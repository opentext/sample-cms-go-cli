// The jsonutil package provides general reusable utils for working with JSON data.
package jsonutil

import (
	"encoding/json"
	logutil "ocp/sample/planets/internal/util/log"
)

func ToJSON(v any) (jsonString string, err error) {
	vBytes, err := json.Marshal(v)

	if err != nil {
		logutil.LogError(err)
	}

	return string(vBytes), err
}
