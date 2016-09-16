package repos

import "github.com/google/go-github/github"

// Get all user's repos
func Get(client *github.Client) (result []*github.Repository, err error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	for {
		repos, nextPage, err := getPage(opt, client)
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
	opt *github.RepositoryListOptions, client *github.Client,
) (repos []*github.Repository, nextPage int, err error) {
	repos, resp, err := client.Repositories.List("", opt)
	if err != nil {
		return repos, 0, err
	}
	return repos, resp.NextPage, err
}
