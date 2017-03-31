package datastore

import "github.com/caarlos0/watchub/config/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
