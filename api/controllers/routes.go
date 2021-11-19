package controllers

import "github.com/Parapheen/skillreview-backend/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.
		HandleFunc("/v1/auth/{provider}",
			middlewares.SetMiddlewareJSON(s.Login)).Methods("GET")
	s.Router.
		HandleFunc("/v1/auth/{provider}/callback",
			middlewares.SetMiddlewareJSON(s.LoginCallback)).
		Methods("GET")

	// Users routes
	s.Router.HandleFunc("/v1/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/v1/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/v1/users/me", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.GetLoggedUser, s.DB))).Methods("GET")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.SetMiddlewareJSON(middlewares.AdminAuthentication(s.UpdateUser, s.DB))).Methods("PUT")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.AdminAuthentication(s.DeleteUser, s.DB)).Methods("DELETE")
	s.Router.HandleFunc("/v1/users/{id}/matches", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.GetRecentMatches, s.DB))).Methods("GET")

	// Matches
	s.Router.HandleFunc("/v1/matches/{id}", middlewares.SetMiddlewareJSON(s.GetMatch)).Methods("GET")

	// ReviewRequests routes
	s.Router.HandleFunc("/v1/requests", middlewares.SetMiddlewareJSON(s.CreateReviewRequest)).Methods("POST")
	s.Router.HandleFunc("/v1/requests", middlewares.SetMiddlewareJSON(s.GetReviewRequests)).Methods("GET")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.SetMiddlewareJSON(s.GetReviewRequest)).Methods("GET")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.UpdateReviewRequest, s.DB))).Methods("PUT")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.Authenticate(s.DeleteReviewRequest, s.DB)).Methods("DELETE")

	// Reviews routes
	s.Router.HandleFunc("/v1/reviews", middlewares.SetMiddlewareJSON(s.CreateReview)).Methods("POST")
	s.Router.HandleFunc("/v1/reviews", middlewares.SetMiddlewareJSON(s.GetReviews)).Methods("GET")
	s.Router.HandleFunc("/v1/reviews/{id}", middlewares.SetMiddlewareJSON(s.GetReview)).Methods("GET")
	s.Router.HandleFunc("/v1/reviews/{id}", middlewares.SetMiddlewareJSON(s.UpdateReview)).Methods("PUT")
	s.Router.HandleFunc("/v1/reviews/{id}", middlewares.Authenticate(s.DeleteReview, s.DB)).Methods("DELETE")

	// Applications routes
	s.Router.HandleFunc("/v1/applications", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.CreateApplication, s.DB))).Methods("POST")
	s.Router.HandleFunc("/v1/applications/{id}", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.GetApplication, s.DB))).Methods("GET")
	s.Router.HandleFunc("/v1/applications/{id}", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.UpdateApplication, s.DB))).Methods("PUT")

	// Internal routes
	// s.Router.HandleFunc("/v1/internal/applications", middlewares.SetMiddlewareJSON(middlewares.InternalAuthenticate(s.GetApplications, s.DB))).Methods("GET")
	// s.Router.HandleFunc("/v1/internal/applications/{id}", middlewares.SetMiddlewareJSON(middlewares.Authenticate(s.UpdateApplication, s.DB))).Methods("PUT")
}
