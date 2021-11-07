package stats_man_provider_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/clients"
	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
	"github.com/Parapheen/skillreview-backend/api/providers"

	"gopkg.in/go-playground/assert.v1"
)

var (
	getRequestFunc func(url string) (*http.Response, error)
)

type getClientMock struct{}

//We are mocking the client method "Get"
func (cm *getClientMock) Get(request string) (*http.Response, error) {
	return getRequestFunc(request)
}

const (
	recentMatchesResponse = `[{"hero_id":100,"match_id":6259722031,"match_timestamp":1636134781,"performance_rating":1,"won_match":true},{"hero_id":69,"match_id":6259670397,"match_timestamp":1636131886,"performance_rating":2,"won_match":false},{"hero_id":9,"match_id":6259581688,"match_timestamp":1636130110,"performance_rating":2,"won_match":false},{"hero_id":69,"match_id":6256304608,"match_timestamp":1635957765,"performance_rating":2,"won_match":false},{"hero_id":69,"match_id":6256208144,"match_timestamp":1635955850,"performance_rating":1,"won_match":true},{"hero_id":100,"match_id":6253191800,"match_timestamp":1635794750,"performance_rating":1,"won_match":true},{"hero_id":9,"match_id":6253121250,"match_timestamp":1635791706,"performance_rating":2,"won_match":false},{"hero_id":21,"match_id":6249625288,"match_timestamp":1635623751,"performance_rating":1,"won_match":true}]`
	profileStatsResponse  = `{"account_id":71935067,"badge_points":7040,"is_plus_subscriber":false,"plus_original_start_date":0,"previous_rank_tier":0,"rank_tier":62,"slots":[{"slot_id":0,"stat":{"stat_id":6,"stat_score":1336491278}}]}	`
	matchResponse         = `{"dire_score":71,"duration":4057,"game_mode":5,"match_id":6249440569,"match_outcome":3,"players":[{"account_id":0,"assists":18,"deaths":12,"hero_id":7,"items":[249,250,108,600,116,50],"kills":21,"player_slot":0,"pro_name":""},{"account_id":0,"assists":23,"deaths":14,"hero_id":86,"items":[1,218,259,108,29,235],"kills":5,"player_slot":1,"pro_name":""},{"account_id":125758825,"assists":15,"deaths":21,"hero_id":119,"items":[0,77,158,63,108,135],"kills":14,"player_slot":2,"pro_name":""},{"account_id":0,"assists":41,"deaths":14,"hero_id":110,"items":[256,267,214,40,119,108],"kills":0,"player_slot":3,"pro_name":""},{"account_id":71935067,"assists":24,"deaths":10,"hero_id":35,"items":[154,263,603,135,63,141],"kills":18,"player_slot":4,"pro_name":""},{"account_id":47959953,"assists":33,"deaths":8,"hero_id":68,"items":[108,100,247,235,29,1466],"kills":17,"player_slot":128,"pro_name":""},{"account_id":191972780,"assists":28,"deaths":15,"hero_id":64,"items":[116,110,48,604,235,277],"kills":16,"player_slot":129,"pro_name":""},{"account_id":107181463,"assists":16,"deaths":12,"hero_id":95,"items":[50,116,174,603,156,208],"kills":13,"player_slot":130,"pro_name":""},{"account_id":207632949,"assists":34,"deaths":10,"hero_id":111,"items":[218,180,254,232,96,108],"kills":9,"player_slot":131,"pro_name":""},{"account_id":180081510,"assists":18,"deaths":13,"hero_id":41,"items":[116,156,158,135,48,160],"kills":15,"player_slot":132,"pro_name":""}],"radiant_score":58,"start_time":1635614213}`
)

