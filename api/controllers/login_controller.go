package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/goth/gothic"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	formaterror "github.com/Parapheen/skillreview-backend/api/utils"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (server *Server) LoginCallback(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	originURL := keys.Get("state")

	authorizedUserInfo, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(err)
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("%v\n", authorizedUserInfo)
	user := models.User{}
	existingUser, err := user.FindUserBySteamID(server.DB, authorizedUserInfo.UserID)

	if err != nil {
		user.Nickname = authorizedUserInfo.NickName
		user.SteamID = authorizedUserInfo.UserID
		err = user.Validate("")
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		userCreated, err := user.SaveUser(server.DB)

		if err != nil {

			formattedError := formaterror.FormatError(err.Error())

			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}

		token, err := auth.CreateToken(userCreated.UUID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s?accessToken=%s", originURL, token))
		responses.JSON(w, http.StatusTemporaryRedirect, token)
	}

	token, err := auth.CreateToken(existingUser.UUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s?accessToken=%s", originURL, token))
	responses.JSON(w, http.StatusTemporaryRedirect, token)
}