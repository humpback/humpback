package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"humpback/api/handle"
	"humpback/api/middleware"
	"humpback/api/static"
	"humpback/config"
	"humpback/types"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine  *gin.Engine
	httpSrv *http.Server
}

func InitRouter(nodeCh chan types.NodeSimpleInfo, serviceCh chan types.ServiceChangeInfo) *Router {
	gin.SetMode(gin.ReleaseMode)
	r := &Router{engine: gin.New()}
	r.setMiddleware(nodeCh, serviceCh)
	r.setRoute()
	return r
}

func (api *Router) Start() {
	go func() {
		listeningAddress := fmt.Sprintf("%s:%s", config.NodeArgs().HostIp, config.NodeArgs().SitePort)
		slog.Info("[Site Api] Listening...", "Address", listeningAddress)
		api.httpSrv = &http.Server{
			Addr:    listeningAddress,
			Handler: api.engine,
		}
		if err := api.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("[Site Api] Listening failed", "Address", listeningAddress, "Error", err)
		}
	}()
}

func (api *Router) Close(c context.Context) error {
	return api.httpSrv.Shutdown(c)
}

func (api *Router) setMiddleware(nodeCh chan types.NodeSimpleInfo, serviceCh chan types.ServiceChangeInfo) {
	api.engine.Use(gin.Recovery(), middleware.Log(), middleware.CorsCheck(), middleware.HandleError(), middleware.SetEventChannel(nodeCh, serviceCh))
}

func (api *Router) setRoute() {
	var routes = map[string]map[string][]any{
		"/webapi": {
			"/common":                 {handle.RouteCommon},
			"/user":                   {handle.RouteUser},
			"/team":                   {middleware.CheckLogin(), handle.RouteTeam},
			"/config":                 {middleware.CheckLogin(), handle.RouteConfig},
			"/registry":               {middleware.CheckLogin(), handle.RouteRegistry},
			"/node":                   {middleware.CheckLogin(), handle.RouteNodes},
			"/group":                  {middleware.CheckLogin(), handle.RouteGroup},
			"/group/:groupId/service": {middleware.CheckLogin(), middleware.CheckInGroup(), handle.RouteService},
			"/group/:groupId/node":    {middleware.CheckLogin(), middleware.CheckInGroup(), handle.RouteGroupNode},
			"/statistics-count":       {middleware.CheckLogin(), handle.RouteStatisticsCount},
			"/activity":               {middleware.CheckLogin(), handle.RouteActivity},
			"/dashboard":              {middleware.CheckLogin(), handle.RouteDashboard},
		},
	}

	for group, list := range routes {
		for path, fList := range list {
			routerGroup := api.engine.Group(fmt.Sprintf("%s%s", group, path), parseSliceAnyToSliceFunc(fList[:len(fList)-1])...)
			groupFunc := fList[len(fList)-1].(func(*gin.RouterGroup))
			groupFunc(routerGroup)
		}
	}
	api.engine.NoRoute(static.WebCache)
}

func parseSliceAnyToSliceFunc(functions []any) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, 0)
	for _, f := range functions {
		if fun, ok := f.(gin.HandlerFunc); ok {
			result = append(result, fun)
		}
	}
	return result
}
