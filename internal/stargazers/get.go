package stargazers

import (
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/google/go-github/github"
)

// Get the list of repos of a given user
func Get(client *github.Client) ([]datastores.Star, error) {
	var result []datastores.Star
	repos, err := repos(client)
	if err != nil {
		return result, err
	}
	for _, repo := range repos {
		stars, err := stars(client, repo)
		if err != nil {
			return result, err
		}
		var stargazers []string
		for _, star := range stars {
			stargazers = append(stargazers, *star.User.Login)
		}
		result = append(
			result,
			datastores.Star{
				RepoID:     int64(*repo.ID),
				RepoName:   *repo.FullName,
				Stargazers: stargazers,
			},
		)
	}
	return result, nil
}

func stars(client *github.Client, repo *github.Repository) ([]*github.Stargazer, error) {
	opt := github.ListOptions{PerPage: 10}
	var allStars []*github.Stargazer
	for {
		stars, resp, err := client.Activity.ListStargazers(
			*repo.Owner.Login, *repo.Name, &opt,
		)
		if err != nil {
			return allStars, err
		}
		allStars = append(allStars, stars...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allStars, nil
}

func repos(client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List("", opt)
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return allRepos, nil
}
