package dto

type GitHubUser struct {
	ID        int
	Login     string
	Email     string
	Followers []string
}
