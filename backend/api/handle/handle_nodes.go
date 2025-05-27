package handle

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"humpback/api/handle/models"
	"humpback/api/middleware"
	"humpback/common/response"
	"humpback/internal/controller"
	"humpback/types"
)

func RouteNodes(router *gin.RouterGroup) {
	router.POST("", middleware.CheckAdminPermissions(), nodesCreate)
	router.PUT("/labels", middleware.CheckAdminPermissions(), nodeUpdateLabels)
	router.PUT("/switch", middleware.CheckAdminPermissions(), nodeUpdateSwitch)
	router.GET("/:id/info", node)
	router.POST("/query", nodesQuery)
	router.GET("/list", nodes)
	router.DELETE("/:id", middleware.CheckAdminPermissions(), nodeDelete)
}

func nodesCreate(c *gin.Context) {
	nodes := make(models.NodesCreateReqInfo, 0)
	if !middleware.BindAndCheckBody(c, &nodes) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	if err := controller.NodeCreate(userInfo, middleware.GetNodeChannel(c), nodes); err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewRespSucceed())
}

func nodeUpdateLabels(c *gin.Context) {
	body := new(models.NodeUpdateLabelReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.NodeUpdateLabel(userInfo, middleware.GetNodeChannel(c), body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func nodeUpdateSwitch(c *gin.Context) {
	body := new(models.NodeUpdateSwitchReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.NodeUpdateSwitch(userInfo, middleware.GetNodeChannel(c), body.NodeId, body.Enable)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func node(c *gin.Context) {
	id := c.Param("id")
	info, err := controller.Node(id)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func nodes(c *gin.Context) {
	list, err := controller.Nodes()
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	slices.SortFunc(list, func(a, b *types.Node) int {
		return types.QuerySortOrder(types.SortOrderAsc, strings.ToLower(a.IpAddress), strings.ToLower(b.IpAddress))
	})
	c.JSON(http.StatusOK, list)
}

func nodesQuery(c *gin.Context) {
	queryInfo := new(models.NodeQueryReqInfo)
	if !middleware.BindAndCheckBody(c, queryInfo) {
		return
	}
	result, err := controller.NodesQuery(queryInfo)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func nodeDelete(c *gin.Context) {
	id := c.Param("id")
	userInfo := middleware.GetUserInfo(c)
	if err := controller.NodeDelete(userInfo, middleware.GetNodeChannel(c), id); err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewRespSucceed())
}
