package services

import (
	"time"

	"github.com/robfig/cron/v3"
)

func CronWinners(procent, maxWinners, hours int) {
	kiev, _ := time.LoadLocation("Europe/Kiev")
	c := cron.New(cron.WithLocation(kiev))

	c.AddFunc("0 0/12 * * ?", func() {
		GetWinners(procent, maxWinners, hours)
	})

	c.Start()

	for {
		time.Sleep(time.Second)
	}
}
