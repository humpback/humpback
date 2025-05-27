package scheduler

import (
	"net/http"
	"strings"
	"time"

	"humpback/internal/db"
	"humpback/pkg/utils"
	"humpback/types"

	"github.com/gin-gonic/gin"
)

func getAllNodes(c *gin.Context) {

	nodes, err := db.GetDataAll[types.Node](db.BucketNodes)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, nodes)
	}
}

func getAllGroups(c *gin.Context) {

	nodes, err := db.GetDataAll[types.NodesGroups](db.BucketNodesGroups)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, nodes)
	}
}

func mockNodes(c *gin.Context) {

	node1 := &types.Node{
		NodeId:    utils.GenerateRandomStringWithLength(8),
		Name:      "hb001",
		IpAddress: "172.30.198.172",
		Port:      8018,
		Status:    "Online",
		IsEnable:  true,
		CreatedAt: time.Now().Unix(),
	}

	db.SaveData(db.BucketNodes, node1.NodeId, node1)

	node2 := &types.Node{
		NodeId:    utils.GenerateRandomStringWithLength(8),
		Name:      "hb002",
		IpAddress: "172.16.41.22",
		Port:      8018,
		Status:    "Online",
		IsEnable:  true,
		CreatedAt: time.Now().Unix(),
	}

	db.SaveData(db.BucketNodes, node2.NodeId, node2)

	node3 := &types.Node{
		NodeId:    utils.GenerateRandomStringWithLength(8),
		Name:      "hb003",
		IpAddress: "172.16.41.23",
		Port:      8018,
		Status:    "Online",
		IsEnable:  true,
		CreatedAt: time.Now().Unix(),
	}

	db.SaveData(db.BucketNodes, node3.NodeId, node3)

	group1 := &types.NodesGroups{
		GroupId:   utils.GenerateRandomStringWithLength(8),
		GroupName: "GroupTest",
		CreatedAt: time.Now().Unix(),
		Nodes:     []string{node1.NodeId, node2.NodeId, node3.NodeId},
	}

	db.SaveData(db.BucketNodesGroups, group1.GroupId, group1)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getAllServices(c *gin.Context) {

	svcs, err := db.GetDataAll[types.Service](db.BucketServices)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, svcs)
	}
}

func getAllConfig(c *gin.Context) {

	svcs, err := db.GetDataAll[types.Config](db.BucketConfigs)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, svcs)
	}
}

func mockGatewayServices(c *gin.Context) {
	gId := c.Param("groupId")
	svc := &types.Service{
		ServiceId:   gId + utils.GenerateRandomStringWithLength(8),
		ServiceName: "Gateway",
		Version:     utils.GenerateRandomStringWithLength(5),
		IsEnabled:   true,
		Status:      types.ServiceStatusNotReady,
		GroupId:     gId,
		CreatedAt:   time.Now().Unix(),
		Deployment: &types.Deployment{
			Type: types.DeployTypeBackground,
			Mode: types.DeployModeGlobal,
		},
	}

	db.SaveData(db.BucketServices, svc.ServiceId, svc)

	sc := c.MustGet("scheduler").(*HumpbackScheduler)

	svcChange := types.ServiceChangeInfo{
		ServiceId: svc.ServiceId,
		Version:   svc.Version,
	}

	sc.ServiceChangeChan <- svcChange

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func mockWebServices(c *gin.Context) {
	gId := c.Param("groupId")
	svc := &types.Service{
		ServiceId:   gId + utils.GenerateRandomStringWithLength(8),
		ServiceName: "Http Web Service",
		Version:     utils.GenerateRandomStringWithLength(5),
		IsEnabled:   true,
		Status:      types.ServiceStatusNotReady,
		GroupId:     gId,
		CreatedAt:   time.Now().Unix(),
		Deployment: &types.Deployment{
			Type:     types.DeployTypeBackground,
			Mode:     types.DeployModeReplicate,
			Replicas: 2,
		},
		Meta: &types.ServiceMetaDocker{
			RegistryDomain: "docker.io",
			Image:          "nginx:latest",
			Network: &types.NetworkInfo{
				Mode: types.NetworkModeBridge,
				Ports: []*types.PortInfo{
					{
						HostPort:      0,
						ContainerPort: 80,
					},
				},
			},
			RestartPolicy: &types.RestartPolicy{
				Mode: types.RestartPolicyModeAlways,
			},
			Envs: []string{"a=b", "name={name}"},
		},
	}

	db.SaveData(db.BucketServices, svc.ServiceId, svc)

	sc := c.MustGet("scheduler").(*HumpbackScheduler)

	svcChange := types.ServiceChangeInfo{
		ServiceId: svc.ServiceId,
		Version:   svc.Version,
	}

	sc.ServiceChangeChan <- svcChange

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func mockScheduleServices(c *gin.Context) {
	gId := c.Param("groupId")
	svc := &types.Service{
		ServiceId:   gId + utils.GenerateRandomStringWithLength(8),
		ServiceName: "Test Schedule Service",
		Version:     utils.GenerateRandomStringWithLength(5),
		IsEnabled:   true,
		Status:      types.ServiceStatusNotReady,
		GroupId:     gId,
		CreatedAt:   time.Now().Unix(),
		Deployment: &types.Deployment{
			Type:     types.DeployTypeSchedule,
			Mode:     types.DeployModeReplicate,
			Replicas: 2,
			Schedule: &types.ScheduleInfo{
				Timeout: "30s",
				Rules:   []string{"*/5 * * * *"},
			},
		},
		Meta: &types.ServiceMetaDocker{
			RegistryDomain: "docker.io",
			Image:          "nginx:latest",
			Network: &types.NetworkInfo{
				Mode: types.NetworkModeBridge,
				Ports: []*types.PortInfo{
					{
						HostPort:      0,
						ContainerPort: 80,
					},
				},
			},
			RestartPolicy: &types.RestartPolicy{
				Mode: types.RestartPolicyModeAlways,
			},
		},
	}

	db.SaveData(db.BucketServices, svc.ServiceId, svc)

	sc := c.MustGet("scheduler").(*HumpbackScheduler)

	svcChange := types.ServiceChangeInfo{
		ServiceId: svc.ServiceId,
		Version:   svc.Version,
	}

	sc.ServiceChangeChan <- svcChange

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func mockServiceAction(c *gin.Context) {
	svcId := c.Param("serviceId")
	action := c.Param("action")

	svc, _ := db.GetDataById[types.Service](db.BucketServices, svcId)

	if svc != nil {

		if strings.EqualFold(action, types.ServiceActionDelete) {
			svc.IsDelete = true
		} else if strings.EqualFold(action, types.ServiceActionDisable) {
			svc.IsEnabled = false
		} else {
			svc.Action = action
		}

		db.SaveData(db.BucketServices, svc.ServiceId, svc)

		svcChange := types.ServiceChangeInfo{
			ServiceId: svc.ServiceId,
			Version:   svc.Version,
			Action:    action,
		}

		sc := c.MustGet("scheduler").(*HumpbackScheduler)
		sc.ServiceChangeChan <- svcChange
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func mockConfigs(c *gin.Context) {
	config1 := &types.Config{
		ConfigId:    utils.GenerateRandomStringWithLength(8),
		ConfigName:  "name",
		ConfigType:  types.ConfigTypeStatic,
		ConfigValue: "james yang",
		CreatedAt:   time.Now().Unix(),
	}

	db.SaveData(db.BucketConfigs, config1.ConfigId, config1)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
