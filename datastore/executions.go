package datastore

import "github.com/Github-Web-Apps/Starhub/shared/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
