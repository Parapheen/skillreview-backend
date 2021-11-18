package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/clients"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (server *Server) CreateApplication(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	application := models.ReviewerApplication{}
	err = json.Unmarshal(body, &application)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	application.Prepare()
	err = application.Validate("")
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
	application.AuthorID = userGotten.ID
	applicationCreated, err := application.SaveApplication(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	if userGotten.Email != "" {
		defer clients.EmailClientStruct.NewApplication(userGotten.Email, userGotten.Nickname)
	}
	msg := tgbotapi.NewMessage(-1001527648645, fmt.Sprintf("New application\nUUID: %s\nUser: %s\nUserUUID: %s\nDescription: %s\nRating: %d\n", application.UUID.String(), application.Author.Nickname, application.Author.UUID.String(), application.Description, application.Rating))
	defer server.Telegram.Send(msg)
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, applicationCreated.ID))
	responses.JSON(w, http.StatusCreated, applicationCreated)
}

func (server *Server) GetApplications(w http.ResponseWriter, r *http.Request) {

	application := models.ReviewerApplication{}

	applications, err := application.FindAllApplications(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, applications)
}

func (server *Server) GetApplication(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	applicationId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	application := models.ReviewerApplication{}

	applicationReceived, err := application.FindApplicationByUUID(server.DB, applicationId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, applicationReceived)
}

func (server *Server) UpdateApplication(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	applicationId, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	application := models.ReviewerApplication{}
	err = server.DB.Model(models.ReviewerApplication{}).Where("uuid = ?", applicationId).Take(&application).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Application not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	applicationUpdate := models.ReviewerApplication{}
	err = json.Unmarshal(body, &applicationUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	applicationUpdate.Prepare()
	err = applicationUpdate.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	applicationUpdated, err := applicationUpdate.UpdateApplication(server.DB, application.UUID)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, applicationUpdated)
}
