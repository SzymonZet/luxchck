package lux

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"szymonzet/luxchck/erroring"
)

const (
	baseUrl string = "https://portalpacjenta.luxmed.pl/"
)

func invokeRequest(url string, requestType string) []byte {
	req := createAuthorizedRequest(url, requestType)
	return getResponse(req)
}

func createAuthorizedRequest(url string, requestType string) *http.Request {
	req, err := http.NewRequest(requestType, url, nil)
	erroring.LogIfError(err, fmt.Sprintf("error when creating %v request for: \n```\n%v\n```\n", requestType, url))
	req.Header.Add("Cookie", CurrentHeaderAuthString)

	return req
}

func getResponse(req *http.Request) []byte {
	resp, err := http.DefaultClient.Do(req)
	erroring.LogIfError(err, fmt.Sprintf("error when invoking request for: \n```\n%v\n```\n", req.URL))

	defer resp.Body.Close()

	if st := resp.StatusCode; st != http.StatusOK {
		erroring.LogIfError(
			fmt.Errorf("response code was %v, expected %v", st, http.StatusOK),
			fmt.Sprintf("problem with response after invoking request for: \n```\n%v\n```\n", req.URL),
		)
	}

	body, err := io.ReadAll(resp.Body)
	erroring.LogIfError(err, fmt.Sprintf("error when reading response from request for: \n```\n%v\n```\n", req.URL))

	return body
}

func addUrlParametersToRequest(req *http.Request, params map[string]string) {
	parametrizedUrl := req.URL.Query()
	for key, val := range params {
		parametrizedUrl.Add(key, val)
	}
	req.URL.RawQuery = parametrizedUrl.Encode()
}

func getFullUrl(endpoint string) string {
	output, err := url.JoinPath(baseUrl, endpoint)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	_, err = url.ParseRequestURI(output)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	return output
}
