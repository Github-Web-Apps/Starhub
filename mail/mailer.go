package mail

import (
	"bytes"
	"html/template"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/dto"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Mailer service
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

// SendChanges report to an existing user
func (mailer *Mailer) SendChanges(data dto.ChangesEmailData) {
	mailer.send(
		data.Login,
		data.Email,
		"Your report from Watchub!",
		"changes",
		data,
	)
}

// SendWelcome email to a new user
func (mailer *Mailer) SendWelcome(data dto.WelcomeEmailData) {
	mailer.send(
		data.Login,
		data.Email,
		"Welcome to Watchub!",
		"welcome",
		data,
	)
}

func (mailer *Mailer) send(name, email, subject, template string, data interface{}) {
	var content bytes.Buffer
	var log = log.WithField("email", email).WithField("template", template)
	var from = mail.NewEmail("Watchub", mailer.from)
	var to = mail.NewEmail(name, email)

	var err = mailer.templates.ExecuteTemplate(&content, template+".html", data)
	if err != nil {
		log.WithError(err).Error("failed parse email template")
		return
	}
	var request = sendgrid.GetRequest(
		mailer.key,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(
		mail.NewV3MailInit(
			from,
			subject,
			to,
			mail.NewContent("text/html", content.String()),
		),
	)
	_, err = sendgrid.API(request)
	if err != nil {
		log.WithError(err).Error("failed to send email")
		return
	}
	log.Info("mail sent")
}
