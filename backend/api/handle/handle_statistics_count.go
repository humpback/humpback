package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"humpback/api/handle/models"
	"humpback/api/middleware"
	"humpback/internal/controller"
)

func RouteStatisticsCount(router *gin.RouterGroup) {
	router.POST("/query", statisticsCountQuery)
}

func statisticsCountQuery(c *gin.Context) {
	body := new(models.StatisticsCountQueryReqInfo)
	if !middleware.BindAndCheckBody(c, body) {
		return
	}
	body.UserInfo = middleware.GetUserInfo(c)
	result, err := controller.StatisticsCountQuery(body)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
