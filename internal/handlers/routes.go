package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"sanservices.git.beanstalkapp.com/goproposal.git/settings"
)

// Middleware that sets requestId and apiKey variables in context
func initBaseContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		requestID := req.Header.Get("Kong-Request-ID")
		apiKey := req.Header.Get("api-key")

		reqTrack := map[string]string{
			"Kong-Request-ID": requestID,
			"api-key":         apiKey,
		}

		ctx := context.WithValue(req.Context(), settings.RequestTracking, reqTrack)
		req = req.WithContext(ctx)

		//call hanlder
		next.ServeHTTP(w, req)
	})
}

// Router returns api router
func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	// set default prefix for Service
	api := router.PathPrefix("/service/mail/api/v1").Subrouter()

	// set endpoints
	api.HandleFunc("/send", h.MailerEndpoint).Methods("POST")

	// sets base context data for all handlers
	api.Use(initBaseContext)

	return router
}
