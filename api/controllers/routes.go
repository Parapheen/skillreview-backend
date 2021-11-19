package controllers

import (
	"github.com/Parapheen/skillreview-backend/api/middlewares"

	sentryhttp "github.com/getsentry/sentry-go/http"
)

func (s *Server) initializeRoutes() {

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Home Route
	s.Router.
		HandleFunc("/",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.Home)),
		).Methods("GET")

	// Login Route
	s.Router.
		HandleFunc("/v1/auth/{provider}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.Login)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/auth/{provider}/callback",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.LoginCallback)),
		).Methods("GET")

	// Users routes
	s.Router.
		HandleFunc("/v1/users",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.CreateUser)),
		).Methods("POST")
	s.Router.
		HandleFunc("/v1/users",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetUsers)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/users/me",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.GetLoggedUser, s.DB))),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/users/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetUser)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/users/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.AdminAuthentication(s.UpdateUser, s.DB))),
		).Methods("PUT")
	s.Router.
		HandleFunc("/v1/users/{id}",
			sentryHandler.HandleFunc(
				middlewares.AdminAuthentication(s.DeleteUser, s.DB)),
		).Methods("DELETE")
	s.Router.
		HandleFunc("/v1/users/{id}/matches",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.GetRecentMatches, s.DB))),
		).Methods("GET")

	// Matches
	s.Router.
		HandleFunc("/v1/matches/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetMatch)),
		).Methods("GET")

	// ReviewRequests routes
	s.Router.
		HandleFunc("/v1/requests",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.CreateReviewRequest)),
		).Methods("POST")
	s.Router.
		HandleFunc("/v1/requests",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetReviewRequests)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/requests/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetReviewRequest)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/requests/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.UpdateReviewRequest, s.DB))),
		).Methods("PUT")
	s.Router.
		HandleFunc("/v1/requests/{id}",
			sentryHandler.HandleFunc(
				middlewares.Authenticate(s.DeleteReviewRequest, s.DB)),
		).Methods("DELETE")

	// Reviews routes
	s.Router.
		HandleFunc("/v1/reviews",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.CreateReview)),
		).Methods("POST")
	s.Router.
		HandleFunc("/v1/reviews",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetReviews)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/reviews/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.GetReview)),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/reviews/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(s.UpdateReview)),
		).Methods("PUT")
	s.Router.
		HandleFunc("/v1/reviews/{id}",
			sentryHandler.HandleFunc(
				middlewares.Authenticate(s.DeleteReview, s.DB)),
		).Methods("DELETE")

	// Applications routes
	s.Router.
		HandleFunc("/v1/applications",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.CreateApplication, s.DB))),
		).Methods("POST")
	s.Router.
		HandleFunc("/v1/applications/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.GetApplication, s.DB))),
		).Methods("GET")
	s.Router.
		HandleFunc("/v1/applications/{id}",
			sentryHandler.HandleFunc(
				middlewares.SetMiddlewareJSON(
					middlewares.Authenticate(s.UpdateApplication, s.DB))),
		).Methods("PUT")

	// Internal routes
	// s.Router.HandleFunc("/v1/internal/applications", middlewares.SetMiddlewareJSON(middlewares.InternalAuthenticate(s.GetApplications, s.DB))).Methods("GET")
	// s.Router.HandleFunc("/v1/internal/applications/{id}", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.UpdateApplication, s.DB))).Methods("PUT")
}
