package server

import (
	"SzymonZet/LuxmedCheck/connection"
	"SzymonZet/LuxmedCheck/erroring"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type citiesEndpointType struct {
	fullEndpointUrl string
}

var CitiesEndpoint citiesEndpointType = citiesEndpointType{
	fullEndpointUrl: connection.GetFullUrl("/PatientPortal/NewPortal/Dictionary/cities"),
}

type cities []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c citiesEndpointType) GetAllCitiesObjects() cities {
	var output cities
	body := invokeRequest(c.fullEndpointUrl, "GET")

	err := json.Unmarshal(body, &output)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))

	return output
}

func (c cities) GetFilteredCities(searchedName string) map[string]int {
	result := make(map[string]int)
	searchedName = strings.ToLower(searchedName)
	for _, val := range c {
		if strings.Contains(strings.ToLower(val.Name), searchedName) {
			result[val.Name] = val.Id
		}
	}

	if len(result) == 0 {
		erroring.QuitIfError(fmt.Errorf("city not found"), fmt.Sprintf("error when trying to find city: %v", searchedName))
	}

	log.Println("at least one city found successfully")

	return result
}
