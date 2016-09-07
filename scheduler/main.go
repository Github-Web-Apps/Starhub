package scheduler

import (
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"github.com/caarlos0/watchub/datastores"
	"github.com/caarlos0/watchub/followers"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

func New(store datastores.Datastore) *cron.Cron {
	c := cron.New()
	fn := process(store)
	c.AddFunc("@every 1h", fn)
	go fn()
	return c
}

func process(store datastores.Datastore) func() {
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
			followers, err := followers.Get(token)
			if err != nil {
				log.Println(err)
				continue
			}
			previousFollowers, err := store.GetFollowers(exec.UserID)
			if err != nil {
				log.Println(err)
				continue
			}
			if err := store.SaveFollowers(
				exec.UserID, toIDArray(followers),
			); err != nil {
				log.Println(err)
				continue
			}
			if len(previousFollowers) == 0 {
				log.Println("First execution for user", exec.UserID)
			}
			log.Println(
				exec.UserID, "had", len(previousFollowers),
				"and now have", len(followers), "followers",
			)
		}
	}
}

func toIDArray(users []*github.User) []int64 {
	var ids []int64
	for _, u := range users {
		ids = append(ids, int64(*u.ID))
	}
	return ids
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
