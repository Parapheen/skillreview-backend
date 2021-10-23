package api

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/Parapheen/skillreview-backend/api/controllers"
	"github.com/Parapheen/skillreview-backend/api/seed"
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

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
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

	runSeeder, _ := strconv.ParseBool(os.Getenv("RUN_SEEDER"))
	if runSeeder {
		fmt.Println("Running seeder")
		seed.Load(server.DB)
	}

	server.Run(":8080")
}