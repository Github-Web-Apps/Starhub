package repos_test

import (
	"os"
	"testing"

	"github.com/caarlos0/watchub/internal/repos"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGet(t *testing.T) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	assert := assert.New(t)

	repos, err := repos.Get(client)
	assert.NotEmpty(repos)
	assert.NoError(err)
}
