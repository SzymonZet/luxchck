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
	fullEndpointUrl string
}

var ServiceVariantsGroupsEndpoint serviceVariantsGroupsEndpointType = serviceVariantsGroupsEndpointType{
	fullEndpointUrl: getFullUrl("/PatientPortal/NewPortal/Dictionary/serviceVariantsGroups"),
}

type serviceVariantsGroups []struct {
	Id       int                        `json:"id"`
	Name     string                     `json:"name"`
	Children []serviceVariantsSubGroups `json:"children"`
}

type serviceVariantsSubGroups struct {
	Id       int                           `json:"id"`
	Name     string                        `json:"name"`
	Children []serviceVariantsSubSubGroups `json:"children"`
}

type serviceVariantsSubSubGroups struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Children []any  `json:"children"`
}

func (s serviceVariantsGroupsEndpointType) GetAllRaw() serviceVariantsGroups {
	var output serviceVariantsGroups
	body := invokeRequest(s.fullEndpointUrl, "GET")

	err := json.Unmarshal(body, &output)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))

	return output
}

func (s serviceVariantsGroups) GetFiltered(searchedName string) map[string]int {
	result := make(map[string]int)
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
							result[fmt.Sprintf("%v -> %v -> %v", lvl1.Name, lvl2.Name, name)] = lvl3.Id
						}
					}
				} else {
					//fmt.Println(lvl2.Name)
					if name := lvl2.Name; strings.Contains(strings.ToLower(name), searchedName) {
						//result[name] = lvl2.Id
						result[fmt.Sprintf("%v -> %v", lvl1.Name, name)] = lvl2.Id
					}
				}
			}
		} else {
			if name := lvl1.Name; strings.Contains(strings.ToLower(name), searchedName) {
				result[name] = lvl1.Id
			}
		}
	}

	if len(result) == 0 {
		erroring.QuitIfError(fmt.Errorf("variant not found"), fmt.Sprintf("error when trying to find variant: %v", searchedName))
	}

	log.Println("at least one variant / visit type found successfully")

	return result
}
