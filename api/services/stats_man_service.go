package services

import (
	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
	"github.com/Parapheen/skillreview-backend/api/providers"
)

type statsManService struct{}

type statsManServiceInterface interface {
	GetUserRecentMatches(input stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, stats_man_domain.StatsManErrorInterface)
	GetUserProfileStats(input stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, stats_man_domain.StatsManErrorInterface)
	GetMatch(input stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, stats_man_domain.StatsManErrorInterface)
}

var (
	StatsManService statsManServiceInterface = &statsManService{}
)

func (s *statsManService) GetUserRecentMatches(input stats_man_domain.ProfileRequest) ([]stats_man_domain.Match, stats_man_domain.StatsManErrorInterface) {
	request := stats_man_domain.ProfileRequest{
		Steam64ID: input.Steam64ID,
	}
	response, err := stats_man_provider.StatsManProvider.GetRecentMatches(request)
	if err != nil {
		return nil, stats_man_domain.NewStatsManError(err.Code, err.ErrorMessage)
	}
	return response, nil
}

func (s *statsManService) GetUserProfileStats(input stats_man_domain.ProfileRequest) (*stats_man_domain.UserStatsResponse, stats_man_domain.StatsManErrorInterface) {
	request := stats_man_domain.ProfileRequest{
		Steam64ID: input.Steam64ID,
	}
	response, err := stats_man_provider.StatsManProvider.GetProfileStats(request)
	if err != nil {
		return nil, stats_man_domain.NewStatsManError(err.Code, err.ErrorMessage)
	}
	return response, nil
}

func (s *statsManService) GetMatch(input stats_man_domain.MatchRequest) (*stats_man_domain.MinimalMatch, stats_man_domain.StatsManErrorInterface) {
	request := stats_man_domain.MatchRequest{
		MatchId: input.MatchId,
	}
	response, err := stats_man_provider.StatsManProvider.GetMatch(request)
	if err != nil {
		return nil, stats_man_domain.NewStatsManError(err.Code, err.ErrorMessage)
	}
	return response, nil
}
