package followers

import "github.com/google/go-github/github"

// Get the given user followers
func Get(user string, client *github.Client) ([]*github.User, error) {
	opt := &github.ListOptions{PerPage: 10}

	var allFollowers []*github.User
	for {
		followers, resp, err := client.Users.ListFollowers(user, opt)
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
