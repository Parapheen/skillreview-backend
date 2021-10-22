package controllers

import (
	"net/http"

	"github.com/Parapheen/skillreview-backend/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "SkillReview API")
}