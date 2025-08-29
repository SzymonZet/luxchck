package lux

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"szymonzet/luxchck/erroring"
	"time"

	"github.com/google/uuid"
)

const dateTimeLayout string = "2006-01-02T15:04:05"

type termsEndpointType struct {
	urlIndex       string
	urlNextTerms   string
	urlOneDayTerms string
}

var TermsEndpoint termsEndpointType = termsEndpointType{
	urlIndex:       getFullUrl("/PatientPortal/NewPortal/terms/index"),
	urlNextTerms:   getFullUrl("/PatientPortal/NewPortal/terms/nextTerms"),
	urlOneDayTerms: getFullUrl("/PatientPortal/NewPortal/terms/oneDayTerms"),
}

type termsIndexResponse struct {
	CorrelationId   string `json:"correlationId"`
	TermsForService struct {
		TermsForDays     []termsforDays `json:"termsForDays"`
		TermsInfoForDays []struct {
			Day          string `json:"day"`
			TermsStatus  int    `json:"termsStatus"` // 0 => potential visits exist that day
			Message      string `json:"message"`
			TermsCounter struct {
				TermsNumber int `json:"termsNumber"`
			} `json:"termsCounter"`
		} `json:"termsInfoForDays"`
	} `json:"termsForService"`
}

type termsOneDayResponse struct {
	CorrelationId string       `json:"correlationId"`
	TermsForDay   termsforDays `json:"termsForDay"`
}

type termsforDays struct {
	Day   string `json:"day"`
	Terms []struct {
		DateTimeFrom string `json:"dateTimeFrom"`
		DateTimeTo   string `json:"dateTimeTo"`
		Doctor       struct {
			AcademicTitle string `json:"academicTitle"`
			FirstName     string `json:"firstName"`
			LastName      string `json:"lastName"`
		} `json:"doctor"`
		Clinic         string `json:"clinic"`
		ClinicGroup    string `json:"clinicGroup"`
		IsTelemedicine bool   `json:"isTelemedicine"`
	} `json:"terms"`
}

type termsResponseExtended struct {
	City           string
	ServiceVariant string
	TermsResponse  []termsforDays
}

type termsResponsesExtended []termsResponseExtended

type termsTarget struct {
	Title string
	Desc  string
	Terms []TermFlatten
}

type TermFlatten struct {
	Day      string
	TimeFrom string
	TimeTo   string
	Clinic   string
	Doctor   string
}

type TermsTargets []termsTarget

