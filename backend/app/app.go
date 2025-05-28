package app

import (
	"context"
	"log/slog"
	"time"

	"humpback/api"
	"humpback/api/static"
	"humpback/config"
	"humpback/internal/controller"
	"humpback/internal/db"
	"humpback/pkg/crypto"
	"humpback/pkg/utils"
	"humpback/scheduler"
)

type App struct {
	webSite   *api.Router
	scheduler *scheduler.HumpbackScheduler
	stopCh    chan struct{}
}

func InitApp() (*App, error) {
	scheduler := scheduler.NewHumpbackScheduler()
	app := &App{
		webSite:   api.InitRouter(scheduler.NodeHeartbeatChan, scheduler.ServiceChangeChan),
		scheduler: scheduler,
		stopCh:    make(chan struct{}),
	}
	slog.Info("[Init DB] Init DB driver...")
	if err := db.InitDB(); err != nil {
		return nil, err
	}
	slog.Info("[Init DB] Init DB driver completed.")
	if err := static.InitStaticsResource(); err != nil {
		return nil, err
	}
	if err := controller.InitData(); err != nil {
		return nil, err
	}
	return app, nil
}

func (app *App) CheckCerts() bool {

	certFile := config.CertArgs().CertFile
	keyFile := config.CertArgs().KeyFile

	return utils.FileExist(certFile) && utils.FileExist(keyFile)
}

func (app *App) Startup() {

	if !app.CheckCerts() {
		ip := config.NodeArgs().HostIp
		host := config.CertArgs().Host
		if host == "" {
			host = "localhost"
		}
		if ip == "" {
			ip = "0.0.0.0"
		}
		if err := crypto.GenerateCertsForHost(host, ip, config.CertArgs().CertFile, config.CertArgs().KeyFile, time.Now().AddDate(5, 0, 0)); err != nil {
			panic(err)
		}
	} else {
		slog.Info("[Cert] Certs already exist, skip generating new certs.")
	}

	controller.Start(app.stopCh)
	app.scheduler.Start()
	app.webSite.Start()
}

func (app *App) Close(c context.Context) error {
	close(app.stopCh)
	if err := app.webSite.Close(c); err != nil {
		return err
	}
	if err := app.scheduler.Close(c); err != nil {
		return err
	}
	return nil
}
