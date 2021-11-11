package controllers

import (
	"net/http"

	"github.com/Parapheen/skillreview-backend/api/clients"
	"github.com/Parapheen/skillreview-backend/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	// err := clients.EmailClientStruct.Greet("parapheen6@gmail.com", "parapheen")
	// if err != nil {
	// 	responses.ERROR(w, http.StatusInternalServerError, err)
	// 	return
	// }
	err := clients.EmailClientStruct.NewReview("parapheen6@gmail.com", "parapheen", "abc", "123")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "SkillReview API")
}
