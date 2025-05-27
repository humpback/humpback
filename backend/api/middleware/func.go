package middleware

import (
	"github.com/gin-gonic/gin"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/types"
)

type CheckInterface interface {
	Check() error
}

func BindAndCheckBody(c *gin.Context, obj CheckInterface) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		AbortErr(c, response.NewBadRequestErr(locales.CodeRequestParamsInvalid, err.Error()))
		return false
	}
	if err := obj.Check(); err != nil {
		AbortErr(c, err)
		return false
	}
	return true
}

func BindBody(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		AbortErr(c, response.NewBadRequestErr(locales.CodeRequestParamsInvalid, err.Error()))
		return false
	}
	return true
}

func CheckBody(c *gin.Context, body CheckInterface) bool {
	if err := body.Check(); err != nil {
		AbortErr(c, err)
		return false
	}
	return true
}

func AbortErr(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

func GetNodeChannel(c *gin.Context) chan types.NodeSimpleInfo {
	ch, exist := c.Get(KeyNodeEventChannel)
	if !exist {
		return nil
	}
	return ch.(chan types.NodeSimpleInfo)
}

func GetServiceChangeChannel(c *gin.Context) chan types.ServiceChangeInfo {
	ch, exist := c.Get(KeyServiceEventChannel)
	if !exist {
		return nil
	}
	return ch.(chan types.ServiceChangeInfo)
}

func GetGroupInfo(c *gin.Context) *types.NodesGroups {
	info, exist := c.Get(KeyGroupInfo)
	if !exist {
		return nil
	}
	return info.(*types.NodesGroups)
}
