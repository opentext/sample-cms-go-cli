// The authutil package handles authentication with OCP
package authutil

import (
	"errors"
	"fmt"
	"net/http"
	"ocp/sample/planets/internal/config"
	ioutil "ocp/sample/planets/internal/util/io"
	logutil "ocp/sample/planets/internal/util/log"
	"strings"

	"github.com/tidwall/gjson"
)

var cachedAccessToken string

// Gets the authentication host from the environment.
func AuthHost() (authHost string, err error) {
	var baseUrl string

	baseUrl, err = config.BaseUrl()

	return fmt.Sprintf("%s/tenants", baseUrl), err
}

// Gets the authentication url from the environment.
func AuthUrl() (url string, err error) {
	var authHost string
	var tenantId string

	authHost, err = AuthHost()

	if err == nil {
		tenantId, err = config.TenantId()
	}

	if err == nil {
		return fmt.Sprintf("%s/%s/oauth2/token", authHost, tenantId), err
	}

	return
}

// Performs an HTTP request with the OCP authentication token.
func DoWithToken(url string, method string) (statusCode int, respBody string, err error) {
	var req *http.Request

	req, err = ioutil.NewRequest(method, url)

	if err == nil {
		err = AddAuthHeader(req)
	}

	if err == nil {
		statusCode, respBody = ioutil.Do(req, false)
	}

	return
}

// Performs an HTTP request with the OCP authentication token and automatically handles retries on failure.
func DoWithTokenAndRetry(url string, method string) (statusCode int, respBody string, err error) {
	var req *http.Request

	req, err = ioutil.NewRequest(method, url)

	if err == nil {
		err = AddAuthHeader(req)
	}

	if err == nil {
		statusCode, respBody = ioutil.Do(req, true)
	}

	return
}

// Performs an HTTP request with the OCP authentication token that requires sending a JSON body.
func DoWithTokenJSONBody(url string, method string, body string) (statusCode int, respBody string, err error) {
	var req *http.Request

	req, err = ioutil.NewRequestJSONBody(method, url, body)

	if err == nil {
		setContentType(req)
		err = AddAuthHeader(req)
	} else {
		logutil.LogError(err)
	}

	if err == nil {
		statusCode, respBody = ioutil.Do(req, false)
	}

	return
}

// Adds the authentication token to CMS requests
func AddAuthHeader(req *http.Request) (err error) {
	var accessToken string

	accessToken, err = AuthToken()

	if err == nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	return
}

// Fetches the authentication token using the configuration store in the environment.
func AuthToken() (accessToken string, err error) {
	if len(cachedAccessToken) > 0 {
		accessToken = cachedAccessToken
	} else {
		var req *http.Request
		var authUrl string
		var authBody string

		if err == nil {
			authUrl, authBody, err = authConfig()
		}

		if err == nil {
			req, err = ioutil.NewRequestJSONBody(http.MethodPost, authUrl, authBody)
		}

		if err == nil {
			accessToken, err = fetchAuthToken(req)
		}
	}

	return
}

func setContentType(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func InvalidateTokenCache() {
	cachedAccessToken = ""
}

// Gets the authentication URL and the auth body to fetch the token
func authConfig() (authUrl string, body string, err error) {
	authUrl, err = AuthUrl()

	if err == nil {
		body, err = authBody()
	}

	if strings.Contains(authUrl, "replace") || strings.Contains(body, "replace") {
		err = errors.New("You need to add your tenant id, client id and client secret to your environment.")
		logutil.LogError(err)
	}

	return
}

// Gets the client ID and secret.
func clientConfig() (confClientId string, clientSecret string, err error) {
	confClientId, err = config.ConfClientId()

	if err == nil {
		clientSecret, err = config.ClientSecret()
	}

	return
}

// Generates the auth body to fetch the token.
func authBody() (authBody string, err error) {
	var confClientId string
	var clientSecret string

	confClientId, clientSecret, err = clientConfig()

	if err == nil {
		authBody = fmt.Sprintf(`{
			"client_id": "%s",
			"client_secret": "%s",
			"grant_type": "client_credentials"
		  }`, confClientId, clientSecret)
	}

	return
}

// Fetches an auth token from OCP
func fetchAuthToken(req *http.Request) (accessToken string, err error) {
	setContentType(req)

	logutil.Log(logutil.INFO_LEVEL, "Fetching access token")
	statusCode, responseBody := ioutil.Do(req, false)
	if statusCode < 400 {
		logutil.Log(logutil.INFO_LEVEL, "Access token fetched successfully")
		accessToken = gjson.Get(string(responseBody), "access_token").String()
		cachedAccessToken = accessToken
	} else {
		err = errors.New("Failed to fetch access token")
		logutil.LogError(err)
	}

	return
}
