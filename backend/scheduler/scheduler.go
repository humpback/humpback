package scheduler

import (
	"context"
	"encoding/pem"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
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
	security            *security.SecurityManager
	workerCerts         map[string]*security.CertificateBundle
	workerTokens        map[string]string // workerID -> token
	sync.RWMutex
}

func NewHumpbackScheduler(sm *security.SecurityManager) *HumpbackScheduler {
	hs := &HumpbackScheduler{}
	hs.NodeHeartbeatChan = make(chan types.NodeSimpleInfo, 100)
	hs.ContainerChangeChan = make(chan types.ContainerStatus, 100)
	hs.ServiceChangeChan = make(chan types.ServiceChangeInfo, 100)
	hs.serviceCtrl = NewServiceController(hs.NodeHeartbeatChan, hs.ContainerChangeChan, hs.ServiceChangeChan)
	hs.nodeCtrl = NewNodeController(hs.NodeHeartbeatChan, hs.ContainerChangeChan)
	hs.security = sm

	node.NewCacheManager()

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

	if !node.RegisterInfo.IsRegister ||
		node.RegisterInfo.Token != payload.Token ||
		time.Now().UnixMilli() > node.RegisterInfo.ExpireAt {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token or expired"})
		return
	}

	sc := c.MustGet("scheduler").(*HumpbackScheduler)

	sc.Lock()
	defer sc.Unlock()

	var err error
	certBundle, ok := sc.workerCerts[nodeId]

	if !ok {
		certBundle, err = sc.security.CreateCertificateBundle(nodeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create certificate"})
			return
		}
	}

	token, ok := sc.workerTokens[nodeId]
	if !ok {
		token, err = sc.security.GenerateWorkerToken(nodeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
	}

	sc.workerCerts[nodeId] = certBundle
	sc.workerTokens[nodeId] = token

	response := struct {
		CertPEM string `json:"certPem"`
		KeyPEM  string `json:"keyPem"`
		Token   string `json:"token"`
		CAPEM   string `json:"caPem"`
	}{
		CertPEM: string(certBundle.CertPEM),
		KeyPEM:  string(certBundle.KeyPEM),
		Token:   token,
		CAPEM: string(pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: sc.security.CACert.Raw,
		})),
	}

	c.JSON(http.StatusOK, response)

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

	sc.Lock()
	defer sc.Unlock()

	token := sc.workerTokens[nodeId]
	var newToken string
	claims, err := sc.security.VerifyWorkerToken(token)
	if err == nil {
		expiry := claims.ExpiresAt.Time
		if time.Until(expiry) < tokenRefreshTime {
			newToken, err := sc.security.GenerateWorkerToken(nodeId)
			if err == nil {
				sc.workerTokens[nodeId] = newToken
				token = newToken
				log.Printf("Refreshed token for worker %s\n", nodeId)
			}
		}
	}

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
	sc := c.MustGet("scheduler").(*HumpbackScheduler)
	_, err := sc.security.VerifyWorkerToken(tokenString)
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
			Addr:      listeningAddress,
			Handler:   e,
			TLSConfig: cert.CreateTLSConfig(true), // 后面要改成false
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
