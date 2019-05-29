package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"sanservices.git.beanstalkapp.com/goproposal.git/internal/models"
	"sanservices.git.beanstalkapp.com/goproposal.git/logger"
)

// errors
var (
	errEmptyRecipients    = errors.New("Email recipients empty")
	errInvalidEmailFormat = errors.New("Email format not valid")
	errInvalidEmailDomain = errors.New("Email domain not valid")
	errMailNotSent        = errors.New("Error sending email")
)

func (h *Handler) sendMail(mailData *models.Mail) error {
	var err error

	if len(mailData.Recipients) == 0 {
		return errEmptyRecipients
	}

	// verify emails
	for _, recepient := range mailData.Recipients {
		err = checkmail.ValidateFormat(recepient)
		if err != nil {
			return errInvalidEmailFormat
		}

		err := checkmail.ValidateHost(recepient)
		if err != nil {
			return errInvalidEmailDomain
		}
	}

	// send mails
	err = h.Service.SendMail(mailData)
	if err != nil {
		return errMailNotSent
	}

	return nil
}

// MailerEndpoint handler
func (h *Handler) MailerEndpoint(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	defer logger.TrackExecutionTime(ctx, "handlers.SendMails", time.Now())

	// get data from payload
	mailData := &models.Mail{}
	err := jsonConverter(req.Body, &mailData)

	if err != nil {
		logger.LogError(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonErrorResponse(w, err, nil)
		return
	}

	recipients := strings.Join(mailData.Recipients[:], ", ")

	if err = h.sendMail(mailData); err != nil {
		logMessage := fmt.Sprintf("source=\"http\", mail_Sent=\"false\", email_recipients=\"%s\" error=\"%s\"", recipients, err.Error())
		logger.LogError(ctx, logMessage)

		if err == errInvalidEmailFormat || err == errInvalidEmailDomain || err == errEmptyRecipients {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		jsonErrorResponse(w, err, nil)
		return
	}

	logmessage := fmt.Sprintf("source=\"http\", origin_requestId=\"%s\", mail_sent=\"true\", email_recipients=\"%s\"", mailData.OriginRequestID, recipients)
	logger.Log(ctx, logmessage)
	w.WriteHeader(http.StatusOK)
	jsonResponse(w, "mail(s) sent successfully")
}

// OnMessageReceived handles brocast messages
func (h *Handler) OnMessageReceived(topic string, message string) {

	if topic == "mail_service" {
		mailPayload := &models.Mail{}
		err := json.Unmarshal([]byte(message), &mailPayload)
		if err != nil {
			log.Println(err)
			return
		}

		recipients := strings.Join(mailPayload.Recipients[:], ", ")
		err = h.sendMail(mailPayload)

		if err != nil {
			logMessage := fmt.Sprintf("source=\"brocast\", requestId=\"%s\", mail_Sent=\"false\", email_recipients=\"%s\" error=\"%s\"", mailPayload.OriginRequestID, recipients, err.Error())
			log.Println(logMessage)
			return
		}

		logmessage := fmt.Sprintf("source=\"brocast\", requestId=\"%s\" mail_sent=\"true\", email_recipients=\"%s\"", mailPayload.OriginRequestID, recipients)
		log.Println(logmessage)
	}
}
