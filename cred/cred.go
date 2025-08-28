package cred

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"syscall"
	"szymonzet/luxchck/erroring"
	"time"

	"golang.org/x/term"
)

// there are some bad hardcodes + logic overlaps with lux package
// this is somewhat intentional, to avoid exposing login/password anywhere else
// while making it easy to refresh and obtain header cookie

var loginUrl string = "https://portalpacjenta.luxmed.pl/PatientPortal/Account/LogIn"
var headerCookie string
var login, password string

func SetLoginAndPassword() {
	login = getSecureString("Login")
	password = getSecureString("Password")
}

func GetHeaderCookie() string {
	return headerCookie
}

func RefreshHeaderCookie() {
	loginReqBody := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		login, password,
	}

	loginRespBody := struct {
		Succeeded    bool   `json:"succeded"` // typo in their json ;)
		ErrorMessage string `json:"errorMessage"`
		Token        string `json:"token"`
	}{}

	var content []byte
	var cookieHeaderString strings.Builder

	for {
		log.Printf("trying to (re)log in...")
		loginRequestBody, err := json.Marshal(loginReqBody)
		erroring.QuitIfError(err, "error when trying to marshal auth data")

		payload := strings.NewReader(string(loginRequestBody))

		req, err := http.NewRequest("GET", loginUrl, payload)
		erroring.QuitIfError(err, "error when trying to create auth request")
		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		log.Printf("sending login request...")
		resp, err := http.DefaultClient.Do(req)
		erroring.QuitIfError(err, "error when getting response from the auth request")

		for _, cookie := range resp.Cookies() {
			cookieString := fmt.Sprintf("%v=%v; ", cookie.Name, cookie.Value)
			_, err := cookieHeaderString.WriteString(cookieString)
			erroring.QuitIfError(err, "error when trying to build auth header cookie string")
		}

		content, err = io.ReadAll(resp.Body)

		// repeat in case of maintenance
		if strings.Contains(string(content), "<title>Przerwa techniczna</title>") || strings.Contains(string(content), "<title>Technical maintenance</title>") {
			log.Println("maintenance break detected, trying again in 5 minutes...")
			resp.Body.Close()
			time.Sleep(5 * time.Minute)
		} else if strings.Contains(string(content), "<html>") {
			log.Println("some html returned (possibly maintenance), trying again in 5 minutes..")
			resp.Body.Close()
			time.Sleep(5 * time.Minute)
		} else {
			defer resp.Body.Close()
			break
		}

		erroring.QuitIfError(err, "error when getting body from auth request response")
	}

	err := json.Unmarshal(content, &loginRespBody)
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to unmarshal auth response: \n```\n%v\n```\n", string(content)))

	if loginRespBody.ErrorMessage != "" || !loginRespBody.Succeeded {
		erroring.QuitIfError(fmt.Errorf("login response not as expected, response error message: `%v`", loginRespBody.ErrorMessage), "login error")
	}
	log.Println("(re)logged in successfully")

	headerCookie = cookieHeaderString.String()
}

func getSecureString(name string) string {
	fmt.Printf("Type %v and press enter: \t", name)
	val, err := term.ReadPassword(int(syscall.Stdin))
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to obtain %v", name))
	fmt.Println()

	return string(val)
}
