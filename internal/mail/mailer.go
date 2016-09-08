package mail

import (
	"bytes"
	"log"

	"html/template"

	"github.com/caarlos0/watchub/internal/config"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const welcomeTemplate = `
<h2>Hi @{{.Login}}!</h2>

<p>
	You have <strong>{{.Followers}}</strong> followers.
	We will watch for more changes and let you know!
</p>

<p style="color: #333">
	Sincerely,<br/>
	The Watchub Team
<p>
`

const changesTemplate = `
<h2>Hi @{{.Login}}!</h2>

<p>
	Hey, here what changed in your account recently:
</p>

{{if .NewFollowers }}
	<p>
		<h3>People who are now <strong>following</strong> you:</h3>
		<ul>
			{{range $index, $element := .NewFollowers}}
			<li>
				<a href="http://github.com/{{.}}">{{.}}</a>
			</li>
			{{end}}
		</ul>
	</p>
{{end}}

{{if .Unfollowers }}
	<p>
		<h3>People who <strong>unfollowed</strong> you:</h3>
		<ul>
			{{range $index, $element := .Unfollowers}}
			<li>
				<a href="http://github.com/{{.}}">{{.}}</a>
			</li>
			{{end}}
		</ul>
	</p>
{{end}}

<p>
	You now have <strong>{{.Followers}}</strong> followers.
	We will watch for more changes and let you know!
</p>

<p style="color: #333">
	Sincerely,<br/>
	The Watchub Team
<p>
`

// Mailer
type Mailer struct {
	key         string
	from        string
	welcomeTmpl *template.Template
	changesTmpl *template.Template
}

// New mailer
func New(config config.Config) *Mailer {
	welcomeTmpl, err := template.New("?").Parse(welcomeTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	changesTmpl, err := template.New("?").Parse(changesTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	return &Mailer{
		key:         config.SendgridAPIKey,
		from:        "caarlos0@gmail.com",
		welcomeTmpl: welcomeTmpl,
		changesTmpl: changesTmpl,
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
	err := mailer.changesTmpl.Execute(&mailContent, data)
	if err != nil {
		log.Println("Failed to mail", data, ".", err)
	}
	mailer.send(data.Login, data.Email, mailContent.String(), subject)
}

// SendWelcome email to a new user
func (mailer *Mailer) SendWelcome(data WelcomeData) {
	subject := "Welcome to Watchub!"
	var mailContent bytes.Buffer
	err := mailer.welcomeTmpl.Execute(&mailContent, data)
	if err != nil {
		log.Println("Failed to mail", data, ".", err)
	}
	mailer.send(data.Login, data.Email, mailContent.String(), subject)
}
