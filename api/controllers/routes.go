package controllers

import "github.com/Parapheen/skillreview-backend/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//ReviewRequests routes
	s.Router.HandleFunc("/requests", middlewares.SetMiddlewareJSON(s.CreateReviewRequest)).Methods("POST")
	s.Router.HandleFunc("/requests", middlewares.SetMiddlewareJSON(s.GetReviewRequests)).Methods("GET")
	s.Router.HandleFunc("/requests/{id}", middlewares.SetMiddlewareJSON(s.GetReviewRequest)).Methods("GET")
	s.Router.HandleFunc("/requests/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateReviewRequest))).Methods("PUT")
	s.Router.HandleFunc("/requests/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteReviewRequest)).Methods("DELETE")
}