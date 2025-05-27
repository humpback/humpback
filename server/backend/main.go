package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "time"
    
    "humpback/app"
    "humpback/config"
    "humpback/pkg/glog"
    "humpback/pkg/utils"
)

func init() {
    glog.Open(glog.WithOutputSource(glog.OutputTypeStd), glog.WithOutputFormat(glog.OutputFormatDefault))
    if err := config.InitConfig(); err != nil {
        panic(err)
    }
    slog.Info("[Config] Init completed.")
    utils.PrintJson(config.Config())
}

func main() {
    application, err := app.InitApp()
    if err != nil {
        panic(err)
    }
    
    slog.Info("[APP] Init completed, next step startup.")
    
    // Start server
    application.Startup()
    
    slog.Info("[APP] Startup completed.")
    
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    // Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
    <-ctx.Done()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := application.Close(ctx); err != nil {
        slog.Error(err.Error())
    }
    glog.Close()
    slog.Info("[App] quit...")
}
