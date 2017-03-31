package stargazers

import (
	"context"
	"sync"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/caarlos0/watchub/internal/repos"
	"github.com/google/go-github/github"
	"golang.org/x/sync/errgroup"
)

// Get the list of repos of a given user
func Get(client *github.Client) (result []datastores.Star, err error) {
	repos, err := repos.Get(client)
	if err != nil {
		return
	}

	var g errgroup.Group
	var m sync.Mutex

	for _, repo := range repos {
		repo := repo
		g.Go(func() error {
			r, er := processRepo(client, repo)
			if er != nil {
				return er
			}
			m.Lock()
			defer m.Unlock()
			result = append(result, r)
			return nil
		})
	}
	log.Info("waiting for the goroutines to end")
	err = g.Wait()
	return
}

func processRepo(
	client *github.Client, repo *github.Repository,
) (result datastores.Star, err error) {
	stars, err := stars(client, repo)
	if err != nil {
		return result, err
	}
	var stargazers []string
	for _, star := range stars {
		stargazers = append(stargazers, *star.User.Login)
	}
	return datastores.Star{
		RepoID:     int64(*repo.ID),
		RepoName:   *repo.FullName,
		Stargazers: stargazers,
	}, nil
}

func stars(
	client *github.Client, repo *github.Repository,
) (result []*github.Stargazer, err error) {
	opt := &github.ListOptions{PerPage: 30}
	for {
		repos, nextPage, err := getPage(opt, client, repo)
		if err != nil {
			return result, err
		}
		result = append(result, repos...)
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return result, nil
}

func getPage(
	opt *github.ListOptions, client *github.Client, repo *github.Repository,
) (stars []*github.Stargazer, nextPage int, err error) {
	ctx := context.Background()
	stars, resp, err := client.Activity.ListStargazers(
		ctx, *repo.Owner.Login, *repo.Name, opt,
	)
	if err != nil {
		return stars, 0, err
	}
	return stars, resp.NextPage, nil
}
