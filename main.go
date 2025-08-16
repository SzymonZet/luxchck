package main

import (
	"flag"
	"fmt"
	"log"
	"szymonzet/luxchck/credentials"
	"szymonzet/luxchck/lux"
	"szymonzet/luxchck/publish"
)

func main() {

	fmt.Println("==============================")
	fmt.Println("szymonzet/luxchck")
	fmt.Println("==============================")
	fmt.Println()

	params := make(map[string]*string)

	params["visitType"] = flag.String("visitType", "", "type of a visit, for example: `konsultacja internistyczna`")
	params["city"] = flag.String("city", "", "city, for example: `warszawa`")
	flag.Parse()

	fmt.Printf("%10v | %-20v\n", "PARAMETER", "VALUE")
	for key, val := range params {
		fmt.Printf("%10v | %-20v\n", key, *val)
	}
	fmt.Println()

	login := credentials.GetSecureString("Login")
	pass := credentials.GetSecureString("Password")

	fmt.Println()
	fmt.Println("==============================")
	fmt.Println()

	log.Println("main processing started...")

	lux.RefreshAuthToken(login, pass)
	groups := lux.ServiceVariantsGroupsEndpoint.GetAllServiceVariantsGroupsObjects().GetFilteredVariants(*params["visitType"])
	cities := lux.CitiesEndpoint.GetAllCitiesObjects().GetFilteredCities(*params["city"])
	terms := lux.TermsEndpoint.GetFilteredTermObjects(cities, groups)

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
