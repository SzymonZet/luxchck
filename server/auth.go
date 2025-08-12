package server

import (
	"SzymonZet/LuxmedCheck/connection"
	"SzymonZet/LuxmedCheck/erroring"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type HttpRequest *http.Request

var CurrentHeaderAuthString string
var loginUrl string = connection.GetFullUrl("/PatientPortal/Account/LogIn")

func RefreshAuthToken(login string, password string) {
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

	// fmt.Println("=====")
	// fmt.Println(cookieHeaderString.String())
	// fmt.Println("=====")

	content, err := io.ReadAll(resp.Body)
	erroring.QuitIfError(err, "error when getting body from auth request response")

	err = json.Unmarshal(content, &loginRespBody)
	erroring.QuitIfError(err, "error when trying to unmarshal auth response")

	if loginRespBody.ErrorMessage != "" || !loginRespBody.Succeeded {
		erroring.QuitIfError(fmt.Errorf("login response not as expected, response error message: `%v`", loginRespBody.ErrorMessage), "login error")
	}
	log.Println("logged in successfully")

	CurrentHeaderAuthString = cookieHeaderString.String()
	//CurrentAuthToken = fmt.Sprintf("Authorization-Token=%v; UserAdditionalInfo=", loginRespBody.Token)

}

func AddAuthTokenHeader(httpRequest *http.Request) {
	httpRequest.Header.Add("Cookie", CurrentHeaderAuthString)
}
