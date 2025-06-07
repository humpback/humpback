package app

import (
	"context"
	"fmt"
	"log/slog"

	"humpback/api"
	"humpback/api/static"
	"humpback/config"
	"humpback/internal/controller"
	"humpback/internal/db"
	"humpback/scheduler"
	"humpback/security"
)

type App struct {
	webSite   *api.Router
	scheduler *scheduler.HumpbackScheduler
	stopCh    chan struct{}
}

func InitApp() (*App, error) {

	app := &App{
		stopCh: make(chan struct{}),
	}

	err := security.InitSecurityManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create security manager: %w", err)
	}

	scheduler := scheduler.NewHumpbackScheduler()
	app.scheduler = scheduler

	app.webSite = api.InitRouter(scheduler.NodeHeartbeatChan, scheduler.ServiceChangeChan)

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

func (app *App) Startup() {

	if err := security.GenerateCA(); err != nil {
		panic(fmt.Errorf("failed to generate CA: %w", err))
	}

	websiteBundle, err := security.GenerateWebsiteCert(config.CertArgs().CertFile, config.CertArgs().KeyFile)
	if err != nil {
		panic(fmt.Errorf("failed to create website certificate: %w", err))
	}

	serverBundle, err := security.CreateCertificateBundle("humpback-server")
	if err != nil {
		panic(fmt.Errorf("failed to create master certificate: %w", err))
	}

	controller.Start(app.stopCh)
	app.scheduler.Start(serverBundle)
	app.webSite.Start(websiteBundle)
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
