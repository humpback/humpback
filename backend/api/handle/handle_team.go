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

func RouteTeam(router *gin.RouterGroup) {
	router.POST("", middleware.CheckAdminPermissions(), teamCreate)
	router.PUT("", middleware.CheckAdminPermissions(), teamUpdate)
	router.GET("/:id/info", middleware.CheckAdminPermissions(), team)
	router.GET("/list", teams)
	router.POST("/query", middleware.CheckAdminPermissions(), teamsQuery)
	router.GET("/query-by-user/:userId", middleware.CheckAdminPermissions(), teamsByUserId)
	router.DELETE("/:id", middleware.CheckAdminPermissions(), teamDelete)
}

func teamCreate(c *gin.Context) {
	body := new(models.TeamCreateReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.TeamCreate(userInfo, body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func teamUpdate(c *gin.Context) {
	body := new(models.TeamUpdateReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	userInfo := middleware.GetUserInfo(c)
	id, err := controller.TeamUpdate(userInfo, body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func team(c *gin.Context) {
	id := c.Param("id")
	info, err := controller.Team(id)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func teams(c *gin.Context) {
	result, err := controller.Teams()
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	slices.SortFunc(result, func(a, b *types.Team) int {
		return types.QuerySortOrder(types.SortOrderAsc, strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	c.JSON(http.StatusOK, result)
}

func teamsQuery(c *gin.Context) {
	queryInfo := new(models.TeamQueryReqInfo)
	if !middleware.BindAndCheckBody(c, queryInfo) {
		return
	}
	result, err := controller.TeamQuery(queryInfo)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func teamsByUserId(c *gin.Context) {
	userId := c.Param("userId")
	teams, err := controller.TeamsGetByUserId(userId)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, teams)
}

func teamDelete(c *gin.Context) {
	id := c.Param("id")
	userInfo := middleware.GetUserInfo(c)
	if err := controller.TeamDelete(userInfo, id); err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewRespSucceed())
}
