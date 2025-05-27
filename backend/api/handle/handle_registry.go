package handle

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"humpback/api/handle/models"
	"humpback/api/middleware"
	"humpback/common/response"
	"humpback/internal/controller"
	"humpback/types"
)

func RouteRegistry(router *gin.RouterGroup) {
	router.POST("", middleware.CheckAdminPermissions(), registryCreate)
	router.PUT("", middleware.CheckAdminPermissions(), registryUpdate)
	router.GET("/:id/info", registry)
	router.POST("/query", registryQuery)
	router.GET("/list", registries)
	router.DELETE("/:id", middleware.CheckAdminPermissions(), registryDelete)
}

func registryCreate(c *gin.Context) {
	body := new(models.RegistryCreateReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.RegistryCreate(userInfo, body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func registryUpdate(c *gin.Context) {
	body := new(models.RegistryUpdateReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.RegistryUpdate(userInfo, body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func registry(c *gin.Context) {
	id := c.Param("id")
	hasAuth := c.Query("hasAuth")
	info, err := controller.Registry(id)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	if strings.ToLower(hasAuth) != "true" {
		info.Username = ""
		info.Password = ""
	}
	c.JSON(http.StatusOK, info)
}

func registries(c *gin.Context) {
	list, err := controller.Registries()
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	for _, registiry := range list {
		registiry.Username = ""
		registiry.Password = ""
	}
	c.JSON(http.StatusOK, list)
}

func registryQuery(c *gin.Context) {
	queryInfo := new(models.RegistryQueryReqInfo)
	if !middleware.BindAndCheckBody(c, queryInfo) {
		return
	}
	list, err := controller.RegistryQuery(queryInfo)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	result := &response.QueryResult[types.QueryRegistry]{
		Total: list.Total,
		List:  make([]*types.QueryRegistry, 0),
	}
	for _, g := range list.List {
		hasAuth := g.Username != "" && g.Password != ""
		g.Username = ""
		g.Password = ""
		result.List = append(result.List, &types.QueryRegistry{
			HasAuth:  hasAuth,
			Registry: g,
		})
	}
	c.JSON(http.StatusOK, result)
}

func registryDelete(c *gin.Context) {
	id := c.Param("id")
	userInfo := middleware.GetUserInfo(c)
	if err := controller.RegistryDelete(userInfo, id); err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewRespSucceed())
}
