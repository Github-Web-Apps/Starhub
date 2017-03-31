package dto

// StarEmailData is the DTO representing a repository and the users starring it
type StarEmailData struct {
	Repo  string
	Users []string
}

// ChangesEmailData is the DTO passed down to the daily email
type ChangesEmailData struct {
	Login                 string
	Email                 string
	Followers             int
	Stars                 int
	Repos                 int
	NewFollowers          []string
	Unfollowers           []string
	NewStars              []StarEmailData
	Unstars               []StarEmailData
	ChangeSubscriptionURL string
}
