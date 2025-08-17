package lux

import (
	"encoding/json"
	"fmt"
	"log"
	"szymonzet/luxchck/erroring"
	"time"

	"github.com/google/uuid"
)

type termsEndpointType struct {
	fullEndpointUrl string
}

var TermsEndpoint termsEndpointType = termsEndpointType{
	fullEndpointUrl: getFullUrl("/PatientPortal/NewPortal/terms/index"),
}

type TermsRoot struct {
	CorrelationId   string          `json:"correlationId"`
	TermsForService termsForService `json:"termsForService"`
}

type termsForService struct {
	TermsForDays []termsForDays `json:"termsForDays"`
}

type termsForDays struct {
	Day   string  `json:"day"`
	Terms []terms `json:"terms"`
}

type terms struct {
	DateTimeFrom string `json:"dateTimeFrom"`
	DateTimeTo   string `json:"dateTimeTo"`
	Doctor       doctor `json:"doctor"`
	Clinic       string `json:"clinic"`
	ClinicGroup  string `json:"clinicGroup"`
}

type doctor struct {
	AcademicTitle string `json:"academicTitle"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}

type TermsRootMultiple struct {
	City           string
	ServiceVariant string
	TermsRoot      TermsRoot
}

func (t termsEndpointType) GetAllRaw(cities map[string]int, serviceVariants map[string]int) []TermsRootMultiple {
	var output []TermsRootMultiple
	var termsRoot TermsRoot
	dateFrom := time.Now().Format("2006-01-02")
	dateTo := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	for cityName, cityId := range cities {
		for variantName, variantId := range serviceVariants {
			req := createAuthorizedRequest(t.fullEndpointUrl, "GET")

			params := map[string]string{
				// actual parameters
				"searchPlace.id":   fmt.Sprintf("%v", cityId),
				"searchPlace.name": fmt.Sprintf("%v", cityName),
				"serviceVariantId": fmt.Sprintf("%v", variantId),
				"searchDateFrom":   fmt.Sprintf("%v", dateFrom),
				"searchDateTo":     fmt.Sprintf("%v", dateTo),

				// hardcoded values based on actual calls (not sure what they mean)
				"searchPlace.type":          "0",
				"languageId":                "10",
				"searchDatePreset":          "7",
				"processId":                 uuid.New().String(),
				"nextSearch":                "false",
				"searchByMedicalSpecialist": "false",
				"serviceVariantSource":      "2",
				"locationReplaced":          "false",
				"delocalized":               "false",
			}

			addUrlParametersToRequest(req, params)

			log.Printf("city: %v (%v) | variant: %v (%v) - trying to get a response...", cityName, cityId, variantName, variantId)

			body := getResponse(req)

			err := json.Unmarshal(body, &termsRoot)
			erroring.LogIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n```\n%v\n```\n", string(body)))

			if termsRoot.CorrelationId == "" {
				log.Printf("city: %v (%v) | variant: %v (%v) - no response or no visits could be read from:\n```\n%v\n```\n", cityName, cityId, variantName, variantId, string(body))
			} else {
				daysWithVisits := len(termsRoot.TermsForService.TermsForDays)
				log.Printf("city: %v (%v) | variant: %v (%v) - days with visits read: %v", cityName, cityId, variantName, variantId, daysWithVisits)

				if daysWithVisits > 0 {
					newEntry := TermsRootMultiple{
						City:           fmt.Sprintf("%v (%v)", cityName, cityId),
						ServiceVariant: fmt.Sprintf("%v (%v)", variantName, variantId),
						TermsRoot:      termsRoot,
					}

					output = append(output, newEntry)
				}
			}
		}
	}

	log.Println("list of terms/visits obtained successfully")

	return output
}
