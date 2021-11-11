package stats_man_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Parapheen/skillreview-backend/api/clients"
	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
)

type statsManProvider struct{}

type statsManServiceInterface interface {
	GetRecentMatches(request stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, *stats_man_domain.StatsManError)
	GetProfileStats(request stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, *stats_man_domain.StatsManError)
	GetMatch(request stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, *stats_man_domain.StatsManError)
}

var (
	StatsManProvider statsManServiceInterface = &statsManProvider{}
)

func (p *statsManProvider) GetRecentMatches(request stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, *stats_man_domain.StatsManError) {
	url := fmt.Sprintf("%sprofiles/%s/recent_matches", os.Getenv("STATS_API"), request.Steam64ID)
	response, err := clients.ClientStruct.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get recent matches from stats man api %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()

	//The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		var errResponse stats_man_domain.StatsManError
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &stats_man_domain.StatsManError{
				Code:         http.StatusInternalServerError,
				ErrorMessage: "invalid json response body",
			}
		}
		errResponse.Code = response.StatusCode
		return nil, &errResponse
	}
	var result []stats_man_domain.Match
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal weather successful response: %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling weather fetch response"}
	}
	return result, nil
}

func (p *statsManProvider) GetProfileStats(request stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, *stats_man_domain.StatsManError) {
	url := fmt.Sprintf("%sprofiles/%s/card", os.Getenv("STATS_API"), request.Steam64ID)
	response, err := clients.ClientStruct.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get profile stats from stats man api %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()

	//The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		var errResponse stats_man_domain.StatsManError
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &stats_man_domain.StatsManError{
				Code:         http.StatusInternalServerError,
				ErrorMessage: "invalid json response body",
			}
		}
		errResponse.Code = response.StatusCode
		return nil, &errResponse
	}
	var result stats_man_domain.UserStatsResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal weather successful response: %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling weather fetch response"}
	}
	return &result, nil
}

func (p *statsManProvider) GetMatch(request stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, *stats_man_domain.StatsManError) {
	url := fmt.Sprintf("%smatches/%s", os.Getenv("STATS_API"), request.MatchId)
	response, err := clients.ClientStruct.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get match from stats man api %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &stats_man_domain.StatsManError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()

	//The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		var errResponse stats_man_domain.StatsManError
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &stats_man_domain.StatsManError{
				Code:         http.StatusInternalServerError,
				ErrorMessage: "invalid json response body",
			}
		}
		errResponse.Code = response.StatusCode
		return nil, &errResponse
	}
	var result stats_man_domain.MinimalMatch
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal weather successful response: %s", err.Error()))
		return nil, &stats_man_domain.StatsManError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling weather fetch response"}
	}
	return &result, nil
}
