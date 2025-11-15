package routes

import (
	"crimson-sunrise.site/pkg/common"
	"crimson-sunrise.site/pkg/db"
	"crimson-sunrise.site/pkg/persistence"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetAll(response http.ResponseWriter, request *http.Request) {
	_ = request
	allPosts, err := persistence.GetAllBlogPosts(db.DB)
	if err != nil {
		log.Printf("Failed to fetch blog posts with error %s", err.Error())
		common.WriteJsonError(err,response, http.StatusNotFound)
	}
	if allPosts == nil || len(allPosts) == 0 {
		log.Print("No posts found..")
		common.WriteJsonError(errors.New("no blog posts found"),response, http.StatusNotFound)
	}
	jsonResponse, err := json.Marshal(allPosts)
	if err != nil {
		common.WriteJsonError(err,response, http.StatusInternalServerError)
	}
	response.Header().Add("Content-Type","application/json")
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("Failed to write blog posts json response with error %s", err.Error())
	}
}

func NewPost(response http.ResponseWriter, request *http.Request) {

}


func GetPostByID(response http.ResponseWriter, request *http.Request) {

}

func UpdatePostByID(response http.ResponseWriter, request *http.Request) {

}

func DeletePostByID(response http.ResponseWriter, request *http.Request) {

}