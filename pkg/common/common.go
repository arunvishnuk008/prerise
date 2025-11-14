package common

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func EpochSecondToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, msInt*int64(time.Second)), nil
}

func WriteJsonError(err error, response http.ResponseWriter, statusCode int) {
	jsonError, _ := json.Marshal(map[string]string{
		"error": err.Error(),
		})
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	_, writeError := response.Write(jsonError)
	if writeError != nil {
		log.Print("Failed to write panic as json response")
	}
}

