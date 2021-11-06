package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/utils"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (server *Server) CreateReviewRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reviewReq := models.ReviewRequest{}
	err = json.Unmarshal(body, &reviewReq)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reviewReq.Prepare()
	err = reviewReq.Validate("")
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
	reviewReq.AuthorID = userGotten.ID
	reviewReq.Author = *userGotten
	reviewCreated, err := reviewReq.SaveReviewRequest(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, reviewCreated.ID))
	responses.JSON(w, http.StatusCreated, reviewCreated)
}

func (server *Server) GetReviewRequests(w http.ResponseWriter, r *http.Request) {

	reviewReq := models.ReviewRequest{}

	reviewReqs, err := reviewReq.FindAllReviewRequests(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviewReqs)
}

func (server *Server) GetReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	reviewReqId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	reviewReq := models.ReviewRequest{}

	reviewReqReceived, err := reviewReq.FindReviewRequestByUIID(server.DB, reviewReqId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviewReqReceived)
}

func (server *Server) UpdateReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	reviewReqId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	reviewReq := models.ReviewRequest{}
	reviewReqGotten, err := reviewReq.FindReviewRequestByUIID(server.DB, reviewReqId)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Review Request not found"))
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
	if userGotten.ID != reviewReqGotten.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewReqUpdate := models.ReviewRequest{}
	err = json.Unmarshal(body, &reviewReqUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewReqUpdate.Prepare()
	err = reviewReqUpdate.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewReqUpdate.ID = reviewReqGotten.ID

	postUpdated, err := reviewReqUpdate.UpdateReviewRequest(server.DB)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeleteReviewRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	reviewReqId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	reviewReq := models.ReviewRequest{}
	err = server.DB.Model(models.ReviewRequest{}).Where("id = ?", reviewReqId).Take(&reviewReq).Error
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

	if userGotten.ID != reviewReq.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = reviewReq.DeleteReviewRequest(server.DB, reviewReqId, userUUID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", reviewReqId))
	responses.JSON(w, http.StatusNoContent, "")
}
