package user

import (
	"context"

	"github.com/Intika-Web-Apps/Watchub-Mirror/github/email"
	"github.com/Intika-Web-Apps/Watchub-Mirror/github/followers"
	"github.com/Intika-Web-Apps/Watchub-Mirror/shared/dto"
	"github.com/google/go-github/github"
)

// Info gets a github user info, like login, email and followers
func Info(ctx context.Context, client *github.Client) (user dto.GitHubUser, err error) {
	u, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return user, err
	}
	email, err := email.Get(ctx, client)
	if err != nil {
		return user, err
	}
	followers, err := followers.Get(ctx, client)
	if err != nil {
		return user, err
	}

	user.ID = *u.ID
	user.Login = *u.Login
	user.Email = email
	user.Followers = ToLoginArray(followers)
	return
}
