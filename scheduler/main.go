package scheduler

import (
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"github.com/caarlos0/watchub/datastores"
	"github.com/caarlos0/watchub/diff"
	"github.com/caarlos0/watchub/followers"
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

// New scheduler
func New(store datastores.Datastore, oauth *oauth.Oauth) *cron.Cron {
	c := cron.New()
	fn := process(store, oauth)
	c.AddFunc("@every 1h", fn)
	go fn()
	return c
}

func process(store datastores.Datastore, oauth *oauth.Oauth) func() {
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
			followers, err := followers.Get(token, oauth)
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

			if len(previousFollowers) == 0 {
				log.Println("First execution for user", exec.UserID)
			} else {
				newFollowers := diff.Of(followersLogin, previousFollowers)
				unfollowers := diff.Of(previousFollowers, followersLogin)
				log.Println(
					exec.UserID, "has", len(newFollowers), "new followers and",
					len(unfollowers), "unfollows",
				)
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
