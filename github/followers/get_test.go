// +build integration

package followers_test

import (
	"os"
	"testing"

	"github.com/caarlos0/watchub/github/followers"
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

	stargazers, err := followers.Get(ctx, client)
	assert.NotEmpty(stargazers)
	assert.NoError(err)
}
