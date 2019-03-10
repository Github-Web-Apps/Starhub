package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/Intika-Web-Apps/Starhub-Notifier/config"
	"github.com/Intika-Web-Apps/Starhub-Notifier/shared/dto"
	"github.com/matcornic/hermes"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// nolint: gochecknoglobals
var emailConfig = hermes.Hermes{
	Product: hermes.Product{
		Name:      "Starhub Notifier",
		Link:      "https://starhub-notifier.duckdns.org/",
		Logo:      "https://raw.githubusercontent.com/Intika-Web-Apps/Starhub-Notifier/master/static/logo/logo-144.png",
		Copyright: "Copyright © 2019 Starhub Notifier.",
	},
	Theme: new(hermes.Flat),
}

// nolint: gochecknoglobals
var welcomeIntro = []string{
	"Welcome To Starhub Notifier!",
	"We're very excited to have you on board.",
}

// nolint: gochecknoglobals
var changesIntro = []string{
	"Here is what changed in your account recently:",
}

func New(config config.Config) *Service {
	return &Service{
		hermes:  emailConfig,
		changes: template.Must(template.ParseFiles("static/mail/changes.md")),
		welcome: template.Must(template.ParseFiles("static/mail/welcome.md")),
		config:  config,
	}
}

type Service struct {
	hermes  hermes.Hermes
	config  config.Config
	changes *template.Template
	welcome *template.Template
}

func (s *Service) SendWelcome(data dto.WelcomeEmailData) {
	html, err := s.generate(data.Login, data, s.welcome, welcomeIntro)
	if err != nil {
		log.WithError(err).Error("failed to generate welcome email")
		return
	}
	s.send(data.Login, data.Email, "Welcome To Starhub Notifier!", html)
}

func (s *Service) SendChanges(data dto.ChangesEmailData) {
	html, err := s.generate(data.Login, data, s.changes, changesIntro)
	if err != nil {
		log.WithError(err).Error("failed to generate changes email")
		return
	}
	s.send(data.Login, data.Email, "Your Report From Starhub Notifier!", html)
}

func (s *Service) generate(login string, data interface{}, tmpl *template.Template, intros []string) (string, error) {
	var wr bytes.Buffer
	if err := tmpl.Execute(&wr, data); err != nil {
		return "", err
	}
	return s.hermes.GenerateHTML(
		hermes.Email{
			Body: hermes.Body{
				Name:   login,
				Intros: intros,
				FreeMarkdown: hermes.Markdown(
					strings.Join(
						[]string{
							wr.String(),
							"\n\n",
							"We will continue to watch for changes and let you know!",
							"\n\n---\n\n",
							"<small>",
							`Liking our service? Maybe you'll consider [make a donation](https://starhub-notifier.duckdns.org/donate) to help with the server fees.`,
							fmt.Sprintf(
								`You might also want to change [your settings](%s).`,
								"https://github.com/settings/connections/applications/"+s.config.ClientID,
							),
							"</small>",
						},
						" ",
					),
				),
			},
		},
	)
}

func (s *Service) send(name, email, subject, html string) {
	var log = log.WithField("email", email)
	var from = mail.NewEmail("Starhub-Notifier", "noreply@starhub-notifier.duckdns.org")
	var to = mail.NewEmail(name, email)
	var request = sendgrid.GetRequest(
		s.config.SendgridAPIKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	// prevent grouping in gmail
	request.Headers["X-Entity-Ref-ID"] = "starhub-" + time.Now().String()
	request.Method = "POST"
	request.Body = mail.GetRequestBody(
		mail.NewV3MailInit(
			from,
			subject,
			to,
			mail.NewContent("text/html", html),
		),
	)
	resp, err := sendgrid.API(request)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("failed to send email")
		return
	}
	log.WithField("status", resp.StatusCode).Info("email request posted")
}
