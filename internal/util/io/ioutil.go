// The ioutil package provides general reusable utils for http and file interactions.
package ioutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	logutil "ocp/sample/planets/internal/util/log"
	"os"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
)

// Create simple requests with no body
func NewRequest(method string, url string) (req *http.Request, err error) {
	return http.NewRequest(method, url, nil)
}

// Create requests with a JSON body
func NewRequestJSONBody(method string, url string, body string) (req *http.Request, err error) {
	return http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
}

// Wrapper for the native http vs the third-party retry http clients.
// Also reads the HTTP response body into a string.
func Do(req *http.Request, withRetry bool) (statusCode int, respBody string) {
	var resp *http.Response
	var err error

	if withRetry == true {
		resp, err = doWithRetry(req)
	} else {
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	if err == nil {
		respBody, err = readerAsString(resp.Body)
	} else {
		logutil.LogError(err)
	}

	logResponseWithError(req, resp, respBody)

	defer resp.Body.Close()
	return resp.StatusCode, respBody
}

// If we receive an error status code then log the result
func logResponseWithError(req *http.Request, resp *http.Response, respBody string) {
	if resp.StatusCode >= 400 {
		logutil.Log(logutil.ERROR_LEVEL, fmt.Sprintf("Url: %s", req.URL))
		logutil.Log(logutil.ERROR_LEVEL, fmt.Sprintf("HTTP status code: %s", strconv.Itoa(resp.StatusCode)))
		logutil.Log(logutil.ERROR_LEVEL, fmt.Sprintf("Response body: %s", respBody))
	}
}

// Use the go-retryablehttp module to handle the request
func doWithRetry(req *http.Request) (resp *http.Response, err error) {
	var retryableRequest *retryablehttp.Request

	retryableClient := retryablehttp.NewClient()
	retryableRequest, err = retryablehttp.FromRequest(req)

	if err == nil {
		return retryableClient.Do(retryableRequest)
	}

	return
}

// Reads in a file and converts the contents to a string
func ReadFileAsString(path string) (jsonBody string, err error) {
	var fileContents *os.File

	fileContents, err = os.Open(path)

	if err == nil {
		jsonBody, err = readerAsString(fileContents)
	} else {
		logutil.LogError(err)
	}

	defer fileContents.Close()

	return jsonBody, err
}

func readerAsString(reader io.Reader) (respBody string, err error) {
	var bodyBytes []byte

	bodyBytes, err = io.ReadAll(reader)

	if err == nil {
		respBody = string(bodyBytes)
	} else {
		logutil.LogError(err)
	}

	return
}
