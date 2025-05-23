package util

import (
	"errors"
	"io"
	"net/http"
)

func Post(fullURL string, formBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", fullURL, formBody)
	if err != nil {
		return nil, errors.Join(err, errors.New("error creating request"))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error executing request"))
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(err, errors.New("error reading response body"))
	}

	return body, nil
}

func Get(fullURL string) ([]byte, error) {
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, errors.New("error executing GET request")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	return bodyBytes, nil
}
