package responses

type Match struct {
	HeroID           int  `json:"hero_id"`
	MatchID          int  `json:"match_id"`
	MatchTimestamp   int  `json:"match_timestamp"`
	PerfomanceRating int  `json:"perfomance_rating"`
	WonMatch         bool `json:"won_match"`
}

type MatchPlayer struct {
	AccountID  int    `json:"account_id"`
	Assists    int    `json:"assists"`
	Deaths     int    `json:"deaths"`
	HeroID     int    `json:"hero_id"`
	Items      []int  `json:"items"`
	Kills      int    `json:"kills"`
	PlayerSlot int    `json:"player_slot"`
	ProName    string `json:"pro_name"`
}

type MinimalMatch struct {
	DireScore    int           `json:"dire_score"`
	Duration     int           `json:"duration"`
	GameMode     int           `json:"game_mode"`
	MatchID      int           `json:"match_id"`
	MatchOutcome int           `json:"match_outcome"`
	Players      []MatchPlayer `json:"players"`
	RadiantScore int           `json:"radiant_score"`
	StartTime    int           `json:"start_time"`
}

type UserStatsResponse struct {
	UserID                int  `json:"account_id"`
	BadgePoints           int  `json:"badge_points"`
	IsPlusSubscriber      bool `json:"is_plus_subscriber"`
	PlusOriginalStartDate int  `json:"plus_original_start_date"`
	PreviousRankTier      int  `json:"previous_rank_tier"`
	RankTier              int  `json:"rank_tier"`
}

type Player struct {
	UserID              string `json:"steamid"`
	NickName            string `json:"personaname"`
	Name                string `json:"realname"`
	AvatarURL           string `json:"avatarfull"`
	LocationCountryCode string `json:"loccountrycode"`
	LocationStateCode   string `json:"locstatecode"`
}

type SteamResponse struct {
	Response struct {
		Players []Player `json:"players"`
	} `json:"response"`
}
