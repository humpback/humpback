package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"humpback/config"
	"humpback/internal/db"
	"humpback/internal/node"
	"humpback/security"
	"humpback/types"

	"github.com/gin-gonic/gin"
)

const tokenRefreshTime = 1 * time.Hour // 在token还有1小时过期时刷新

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
	node.NewAgentManager()

	return hs
}

func doRegister(c *gin.Context) {

	payload := types.RegisterInfo{}

	if c.BindJSON(&payload) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	nodeId, _ := node.MatchNodeWithIpAddress(payload.IpAddress)
	if nodeId == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	node := node.GetNodeInfo(nodeId)
	if node == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	if !node.RegisterInfo.IsRegister &&
		(node.RegisterInfo.Token != payload.Token ||
			time.Now().Unix() > node.RegisterInfo.ExpireAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token or expired"})
		return
	}

	if node.RegisterInfo.IsRegister && node.RegisterInfo.Token != payload.Token {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token or expired"})
		return
	}

	sc := c.MustGet("scheduler").(*HumpbackScheduler)

	response, err := sc.nodeCtrl.HandlerNodeRegister(nodeId, payload.IpAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create certificate"})
	} else {
		c.JSON(http.StatusOK, response)
	}
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

	newToken := sc.nodeCtrl.RefreshNodeToken(nodeId)

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

func tokenAuthMiddleware(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		c.Abort()
		return
	}

	// 验证Bearer token格式
	if len(token) < 7 || token[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	tokenString := token[7:]
	_, err := security.VerifyWorkerToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func (scheduler *HumpbackScheduler) Start(cert *security.CertificateBundle) {
	scheduler.serviceCtrl.RestoreServiceManager()
	scheduler.nodeCtrl.RestoreNodes()
	go func() {
		e := gin.Default()

		e.Use(func(c *gin.Context) {
			c.Set("scheduler", scheduler)
			c.Next()
		})

		e.POST("/api/register", doRegister)

		e.POST("/api/health", tokenAuthMiddleware, doHealth)

		e.GET("/api/config/:name", tokenAuthMiddleware, getConfigByName)

		/*
			e.GET("/mock/nodes", mockNodes)

			e.GET("/nodes", getAllNodes)

			e.GET("/groups", getAllGroups)

			e.GET("/services", getAllServices)

			e.GET("/configs", getAllConfig)

			e.GET("/mock/node/:nodeId/container/:cid/stats", mockContainerStats)

			e.GET("/mock/service/:groupId/gateway", mockGatewayServices)

			e.GET("/mock/service/:groupId/web", mockWebServices)

			e.GET("/mock/configs", mockConfigs)

			e.GET("/mock/service/:groupId/schedule", mockScheduleServices)

			e.GET("/mock/action/:serviceId/:action", mockServiceAction)
		*/

		listeningAddress := fmt.Sprintf(":%s", config.BackendArgs().BackendPort)
		slog.Info("[Scheduler Api] Listening...", "Address", listeningAddress)
		scheduler.httpSrv = &http.Server{
			Addr:      listeningAddress,
			Handler:   e,
			TLSConfig: cert.CreateTLSConfig(true),
		}
		if err := scheduler.httpSrv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
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
