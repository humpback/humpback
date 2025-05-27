package handle

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "humpback/api/handle/models"
    "humpback/api/middleware"
    "humpback/common/response"
    "humpback/internal/controller"
)

func RouteService(router *gin.RouterGroup) {
    router.POST("/query", serviceQuery)
    router.GET("/total", serviceTotal)
    router.POST("", serviceCreate)
    router.POST("/clone", serviceClone)
    router.PUT("", serviceUpdate)
    router.GET("/:serviceId/info", serviceInfo)
    router.PUT("/operate", serviceOperate)
    router.DELETE("/:serviceId", serviceDelete)
    router.PUT("/:serviceId/instance/operate", instanceOperate)
    router.POST("/:serviceId/instance/logs", instanceQueryLogs)
    router.POST("/:serviceId/instance/performances", instancePerformances)
}

func serviceQuery(c *gin.Context) {
    queryInfo := new(models.ServiceQueryReqInfo)
    if !middleware.BindAndCheckBody(c, queryInfo) {
        return
    }
    groupId := c.Param("groupId")
    queryInfo.UserInfo = middleware.GetUserInfo(c)
    result, err := controller.ServiceQuery(groupId, queryInfo)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceTotal(c *gin.Context) {
    groupId := c.Param("groupId")
    result, err := controller.ServiceTotal(groupId)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceInfo(c *gin.Context) {
    groupId := c.Param("groupId")
    serviceId := c.Param("serviceId")
    result, err := controller.Service(groupId, serviceId)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceCreate(c *gin.Context) {
    body := new(models.ServiceCreateReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    groupInfo := middleware.GetGroupInfo(c)
    result, err := controller.ServiceCreate(userInfo, groupInfo, body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceClone(c *gin.Context) {
    body := new(models.ServiceCloneReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    groupInfo := middleware.GetGroupInfo(c)
    result, err := controller.ServiceClone(userInfo, groupInfo, body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceUpdate(c *gin.Context) {
    body := new(models.ServiceUpdateReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    groupInfo := middleware.GetGroupInfo(c)
    result, err := controller.ServiceUpdate(userInfo, groupInfo, middleware.GetServiceChangeChannel(c), body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceOperate(c *gin.Context) {
    body := new(models.ServiceOperateReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    groupInfo := middleware.GetGroupInfo(c)
    result, err := controller.ServiceOperate(userInfo, groupInfo, middleware.GetServiceChangeChannel(c), body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func serviceDelete(c *gin.Context) {
    serviceId := c.Param("serviceId")
    groupId := c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    groupInfo := middleware.GetGroupInfo(c)
    if err := controller.ServiceSoftDelete(userInfo, groupInfo, middleware.GetServiceChangeChannel(c), groupId, serviceId); err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, response.NewRespSucceed())
}

func instanceOperate(c *gin.Context) {
    body := new(models.InstanceOperateReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    body.ServiceId = c.Param("serviceId")
    groupInfo := middleware.GetGroupInfo(c)
    userInfo := middleware.GetUserInfo(c)
    if err := controller.InstanceOperate(userInfo, groupInfo, body); err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, response.NewRespSucceed())
}

func instanceQueryLogs(c *gin.Context) {
    body := new(models.InstanceLogsReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    body.ServiceId = c.Param("serviceId")
    groupInfo := middleware.GetGroupInfo(c)
    logs, err := controller.InstanceQueryLogs(groupInfo, body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, logs)
}

func instancePerformances(c *gin.Context) {
    body := new(models.InstancesPerformanceReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    body.ServiceId = c.Param("serviceId")
    groupInfo := middleware.GetGroupInfo(c)
    result, err := controller.InsatncePerformances(groupInfo, body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}
