package controllers

import "github.com/Parapheen/skillreview-backend/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.
		HandleFunc("/v1/auth/{provider}/callback",
			middlewares.SetMiddlewareJSON(s.LoginCallback)).
		Methods("GET")

		
	// Users routes
	s.Router.HandleFunc("/v1/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/v1/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/v1/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")
	s.Router.HandleFunc("/v1/users/{id}/matches", middlewares.SetMiddlewareAuthentication(s.GetRecentMatches)).Methods("GET")

	// ReviewRequests routes
	s.Router.HandleFunc("/v1/requests", middlewares.SetMiddlewareJSON(s.CreateReviewRequest)).Methods("POST")
	s.Router.HandleFunc("/v1/requests", middlewares.SetMiddlewareJSON(s.GetReviewRequests)).Methods("GET")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.SetMiddlewareJSON(s.GetReviewRequest)).Methods("GET")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateReviewRequest))).Methods("PUT")
	s.Router.HandleFunc("/v1/requests/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteReviewRequest)).Methods("DELETE")
}