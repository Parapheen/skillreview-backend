package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/utils"
)

func (server *Server) CreateReviewRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	review := models.ReviewRequest{}
	err = json.Unmarshal(body, &review)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	review.Prepare()
	err = review.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != review.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	reviewCreated, err := review.SaveReviewRequest(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, reviewCreated.ID))
	responses.JSON(w, http.StatusCreated, reviewCreated)
}

func (server *Server) GetReviewRequests(w http.ResponseWriter, r *http.Request) {

	review := models.ReviewRequest{}

	reviews, err := review.FindAllReviewRequests(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviews)
}

func (server *Server) GetReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	review_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	review := models.ReviewRequest{}

	postReceived, err := review.FindReviewRequestByID(server.DB, review_id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdateReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	review_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	review := models.ReviewRequest{}
	err = server.DB.Model(models.ReviewRequest{}).Where("id = ?", review_id).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	if uid != review.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewUpdate := models.ReviewRequest{}
	err = json.Unmarshal(body, &reviewUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != reviewUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	reviewUpdate.Prepare()
	err = reviewUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewUpdate.ID = review.ID

	postUpdated, err := reviewUpdate.UpdateReviewRequest(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeleteReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	review_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	review := models.ReviewRequest{}
	err = server.DB.Model(models.ReviewRequest{}).Where("id = ?", review_id).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	if uid != review.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = review.DeleteReviewRequest(server.DB, review_id, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", review_id))
	responses.JSON(w, http.StatusNoContent, "")
}