package mail

import (
	"html/template"
	"io/ioutil"
	"testing"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/stretchr/testify/assert"
)

func TestWelcomeMail(t *testing.T) {
	s := MailSvc{
		hermes: emailConfig,
		config: config.Config{
			ClientID: "1",
		},
		welcome: template.Must(template.ParseFiles("../static/mail/welcome.md")),
	}
	data := dto.WelcomeEmailData{
		Email:     "caarlos0@gmail.com",
		Followers: 10,
		Login:     "caarlos0",
		Repos:     5,
		Stars:     1,
	}
	html, err := s.generate(data.Login, data, s.welcome, welcomeIntro)
	assert.NoError(t, err)
	// TODO: make this flag-toggleable
	// assert.NoError(t, ioutil.WriteFile("testdata/welcome.html", []byte(html), 0644))
	bts, err := ioutil.ReadFile("testdata/welcome.html")
	assert.NoError(t, err)
	assert.Equal(t, string(bts), html)
}

func TestChangesMail(t *testing.T) {
	s := MailSvc{
		hermes: emailConfig,
		config: config.Config{
			ClientID: "1",
		},
		changes: template.Must(template.ParseFiles("../static/mail/changes.md")),
	}
	data := dto.ChangesEmailData{
		Email:        "caarlos0@gmail.com",
		Followers:    10,
		Login:        "caarlos0",
		Repos:        5,
		Stars:        1,
		NewFollowers: []string{"juvenal", "moises"},
		NewStars: []dto.StarEmailData{
			{
				Repo:  "test/test",
				Users: []string{"juvenal", "moises"},
			},
		},
		Unfollowers: []string{"outro-juvenal", "outro-moises"},
		Unstars: []dto.StarEmailData{
			{
				Repo:  "test/test",
				Users: []string{"outro-juvenal", "outro-moises"},
			},
		},
	}
	html, err := s.generate(data.Login, data, s.changes, changesIntro)
	assert.NoError(t, err)
	// TODO: make this flag-toggleable
	// assert.NoError(t, ioutil.WriteFile("testdata/changes.html", []byte(html), 0644))
	bts, err := ioutil.ReadFile("testdata/changes.html")
	assert.NoError(t, err)
	assert.Equal(t, string(bts), html)
}
