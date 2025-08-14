package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"
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

func convertNewLine(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}

func sendMail(recipient, subject, body string) error {
	// client := sendgrid.NewSendClient(setting.Sendgrid.ApiKey)
	// from := mail.NewEmail("中村雀右衛門オフィシャルウェブサイト", setting.Sendgrid.Sender)
	// to := mail.NewEmail("", recipient)
	// content := mail.NewContent("text/plain", plainTextContent)
	// message := mail.NewV3MailInit(from, subject, to, content)
	// response, err := client.Send(message)
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", setting.Brevo.ApiKey)
	client := brevo.NewAPIClient(cfg)
	email := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "中村雀右衛門オフィシャルウェブサイト",
			Email: setting.Brevo.Sender,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: recipient,
			},
		},
		HtmlContent: body,
		Subject:     subject,
	}
	// email.Subject = subject
	// email.HtmlContent = plainTextContent
	// email.TextContent = convertNewLine(plainTextContent, "\r\n")
	// email.Sender = &brevo.SendSmtpEmailSender{
	// 	Email: setting.Brevo.Sender,
	// }
	_, response, err := client.TransactionalEmailsApi.SendTransacEmail(context.Background(), email)
	if err != nil {
		return err
	}
	if response.StatusCode/100 >= 4 {
		return fmt.Errorf("send mail failed: %v", response.Body)
	}
	return nil
}
