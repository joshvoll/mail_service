package handlers

import (
	"sanservices.git.beanstalkapp.com/goproposal.git/internal/services"
)

// Handler main structure with handlers functionality
type Handler struct {
	Service *services.Service
}

// New initialize main handler
func New(s *services.Service) *Handler {
	return &Handler{
		Service: s,
	}
}
