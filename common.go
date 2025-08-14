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
	_, response, err := client.TransactionalEmailsApi.SendTransacEmail(context.Background(), email)
	if err != nil {
		return err
	}
	if response.StatusCode/100 >= 4 {
		return fmt.Errorf("send mail failed: %v", response.Body)
	}
	return nil
}
