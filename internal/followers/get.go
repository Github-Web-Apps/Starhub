package followers

import "github.com/google/go-github/github"

// Get the list of followers of a given user
func Get(client *github.Client) ([]*github.User, error) {
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
