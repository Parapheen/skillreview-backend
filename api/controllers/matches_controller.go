package controllers

import (
	"net/http"

	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/services"
	"github.com/gorilla/mux"
)

func (server *Server) GetMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID := vars["id"]

	request := stats_man_domain.MatchRequest{
		MatchId: matchID,
	}
	match, apiError := services.StatsManService.GetMatch(request)
	if apiError != nil {
		responses.JSON(w, apiError.Status(), apiError)
		return
	}

	responses.JSON(w, http.StatusOK, match)
}
