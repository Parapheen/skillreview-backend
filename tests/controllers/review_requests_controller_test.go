package controllertests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/middlewares"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetReviewRequestByID(t *testing.T) {

	err := refreshUserAndReviewRequestTables()
	if err != nil {
		log.Fatal(err)
	}
	reqr, _, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatal(err)
	}
	rrSample := []struct {
		id                   string
		statusCode           int
		description          string
		self_rate_laning     int
		self_rate_teamfights int
		self_rate_overall    int
		state                string
		errorMessage         string
	}{
		{
			id:                   reqr.UUID.String(),
			statusCode:           200,
			description:          reqr.Description,
			self_rate_laning:     reqr.SelfRateLaning,
			self_rate_teamfights: reqr.SelfRateTeamfights,
			self_rate_overall:    reqr.SelfRateOverall,
			state:                string(reqr.State),
		},
	}
	for _, v := range rrSample {
		req, err := http.NewRequest("GET", "/v1/requests", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middlewares.Authenticate(server.GetReviewRequest, server.DB))

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		if v.statusCode == 200 {
			assert.Equal(t, reqr.Description, responseMap["description"])
			assert.Equal(t, reqr.SelfRateLaning, responseMap["self_rate_laning"])
			assert.Equal(t, reqr.SelfRateTeamfights, responseMap["self_rate_teamfights"])
			assert.Equal(t, reqr.SelfRateOverall, responseMap["self_rate_overall"])
			assert.Equal(t, reqr.UUID.String(), responseMap["id"])
		}
	}
}
