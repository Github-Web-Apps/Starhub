package followers

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// Get the list of followers of a given user
func Get(
	ctx context.Context,
	client *github.Client,
) (result []*github.User, err error) {
	var opt = &github.ListOptions{PerPage: 30}
	for {
		followers, nextPage, err := getPage(ctx, client, opt)
		if err != nil {
			return result, errors.Wrap(err, "failed to get followers")
		}
		result = append(result, followers...)
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return result, nil
}

func getPage(
	ctx context.Context,
	client *github.Client,
	opt *github.ListOptions,
) (followers []*github.User, nextPage int, err error) {
	followers, resp, err := client.Users.ListFollowers(ctx, "", opt)
	if err != nil {
		return
	}
	return followers, resp.NextPage, err
}
