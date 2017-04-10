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
		key:  config.SendgridAPIKey,
		from: "noreply@watchub.pw",
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

func (mailer *Mailer) send(name, email, subject, templateName string, data interface{}) {
	var content bytes.Buffer
	var log = log.WithField("email", email).WithField("template", templateName)
	var from = mail.NewEmail("Watchub", mailer.from)
	var to = mail.NewEmail(name, email)
	if err := executeTemplate(&content, templateName, data); err != nil {
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
	if _, err := sendgrid.API(request); err != nil {
		log.WithError(err).Error("failed to send email")
		return
	}
	log.Info("mail sent")
}

func executeTemplate(content *bytes.Buffer, name string, data interface{}) error {
	templates, _ := template.ParseFiles(
		"static/mail/layout.html",
		"static/mail/"+name+".html",
	)
	return templates.ExecuteTemplate(content, "layout", data)
}
