package scheduler

import (
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"github.com/caarlos0/watchub/internal/config"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/caarlos0/watchub/internal/diff"
	"github.com/caarlos0/watchub/internal/followers"
	"github.com/caarlos0/watchub/internal/mail"
	"github.com/caarlos0/watchub/internal/oauth"
	"github.com/caarlos0/watchub/internal/stargazers"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

// New scheduler
func New(
	config config.Config, store datastores.Datastore, oauth *oauth.Oauth,
) *cron.Cron {
	c := cron.New()
	fn := process(config, store, oauth)
	c.AddFunc(config.Schedule, fn)
	go fn()
	return c
}

func process(
	config config.Config, store datastores.Datastore, oauth *oauth.Oauth,
) func() {
	return func() {
		execs, err := store.Executions()
		if err != nil {
			log.Println(err)
			return
		}
		mailer := mail.New(config)
		for _, exec := range execs {
			log.Println("Processing", exec.UserID)
			token, err := tokenFromJSON(exec.Token)
			if err != nil {
				log.Println(err)
				return
			}
			client := oauth.Client(token)
			go doProcess(client, mailer, store, exec)
		}
	}
}

func doProcess(
	client *github.Client,
	mailer *mail.Mailer,
	store datastores.Datastore,
	exec datastores.Execution,
) {
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Println(err)
		return
	}
	followers, err := followers.Get(client)
	if err != nil {
		log.Println(err)
		return
	}
	previousFollowers, err := store.GetFollowers(exec.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	followersLogin := toLoginArray(followers)
	if err := store.SaveFollowers(exec.UserID, followersLogin); err != nil {
		log.Println(err)
		return
	}

	stars, err := stargazers.Get(client)
	if err != nil {
		log.Println(err)
		return
	}
	previousStars, err := store.GetStars(exec.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	if err := store.SaveStars(exec.UserID, stars); err != nil {
		log.Println(err)
		return
	}

	if len(previousFollowers) == 0 && len(previousStars) == 0 {
		starCount := 0
		for _, star := range stars {
			starCount += len(star.Stargazers)
		}
		mailer.SendWelcome(
			mail.WelcomeData{
				Login:     *user.Login,
				Email:     *user.Email,
				Followers: len(followers),
				Stars:     starCount,
				Repos:     len(stars),
			},
		)
	} else {
		newFollowers := diff.Of(followersLogin, previousFollowers)
		unfollowers := diff.Of(previousFollowers, followersLogin)
		if len(newFollowers) > 0 || len(unfollowers) > 0 {
			mailer.SendChanges(
				mail.ChangesData{
					Login:        *user.Login,
					Email:        *user.Email,
					Followers:    len(followers),
					NewFollowers: newFollowers,
					Unfollowers:  unfollowers,
				},
			)
		}
	}
}

func toLoginArray(users []*github.User) []string {
	var logins []string
	for _, u := range users {
		logins = append(logins, *u.Login)
	}
	return logins
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
