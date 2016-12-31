package stargazers

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/caarlos0/watchub/internal/repos"
	"github.com/google/go-github/github"
)

// Get the list of repos of a given user
func Get(client *github.Client) (result []datastores.Star, err error) {
	repos, err := repos.Get(client)
	if err != nil {
		return result, err
	}
	var wg sync.WaitGroup
	results := make(chan datastores.Star, len(repos))
	errors := make(chan error)
	for _, repo := range repos {
		wg.Add(1)
		go func(client *github.Client, repo *github.Repository) {
			r, err := processRepo(client, repo)
			if err != nil {
				errors <- err
			} else {
				results <- r
			}
		}(client, repo)
	}
	go func() {
		for {
			select {
			case r := <-results:
				result = append(result, r)
				wg.Done()
			case e := <-errors:
				err = e
				wg.Done()
			}
		}
	}()
	log.Println("Waiting for the goroutines to end")
	wg.Wait()
	return result, err
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
	stars, resp, err := client.Activity.ListStargazers(
		*repo.Owner.Login, *repo.Name, opt,
	)
	if err != nil {
		return stars, 0, err
	}
	return stars, resp.NextPage, nil
}
