package scheduler

import (
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastores"
	"github.com/caarlos0/watchub/diff"
	"github.com/caarlos0/watchub/followers"
	"github.com/caarlos0/watchub/mail"
	"github.com/caarlos0/watchub/oauth"
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
		for _, exec := range execs {
			log.Println("Processing", exec.UserID)
			token, err := tokenFromJSON(exec.Token)
			if err != nil {
				log.Println(err)
				continue
			}
			client := oauth.Client(token)
			user, _, err := client.Users.Get("")
			if err != nil {
				log.Println(err)
				continue
			}
			followers, err := followers.Get(client)
			if err != nil {
				log.Println(err)
				continue
			}
			previousFollowers, err := store.GetFollowers(exec.UserID)
			if err != nil {
				log.Println(err)
				continue
			}
			followersLogin := toLoginArray(followers)
			if err := store.SaveFollowers(
				exec.UserID, followersLogin,
			); err != nil {
				log.Println(err)
				continue
			}

			m := mail.New(config)
			if len(previousFollowers) == 0 {
				m.SendWelcome(
					mail.WelcomeData{
						Login:     *user.Login,
						Email:     *user.Email,
						Followers: len(followers),
					},
				)
			} else {
				newFollowers := diff.Of(followersLogin, previousFollowers)
				unfollowers := diff.Of(previousFollowers, followersLogin)
				if len(newFollowers) > 0 || len(unfollowers) > 0 {
					m.SendChanges(
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
