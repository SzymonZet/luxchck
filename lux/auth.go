package lux

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"szymonzet/luxchck/erroring"
)

var HeaderCookie string
var loginUrl string = getFullUrl("/PatientPortal/Account/LogIn")
var login, password string

func RefreshHeaderCookie(login string, password string) {
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

	loginRequestBody, err := json.Marshal(loginReqBody)
	erroring.QuitIfError(err, "error when trying to marshal auth data")

	payload := strings.NewReader(string(loginRequestBody))

	req, err := http.NewRequest("GET", loginUrl, payload)
	erroring.QuitIfError(err, "error when trying to create auth request")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	erroring.QuitIfError(err, "error when getting response from the auth request")
	defer resp.Body.Close()

	var cookieHeaderString strings.Builder
	for _, cookie := range resp.Cookies() {
		cookieString := fmt.Sprintf("%v=%v; ", cookie.Name, cookie.Value)
		_, err := cookieHeaderString.WriteString(cookieString)
		erroring.QuitIfError(err, "error when trying to build auth header cookie string")
	}

	content, err := io.ReadAll(resp.Body)
	erroring.QuitIfError(err, "error when getting body from auth request response")

	err = json.Unmarshal(content, &loginRespBody)
	erroring.QuitIfError(err, "error when trying to unmarshal auth response")

	if loginRespBody.ErrorMessage != "" || !loginRespBody.Succeeded {
		erroring.QuitIfError(fmt.Errorf("login response not as expected, response error message: `%v`", loginRespBody.ErrorMessage), "login error")
	}
	log.Println("logged in successfully")

	HeaderCookie = cookieHeaderString.String()
}
