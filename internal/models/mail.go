package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Mail object struct
type Mail struct {
	OriginRequestID string   `json:"origin_requestId"`
	Recipients      []string `json:"recipients"`
	Subject         string   `json:"subject"`
	Body            string   `json:"message_body"`
}

// BuildMessage creates the message to be sent
func (mail *Mail) BuildMessage() string {

	var message string
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	if len(mail.Recipients) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.Recipients, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += mime
	message += "\r\n" + mail.Body

	return message
}

// ToString returns string representation of struct
func (mail Mail) ToString() string {
	str, err := json.Marshal(mail)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(str)
}
