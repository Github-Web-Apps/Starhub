package datastores

type Execution struct {
	UserID int    `db:"user_id"`
	Token  string `db:"token" json:"-"`
}