//When the everything is good
func TestGetRecentMatchesNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(profileStatsResponse)),
		}, nil
	}
	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

	response, err := stats_man_provider.StatsManProvider.GetProfileStats(stats_man_domain.ProfileRequest{Steam64ID: "123456789"})
	assert.NotEqual(t, response, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.UserID, 71935067)
	assert.Equal(t, response.RankTier, 62)
	assert.Equal(t, response.PreviousRankTier, 0)
}

func TestGetProfileStatsNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(recentMatchesResponse)),
		}, nil
	}
	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

	response, err := stats_man_provider.StatsManProvider.GetRecentMatches(stats_man_domain.ProfileRequest{Steam64ID: "123456789"})
	assert.NotEqual(t, response, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, 8, len(response))
	assert.Equal(t, 6259722031, response[0].MatchID)
	assert.Equal(t, 6259670397, response[1].MatchID)
}

func TestGetMatchNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(matchResponse)),
		}, nil
	}
	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

	response, err := stats_man_provider.StatsManProvider.GetMatch(stats_man_domain.MatchRequest{MatchId: "123456789"})
	assert.NotEqual(t, response, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.MatchID, 6249440569)
	assert.Equal(t, response.DireScore, 71)
	assert.Equal(t, response.RadiantScore, 58)
}

// func TestGetWeatherInvalidLatitude(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 34223.3445, -71.0589})

// 	assert.NotNil(t, err)
// 	assert.Nil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "The given location is invalid", err.ErrorMessage)
// }

// func TestGetWeatherInvalidLongitude(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 44.3601, -74331.0589})

// 	assert.NotNil(t, err)
// 	assert.Nil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "The given location is invalid", err.ErrorMessage)
// }

// func TestGetWeatherInvalidFormat(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "Poorly formatted request"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 0, -74331.0589})

// 	assert.NotNil(t, err)
// 	assert.Nil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "Poorly formatted request", err.ErrorMessage)
// }

// //When no body is provided
// func TestGetWeatherInvalidRestClient(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "invalid rest client response"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 0, -74331.0589})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "invalid rest client response", err.ErrorMessage)
// }

// func TestGetWeatherInvalidResponseBody(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "Invalid response body"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"wrong_anything", 44.3601, -71.0589})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "Invalid response body", err.ErrorMessage)
// }

// func TestGetWeatherInvalidRequest(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		invalidCloser, _ := os.Open("-asf3")
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       invalidCloser,
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"wrong_anything", 44.3601, -71.0589})
// 	assert.Nil(t, response)
// 	assert.NotNil(t, err)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Code)
// 	assert.EqualValues(t, "invalid argument", err.ErrorMessage)
// }

// //When the error response is invalid, here the code is supposed to be an integer, but a string was given.
// //This can happen when the api owner changes some data types in the api
// func TestGetWeatherInvalidErrorInterface(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"code": "string code"}`)),
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 44.3601, -71.0589})
// 	assert.Nil(t, response)
// 	assert.NotNil(t, err)
// 	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
// 	assert.EqualValues(t, "invalid json response body", err.ErrorMessage)
// }

// //We are getting a postive response from the api, but, the datatype of the response returned does not match the struct datatype we have defined (does not match the struct type we want to unmarshal this response into).
// func TestGetWeatherInvalidResponseInterface(t *testing.T) {
// 	getRequestFunc = func(url string) (*http.Response, error) {
// 		return &http.Response{
// 			StatusCode: http.StatusOK,
// 			Body:       ioutil.NopCloser(strings.NewReader(`{"latitude": "string latitude", "longitude": -71.0589, "timezone": "America/New_York"}`)), //when we use string for latitude instead of float
// 		}, nil
// 	}
// 	restclient.ClientStruct = &getClientMock{} //without this line, the real api is fired

// 	response, err := WeatherProvider.GetWeather(weather_domain.WeatherRequest{"anything", 44.3601, -71.0589})
// 	assert.Nil(t, response)
// 	assert.NotNil(t, err)
// 	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
// 	assert.EqualValues(t, "error unmarshaling weather fetch response", err.ErrorMessage)
// }
