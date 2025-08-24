package publish

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"szymonzet/luxchck/lux"
)

func StartPublishServer(params map[string]*string, terms lux.TermsTargets) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", generateHtml(params, terms))
	})
	url := "127.0.0.1:8090"
	log.Printf("server started on: %v\n", url)
	http.ListenAndServe(url, nil)
}

func generateHtml(params map[string]*string, terms lux.TermsTargets) string {
	var htmlBuilder strings.Builder

	htmlBuilder.WriteString(`<html><head><title>luxchck - Terms</title><style>table, tr, td{border: 1px solid black} tr, td{padding: 3px}</style></head><body>`)

	htmlBuilder.WriteString("<h1>Defined search parameters</h1><ul>")
	for key, val := range params {
		htmlBuilder.WriteString(fmt.Sprintf("<li><b>%v</b>: <u>%v</u></li>", key, *val))
	}
	htmlBuilder.WriteString("</ul>")

	for _, term := range terms {
		htmlBuilder.WriteString(fmt.Sprintf(`<h1>%v</h1>`, term.Title))
		htmlBuilder.WriteString(generateTermsTableHtml(term.Terms))
	}

	htmlBuilder.WriteString(`</body></html>`)

	return htmlBuilder.String()
}

func generateTermsTableHtml(termsFlatten []lux.TermFlatten) string {
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

	for _, term := range termsFlatten {
		htmlBuilder.WriteString(
			fmt.Sprintf(
				"<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>",
				term.Day,
				term.TimeFrom,
				term.TimeTo,
				term.Clinic,
				term.Doctor,
			),
		)
	}

	htmlBuilder.WriteString("</table>")

	return htmlBuilder.String()
}
