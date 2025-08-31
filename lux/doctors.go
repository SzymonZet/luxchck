package lux

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"szymonzet/luxchck/erroring"
)

type doctorsEndpointType struct {
	doctorsUrl              string
	facilitiesAndDoctorsUrl string
}

var DoctorsEndpoint doctorsEndpointType = doctorsEndpointType{
	doctorsUrl:              getFullUrl("/PatientPortal/NewPortal/Dictionary/Doctors"),
	facilitiesAndDoctorsUrl: getFullUrl("/PatientPortal/NewPortal/Dictionary/facilitiesAndDoctors"),
}

type facilitiesAndDoctorsResponse struct {
	Doctors doctorsResponse `json:"doctors"`
}

type doctorsResponse []struct {
	Id            int    `json:"id"`
	AcademicTitle string `json:"academicTitle"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}

func (d doctorsEndpointType) GetAllDoctorsMap(serviceVariants []ServiceVariantsGroupsTarget, cities map[string]int) map[int]string {
	output := make(map[int]string)
	for _, variant := range serviceVariants {
		log.Printf("variant: %v (%v) - getting all the doctors...\n", variant.FullName, variant.ChildId)
		newDoctors := d.getAllRaw(variant, cities).asTargetMap()
		log.Printf("variant: %v (%v) - found %v doctors", variant.FullName, variant.ChildId, len(newDoctors))
		maps.Copy(output, newDoctors)
	}

	log.Printf("overall %v doctors", len(output))

	return output
}

func (d doctorsEndpointType) getAllRaw(serviceVariant ServiceVariantsGroupsTarget, cities map[string]int) doctorsResponse {

	params := map[string]string{
		"serviceVariantId": fmt.Sprint(serviceVariant.ChildId),
		"visitLanguageId":  "10",
	}

	getDoctorsResponseFromUrl := func(url string) []byte {
		req := createAuthorizedRequest(url, "GET")
		addUrlParametersToRequest(req, params)
		return getResponse(req)
	}

	var output doctorsResponse
	body := getDoctorsResponseFromUrl(d.doctorsUrl)
	err := json.Unmarshal(body, &output)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))

	if len(output) == 0 {
		var fadResponse facilitiesAndDoctorsResponse
		for _, cityId := range cities {
			params["cityId"] = fmt.Sprint(cityId)
			body := getDoctorsResponseFromUrl(d.facilitiesAndDoctorsUrl)
			err := json.Unmarshal(body, &fadResponse)
			erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n%v", string(body)))
			output = append(output, fadResponse.Doctors...)
		}
	}

	return output
}

func (d doctorsResponse) asTargetMap() map[int]string {
	output := make(map[int]string)

	for _, doctor := range d {
		output[doctor.Id] = fmt.Sprintf("%v %v %v", doctor.AcademicTitle, doctor.FirstName, doctor.LastName)
	}

	return output
}
