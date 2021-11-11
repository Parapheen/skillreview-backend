package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/middlewares"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/assert.v1"
	"syreclabs.com/go/faker"
)

func TestGetUserByID(t *testing.T) {

	err := refreshUserAndReviewRequestTables()
	if err != nil {
		log.Fatal(err)
	}
	_, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatal(err)
	}
	token, err := CreateToken(user.UUID)
	tokenString := fmt.Sprintf("Bearer %v", token)
	if err != nil {
		log.Fatalf("cannot generate token: %v\n", err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		nickname     string
		steam64_id   string
		steam32_id   string
		avatar       string
		rank         string
		email        string
		tokenGiven   string
		errorMessage string
	}{
		{
			id:         user.UUID.String(),
			statusCode: 200,
			nickname:   user.Nickname,
			email:      user.Email,
			steam64_id: user.Steam64ID,
			steam32_id: user.Steam32ID,
			avatar:     user.Avatar,
			rank:       user.Rank,
			tokenGiven: tokenString,
		},
	}
	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/v1/users", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middlewares.Authenticate(server.GetUser, server.DB))

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		if v.statusCode == 200 {
			assert.Equal(t, user.Nickname, responseMap["nickname"])
			assert.Equal(t, user.Email, responseMap["email"])
			assert.Equal(t, user.UUID.String(), responseMap["id"])
			assert.Equal(t, len(responseMap["reviews"].([]interface{})), 0)
			assert.NotEqual(t, len(responseMap["review_requests"].([]interface{})), 0)
		}
	}
}

func TestUpdateUser(t *testing.T) {

	err := refreshUserAndReviewRequestTables()
	if err != nil {
		log.Fatal(err)
	}
	_, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatal(err)
	}
	token, err := CreateToken(user.UUID)
	if err != nil {
		log.Fatalf("cannot generate token: %v\n", err)
	}
	randomToken, err := CreateToken(uuid.NewV4())
	if err != nil {
		log.Fatalf("cannot generate token: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)
	randomTokenString := fmt.Sprintf("Bearer %v", randomToken)

	email := faker.Internet().Email()
	nickname := faker.Internet().UserName()

	userSample := []struct {
		inputJSON    string
		statusCode   int
		userId       string
		nickname     string
		email        string
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON: fmt.Sprintf(
				`{"email": "%s", "nickname": "%s"}`,
				email,
				nickname,
			),
			statusCode:   200,
			tokenGiven:   tokenString,
			email:        email,
			nickname:     nickname,
			userId:       user.UUID.String(),
			errorMessage: "",
		},
		{
			inputJSON: fmt.Sprintf(
				`{"email": "%s", "nickname": "%s"}`,
				email,
				nickname,
			),
			statusCode:   401,
			tokenGiven:   randomTokenString,
			email:        email,
			nickname:     nickname,
			userId:       user.UUID.String(),
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range userSample {

		req, err := http.NewRequest("PUT", "/v1/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.userId})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(
			middlewares.AdminAuthentication(server.UpdateUser, server.DB),
		)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, email, responseMap["email"])
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