func (t termsEndpointType) GetAllRaw(cities map[string]int, serviceVariants []ServiceVariantsGroupsTarget) termsResponsesExtended {
	var output termsResponsesExtended
	var termsIndexResp termsIndexResponse
	dateFrom := time.Now().Format("2006-01-02")
	dateTo := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	// horribly written; to be re-factored into smaller bits
	for cityName, cityId := range cities {
		for _, variant := range serviceVariants {
			req := createAuthorizedRequest(t.urlIndex, "GET")

			params := map[string]string{
				// actual parameters
				"serviceVariantId": fmt.Sprintf("%v", variant.ChildId),
				"searchPlace.id":   fmt.Sprintf("%v", cityId),
				"searchPlace.name": fmt.Sprintf("%v", cityName),
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

			// if !variant.IsTelemedicine {
			// 	params["searchPlace.id"] = fmt.Sprintf("%v", cityId)
			// 	params["searchPlace.name"] = fmt.Sprintf("%v", cityName)
			// }

			addUrlParametersToRequest(req, params)

			log.Printf("city: %v (%v) | variant: %v (%v) - index - trying to get a response...", cityName, cityId, variant.FullName, variant.ChildId)

			body := getResponse(req)

			err := json.Unmarshal(body, &termsIndexResp)
			erroring.LogIfError(err, fmt.Sprintf("error when trying to unmarshal response from:\n```\n%v\n```\n", string(body)))

			if termsIndexResp.CorrelationId == "" {
				log.Printf("city: %v (%v) | variant: %v (%v) - index - no response or no visits could be read from:\n```\n%v\n```\n", cityName, cityId, variant.FullName, variant.ChildId, string(body))
			} else {
				daysWithVisits := len(termsIndexResp.TermsForService.TermsForDays)
				log.Printf("city: %v (%v) | variant: %v (%v) - index - days with visits read: %v", cityName, cityId, variant.FullName, variant.ChildId, daysWithVisits)

				if daysWithVisits > 0 {
					newEntry := termsResponseExtended{
						City:           fmt.Sprintf("%v (%v)", cityName, cityId),
						ServiceVariant: fmt.Sprintf("%v (%v)", variant.FullName, variant.ChildId),
						TermsResponse:  termsIndexResp.TermsForService.TermsForDays,
					}

					output = append(output, newEntry)
				} else {
					for _, info := range termsIndexResp.TermsForService.TermsInfoForDays {
						if info.TermsCounter.TermsNumber > 0 {
							day := extractDate(info.Day)
							var termsOneDayResp termsOneDayResponse
							req = createAuthorizedRequest(t.urlOneDayTerms, "GET")
							params["searchDateFrom"] = day
							params["searchDateTo"] = day
							params["expectedTermsNumber"] = fmt.Sprint(info.TermsCounter.TermsNumber)

							// delete(params, "serviceVariantSource")
							// delete(params, "nextSearch")
							// delete(params, "locationReplaced")
							// delete(params, "searchDatePreset")
							// if variant.IsTelemedicine {
							// 	delete(params, "searchPlace.id")
							// 	delete(params, "searchPlace.name")
							// }

							addUrlParametersToRequest(req, params)

							log.Printf("city: %v (%v) | variant: %v (%v) - oneDayTerms - %v - trying to get a response...", cityName, cityId, variant.FullName, variant.ChildId, day)

							body := getResponse(req)

							err = json.Unmarshal(body, &termsOneDayResp)

							// debug
							fmt.Println()
							// fmt.Println(req.URL)
							// fmt.Println(string(body))
							fmt.Println(termsOneDayResp)
							fmt.Println()

							if termsIndexResp.CorrelationId == "" {
								log.Printf("city: %v (%v) | variant: %v (%v) - oneDayTerms - %v - no response or no visits could be read from:\n```\n%v\n```\n", cityName, cityId, variant.FullName, variant.ChildId, string(body), day)
							} else {
								daysWithVisits := len(termsIndexResp.TermsForService.TermsForDays)
								log.Printf("city: %v (%v) | variant: %v (%v) - oneDayTerms - %v - days with visits read: %v", cityName, cityId, variant.FullName, variant.ChildId, daysWithVisits, day)

								if daysWithVisits > 0 {
									newEntry := termsResponseExtended{
										City:           fmt.Sprintf("%v (%v)", cityName, cityId),
										ServiceVariant: fmt.Sprintf("%v (%v)", variant.FullName, variant.ChildId),
										TermsResponse:  termsIndexResp.TermsForService.TermsForDays,
									}

									output = append(output, newEntry)
								}
							}

						}
					}
				}
			}
		}
	}

	log.Println("list of terms/visits obtained successfully")

	return output
}

func (t termsResponsesExtended) FilterAndClean(clinics []string, doctors []string) TermsTargets {
	var output TermsTargets
	for _, termRootMultiple := range t {
		termsTarget := termsTarget{
			Title: fmt.Sprintf("%v | %v", termRootMultiple.City, termRootMultiple.ServiceVariant),
		}
		for _, termForDay := range termRootMultiple.TermsResponse {
			for _, term := range termForDay.Terms {
				var isClinic, isDoctor bool

				// this is suboptimal, as it doesn't really make sense to process the result BEFORE filtering
				// kept like this for now, for accurate logging (to be changed once properly tested / confidence gained)
				fullDoctorName := fmt.Sprintf("%v %v %v", term.Doctor.AcademicTitle, term.Doctor.FirstName, term.Doctor.LastName)
				new := TermFlatten{
					Day:      extractDate(termForDay.Day),
					TimeFrom: extractTime(term.DateTimeFrom),
					TimeTo:   extractTime(term.DateTimeTo),
					Clinic:   term.Clinic,
					Doctor:   fullDoctorName,
				}

				if term.IsTelemedicine {
					isClinic = true
				} else {
					for _, clinic := range clinics {
						if strings.Contains(strings.ToLower(term.Clinic), strings.ToLower(clinic)) {
							isClinic = true
							break
						}
					}
				}

				if isClinic || len(clinics) == 0 {
					for _, doctor := range doctors {
						if strings.Contains(strings.ToLower(fullDoctorName), strings.ToLower(doctor)) {
							isDoctor = true
							break
						}
					}
				}

				if isDoctor || len(doctors) == 0 {
					termsTarget.Terms = append(termsTarget.Terms, new)
					log.Printf("term %v met the filter conditions", new)
				} else {
					log.Printf("term %v did NOT meet filter conditions (clinics: %v | doctors: %v)", new, isClinic, isDoctor)
				}

			}
		}

		output = append(output, termsTarget)
	}

	return output
}

func extractTime(dateTime string) string {
	parsed, err := time.Parse(dateTimeLayout, dateTime)
	erroring.LogIfError(err, fmt.Sprintf("error when trying parse dateTime:%v", dateTime))
	return parsed.Format("15:04:05")
}

func extractDate(dateTime string) string {
	parsed, err := time.Parse(dateTimeLayout, dateTime)
	erroring.LogIfError(err, fmt.Sprintf("error when trying parse dateTime:%v", dateTime))
	return parsed.Format("2006-01-02")
}
