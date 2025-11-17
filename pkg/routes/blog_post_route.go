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

func GetAll(response http.ResponseWriter, request *http.Request) {
	_ = request
	allPosts, err := persistence.GetAllBlogPosts(db.DB)
	if err != nil {
		log.Printf("Failed to fetch blog posts with error %s", err.Error())
		common.WriteJsonError(err,response, http.StatusNotFound)
		return
	}
	if allPosts == nil || len(allPosts) == 0 {
		log.Print("No posts found..")
		common.WriteJsonError(errors.New("no blog posts found"),response, http.StatusNotFound)
		return
	}
	jsonResponse, err := json.Marshal(allPosts)
	if err != nil {
		common.WriteJsonError(err,response, http.StatusInternalServerError)
		return
	}
	response.Header().Add("Content-Type","application/json")
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("Failed to write blog posts json response with error %s", err.Error())
	}
}

func NewPost(response http.ResponseWriter, request *http.Request) {
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusBadRequest)
		return
	}
	var post model.BlogPost
	err = json.Unmarshal(requestBody,&post)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusBadRequest)
		return
	}
	createdPost, err := persistence.AddNewPost(db.DB,post)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusUnprocessableEntity)
		return
	}
	if createdPost == nil {
		common.WriteJsonError(err, response,http.StatusUnprocessableEntity)
		return
	}
	jsonResponse, err := json.Marshal(createdPost)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusInternalServerError)
		return
	}
	response.Header().Add("Content-Type","application/json")
	response.WriteHeader(http.StatusCreated)
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("failed to write blog post json response due to error :: %s",err.Error())
	}
}


func GetPostByID(response http.ResponseWriter, request *http.Request) {
	idParam := request.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.WriteJsonError(err,response, http.StatusBadRequest)
		return
	}
	post, err := persistence.GetBlogPostByID(db.DB, int64(id))
	if err != nil {
		common.WriteJsonError(err,response,http.StatusNotFound)
		return
	}
	if post == nil {
		common.WriteJsonError(errors.New("blog post not found"), response, http.StatusNotFound)
		return
	}
	jsonResponse, err := json.Marshal(post)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusInternalServerError)
	}
	response.Header().Add("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("failed to write get post by id json response due to error :: %s", err.Error())
	}
}

func UpdatePostByID(response http.ResponseWriter, request *http.Request) {
	idParam := request.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.WriteJsonError(err,response, http.StatusBadRequest)
		return
	}
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusBadRequest)
		return
	}
	var post model.BlogPost
	err = json.Unmarshal(requestBody,&post)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusBadRequest)
		return
	}
	post.ID = int64(id)
	updatedPost, err := persistence.UpdateBlogPost(db.DB, post)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusUnprocessableEntity)
		return
	}
	if updatedPost == nil {
		common.WriteJsonError(errors.New("failed to update blog post"),response,http.StatusUnprocessableEntity)
		return
	}
	jsonResponse, err := json.Marshal(updatedPost)
	if err != nil {
		common.WriteJsonError(err,response,http.StatusInternalServerError)
	}
	response.Header().Add("Content-Type","application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Printf("failed to write update post by id json response due to error :: %s", err.Error())
	}
}

func DeletePostByID(response http.ResponseWriter, request *http.Request) {
	idParam := request.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.WriteJsonError(err,response, http.StatusBadRequest)
		return
	}
	ok, err := persistence.DeleteBlogPostByID(db.DB,int64(id))
	if err != nil {
		common.WriteJsonError(err,response,http.StatusUnprocessableEntity)
		return
	}
	if !ok {
		common.WriteJsonError(errors.New("failed to delete blog post"), response, http.StatusUnprocessableEntity)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}