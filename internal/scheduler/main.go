package scheduler

import (
	"encoding/json"
	"errors"
	"log"
	"time"

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
	c.AddFunc(config.Schedule, process(config, store, oauth))
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
				log.Println("Failed to get the token", exec.UserID, err)
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
	start := time.Now()

	// user info
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Println("Failed to get user data", exec.UserID, err)
		return
	}
	email, err := getEmail(client)
	if err != nil {
		log.Println("Failed to get user email", exec.UserID, err)
		return
	}

	// followers
	followers, err := followers.Get(client)
	if err != nil {
		log.Println(
			"Failed to store user followers from github", exec.UserID, err,
		)
		return
	}
	previousFollowers, err := store.GetFollowers(exec.UserID)
	if err != nil {
		log.Println("Failed to get user followers from db", exec.UserID, err)
		return
	}
	followersLogin := toLoginArray(followers)
	if err := store.SaveFollowers(exec.UserID, followersLogin); err != nil {
		log.Println("Failed to store user followers to db", exec.UserID, err)
		return
	}

	// stars
	stars, err := stargazers.Get(client)
	if err != nil {
		log.Println(
			"Failed to get user repos stargazers from github", exec.UserID, err,
		)
		return
	}
	previousStars, err := store.GetStars(exec.UserID)
	if err != nil {
		log.Println(
			"Failed to get user repos stargazers from db", exec.UserID, err,
		)
		return
	}
	if err := store.SaveStars(exec.UserID, stars); err != nil {
		log.Println(
			"Failed to store user repos stargazers to db", exec.UserID, err,
		)
		return
	}

	// send email
	if len(previousFollowers)+len(previousStars) == 0 {
		mailer.SendWelcome(
			mail.WelcomeData{
				Login:     *user.Login,
				Email:     email,
				Followers: len(followers),
				Stars:     countStars(stars),
				Repos:     len(stars),
			},
		)
	} else {
		newFollowers := diff.Of(followersLogin, previousFollowers)
		unfollowers := diff.Of(previousFollowers, followersLogin)
		newStars, unstars := stargazerStatistics(stars, previousStars)
		if len(newFollowers)+len(unfollowers)+len(newStars)+len(unstars) > 0 {
			mailer.SendChanges(
				mail.ChangesData{
					Login:        *user.Login,
					Email:        email,
					Followers:    len(followers),
					Stars:        countStars(stars),
					Repos:        len(stars),
					NewFollowers: newFollowers,
					Unfollowers:  unfollowers,
					NewStars:     newStars,
					Unstars:      unstars,
				},
			)
		}
	}
	log.Println(
		"Processing", exec.UserID, "=", email,
		"took", time.Since(start).Seconds(), "seconds",
	)
}

func countStars(stars []datastores.Star) int {
	starCount := 0
	for _, star := range stars {
		starCount += len(star.Stargazers)
	}
	return starCount
}

func getEmail(client *github.Client) (email string, err error) {
	emails, _, err := client.Users.ListEmails(&github.ListOptions{PerPage: 10})
	if err != nil {
		return email, err
	}
	for _, e := range emails {
		if *e.Primary && *e.Verified {
			return *e.Email, err
		}
	}
	return email, errors.New("No email found!")
}

func stargazerStatistics(
	stars, previousStars []datastores.Star,
) (newStars, unstars []mail.StarData) {
	for _, s := range stars {
		for _, os := range previousStars {
			if s.RepoID != os.RepoID {
				continue
			}
			if d := getDiff(s.RepoName, s.Stargazers, os.Stargazers); d != nil {
				newStars = append(newStars, *d)
			}
			if d := getDiff(s.RepoName, os.Stargazers, s.Stargazers); d != nil {
				unstars = append(unstars, *d)
			}
			break
		}
	}
	return newStars, unstars
}

func getDiff(name string, x, xs []string) *mail.StarData {
	d := diff.Of(x, xs)
	if len(d) > 0 {
		return &mail.StarData{
			Repo:  name,
			Users: d,
		}
	}
	return nil
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
