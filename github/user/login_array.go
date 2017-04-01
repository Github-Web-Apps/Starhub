package user

import "github.com/google/go-github/github"

// ToLoginArray maps a list of github users to a list of their logins
func ToLoginArray(users []*github.User) (logins []string) {
	for _, u := range users {
		logins = append(logins, *u.Login)
	}
	return
}
