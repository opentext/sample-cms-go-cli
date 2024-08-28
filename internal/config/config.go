// The config package provides a central place to get information from the environment
package config

import (
	"errors"
	"fmt"
	"net/url"
	logutil "ocp/sample/planets/internal/util/log"
	"os"
)

const (
	VAR_BASE_URL         = "CMS_DEMO_BASE_URL"
	VAR_TENANT_ID        = "CMS_DEMO_TENANT_ID"
	VAR_CONF_CLIENT_ID   = "CMS_DEMO_CONF_CLIENT_ID"
	VAR_CLIENT_SECRET    = "CMS_DEMO_CLIENT_SECRET"
	VAR_SAMPLE_DATA_PATH = "CMS_DEMO_SAMPLE_DATA_PATH"
)

// The base url for the OCP environment
func BaseUrl() (baseUrl string, err error) {
	baseUrl, err = envVar(VAR_BASE_URL)

	if err == nil {
		_, err = url.ParseRequestURI(baseUrl)
		if err != nil {
			logutil.Log(logutil.ERROR_LEVEL, "CMS_DEMO_BASE_URL environment variable is not a valid url")
		}
	}

	return
}

// The CMS host for the OCP environment
func CMSHost() (cmsHost string, err error) {
	var baseUrl string

	baseUrl, err = BaseUrl()

	if err == nil {
		cmsHost = fmt.Sprintf("%s/cms", baseUrl)
	}

	return
}

// The Tenant ID for the OCP environment
func TenantId() (tenantId string, err error) {
	return envVar(VAR_TENANT_ID)
}

// The Confidential Client ID for the OCP environment
func ConfClientId() (confClientId string, err error) {
	return envVar(VAR_CONF_CLIENT_ID)
}

// The Client Secret for the OCP environment
func ClientSecret() (clientSecret string, err error) {
	return envVar(VAR_CLIENT_SECRET)
}

// The path to the sample planet data
func SampleDataPath() (sampleDataPath string, err error) {
	sampleDataPath, err = envVar(VAR_SAMPLE_DATA_PATH)
	if err == nil {
		_, err = os.Stat(sampleDataPath)
		if err != nil {
			logutil.Log(logutil.ERROR_LEVEL, fmt.Sprintf("Sample data file is not present at %s", sampleDataPath))
		}
	}
	return
}

// Gets an environment variable and returns an error if no value is set.
func envVar(key string) (val string, err error) {
	val = os.Getenv(key)
	if len(val) == 0 {
		errStr := fmt.Sprintf("%s environment variable is missing.", key)
		err = errors.New(errStr)
		logutil.Log(logutil.ERROR_LEVEL, errStr)
	}
	return
}
