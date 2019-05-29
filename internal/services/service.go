package services

import (
	"sanservices.git.beanstalkapp.com/goproposal.git/settings"
)

// Service struct definition containing all services functionalties
type Service struct {
}

// New initialize service and database connections
func New(smtpServer settings.SMTPServer) (*Service, error) {

	s := &Service{}
	return s, nil
}
