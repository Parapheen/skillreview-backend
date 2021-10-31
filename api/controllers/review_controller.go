package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/utils"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (server *Server) CreateReview(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	review := models.Review{}
	err = json.Unmarshal(body, &review)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	review.Prepare()
	err = review.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	userUUID, err := uuid.FromString(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByUIID(server.DB, userUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	requestUUID, err := uuid.FromString(review.ReviewRequestUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	request := models.ReviewRequest{}
	requestGotten, err := request.FindReviewRequestByUIID(server.DB, requestUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	review.AuthorID = userGotten.ID
	review.ReviewRequestID = requestGotten.ID
	reviewCreated, err := review.SaveReview(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, reviewCreated.ID))
	responses.JSON(w, http.StatusCreated, reviewCreated)
}

func (server *Server) GetReviews(w http.ResponseWriter, r *http.Request) {

	review := models.Review{}

	reviews, err := review.FindAllReviews(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviews)
}

func (server *Server) GetReview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	reviewId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	review := models.Review{}

	reviewReceived, err := review.FindReviewByUIID(server.DB, reviewId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviewReceived)
}

func (server *Server) UpdateReview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	reviewId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	review := models.Review{}
	err = server.DB.Model(models.ReviewRequest{}).Where("uuid = ?", reviewId).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Review not found"))
		return
	}
	userUUID, err := uuid.FromString(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByUIID(server.DB, userUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	reviewRequestUUID, err := uuid.FromString(review.ReviewRequestUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	request := models.ReviewRequest{}
	requestGotten, err := request.FindReviewRequestByUIID(server.DB, reviewRequestUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	
	reviewUpdate := models.Review{}
	err = json.Unmarshal(body, &reviewUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if userGotten.ID != review.AuthorID || userGotten.ID != requestGotten.Author.ID || userGotten.Steam64ID != os.Getenv("SUPER_ADMIN_STEAM64ID") {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	reviewUpdate.Prepare()
	err = reviewUpdate.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewUpdated, err := reviewUpdate.UpdateReview(server.DB, review.UUID)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, reviewUpdated)
}

func (server *Server) DeleteReview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	reviewId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	review := models.Review{}
	err = server.DB.Model(models.Review{}).Where("id = ?", reviewId).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}
	userUUID, err := uuid.FromString(uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByUIID(server.DB, userUUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userGotten.ID != review.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = review.DeleteReview(server.DB, reviewId, userUUID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", reviewId))
	responses.JSON(w, http.StatusNoContent, "")
}
