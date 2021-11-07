package controllers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/domain/stats_man_domain"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	"github.com/Parapheen/skillreview-backend/api/services"
	"github.com/Parapheen/skillreview-backend/api/utils"
	formaterror "github.com/Parapheen/skillreview-backend/api/utils"
	"github.com/luanruisong/g-steam"
	"github.com/markbates/goth/gothic"
)

func FetchUser(steamID string) (responses.Player, error) {
	resp := responses.SteamResponse{}
	client := steam.NewClient(os.Getenv("STEAM_API_KEY"))
	api := client.Api()
	_, err := api.Server("ISteamUser"). // Set up service interface
						Method("GetPlayerSummaries").  // Set access function
						Version("v0002").              // Set version
						AddParam("steamids", steamID). // Setting parameters (If the key parameter is not set, the client's appKey will be added by default)
						Get(&resp)
	if err != nil {
		return responses.Player{}, err
	}
	player := resp.Response.Players[0]
	return player, nil
}

func ConvertRankToMedal(rankTier int) (string, error) {
	m := map[int]string{
		80: "Immortal",
		70: "Divine",
		60: "Ancient",
		50: "Legend",
		40: "Archon",
		30: "Crusader",
		20: "Guardian",
		10: "Herald",
		0:  "Uncalibrated",
	}
	rank := int(math.Floor(float64(rankTier)/10) * 10)
	tier := rankTier % 10
	medal := "Unknown"
	if rank >= 0 && rank < 90 {
		medal = m[rank]
		if rank != 80 && rank > 0 && tier > 0 {
			medal = fmt.Sprintf("%s %d", medal, tier)
		}
	}
	return medal, nil
}

func FetchUserRank(steamID string) (string, error) {
	request := stats_man_domain.ProfileRequest{
		Steam64ID: steamID,
	}
	stats, apiError := services.StatsManService.GetUserProfileStats(request)
	if apiError != nil {
		return "", errors.New(apiError.Message())
	}

	rank, err := ConvertRankToMedal(stats.RankTier)
	if err != nil {
		return "", err
	}
	return rank, nil
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (server *Server) LoginCallback(w http.ResponseWriter, r *http.Request) {
	authorizedUserInfo, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(err)
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserBySteamID(server.DB, authorizedUserInfo.UserID)
	if err != nil {
		steamUser, err := FetchUser(authorizedUserInfo.UserID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		user.Nickname = steamUser.NickName
		user.Steam64ID = steamUser.UserID
		user.Steam32ID, err = utils.Steam64toSteam32(user.Steam64ID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		user.Avatar = steamUser.AvatarURL
		user.Rank, err = FetchUserRank(user.Steam64ID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		user.Prepare()
		err = user.Validate("")
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		userCreated, err := user.SaveUser(server.DB)

		if err != nil {

			formattedError := formaterror.FormatError(err.Error())

			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}

		token, err := auth.CreateToken(userCreated.UUID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w.Header().Set("Location", fmt.Sprintf("%s?accessToken=%s", os.Getenv("FRONTEND_URL"), token))
		responses.JSON(w, http.StatusTemporaryRedirect, token)
		return
	}

	token, err := auth.CreateToken(userGotten.UUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Location", fmt.Sprintf("%s?accessToken=%s", os.Getenv("FRONTEND_URL"), token))
	responses.JSON(w, http.StatusTemporaryRedirect, token)
}
