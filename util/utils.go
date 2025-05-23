package util

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func GenerateOrderNumber(minLen, maxLen int) string {
	rand.NewSource(time.Now().UnixNano())
	length := rand.Intn(maxLen-minLen+1) + minLen
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var orderNumber string
	for i := 0; i < length; i++ {
		orderNumber += string(chars[rand.Intn(len(chars))])
	}
	return orderNumber
}

func FindRequestId(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "input" {
		inputType, id, name, value := getInputAttributes(n)
		if inputType == "hidden" && id == "request_id" && name == "request_id" {
			return value
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if value := FindRequestId(c); value != "" {
			return value
		}
	}

	return ""
}
func FindPaRes(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "input" {
		inputType, _, name, value := getInputAttributes(n)
		if inputType == "hidden" && name == "PaRes" {
			return value
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if value := FindPaRes(c); value != "" {
			return value
		}
	}

	return ""
}

// getInputAttributes extracts relevant attributes from an <input> element.
func getInputAttributes(n *html.Node) (inputType, id, name, value string) {
	for _, attr := range n.Attr {
		switch attr.Key {
		case "type":
			inputType = attr.Val
		case "id":
			id = attr.Val
		case "name":
			name = attr.Val
		case "value":
			value = attr.Val
		}
	}
	return
}

// StructToURLParams converts a struct into URL query parameters
func StructToURLParams(data interface{}) string {
	var queryParams []string

	// Get the type and value of the struct
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Get the JSON tag of the field (or use field name if no JSON tag)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		} else {
			// Strip ",omitempty" from the JSON tag
			tag = strings.Split(tag, ",")[0]
		}
		if tag == "api_client" {
			continue
		}
		// Check if the field value is non-empty
		if !isEmpty(value) {
			// Format the query parameter and add it to the slice
			encodedValue := url.QueryEscape(fmt.Sprintf("%v", value))
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", tag, encodedValue))
		}
	}

	// Join all query parameters with "&" to form the query string
	queryString := strings.Join(queryParams, "&")
	return queryString
}

// isEmpty checks if a value is nil or an empty string
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	default:
		return false
	}
}
