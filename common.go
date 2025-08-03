package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type errorResponse struct {
	Message string `json:"message"`
	Debug   any    `json:"debug"`
}

func respond(w http.ResponseWriter, status int, attributes any) {
	if status/100 >= 5 {
		log.Printf("%#v", attributes)
	} else if status/100 >= 3 {
		log.Printf("%#v", attributes)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(attributes)
}

func sendMail(recipient, subject, plainTextContent string) error {
	client := sendgrid.NewSendClient(setting.Sendgrid.ApiKey)
	from := mail.NewEmail("中村雀右衛門オフィシャルウェブサイト", setting.Sendgrid.Sender)
	to := mail.NewEmail("", recipient)
	content := mail.NewContent("text/plain", plainTextContent)
	message := mail.NewV3MailInit(from, subject, to, content)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode/100 >= 4 {
		return fmt.Errorf("send mail failed: %v", response.Body)
	}
	return nil
}
