package followers

import (
	"github.com/caarlos0/watchub/oauth"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Get the list of followers of a given user
func Get(token *oauth2.Token, oauth *oauth.Oauth) ([]*github.User, error) {
	client := oauth.Client(token)

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
