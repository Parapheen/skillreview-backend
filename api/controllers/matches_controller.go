package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/gorilla/mux"
)

func (server *Server) GetMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	match := responses.MinimalMatch{}

	matchID := vars["id"]

	client := http.DefaultClient
	resp, err := client.Get(fmt.Sprintf("%smatches/%s", os.Getenv("STATS_API"), matchID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(content, &match)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if match.MatchID == 0 {
		responses.JSON(w, http.StatusNotFound, err)
	}
	responses.JSON(w, http.StatusOK, match)
}
