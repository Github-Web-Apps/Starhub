package scheduler

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/github/repos"
	"github.com/caarlos0/watchub/github/stargazers"
	"github.com/caarlos0/watchub/github/user"
	"github.com/caarlos0/watchub/mail"
	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/shared/diff"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/model"
	"github.com/robfig/cron"
)

// New scheduler
func New(
	config config.Config,
	store datastore.Datastore,
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
			go process(exec, config, store, oauth)
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
	oauth *oauth.Oauth,
) {
	var start = time.Now()
	var log = log.WithField("id", exec.UserID)
	var ctx = context.Background()
	client, err := oauth.ClientFrom(ctx, exec.Token)
	if err != nil {
		log.WithError(err).Error("failed to authenticate")
		return
	}

	log.Info("started processing...")
	user, err := user.Info(ctx, client)
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

	repos, err := repos.Get(ctx, client)
	if err != nil {
		log.WithError(err).Error("failed to get user repos from github")
		return
	}
	stars, err := stargazers.Get(ctx, client, repos)
	if err != nil {
		log.WithError(err).Error("failed to get stargazers from github")
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
	var mailer = mail.New(config)
	if len(followers)+len(previousStars) == 0 {
		mailer.SendWelcome(
			dto.WelcomeEmailData{
				Login:     user.Login,
				Email:     user.Email,
				Followers: len(followers),
				Stars:     countStars(stars),
				Repos:     len(stars),
				ClientID:  config.ClientID,
			},
		)
	} else {
		newFollowers := diff.Of(user.Followers, followers)
		unfollowers := diff.Of(followers, user.Followers)
		newStars, unstars := stargazerStatistics(stars, previousStars)
		if len(newFollowers)+len(unfollowers)+len(newStars)+len(unstars) > 0 {
			mailer.SendChanges(
				dto.ChangesEmailData{
					Login:        user.Login,
					Email:        user.Email,
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

func countStars(stars []model.Star) (count int) {
	for _, star := range stars {
		count += len(star.Stargazers)
	}
	return
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
		return &dto.StarEmailData{
			Repo:  name,
			Users: d,
		}
	}
	return nil
}
