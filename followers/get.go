package followers

import (
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Get(token *oauth2.Token) ([]*github.User, error) {
	client := github.NewClient(oauth.Config.Client(oauth2.NoContext, token))

	opt := &github.ListOptions{PerPage: 10}

	var allFollowers []*github.User
	for {
		followers, resp, err := client.Users.ListFollowers("", opt)
		if err != nil {
			return allFollowers, err
		}
		allFollowers = append(allFollowers, followers...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allFollowers, nil
}
