package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

type MailtrapMailer struct {
	apiKey  string
	inboxID string
	from    string
}

type mailtrapRequest struct {
	From     mailtrapFrom `json:"from"`
	To       []mailtrapTo `json:"to"`
	Subject  string       `json:"subject"`
	HTML     string       `json:"html"`
	Category string       `json:"category"`
}

type mailtrapFrom struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type mailtrapTo struct {
	Email string `json:"email"`
}

func NewMailtrapMailer(apiKey, fromEmail, inboxID string) *MailtrapMailer {
	return &MailtrapMailer{
		apiKey:  apiKey,
		from:    fromEmail,
		inboxID: inboxID,
	}
}

func (m *MailtrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	// Parse template
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	payload := mailtrapRequest{
		From: mailtrapFrom{
			Email: m.from,
			Name:  FromName,
		},
		To: []mailtrapTo{
			{
				Email: email,
			},
		},
		Subject:  subject.String(),
		HTML:     body.String(),
		Category: "User Activation",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("https://sandbox.api.mailtrap.io/api/send/%s", m.inboxID)
	fmt.Println(url)

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			if i == maxRetries-1 {
				return fmt.Errorf("failed to create request: %v", err)
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
		req.Header.Set("Content-Type", "application/json")
		fmt.Println(req.Header)
		// req.Header.Set("Api-Token", m.apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			if i == maxRetries-1 {
				return fmt.Errorf("failed to send request: %v", err)
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		if i == maxRetries-1 {
			return fmt.Errorf("failed to send email, status: %d, body: %s", resp.StatusCode, string(body))
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to send email after %d retries", maxRetries)
}
