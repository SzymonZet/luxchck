package lux

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"szymonzet/luxchck/erroring"
)

type serviceVariantsGroupsEndpointType struct {
	url string
}

var ServiceVariantsGroupsEndpoint serviceVariantsGroupsEndpointType = serviceVariantsGroupsEndpointType{
	url: getFullUrl("/PatientPortal/NewPortal/Dictionary/serviceVariantsGroups"),
}

type serviceVariantsGroupsResponse []struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	IsTelemedicine bool   `json:"isTelemedicine"`
	Children       []struct {
		Id             int    `json:"id"`
		Name           string `json:"name"`
		IsTelemedicine bool   `json:"isTelemedicine"`
		Children       []struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			Children       []any  `json:"children"`
			IsTelemedicine bool   `json:"isTelemedicine"`
		} `json:"children"`
	} `json:"children"`
}

type ServiceVariantsGroupsTarget struct {
	ChildId        int
	FullName       string
	IsTelemedicine bool
}

func (s serviceVariantsGroupsEndpointType) GetAllRaw() serviceVariantsGroupsResponse {
	var output serviceVariantsGroupsResponse
	body := invokeRequest(s.url, "GET")

	err := json.Unmarshal(body, &output)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))

	return output
}

func (s serviceVariantsGroupsResponse) GetFiltered(searchedName string) []ServiceVariantsGroupsTarget {
	var result []ServiceVariantsGroupsTarget
	searchedName = strings.ToLower(searchedName)

	// todo: probably can be done via recurrence (?)
	for _, lvl1 := range s {
		if len(lvl1.Children) != 0 {
			for _, lvl2 := range lvl1.Children {
				if len(lvl2.Children) != 0 {
					for _, lvl3 := range lvl2.Children {
						if len(lvl3.Children) != 0 {
							erroring.LogIfError(errors.New("more serviceVariantsGroups levels than predicted"), "error when trying to extract serviceVariantsGroups")
						}
						if name := lvl3.Name; strings.Contains(strings.ToLower(name), searchedName) {
							//result[fmt.Sprintf("%v -> %v -> %v", lvl1.Name, lvl2.Name, name)] = lvl3.Id
							newEntry := ServiceVariantsGroupsTarget{
								ChildId:        lvl3.Id,
								FullName:       fmt.Sprintf("%v -> %v -> %v", lvl1.Name, lvl2.Name, name),
								IsTelemedicine: lvl3.IsTelemedicine,
							}
							result = append(result, newEntry)
						}
					}
				} else {
					//fmt.Println(lvl2.Name)
					if name := lvl2.Name; strings.Contains(strings.ToLower(name), searchedName) {
						//result[fmt.Sprintf("%v -> %v", lvl1.Name, name)] = lvl2.Id
						newEntry := ServiceVariantsGroupsTarget{
							ChildId:        lvl2.Id,
							FullName:       fmt.Sprintf("%v -> %v", lvl1.Name, name),
							IsTelemedicine: lvl2.IsTelemedicine,
						}
						result = append(result, newEntry)
					}
				}
			}
		} else {
			if name := lvl1.Name; strings.Contains(strings.ToLower(name), searchedName) {
				newEntry := ServiceVariantsGroupsTarget{
					ChildId:        lvl1.Id,
					FullName:       fmt.Sprintf("%v", name),
					IsTelemedicine: lvl1.IsTelemedicine,
				}
				result = append(result, newEntry)
			}
		}
	}

	if len(result) == 0 {
		erroring.QuitIfError(fmt.Errorf("variant not found"), fmt.Sprintf("error when trying to find variant: %v", searchedName))
	}

	log.Println("at least one variant / visit type found successfully")

	return result
}
