package datastores

type Userdatastore interface {
	GetFollowers(userID int64) ([]string, error)
	SaveFollowers(userID int64, followers []string) error
}
