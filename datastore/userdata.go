package datastore

import "github.com/caarlos0/watchub/shared/model"

type Userdatastore interface {
	GetFollowers(userID int64) ([]string, error)
	SaveFollowers(userID int64, followers []string) error
	GetStars(userID int64) ([]model.Star, error)
	SaveStars(userID int64, stars []model.Star) error
}
