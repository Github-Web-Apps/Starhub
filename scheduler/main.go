package scheduler

import (
	"log"

	"github.com/caarlos0/watchub/datastores"
	"github.com/robfig/cron"
)

func New(store datastores.Datastore) *cron.Cron {
	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		execs, err := store.Executions()
		if err != nil {
			log.Println(err)
		}
		for _, exec := range execs {
			log.Println(exec.UserID)
		}
	})
	return c
}
