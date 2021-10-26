package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Parapheen/skillreview-backend/api/auth"
	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/Parapheen/skillreview-backend/api/responses"
	formaterror "github.com/Parapheen/skillreview-backend/api/utils"
	"github.com/luanruisong/g-steam"
)

type Player struct {
	UserID              string `json:"steamid"`
	NickName            string `json:"personaname"`
	Name                string `json:"realname"`
	AvatarURL           string `json:"avatarfull"`
	LocationCountryCode string `json:"loccountrycode"`
	LocationStateCode   string `json:"locstatecode"`
}

type steamResponse struct {
	Response struct {
		Players []Player `json:"players"`
	} `json:"response"`
}

func FetchUser(steamID string) (Player, error) {
	resp := steamResponse{}
	client := steam.NewClient(os.Getenv("STEAM_API_KEY"))
	api := client.Api()
	_, err := api.Server("ISteamUser"). // Set up service interface
		Method("GetPlayerSummaries"). // Set access function
		Version("v0002"). // Set version
		AddParam("steamids", steamID). // Setting parameters (If the key parameter is not set, the client's appKey will be added by default)
		Get(&resp)
	if err != nil {
		return Player{}, err
	}
	player := resp.Response.Players[0]
	return player, nil
}

func CompleteAuth(vars url.Values) (string, error){
	v := make(url.Values)
	v.Set("openid.assoc_handle", vars.Get("openid.assoc_handle"))
	v.Set("openid.signed", vars.Get("openid.signed"))
	v.Set("openid.sig", vars.Get("openid.sig"))
	v.Set("openid.ns", vars.Get("openid.ns"))

	split := strings.Split(vars.Get("openid.signed"), ",")
	for _, item := range split {
		v.Set("openid."+item, vars.Get("openid."+item))
	}
	v.Set("openid.mode", "check_authentication")
	client := http.DefaultClient

	resp, err := client.PostForm("https://steamcommunity.com/openid/login", v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := strings.Split(string(content), "\n")
	if response[0] != "ns:"+"http://specs.openid.net/auth/2.0" {
		return "", errors.New("Wrong ns in the response.")
	}
	
	if response[1] == "is_valid:false" {
		return "", errors.New("Unable validate openId.")
	}

	openIDURL := vars.Get("openid.claimed_id")
	validationRegExp := regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	if !validationRegExp.MatchString(openIDURL) {
		return "", errors.New("Invalid Steam ID pattern.")
	}

	steamID := regexp.MustCompile("\\D+").ReplaceAllString(openIDURL, "")
	return steamID, nil
}


func (server *Server) LoginCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	vars := r.URL.Query()
	// verify openid info
	steamID, err := CompleteAuth(vars)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserBySteamID(server.DB, steamID)

	if err != nil {
		steamUser, err := FetchUser(steamID)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		user.Nickname = steamUser.NickName
		user.SteamID = steamUser.UserID
		user.Avatar = steamUser.AvatarURL
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
		responses.JSON(w, http.StatusOK, userCreated)
	}

	token, err := auth.CreateToken(userGotten.UUID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	responses.JSON(w, http.StatusOK, userGotten)
}