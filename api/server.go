package api

import (
	"log"
	"os"
	"sort"

	"github.com/Parapheen/skillreview-backend/api/controllers"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/steam"
)

var server = controllers.Server{}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error getting env, not comming through %v", err)
	}

	goth.UseProviders(steam.New(
		os.Getenv("STEAM_API_KEY"),
		os.Getenv("STEAM_CALLBACK"),
	))

	m := make(map[string]string)
	m["steam"] = "Steam"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	_ = &ProviderIndex{
		Providers:    keys,
		ProvidersMap: m,
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":8080")
}
