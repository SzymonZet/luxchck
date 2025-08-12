package connection

import (
	"SzymonZet/LuxmedCheck/erroring"
	"net/url"
)

const (
	baseUrl string = "https://portalpacjenta.luxmed.pl/"
)

func GetFullUrl(endpoint string) string {
	output, err := url.JoinPath(baseUrl, endpoint)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	_, err = url.ParseRequestURI(output)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	return output
}
