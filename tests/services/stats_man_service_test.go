package services_test

import (
	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
	"github.com/Parapheen/skillreview-backend/api/providers"
	"github.com/Parapheen/skillreview-backend/api/services"

	"testing"

	"gopkg.in/go-playground/assert.v1"
)

var (
	getRecentMatchesProviderFunc func(request stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, *stats_man_domain.StatsManError)
	getProfileStatsProviderFunc  func(request stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, *stats_man_domain.StatsManError)
	getMatchProviderFunc         func(request stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, *stats_man_domain.StatsManError)
)

type getProviderMock struct{}

func (c *getProviderMock) GetRecentMatches(request stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, *stats_man_domain.StatsManError) {
	return getRecentMatchesProviderFunc(request)
}

func (c *getProviderMock) GetProfileStats(request stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, *stats_man_domain.StatsManError) {
	return getProfileStatsProviderFunc(request)
}

func (c *getProviderMock) GetMatch(request stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, *stats_man_domain.StatsManError) {
	return getMatchProviderFunc(request)
}

// func TestStatsManServiceRecentMatchesWrongSteamID(t *testing.T) {
// 	getWeatherProviderFunc = func(request weather_domain.WeatherRequest) (*weather_domain.Weather, *weather_domain.WeatherError) {
// 		return nil, &weather_domain.WeatherError{
// 			Code:         400,
// 			ErrorMessage: "The given location is invalid",
// 		}
// 	}
// 	weather_provider.WeatherProvider = &getProviderMock{} //without this line, the real api is fired

// 	request := weather_domain.WeatherRequest{ApiKey: "api_key", Latitude: 123443, Longitude: -71.0589}
// 	result, err := WeatherService.GetWeather(request)
// 	assert.Nil(t, result)
// 	assert.NotNil(t, err)
// 	assert.EqualValues(t, http.StatusBadRequest, err.Status())
// 	assert.EqualValues(t, "The given location is invalid", err.Message())
// }

func TestStatsManServiceRecentMatches(t *testing.T) {
	getRecentMatchesProviderFunc = func(request stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, *stats_man_domain.StatsManError) {
		res := []stats_man_domain.Match{
			{
				HeroID:           12,
				MatchID:          123456789,
				MatchTimestamp:   123213123,
				PerfomanceRating: 3,
				WonMatch:         true,
			},
			{
				HeroID:           24,
				MatchID:          12345,
				MatchTimestamp:   3123123123,
				PerfomanceRating: 5,
				WonMatch:         false,
			},
			{
				HeroID:           10,
				MatchID:          52414123,
				MatchTimestamp:   532546,
				PerfomanceRating: 1,
				WonMatch:         true,
			},
		}
		return res, nil
	}
	stats_man_provider.StatsManProvider = &getProviderMock{} //without this line, the real api is fired

	request := stats_man_domain.ProfileRequest{Steam64ID: "123456789"}
	result, err := services.StatsManService.GetUserRecentMatches(request)
	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, 123456789, result[0].MatchID)
	assert.Equal(t, 12345, result[1].MatchID)
	assert.Equal(t, 52414123, result[2].MatchID)
	assert.Equal(t, true, result[0].WonMatch)
	assert.Equal(t, false, result[1].WonMatch)
	assert.Equal(t, true, result[2].WonMatch)
}
func TestStatsManServiceUserStats(t *testing.T) {
	getProfileStatsProviderFunc = func(request stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, *stats_man_domain.StatsManError) {
		res := stats_man_domain.UserStatsResponse{
			UserID:                12345,
			BadgePoints:           123,
			IsPlusSubscriber:      false,
			PlusOriginalStartDate: 12312312,
			PreviousRankTier:      0,
			RankTier:              64,
		}
		return &res, nil
	}
	stats_man_provider.StatsManProvider = &getProviderMock{} //without this line, the real api is fired

	request := stats_man_domain.ProfileRequest{Steam64ID: "123456789"}
	result, err := services.StatsManService.GetUserProfileStats(request)
	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, result.UserID, 12345)
	assert.Equal(t, result.BadgePoints, 123)
	assert.Equal(t, result.IsPlusSubscriber, false)
	assert.Equal(t, result.PreviousRankTier, 0)
	assert.Equal(t, result.RankTier, 64)
}
func TestStatsManServiceGetMatch(t *testing.T) {
	getMatchProviderFunc = func(request stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, *stats_man_domain.StatsManError) {
		res := stats_man_domain.MinimalMatch{
			DireScore:    123,
			Duration:     122112,
			GameMode:     2,
			MatchID:      123123123,
			MatchOutcome: 1,
			Players:      []stats_man_domain.MatchPlayer{},
			RadiantScore: 123,
		}
		return &res, nil
	}
	stats_man_provider.StatsManProvider = &getProviderMock{} //without this line, the real api is fired

	request := stats_man_domain.MatchRequest{MatchId: "123123123"}
	result, err := services.StatsManService.GetMatch(request)
	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, result.DireScore, 123)
	assert.Equal(t, result.Duration, 122112)
	assert.Equal(t, result.GameMode, 2)
	assert.Equal(t, result.MatchID, 123123123)
	assert.Equal(t, result.MatchOutcome, 1)
	assert.Equal(t, len(result.Players), 0)
	assert.Equal(t, result.RadiantScore, 123)
}
