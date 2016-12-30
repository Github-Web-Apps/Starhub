package mail

import (
	"bytes"
	"html/template"
	"log"

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
		from:      "noreply@watchub.pw",
		templates: template.Must(template.ParseGlob("static/mail/*.html")),
	}
}

type StarData struct {
	Repo  string
	Users []string
}

type WelcomeData struct {
	Login                 string
	Email                 string
	Followers             int
	Stars                 int
	Repos                 int
	ChangeSubscriptionURL string
}

type ChangesData struct {
	Login                 string
	Email                 string
	Followers             int
	Stars                 int
	Repos                 int
	NewFollowers          []string
	Unfollowers           []string
	NewStars              []StarData
	Unstars               []StarData
	ChangeSubscriptionURL string
}

func (mailer *Mailer) send(name, email, subject, template string, data interface{}) {
	from := mail.NewEmail("Watchub", mailer.from)
	to := mail.NewEmail(name, email)

	var content bytes.Buffer
	if err := mailer.templates.ExecuteTemplate(
		&content, template+".html", data,
	); err != nil {
		log.Println("Failed to mail", data, ".", err)
	}
	m := mail.NewV3MailInit(
		from, subject, to, mail.NewContent("text/html", content.String()),
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
	mailer.send(
		data.Login, data.Email, "Your report from Watchub!", "changes", data,
	)
}

// SendWelcome email to a new user
func (mailer *Mailer) SendWelcome(data WelcomeData) {
	mailer.send(data.Login, data.Email, "Welcome to Watchub!", "welcome", data)
}
