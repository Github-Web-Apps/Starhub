package model

// Execution model
type Execution struct {
	UserID int64  `db:"user_id"`
	Token  string `db:"token" json:"-"`
}
