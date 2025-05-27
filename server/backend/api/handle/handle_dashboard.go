package handle

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "humpback/api/middleware"
    "humpback/internal/controller"
)

func RouteDashboard(router *gin.RouterGroup) {
    router.GET("/resource-statistics", dashboardResourceStatistics)
}

func dashboardResourceStatistics(c *gin.Context) {
    userInfo := middleware.GetUserInfo(c)
    result, err := controller.DashboardResourceStatistics(userInfo)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}
