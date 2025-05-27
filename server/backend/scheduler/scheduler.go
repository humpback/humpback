package scheduler

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    
    "humpback/config"
    "humpback/internal/db"
    "humpback/internal/node"
    "humpback/types"
    
    "github.com/gin-gonic/gin"
)

type HumpbackScheduler struct {
    httpSrv             *http.Server
    nodeCtrl            *NodeController
    serviceCtrl         *ServiceController
    NodeHeartbeatChan   chan types.NodeSimpleInfo
    ContainerChangeChan chan types.ContainerStatus
    ServiceChangeChan   chan types.ServiceChangeInfo
}

func NewHumpbackScheduler() *HumpbackScheduler {
    hs := &HumpbackScheduler{}
    hs.NodeHeartbeatChan = make(chan types.NodeSimpleInfo, 100)
    hs.ContainerChangeChan = make(chan types.ContainerStatus, 100)
    hs.ServiceChangeChan = make(chan types.ServiceChangeInfo, 100)
    hs.serviceCtrl = NewServiceController(hs.NodeHeartbeatChan, hs.ContainerChangeChan, hs.ServiceChangeChan)
    hs.nodeCtrl = NewNodeController(hs.NodeHeartbeatChan, hs.ContainerChangeChan)
    
    node.NewCacheManager()
    
    return hs
}

func doHealth(c *gin.Context) {
    payload := types.HealthInfo{}
    
    if c.BindJSON(&payload) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
        return
    }
    nodeId, ip := node.MatchNodeWithIpAddress(payload.HostInfo.IpAddress)
    if nodeId == "" {
        c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
        return
    }
    payload.NodeId = nodeId
    payload.IpAddress = ip
    sc := c.MustGet("scheduler").(*HumpbackScheduler)
    sc.nodeCtrl.HeartBeat(payload)
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (scheduler *HumpbackScheduler) Start() {
    scheduler.serviceCtrl.RestoreServiceManager()
    scheduler.nodeCtrl.RestoreNodes()
    go func() {
        e := gin.Default()
        
        e.Use(func(c *gin.Context) {
            c.Set("scheduler", scheduler)
            c.Next()
        })
        
        e.POST("/api/health", doHealth)
        
        e.GET("/api/config/:name", getConfigByName)
        
        e.GET("/mock/nodes", mockNodes)
        
        e.GET("/nodes", getAllNodes)
        
        e.GET("/groups", getAllGroups)
        
        e.GET("/services", getAllServices)
        
        e.GET("/configs", getAllConfig)
        
        e.GET("/mock/service/:groupId/gateway", mockGatewayServices)
        
        e.GET("/mock/service/:groupId/web", mockWebServices)
        
        e.GET("/mock/configs", mockConfigs)
        
        e.GET("/mock/service/:groupId/schedule", mockScheduleServices)
        
        e.GET("/mock/action/:serviceId/:action", mockServiceAction)
        
        listeningAddress := fmt.Sprintf("%s:%s", config.NodeArgs().HostIp, config.BackendArgs().BackendPort)
        slog.Info("[Scheduler Api] Listening...", "Address", listeningAddress)
        scheduler.httpSrv = &http.Server{
            Addr:    listeningAddress,
            Handler: e,
        }
        if err := scheduler.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("[Scheduler Api] Listening failed", "Address", listeningAddress, "Error", err)
        }
    }()
}

func (scheduler *HumpbackScheduler) Close(c context.Context) error {
    return scheduler.httpSrv.Shutdown(c)
}

func getConfigByName(c *gin.Context) {
    configName := c.Param("name")
    configValue, err := db.ConfigsGetByName(configName, true)
    if err != nil {
        c.String(http.StatusInternalServerError, err.Error())
    } else if len(configValue) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
    } else {
        c.String(http.StatusOK, configValue[0].ConfigValue)
    }
}
