package datastores

type Execution struct {
	UserID int    `db:"user_id"`
	Token  string `db:"token" json:"-"`
}

type Execstore interface {
	Executions() ([]Execution, error)
}
