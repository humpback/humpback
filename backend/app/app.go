package app

import (
    "context"
    "fmt"
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

func (app *App) CheckCerts() error {
    certFile := config.CertArgs().CertFile
    keyFile := config.CertArgs().KeyFile
    
    if err := utils.EnsureDirAndForFile(certFile); err != nil {
        return fmt.Errorf("[Cert] cert file: %s", err)
    }
    if err := utils.EnsureDirAndForFile(keyFile); err != nil {
        return fmt.Errorf("[Cert] key file: %s", err)
    }
    
    if utils.FileExist(certFile) && utils.FileExist(keyFile) {
        slog.Info("[Cert] Certs already exist, skip generating new certs.")
        return nil
    }
    
    ip := config.NodeArgs().HostIp
    host := config.CertArgs().Host
    if host == "" {
        host = "localhost"
    }
    if ip == "" {
        ip = "0.0.0.0"
    }
    if err := crypto.GenerateCertsForHost(host, ip, certFile, keyFile, time.Now().AddDate(5, 0, 0)); err != nil {
        return fmt.Errorf("[Cert] generate cert and key file failed: %s", err)
    }
    slog.Info("[Cert] Certs not exist, generating new certs succeed.")
    return nil
}

func (app *App) Startup() {
    if err := app.CheckCerts(); err != nil {
        panic(err)
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
