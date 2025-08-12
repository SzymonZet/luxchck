package main

import (
	"SzymonZet/LuxmedCheck/credentials"
	"SzymonZet/LuxmedCheck/publish"
	"SzymonZet/LuxmedCheck/server"
	"flag"
	"fmt"
	"log"
)

func main() {

	fmt.Println("==============================")
	fmt.Println("SzymonZet/LuxmedCheck")
	fmt.Println("==============================")
	fmt.Println()

	params := make(map[string]*string)

	params["visitType"] = flag.String("visitType", "", "type of a visit, for example: `konsultacja internistyczna`")
	params["city"] = flag.String("city", "", "city, for example: `warszawa`")
	flag.Parse()

	fmt.Printf("%10v : %-20v\n", "PARAMETER", "VALUE")
	for key, val := range params {
		fmt.Printf("%10v | %-20v\n", key, *val)
	}
	fmt.Println()

	login := credentials.GetSecureString("Luxmed Login")
	pass := credentials.GetSecureString("Luxmed Password")

	fmt.Println()
	fmt.Println("==============================")
	fmt.Println()

	log.Println("main processing started...")

	server.RefreshAuthToken(login, pass)
	groups := server.ServiceVariantsGroupsEndpoint.GetAllServiceVariantsGroupsObjects().GetFilteredVariants(*params["visitType"])
	cities := server.CitiesEndpoint.GetAllCitiesObjects().GetFilteredCities(*params["city"])
	terms := server.TermsEndpoint.GetFilteredTermObjects(cities, groups)

	log.Println("main processing completed successfully")

	// debug
	// fmt.Println(groups)
	// fmt.Println(cities)
	// fmt.Println("Day: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Day)
	// fmt.Println("Clinic: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Terms[0].Clinic)
	// fmt.Println("ClinicGroup: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Terms[0].ClinicGroup)
	// fmt.Println("DateTimeFrom: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Terms[0].DateTimeFrom)
	// fmt.Println("DateTimeTo: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Terms[0].DateTimeTo)
	// fmt.Println("Doctor: ", terms[0].TermsRoot.TermsForService.TermsForDays[0].Terms[0].Doctor)

	publish.StartPublishServer(params, terms)
}
