package app

import (
	"buding-job/common/log"
	"buding-job/job"
	"buding-job/web"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BuDingJobApp struct {
	jobApp *job.JobApp
	webApp *web.WebApp
}

func NewBuDingJobApp() *BuDingJobApp {
	return &BuDingJobApp{
		webApp: web.NewWebApp(),
		jobApp: job.NewJobApp(),
	}
}
func (app *BuDingJobApp) Start() {
	app.jobApp.Start()
	app.webApp.Start()
	log.GetLog().Infoln("SERVER START SUCCESS")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.GetLog().Infoln("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.webApp.Stop(ctx)
	app.jobApp.Stop()
	log.GetLog().Infoln("Server exiting")
}
