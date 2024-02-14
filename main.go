package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type QueryArgs []string

func (a *QueryArgs) String() string {
	var result string
	for index, q := range *a {
		queryFmt := "%s&"
		if index == len(*a)-1 {
			queryFmt = "%s"
		}
		result += fmt.Sprintf(queryFmt, q)
	}
	return result
}
func (i *QueryArgs) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type ContentType struct {
	JSON string
}

var content_type ContentType = ContentType{
	JSON: "application/json",
}

func main() {
	rawURL := flag.String("url", "", "URL to make request to.")
	contentType := flag.String("content-type", "", "Content-Type request header.")
	method := flag.String("method", "", "GET, POST, PUT, PATCH, DELETE, OPTIONS.")
	body := flag.String("body", "", "Request body.")
	var querySlice QueryArgs
	flag.Var(&querySlice, "q", "Request query arguments.")
	flag.Parse()

	if *rawURL == "" {
		log.Fatalln("URL is required.")
	}

	// parse request URI
	parsedURL, err := url.ParseRequestURI(*rawURL)
	if err != nil {
		log.Fatalln("Couldn't parse URL:", *rawURL)
	}

	// add query string
	values := parsedURL.Query()
	for _, qStr := range querySlice {
		qSlice := strings.Split(qStr, "=")
		if len(qSlice) != 2 {
			log.Fatalln("Invalid query string:", qStr)
		}
		values.Add(qSlice[0], qSlice[1])
	}
	if values != nil {
		parsedURL.RawQuery = values.Encode()
	}

	if *body != "" && isJSONContent(*contentType) {
		var temp map[string]interface{}
		if err := json.Unmarshal([]byte(*body), &temp); err != nil {
			log.Fatalln("Couldn't marshal body:", err)
		}
	}

	if *method == "" || !isAllowedMethod(*method) {
		log.Fatalln("Invalid request method:", *method)
	}

	client := &http.Client{}
	request, err := http.NewRequest(*method, parsedURL.String(), bytes.NewBuffer([]byte(*body)))
	if err != nil {
		log.Fatalln("Couldn't create request:", err)
	}

	// check Content-Type
	if *contentType != "" && !isAllowedContentType(*contentType) {
		log.Fatalln("Unallowed content type:", *contentType)
	} else {
		request.Header.Add("Content-Type", *contentType)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln("Request error::", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Couldn't read response body:", err)
	}
	fmt.Println(string(result))
}

var ALLOWED_CONTENT_TYPES = [1]string{"application/json"}

func isAllowedContentType(input string) bool {
	for _, contentType := range ALLOWED_CONTENT_TYPES {
		if contentType == input {
			return true
		}
	}
	return false
}

func isJSONContent(contentType string) bool {
	return contentType == content_type.JSON
}

var METHODS = [6]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

func isAllowedMethod(input string) bool {
	for _, contentType := range METHODS {
		if contentType == input {
			return true
		}
	}
	return false
}
