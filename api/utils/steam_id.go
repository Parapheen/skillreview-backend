package utils

import (
	"errors"
	"strconv"
)

func SteamInt64ToString(steamInt int64) (string, error) {
	if steamInt <= 76561197960265728 {
		return string(""), errors.New("64 bit steamid int should be bigger than 76561197960265728")
	}
	steamInt = steamInt - 76561197960265728
	remainder := steamInt % 2
	steamInt = steamInt / 2
	return "STEAM_0:" + strconv.FormatInt(remainder, 10) + ":" + strconv.FormatInt(steamInt, 10), nil

}

func SteamStringToInt32(steamstring string) (int, error) {
	Y, err := strconv.Atoi(steamstring[8:9])
	if err != nil {
		return int(0), err
	}
	Z, err := strconv.Atoi(steamstring[10:])
	if err != nil {
		return int(0), err
	}
	return (Z * 2) + Y, nil
}

func Steam64toSteam32(steam64Id string) (string, error) {
	steamInt, err := strconv.ParseInt(steam64Id, 10, 64)
	if err != nil {
		return "", err
	}
	stringID, err := SteamInt64ToString(steamInt)
	if err != nil {
		return "", err
	}
	steam32, err := SteamStringToInt32(stringID)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(steam32), nil
}