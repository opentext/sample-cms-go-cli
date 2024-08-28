package cms

import (
	"fmt"
	"ocp/sample/planets/internal/config"
	ioutil "ocp/sample/planets/internal/util/io"
	jsonutil "ocp/sample/planets/internal/util/json"
	logutil "ocp/sample/planets/internal/util/log"

	"github.com/tidwall/gjson"
)

const (
	PlanetType     = "un_planet"
	PlanetCategory = "object"
)

// Model for the Planet type CMS properties
type PlanetProps struct {
	Diameter        int64   `json:"diameter"`
	LengthOfDay     float64 `json:"length_of_day"`
	NumberOfMoons   *int64  `json:"number_of_moons,integer,omitempty"`
	MeanTemperature *int64  `json:"mean_temperature,integer,omitempty"`
}

// Reads in planet data from the json sample data and creates one instance per object
// Deliberately doesn't populate the "number_of_moons" and "mean_temperature" CMS attributes.
func CreatePlanets() (err error) {
	var planetJSON string

	planetJSON, err = readPlanetData()

	if err == nil {
		gjson.Parse(planetJSON).ForEach(func(_, value gjson.Result) bool {
			instanceBody := &InstanceBody{
				Name: value.Get("name").String(),
				Properties: PlanetProps{
					Diameter:    value.Get("diameter").Int(),
					LengthOfDay: value.Get("length_of_day").Float(),
				},
			}
			var postBody string
			postBody, err = jsonutil.ToJSON(instanceBody)

			if err == nil {
				CreateInstance(PlanetCategory, PlanetType, postBody)
			}

			return true
		})
	}

	return
}

// Fetches the existing planets instance from CMS. Loops through and performs an update on each instance.
// CMS type attributes "number_of_moons" and "mean_temperature" that weren't previously set are set now.
func UpdatePlanets() (err error) {
	var id string
	var planetJSON string
	var instances gjson.Result

	planetJSON, err = readPlanetData()

	if err == nil {
		_, instances, err = InstancesByType(PlanetCategory, PlanetType)
	}

	if err == nil {
		gjson.Parse(planetJSON).ForEach(func(_, value gjson.Result) bool {
			name := value.Get("name").String()
			numMoons := value.Get("number_of_moons").Int()
			meanTemp := value.Get("mean_temperature").Int()
			instanceBody := &InstanceBody{
				Name: name,
				Properties: PlanetProps{
					Diameter:        value.Get("diameter").Int(),
					LengthOfDay:     value.Get("length_of_day").Float(),
					NumberOfMoons:   &numMoons,
					MeanTemperature: &meanTemp,
				},
			}

			instances.ForEach(func(_, instance gjson.Result) bool {
				if instance.Get("name").String() == name {
					id = instance.Get("id").String()
					return false
				} else {
					return true
				}
			})

			var postBody string
			postBody, err = jsonutil.ToJSON(instanceBody)

			if err == nil {
				UpdateInstance(PlanetCategory, PlanetType, postBody, id)
			}

			return true
		})
	}
	return
}

// Deletes all planet instances.
func DeletePlanets() (err error) {
	return DeleteInstancesByType(PlanetCategory, PlanetType)
}

// Fetches all planet instances from CMS and logs out some basic information to the console.
func PlanetInfo() (err error) {
	_, instances, err := InstancesByType(PlanetCategory, PlanetType)

	if err == nil {
		instances.ForEach(func(_, value gjson.Result) bool {
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Id: %s", value.Get("id").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Type: %s", value.Get("type").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Name: %s", value.Get("name").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Diameter (km): %s", value.Get("properties.diameter").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Length of day (hours): %s", value.Get("properties.length_of_day").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Number of moons: %s", value.Get("properties.number_of_moons").String()))
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("Mean temperature (°C): %s", value.Get("properties.mean_temperature").String()))
			logutil.Log(logutil.INFO_LEVEL, "-----------------------------------------------------")
			return true
		})

		if len(instances.Array()) == 0 {
			logutil.Log(logutil.INFO_LEVEL, fmt.Sprintf("No instances of type %s found", PlanetType))
		}
	}
	return
}

// Reads planet JSON data from the sample file
func readPlanetData() (planetJSON string, err error) {
	var sampleDataPath string

	sampleDataPath, err = config.SampleDataPath()

	if err == nil {
		planetJSON, err = ioutil.ReadFileAsString(sampleDataPath)
	}

	return
}
