package main

import (
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func fromGithub(user string, client *github.Client) ([]*github.User, error) {
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

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	user := os.Args[1]
	log.Println("Gathering data for", user)

	followers, err := fromGithub(user, client)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("You have a total of", len(followers), "followers!")
	for _, follower := range followers {
		log.Println(*follower.Login)
	}
}
