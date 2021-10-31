package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/Parapheen/skillreview-backend/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.AutoMigrate(&models.User{}, &models.ReviewRequest{}, &models.Review{}) // database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})

	loggedRouter := handlers.LoggingHandler(os.Stdout, server.Router)
	port := Getenv("PORT", "8080")
	log.Printf("Listening to port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(headers, methods, origins)(server.Router)),
		handlers.CompressHandler(server.Router), loggedRouter)
}
