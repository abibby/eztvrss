package app

import (
	"github.com/abibby/eztvrss/app/events"
	"github.com/abibby/eztvrss/app/jobs"
	"github.com/abibby/eztvrss/config"
	"github.com/abibby/eztvrss/database"
	"github.com/abibby/eztvrss/routes"
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
)

var Kernel = kernel.NewDefaultKernel(
	kernel.Config(config.Kernel),
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.FetchRSSEvent{}),
	),
	kernel.Listeners(
		kernel.NewListener(jobs.FetchRSS),
	),
	kernel.InitRoutes(routes.InitRoutes),
)
