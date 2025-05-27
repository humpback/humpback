package middleware

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "strings"
    "time"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    
    "humpback/common/locales"
    "humpback/common/response"
    "humpback/config"
    "humpback/internal/controller"
    "humpback/types"
)

const (
    KeyErrCodeMap          = "KeyErrCodeMap"
    KeyNodeEventChannel    = "KeyNodeEventChannel"
    KeyServiceEventChannel = "KeyServiceEventChannel"
    KeyGroupInfo           = "KeyGroupInfo"
)

func Log() gin.HandlerFunc {
    return func(c *gin.Context) {
        if (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut) && strings.EqualFold(c.GetHeader("Content-Type"), "application/json") {
            data, err := io.ReadAll(c.Request.Body)
            if err != nil {
                AbortErr(c, err)
                return
            }
            c.Set("Body", data)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
        }
        startTime := time.Now()
        c.Next()
        if strings.HasPrefix(c.Request.URL.Path, "/webapi") {
            slog.Info("request", c.Request.Method, c.Request.URL, "T", time.Now().Sub(startTime).String())
            v, ok := c.Get("Body")
            if ok {
                fmt.Printf("%s\n", v)
            }
        }
    }
}

func CorsCheck() gin.HandlerFunc {
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowAllOrigins = true
    return cors.New(corsConfig)
}

func HandleError() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        length := len(c.Errors)
        if length == 0 {
            return
        }
        var result *response.ErrInfo
        errInfo := c.Errors[length-1]
        if errInfo != nil && errInfo.Err != nil {
            if !errors.As(errInfo.Err, &result) {
                result = response.NewRespServerErr(errInfo.Err.Error())
            }
        } else {
            result = response.NewRespServerErr()
        }
        slog.Error("Response", "url", c.Request.URL, "msg", result.String())
        if result.StatusCode == http.StatusBadRequest {
            errCodeMap := GetErrCodeMap(c)
            if code, ok := errCodeMap[result.Code]; ok {
                result.ReplaceCode(code)
            }
        }
        if result.StatusCode == 0 {
            result.StatusCode = 500
        }
        result.ParseCodeMsg(c.GetHeader("Accept-Language"))
        c.JSON(result.StatusCode, result)
    }
}

func SetErrCodeMap(code map[string]string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.Set(KeyErrCodeMap, code)
        }
    }
}

func GetErrCodeMap(c *gin.Context) map[string]string {
    errMap, exist := c.Get(KeyErrCodeMap)
    if !exist {
        return make(map[string]string)
    }
    errCode, ok := errMap.(map[string]string)
    if !ok {
        return make(map[string]string)
    }
    return errCode
}

func CheckLogin() gin.HandlerFunc {
    return func(c *gin.Context) {
        startUp := strings.ToLower(c.Query("startup")) == "true"
        sessionId, err := GetCookieSession(c)
        if err != nil {
            AbortErr(c, response.NewRespUnauthorized(startUp))
            return
        }
        SetSessionId(c, sessionId)
        userInfo, expired, err := controller.SessionGetAndRefresh(sessionId)
        if err != nil {
            AbortErr(c, err)
            return
        }
        if expired {
            SetCookieSession(c, sessionId, 0)
            AbortErr(c, response.NewRespUnauthorized(startUp))
            return
        }
        SetCookieSession(c, sessionId, int(config.DBArgs().SessionTimeout.Seconds()))
        SetUserInfo(c, userInfo)
    }
}

func CheckAdminPermissions() gin.HandlerFunc {
    return func(c *gin.Context) {
        userInfo := GetUserInfo(c)
        if userInfo.Role != types.UserRoleAdmin && userInfo.Role != types.UserRoleSuperAdmin {
            AbortErr(c, response.NewNoPermissionErr(locales.CodeNoPermission))
            return
        }
    }
}

func CheckInGroup() gin.HandlerFunc {
    return func(c *gin.Context) {
        userInfo := GetUserInfo(c)
        groupId := c.Param("groupId")
        if groupId == "" {
            AbortErr(c, response.NewBadRequestErr(locales.CodeRequestParamsInvalid))
            return
        }
        groupInfo, err := controller.Group(userInfo, groupId)
        if err != nil {
            AbortErr(c, err)
            return
        }
        c.Set(KeyGroupInfo, groupInfo)
    }
}

func SetEventChannel(nodeCh chan types.NodeSimpleInfo, serviceCh chan types.ServiceChangeInfo) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set(KeyNodeEventChannel, nodeCh)
        c.Set(KeyServiceEventChannel, serviceCh)
    }
}
