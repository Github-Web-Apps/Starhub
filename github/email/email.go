package email

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)

// Get the primary and verified user email
func Get(ctx context.Context, client *github.Client) (email string, err error) {
	emails, _, err := client.Users.ListEmails(ctx, &github.ListOptions{PerPage: 10})
	if err != nil {
		return
	}
	for _, e := range emails {
		if *e.Primary && *e.Verified {
			return *e.Email, err
		}
	}
	return email, errors.New("no email found")
}
