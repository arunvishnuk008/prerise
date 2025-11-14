package routes

import (
	"crimson-sunrise.site/pkg/model"
	"crimson-sunrise.site/pkg/persistence"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Login(response http.ResponseWriter, request *http.Request) {
	var loginRequest model.LoginRequest
	log.Print("reading login request body")
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatalf("Unable to proceed due to error %s", err.Error())
		return
	}
	log.Printf("login request body is :: %s", string(requestBody))
	err = json.Unmarshal(requestBody, &loginRequest)

	if err != nil {
		log.Fatalf("Unable to deserialize request body due to error : %s", err.Error())
		return
	}

	loginResponse, err := persistence.Login(loginRequest)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "unable to authenticate user", http.StatusUnauthorized)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(loginResponse)
	_, _ = response.Write(jsonResponse)
}

func Logout(response http.ResponseWriter, request *http.Request) {
	//TODO : implement logout by revoking access token in db
}

