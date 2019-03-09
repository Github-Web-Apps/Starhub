// +build integration

package stargazers_test

import (
	"os"
	"testing"

	"github.com/Intika-Web-Apps/Watchub-Mirror/github/stargazers"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGet(t *testing.T) {
	var ctx = oauth2.NoContext
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	assert := assert.New(t)
	var opt = &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	}
	repos, _, err := client.Repositories.List(ctx, "", opt)
	assert.NoError(err)
	stargazers, err := stargazers.Get(ctx, client, repos)
	assert.NotEmpty(stargazers)
	assert.NoError(err)
}
