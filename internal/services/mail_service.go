package services

import (
	"crypto/tls"
	"net/smtp"

	"sanservices.git.beanstalkapp.com/goproposal.git/internal/models"
	"sanservices.git.beanstalkapp.com/goproposal.git/settings"
)

// SendMail ...
func (s *Service) SendMail(mail *models.Mail) error {

	config, err := settings.GetConfiguration()
	if err != nil {
		return err
	}

	//build an auth
	auth := smtp.PlainAuth("",
		config.SMTPServer.SenderEmail,
		config.SMTPServer.SenderPassword,
		config.SMTPServer.Host,
	)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.SMTPServer.Host,
	}

	conn, err := tls.Dial("tcp", config.SMTPServer.ServerName(), tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, config.SMTPServer.Host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(config.SMTPServer.SenderEmail); err != nil {
		return err
	}
	for _, recepient := range mail.Recipients {
		if err = client.Rcpt(recepient); err != nil {
			return err
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	// send message
	messageBody := mail.BuildMessage()
	_, err = w.Write([]byte(messageBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}
