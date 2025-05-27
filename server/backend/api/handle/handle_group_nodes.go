package handle

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "humpback/api/handle/models"
    "humpback/api/middleware"
    "humpback/internal/controller"
)

func RouteGroupNode(router *gin.RouterGroup) {
    router.GET("/list", groupNodes)
    router.PUT("", groupNodeUpdate)
    router.POST("/query", groupNodeQuery)
}

func groupNodes(c *gin.Context) {
    groupId := c.Param("groupId")
    result, err := controller.GroupNodes(groupId, middleware.GetUserInfo(c))
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}

func groupNodeUpdate(c *gin.Context) {
    body := new(models.GroupUpdateNodesReqInfo)
    if !middleware.BindAndCheckBody(c, body) {
        return
    }
    body.GroupId = c.Param("groupId")
    userInfo := middleware.GetUserInfo(c)
    id, err := controller.GroupUpdateNodes(userInfo, middleware.GetNodeChannel(c), middleware.GetGroupInfo(c), body)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, id)
}

func groupNodeQuery(c *gin.Context) {
    groupId := c.Param("groupId")
    queryInfo := new(models.GroupQueryNodesReqInfo)
    if !middleware.BindAndCheckBody(c, queryInfo) {
        return
    }
    result, err := controller.GroupNodesQuery(groupId, middleware.GetUserInfo(c), queryInfo)
    if err != nil {
        middleware.AbortErr(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
}
