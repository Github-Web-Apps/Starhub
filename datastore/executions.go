package datastore

import "github.com/Intika-Web-Apps/Watchub-Mirror/shared/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
