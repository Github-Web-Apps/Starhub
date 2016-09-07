package mail

import (
	"bytes"
	"fmt"
	"log"

	"html/template"

	"github.com/caarlos0/watchub/config"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const welcomeTemplate = `
<h4>Hi @{{.Login}}!</h4>

<p>
You have {{.Followers}} followers. We will watch for changes and let you know!
</p>

<p>
Sincerely,<br/>
The Watchub Team
<p>
`

// Mailer
type Mailer struct {
	key     string
	from    string
	welcome *template.Template
}

// New mailer
func New(config config.Config) *Mailer {
	tmpl, err := template.New("?").Parse(welcomeTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	return &Mailer{
		key:     config.SendgridAPIKey,
		from:    "caarlos0@gmail.com",
		welcome: tmpl,
	}
}

type WelcomeData struct {
	Login     string
	Email     string
	Followers int
}

// SendWelcome email to a new user :)
func (mailer *Mailer) SendWelcome(data WelcomeData) {
	from := mail.NewEmail("Watchhub", mailer.from)
	subject := "Welcome to Watchub!"
	to := mail.NewEmail(data.Login, data.Email)

	var mailContent bytes.Buffer
	err := mailer.welcome.Execute(&mailContent, data)
	if err != nil {
		fmt.Println("Failed to mail", data, ".", err)
	}
	content := mail.NewContent("text/html", mailContent.String())
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(
		mailer.key, "/v3/mail/send", "https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err = sendgrid.API(request)
	if err != nil {
		fmt.Println("Failed to mail", data, ".", err)
	} else {
		fmt.Println("Mail sent to", data.Email)
	}
}
