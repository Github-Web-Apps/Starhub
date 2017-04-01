package datastore

import "github.com/caarlos0/watchub/shared/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
