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
	security  *security.SecurityManager
	stopCh    chan struct{}
}

func InitApp() (*App, error) {
	scheduler := scheduler.NewHumpbackScheduler()
	app := &App{
		webSite:   api.InitRouter(scheduler.NodeHeartbeatChan, scheduler.ServiceChangeChan),
		scheduler: scheduler,
		stopCh:    make(chan struct{}),
	}

	sm, err := security.NewSecurityManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create security manager: %w", err)
	}
	app.security = sm

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

	if err := app.security.GenerateCA(config.CertArgs().CertFile, config.CertArgs().KeyFile); err != nil {
		panic(fmt.Errorf("failed to generate CA: %w", err))
	}

	serverBundle, err := app.security.CreateCertificateBundle("humpback-server")
	if err != nil {
		panic(fmt.Errorf("failed to create master certificate: %w", err))
	}

	controller.Start(app.stopCh)
	app.scheduler.Start(serverBundle)
	app.webSite.Start(serverBundle)
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
