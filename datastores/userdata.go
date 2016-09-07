package datastores

type Userdatastore interface {
	GetFollowers(userID int64) ([]int64, error)
	SaveFollowers(userID int64, followers []int64) error
}
