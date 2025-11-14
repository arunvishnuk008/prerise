package routes

import (
	"crimson-sunrise.site/pkg/common"
	"crimson-sunrise.site/pkg/db"
	"crimson-sunrise.site/pkg/model"
	"crimson-sunrise.site/pkg/persistence"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetAllHotTakes(response http.ResponseWriter, request *http.Request ) {
	_ = request
	log.Printf("Fetching all hot takes..")
	allHotTakes, err := persistence.GetAllHotTakes(db.DB)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	if allHotTakes == nil {
		common.WriteJsonError(errors.New("no hot takes found"), response, http.StatusNotFound)
	}
	response.Header().Add("Content-Type","application/json")
	jsonResponse, _ := json.Marshal(allHotTakes)
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("Error Writing json response!! %s", err.Error())
	}
}


func GetHotTakeByID(response http.ResponseWriter, request *http.Request) {
	log.Printf("Fetching hot take by ID..")
	idParam := request.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.WriteJsonError(err, response,http.StatusUnprocessableEntity)
	}
	hotTake, err := persistence.GetHotTakeByID(db.DB,int64(id))
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	if hotTake == nil {
		common.WriteJsonError(errors.New("hot-take not found"), response, http.StatusNotFound)
	}
	jsonResponse, err := json.Marshal(hotTake)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	response.Header().Add("Content-Type","application/json")
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("Error writing json : %s", err.Error())
	}
}

func NewHotTake(response http.ResponseWriter, request *http.Request) {
	log.Printf("Adding a new hot take to the database")
	var hotTake model.HotTake
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	err = json.Unmarshal(requestBody, &hotTake)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	createdHotTake, err := persistence.AddHotTake(db.DB,hotTake)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	if createdHotTake == nil {
		common.WriteJsonError(errors.New("unable to save hot take"), response, http.StatusUnprocessableEntity)
	}
	jsonResponse, err := json.Marshal(createdHotTake)
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	response.Header().Add("Content-Type","application/json")
	response.WriteHeader(http.StatusCreated)
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("Error writing hot take json : %s", err.Error())
	}
}

func DeleteHotTakeByID(response http.ResponseWriter, request *http.Request) {
	log.Printf("Deleting hot take by ID..")
	idParam := request.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.WriteJsonError(err, response,http.StatusUnprocessableEntity)
	}
	ok, err := persistence.DeleteHotTake(db.DB,int64(id))
	if err != nil {
		common.WriteJsonError(err, response, http.StatusUnprocessableEntity)
	}
	if !ok {
		common.WriteJsonError(errors.New("unable to delete hot-take"), response, http.StatusUnprocessableEntity)
	}
	response.Header().Add("Content-Type","application/json")
	response.WriteHeader(http.StatusNoContent)
}