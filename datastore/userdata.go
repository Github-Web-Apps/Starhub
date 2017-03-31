package datastore

type Userdatastore interface {
	GetFollowers(userID int64) ([]string, error)
	SaveFollowers(userID int64, followers []string) error
	GetStars(userID int64) ([]Star, error)
	SaveStars(userID int64, stars []Star) error
}
