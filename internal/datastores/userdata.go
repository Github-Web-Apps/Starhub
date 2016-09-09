package datastores

type Userdatastore interface {
	GetFollowers(userID int64) ([]string, error)
	SaveFollowers(userID int64, followers []string) error
	GetStars(userID int64) ([]Stars, error)
	SaveStars(userID int64, stars []Stars) error
}

type Stars struct {
	RepoID     int64    `json:"repo_id"`
	RepoName   string   `json:"repo_name"`
	Stargazers []string `json:"stargazers"`
}
