package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"szymonzet/luxchck/cred"
	"szymonzet/luxchck/lux"
	"szymonzet/luxchck/publish"
)

func main() {

	fmt.Println("==============================")
	fmt.Println("szymonzet/luxchck")
	fmt.Println("==============================")
	fmt.Println()

	params := make(map[string]*string)

	params["visitTypes"] = flag.String("visitTypes", "", "types of a visit (matching phrase); example: `hemat`")
	params["city"] = flag.String("city", "", "city (has to be explicit, only one at the time); example: `wrocław`")
	params["clinics"] = flag.String("clinics", "", "colon-separated list of clinics addresses/streets (matching phrase); optional; examples: `Legnicka;Swobodna`, `legn`, ``")
	params["doctors"] = flag.String("doctors", "", "colon-separated list of doctors (matching phrase); optional; examples: `Jan Kowalski;Jane Doe`, `kow`, ``")
	flag.Parse()

	fmt.Printf("%10v | %-20v\n", "PARAMETER", "VALUE")
	for key, val := range params {
		fmt.Printf("%10v | %-20v\n", key, *val)
	}
	fmt.Println()

	cred.SetLoginAndPassword()
	cred.RefreshHeaderCookie()

	fmt.Println()
	fmt.Println("==============================")
	fmt.Println()

	log.Println("main processing started...")

	variants := lux.ServiceVariantsGroupsEndpoint.GetAllRaw().GetFiltered(*params["visitTypes"])
	cities := lux.CitiesEndpoint.GetAllRaw().GetFiltered(*params["city"])
	clinics := strings.Split(*params["clinics"], ";")
	doctors := strings.Split(*params["doctors"], ";")

	reqCombinationsCount := len(cities) * len(variants)
	log.Printf("cities: %v | variants: %v | total: %v\n", len(cities), len(variants), reqCombinationsCount)
	if reqCombinationsCount > 5 {
		log.Printf("WARN | too many potential request combinations: %v | you may exceed rates / encounter error 429, please consider narrowing down search parameters!", reqCombinationsCount)
	}

	terms := lux.TermsEndpoint.GetAllRaw(cities, variants).FilterAndClean(clinics, doctors)

	log.Println("main processing completed successfully")

	// debug
	//filteredTerms := terms.FilterAndClean(clinics, doctors)
	//fmt.Println(filteredTerms)
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
