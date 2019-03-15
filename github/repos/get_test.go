// +build integration

package repos_test

import (
	"os"
	"testing"

	"github.com/Github-Web-Apps/Starhub/github/repos"
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

	repos, err := repos.Get(ctx, client)
	assert.NotEmpty(repos)
	assert.NoError(err)
}
