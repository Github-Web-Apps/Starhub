package scheduler

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/github/stargazers"
	"github.com/caarlos0/watchub/github/user"
	"github.com/caarlos0/watchub/shared/diff"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/model"
	"github.com/robfig/cron"
)

// New scheduler
func New(
	config config.Config,
	store model.Datastore,
	oauth *oauth.Oauth,
) *cron.Cron {
	var fn = func() {
		execs, err := store.Executions()
		if err != nil {
			log.WithError(err).Error("failed to get executions")
			return
		}
		for _, exec := range execs {
			exec := exec
			go process(exec, config, store)
		}
	}

	var cron = cron.New()
	if err := cron.AddFunc(config.Schedule, fn); err != nil {
		log.WithError(err).Fatal("failed to start cron service")
	}
	return cron
}

func process(
	exec model.Execution,
	config config.Config,
	store datastore.Datastore,
) {
	var start = time.Now()
	var log = log.WithField("id", exec.UserID)
	var ctx = context.Background()
	var client = oauth.ClientFrom(ctx, exec.Token)

	log.Info("started processing...")
	var user, err = user.Info(ctx, client)
	if err != nil {
		log.WithError(err).Error("failed to get user info")
	}
	log = log.WithField("email", user.Email)

	followers, err := store.GetFollowers(exec.UserID)
	if err != nil {
		log.WithError(err).Error("failed to get user followers from db")
		return
	}
	if err = store.SaveFollowers(exec.UserID, user.Followers); err != nil {
		log.WithError(err).Error("failed to store user followers to db")
		return
	}

	// stars
	stars, err := stargazers.Get(ctx, client)
	if err != nil {
		log.WithError(err).Error("failed to get user repos stargazers from github")
		return
	}
	previousStars, err := store.GetStars(exec.UserID)
	if err != nil {
		log.WithError(err).Error("failed to get user repos stargazers from db")
		return
	}
	if err := store.SaveStars(exec.UserID, stars); err != nil {
		log.WithError(err).Error("failed to store user repos stargazers to db")
		return
	}

	// send email
	if len(followers)+len(previousStars) == 0 {
		mailer.SendWelcome(
			mail.WelcomeData{
				Login:     *user.Login,
				Email:     email,
				Followers: len(followers),
				Stars:     countStars(stars),
				Repos:     len(stars),
				ClientID:  config.ClientID,
			},
		)
	} else {
		newFollowers := diff.Of(followersLogin, followers)
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
					ClientID:     config.ClientID,
				},
			)
		}
	}
	log.WithField("time_taken", time.Since(start).Seconds()).Info("successfully processed")
}

func countStars(stars []model.Star) int {
	starCount := 0
	for _, star := range stars {
		starCount += len(star.Stargazers)
	}
	return starCount
}

func stargazerStatistics(stars, previousStars []model.Star) (newStars, unstars []dto.StarEmailData) {
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

func getDiff(name string, x, xs []string) *dto.StarEmailData {
	var d = diff.Of(x, xs)
	if len(d) > 0 {
		return &mail.StarData{
			Repo:  name,
			Users: d,
		}
	}
	return nil
}
