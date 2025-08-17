package lux

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"szymonzet/luxchck/cred"
	"szymonzet/luxchck/erroring"
	"time"
)

const (
	baseUrl string = "https://portalpacjenta.luxmed.pl/"
)

func getFullUrl(endpoint string) string {
	output, err := url.JoinPath(baseUrl, endpoint)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	_, err = url.ParseRequestURI(output)
	erroring.QuitIfError(err, "error when trying to parse %v url")
	return output
}

func invokeRequest(url string, requestType string) []byte {
	req := createAuthorizedRequest(url, requestType)
	return getResponse(req)
}

func createAuthorizedRequest(url string, requestType string) *http.Request {
	req, err := http.NewRequest(requestType, url, nil)
	erroring.LogIfError(err, fmt.Sprintf("error when creating %v request for: \n```\n%v\n```\n", requestType, url))
	addHeaderCookie(req)

	return req
}

func getResponse(req *http.Request) []byte {
	var output []byte
	maxAttempts := 3
	timeoutOnTooManyRequestsSecs := 60
	timeoutBeforeRequest := 2

	for currentAttempt := 1; currentAttempt <= maxAttempts; currentAttempt++ {
		if currentAttempt > 1 {
			cred.RefreshHeaderCookie()
			addHeaderCookie(req)
		}

		time.Sleep(time.Duration(timeoutBeforeRequest) * time.Second)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			erroring.LogIfError(err, fmt.Sprintf("[attempt %v/%v] error when invoking request for: \n```\n%v\n```\n", currentAttempt, maxAttempts, req.URL))
			continue
		}
		defer resp.Body.Close()

		if st := resp.StatusCode; st != http.StatusOK {
			if st == http.StatusTooManyRequests {
				erroring.LogIfError(
					fmt.Errorf("response code was %v, expected %v", st, http.StatusOK),
					fmt.Sprintf("[attempt %v/%v] trying again in %v seconds...", currentAttempt, maxAttempts, timeoutOnTooManyRequestsSecs),
				)
				time.Sleep(time.Duration(timeoutOnTooManyRequestsSecs) * time.Second)
				resp.Body.Close()
				continue
			} else {
				erroring.LogIfError(
					fmt.Errorf("response code was %v, expected %v", st, http.StatusOK),
					fmt.Sprintf("[attempt %v/%v] problem with response after invoking request for: \n```\n%v\n```\n", currentAttempt, maxAttempts, req.URL),
				)
				resp.Body.Close()
				continue
			}
		} else {
			output, err = io.ReadAll(resp.Body)
			if err != nil {
				erroring.LogIfError(err, fmt.Sprintf("[attempt %v/%v] error when reading response from request for: \n```\n%v\n```\n", currentAttempt, maxAttempts, req.URL))
				resp.Body.Close()
				continue
			}

			break
		}
	}

	if len(output) == 0 {
		log.Printf("out of attempts / never got a proper response for:\n```\n%v\n```\n", req.URL)
	}

	return output
}

func addUrlParametersToRequest(req *http.Request, params map[string]string) {
	parametrizedUrl := req.URL.Query()
	for key, val := range params {
		parametrizedUrl.Add(key, val)
	}
	req.URL.RawQuery = parametrizedUrl.Encode()
}

func addHeaderCookie(req *http.Request) {
	req.Header.Del("Cookie")
	req.Header.Add("Cookie", cred.GetHeaderCookie())
}
