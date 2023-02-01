package main

import (
	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
)

func main() {
	// Seconds field, required
	cron.New(cron.WithSeconds())

	// Seconds field, optional
	cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
}
