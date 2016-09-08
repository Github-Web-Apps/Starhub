package mail

import (
	"bytes"
	"log"

	"html/template"

	"github.com/caarlos0/watchub/internal/config"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Mailer
type Mailer struct {
	key       string
	from      string
	templates *template.Template
}

// New mailer
func New(config config.Config) *Mailer {
	return &Mailer{
		key:       config.SendgridAPIKey,
		from:      "caarlos0@gmail.com",
		templates: template.Must(template.ParseGlob("static/mail/*.html")),
	}
}

type WelcomeData struct {
	Login     string
	Email     string
	Followers int
}

type ChangesData struct {
	Login        string
	Email        string
	Followers    int
	NewFollowers []string
	Unfollowers  []string
}

func (mailer *Mailer) send(name, email, content, subject string) {
	from := mail.NewEmail("Watchhub", mailer.from)
	to := mail.NewEmail(name, email)

	m := mail.NewV3MailInit(
		from, subject, to, mail.NewContent("text/html", content),
	)

	request := sendgrid.GetRequest(
		mailer.key, "/v3/mail/send", "https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println("Failed to mail to", email, ".", err)
	} else {
		log.Println("Mail sent to", email)
	}
}

// SendChanges report to an existing user
func (mailer *Mailer) SendChanges(data ChangesData) {
	subject := "Your report from Watchub!"
	var mailContent bytes.Buffer
	err := mailer.templates.ExecuteTemplate(&mailContent, "changes.html", data)
	if err != nil {
		log.Println("Failed to mail", data, ".", err)
	}
	mailer.send(data.Login, data.Email, mailContent.String(), subject)
}

// SendWelcome email to a new user
func (mailer *Mailer) SendWelcome(data WelcomeData) {
	subject := "Welcome to Watchub!"
	var mailContent bytes.Buffer
	err := mailer.templates.ExecuteTemplate(&mailContent, "welcome.html", data)
	if err != nil {
		log.Println("Failed to mail", data, ".", err)
	}
	mailer.send(data.Login, data.Email, mailContent.String(), subject)
}
