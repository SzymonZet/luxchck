package publish

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"szymonzet/luxchck/lux"
)

func StartPublishServer(params map[string]*string, terms []lux.TermsRootMultiple) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", generateHtml(params, terms))
	})
	url := "127.0.0.1:8090"
	log.Printf("server started on: %v\n", url)
	http.ListenAndServe(url, nil)
}

func generateHtml(params map[string]*string, terms []lux.TermsRootMultiple) string {
	var htmlBuilder strings.Builder

	htmlBuilder.WriteString(`<html><head><title>luxchck - Terms</title><style>table, tr, td{border: 1px solid black} tr, td{padding: 3px}</style></head><body>`)

	htmlBuilder.WriteString("<h1>Defined search parameters</h1><ul>")
	for key, val := range params {
		htmlBuilder.WriteString(fmt.Sprintf("<li><b>%v</b>: <u>%v</u></li>", key, *val))
	}
	htmlBuilder.WriteString("</ul>")

	for _, term := range terms {
		htmlBuilder.WriteString(fmt.Sprintf(`<h1>%v | %v</h1>`, term.City, term.ServiceVariant))
		htmlBuilder.WriteString(generateTermsTableHtml(term.TermsRoot))
	}

	htmlBuilder.WriteString(`</body></html>`)

	return htmlBuilder.String()
}

func generateTermsTableHtml(termsRoot lux.TermsRoot) string {
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<table>")

	// header
	htmlBuilder.WriteString(
		fmt.Sprintf(
			"<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>",
			"<b>Day</b>",
			"<b>TimeFrom</b>",
			"<b>TimeTo</b>",
			"<b>Clinic</b>",
			"<b>Doctor</b>",
		),
	)

	for _, termDay := range termsRoot.TermsForService.TermsForDays {
		for _, term := range termDay.Terms {
			htmlBuilder.WriteString(
				fmt.Sprintf(
					"<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v %v %v</td></tr>",
					termDay.Day,
					term.DateTimeFrom,
					term.DateTimeTo,
					//term.Clinic,
					term.ClinicGroup,
					term.Doctor.AcademicTitle,
					term.Doctor.FirstName,
					term.Doctor.LastName,
				),
			)
		}
	}

	htmlBuilder.WriteString("</table>")

	return htmlBuilder.String()
}
