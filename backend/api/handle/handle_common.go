package handle

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"humpback/api/middleware"
	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/controller"
)

func RouteCommon(router *gin.RouterGroup) {
	router.GET("/config", ruleConfig)
	router.GET("/search/group-service/by-name", middleware.CheckLogin(), searchGroupServiceByName)
}

func ruleConfig(c *gin.Context) {
	data := map[string]any{
		"RuleLengthLimit": enum.RuleLength,
		"RuleFormat":      enum.RuleFormat,
		"EncryptionKey":   enum.PublicKey,
	}
	c.JSON(http.StatusOK, data)
}

func searchGroupServiceByName(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		middleware.AbortErr(c, response.NewBadRequestErr(locales.CodeNameNotEmpty))
		return
	}
	userInfo := middleware.GetUserInfo(c)
	result, err := controller.SearchGroupAndServcieByName(userInfo, name)
	if err != nil {
		middleware.AbortErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
