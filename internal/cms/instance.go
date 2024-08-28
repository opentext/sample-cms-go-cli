package cms

import (
	"fmt"
	"net/http"
	"ocp/sample/planets/internal/config"
	authutil "ocp/sample/planets/internal/util/auth"
	logutil "ocp/sample/planets/internal/util/log"

	"github.com/tidwall/gjson"
)

type InstanceBody struct {
	Name       string      `json:"name"`
	Properties interface{} `json:"properties,omitempty"`
}

const (
	embeddedCollectionKey = "_embedded.collection"
)

// Returns the GET /instances URL for a given category and type.
func InstancesUrl(category string, systemTypeName string) (instancesUrl string, err error) {
	var cmsHost string

	cmsHost, err = config.CMSHost()

	if err == nil {
		instancesUrl = fmt.Sprintf("%s/instances/%s/%s", cmsHost, category, systemTypeName)
	}

	return
}

// Gets instances from CMS for a given category and type.
func InstancesByType(category string, systemTypeName string) (statusCode int, instances gjson.Result, err error) {
	var respBody string
	var instancesUrl string

	instancesUrl, err = InstancesUrl(category, systemTypeName)

	if err == nil {
		statusCode, respBody, err = authutil.DoWithToken(instancesUrl, http.MethodGet)
	}

	if err == nil {
		instances = gjson.Get(respBody, embeddedCollectionKey)
	}

	return
}

// Creates instances in CMS for a given category and type.
func CreateInstance(category string, systemTypeName string, jsonString string) (statusCode int, respBody string, err error) {
	var instancesUrl string

	logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Creating instance of type %s with name: %s", systemTypeName, gjson.Get(jsonString, "name")))

	instancesUrl, err = InstancesUrl(category, systemTypeName)

	if err == nil {
		return authutil.DoWithTokenJSONBody(instancesUrl, http.MethodPost, jsonString)
	}

	return
}

// Updates instances in CMS for a given category, type and id.
func UpdateInstance(category string, systemTypeName string, jsonString string, id string) (statusCode int, respBody string, err error) {
	var instancesUrl string

	logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Updating instance of type %s with name: %s", systemTypeName, gjson.Get(jsonString, "name")))

	instancesUrl, err = InstancesUrl(category, systemTypeName)

	if err == nil {
		return authutil.DoWithTokenJSONBody(fmt.Sprintf("%s/%s", instancesUrl, id), http.MethodPut, jsonString)
	}

	return
}

// Deletes instances from CMS for a given category and type.
// Runs deletes in parallel using channels with automatic retry handling.
func DeleteInstancesByType(category string, systemTypeName string) (err error) {
	var instances gjson.Result
	var statusCode int

	statusCode, instances, err = InstancesByType(category, systemTypeName)

	if statusCode < 400 && err == nil {
		c := make(chan string)

		instances.ForEach(func(_, value gjson.Result) bool {
			id := value.Get("id").String()
			name := value.Get("name").String()
			cmsType := value.Get("type").String()
			deleteUrl := value.Get("_links.urn:eim:linkrel:delete.href").String()

			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Deleting instance of type %s with id: %s and name: %s", cmsType, id, name))

			go deleteInstance(id, cmsType, deleteUrl, c)

			return true
		})

		deleteResponses := make([]string, len(instances.Array()))

		for i := range deleteResponses {
			deleteResponses[i] = <-c
			logutil.Log(logutil.INFO_LEVEL, deleteResponses[i])
		}

		logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Finished deleting instances for type %s", systemTypeName))
	}

	return
}

// Deletes an individual instance from CMS.
func deleteInstance(id string, cmsType string, deleteUrl string, c chan string) {
	var deleteMessage string

	statusCode, _, err := authutil.DoWithTokenAndRetry(deleteUrl, http.MethodDelete)

	if statusCode < 400 && err == nil {
		deleteMessage = fmt.Sprintf("Item with id %s deleted", id)
	} else {
		deleteMessage = fmt.Sprintf("Unable to delete item with id %s", id)
	}

	c <- deleteMessage
}
