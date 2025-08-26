package lux

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"szymonzet/luxchck/erroring"
)

type citiesEndpointType struct {
	url string
}

var CitiesEndpoint citiesEndpointType = citiesEndpointType{
	url: getFullUrl("/PatientPortal/NewPortal/Dictionary/cities"),
}

type citiesResponse []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c citiesEndpointType) GetAllRaw() citiesResponse {
	var output citiesResponse
	body := invokeRequest(c.url, "GET")

	err := json.Unmarshal(body, &output)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))

	return output
}

func (c citiesResponse) GetFiltered(searchedName string) map[string]int {
	result := make(map[string]int)
	searchedName = strings.ToLower(searchedName)
	for _, val := range c {
		if strings.ToLower(val.Name) == searchedName {
			result[val.Name] = val.Id
		}
	}

	if len(result) == 0 {
		erroring.QuitIfError(fmt.Errorf("city not found"), fmt.Sprintf("error when trying to find city: %v", searchedName))
	}

	log.Println("at least one city found successfully")

	return result
}
