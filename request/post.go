package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func Post(fullURL string, formBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", fullURL, formBody)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error executing request")
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	return body, nil
}
