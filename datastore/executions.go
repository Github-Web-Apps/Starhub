package datastore

import "github.com/Intika-Web-Apps/Starhub-Notifier/shared/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
