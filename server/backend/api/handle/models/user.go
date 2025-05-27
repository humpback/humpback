package models

import (
    "slices"
    "strings"
    
    "github.com/jinzhu/copier"
    "humpback/common/enum"
    "humpback/common/locales"
    "humpback/common/response"
    "humpback/common/verify"
    "humpback/pkg/utils"
    "humpback/types"
)

type UserLoginReqInfo struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (u *UserLoginReqInfo) Check() error {
    u.Username = utils.RSADecrypt(u.Username)
    u.Password = utils.RSADecrypt(u.Password)
    
    if err := verify.CheckRequiredAndLengthLimit(u.Username, enum.LimitUsername.Min, enum.LimitUsername.Max, locales.CodeUserNameNotEmpty, locales.CodeUserNameLimitLength); err != nil {
        return err
    }
    if err := verify.CheckRequiredAndLengthLimit(u.Password, enum.LimitPassword.Min, enum.LimitPassword.Max, locales.CodePasswordNotEmpty, locales.CodePasswordLimitLength); err != nil {
        return err
    }
    return nil
}

type MeUpdateReqInfo struct {
    Username    string `json:"username"`
    Email       string `json:"email"`
    Phone       string `json:"phone"`
    Description string `json:"description"`
}

func (u *MeUpdateReqInfo) Check() error {
    if err := verify.CheckRequiredAndLengthLimit(u.Username, enum.LimitUsername.Min, enum.LimitUsername.Max, locales.CodeUserNameNotEmpty, locales.CodeUserNameLimitLength); err != nil {
        return err
    }
    if err := verify.CheckLengthLimit(u.Email, 0, enum.LimitEmail.Max, locales.CodeEmailLimitLength); err != nil {
        return err
    }
    if u.Email != "" {
        if err := verify.CheckEmail(u.Email); err != nil {
            return err
        }
    }
    if err := verify.CheckLengthLimit(u.Phone, 0, enum.LimitPhone.Max, locales.CodePhoneLimitLength); err != nil {
        return err
    }
    if u.Phone != "" {
        if err := verify.CheckPhone(u.Phone); err != nil {
            return err
        }
    }
    if err := verify.CheckLengthLimit(u.Description, 0, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
        return err
    }
    return nil
}

func (u *MeUpdateReqInfo) NewUserInfo(userInfo *types.User) *types.User {
    newUser := new(types.User)
    copier.Copy(newUser, userInfo)
    newUser.Username = u.Username
    newUser.Email = u.Email
    newUser.Phone = u.Phone
    newUser.Description = u.Description
    newUser.UpdatedAt = utils.NewActionTimestamp()
    return newUser
}

type MeChangePasswordReqInfo struct {
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword"`
}

func (u *MeChangePasswordReqInfo) Check() error {
    u.OldPassword = utils.RSADecrypt(u.OldPassword)
    u.NewPassword = utils.RSADecrypt(u.NewPassword)
    if err := verify.CheckRequiredAndLengthLimit(u.OldPassword, enum.LimitPassword.Min, enum.LimitPassword.Max, locales.CodeOldPasswordNotEmpty, locales.CodeOldPasswordLimitLength); err != nil {
        return err
    }
    if err := verify.CheckRequiredAndLengthLimit(u.NewPassword, enum.LimitPassword.Min, enum.LimitPassword.Max, locales.CodeNewPasswordNotEmpty, locales.CodeNewPasswordLimitLength); err != nil {
        return err
    }
    return nil
}

type UserCreateReqInfo struct {
    Username    string         `json:"username"`
    Password    string         `json:"password"`
    Role        types.UserRole `json:"role"`
    Email       string         `json:"email"`
    Phone       string         `json:"phone"`
    Description string         `json:"description"`
    Teams       []string       `json:"teams"`
}

func (u *UserCreateReqInfo) Check() error {
    u.Password = utils.RSADecrypt(u.Password)
    
    if err := verify.CheckRequiredAndLengthLimit(u.Username, enum.LimitUsername.Min, enum.LimitUsername.Max, locales.CodeUserNameNotEmpty, locales.CodeUserNameLimitLength); err != nil {
        return err
    }
    if err := verify.CheckLengthLimit(u.Email, 0, enum.LimitEmail.Max, locales.CodeEmailLimitLength); err != nil {
        return err
    }
    if u.Email != "" {
        if err := verify.CheckEmail(u.Email); err != nil {
            return err
        }
    }
    if err := verify.CheckLengthLimit(u.Phone, 0, enum.LimitPhone.Max, locales.CodePhoneLimitLength); err != nil {
        return err
    }
    if u.Phone != "" {
        if err := verify.CheckPhone(u.Phone); err != nil {
            return err
        }
    }
    if err := verify.CheckLengthLimit(u.Description, 0, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
        return err
    }
    
    if err := verify.CheckRequiredAndLengthLimit(u.Password, enum.LimitPassword.Min, enum.LimitPassword.Max, locales.CodePasswordNotEmpty, locales.CodePasswordLimitLength); err != nil {
        return err
    }
    
    if u.Role != types.UserRoleAdmin && u.Role != types.UserRoleUser {
        return response.NewBadRequestErr(locales.CodeUserRoleIsInvalid)
    }
    if len(u.Teams) == 0 {
        u.Teams = make([]string, 0)
    }
    return nil
}

func (u *UserCreateReqInfo) CheckCreateRole(operator *types.User) error {
    if operator.Role == types.UserRoleAdmin && u.Role != types.UserRoleUser {
        return response.NewBadRequestErr(locales.CodeUserRoleIsInvalid)
    }
    return nil
}

func (u *UserCreateReqInfo) NewUserInfo() *types.User {
    t := utils.NewActionTimestamp()
    return &types.User{
        UserId:      utils.NewGuidStr(),
        Username:    u.Username,
        Password:    u.Password,
        Role:        u.Role,
        Email:       u.Email,
        Phone:       u.Phone,
        Description: u.Description,
        CreatedAt:   t,
        UpdatedAt:   t,
        Teams:       u.Teams,
    }
}

type UserUpdateReqInfo struct {
    UserId string `json:"userId"`
    UserCreateReqInfo
}

func (u *UserUpdateReqInfo) Check() error {
    if err := verify.CheckIsEmpty(u.UserId, locales.CodeUserIdNotEmpty); err != nil {
        return err
    }
    return u.UserCreateReqInfo.Check()
}

func (u *UserUpdateReqInfo) CheckUpdateRole(operator *types.User) error {
    if u.UserId == operator.UserId {
        return response.NewBadRequestErr(locales.CodeUserIsOwner)
    }
    if operator.Role == types.UserRoleAdmin && u.Role != types.UserRoleUser {
        return response.NewBadRequestErr(locales.CodeUserRoleIsInvalid)
    }
    return nil
}

func (u *UserUpdateReqInfo) NewUserInfo(oldUserInfo *types.User) (*types.User, bool) {
    userInfo := &types.User{
        UserId:      u.UserId,
        Username:    u.Username,
        Password:    u.Password,
        Role:        u.Role,
        Email:       u.Email,
        Phone:       u.Phone,
        Description: u.Description,
        CreatedAt:   oldUserInfo.CreatedAt,
        UpdatedAt:   utils.NewActionTimestamp(),
        Teams:       u.Teams,
    }
    if u.Username != oldUserInfo.Username ||
        u.Password != oldUserInfo.Password ||
        u.Email != oldUserInfo.Email ||
        u.Phone != oldUserInfo.Phone ||
        u.Description != oldUserInfo.Description ||
        u.Role != oldUserInfo.Role ||
        len(u.Teams) != len(oldUserInfo.Teams) {
        return userInfo, true
    }
    for _, teamId := range u.Teams {
        if index := slices.Index(oldUserInfo.Teams, teamId); index == -1 {
            return userInfo, true
        }
    }
    return userInfo, false
}

type UserQueryFilterInfo struct {
    Role int `json:"role"`
}

type UserQueryReqInfo struct {
    types.QueryInfo
    FilterInfo *UserQueryFilterInfo `json:"-"`
}

func (u *UserQueryReqInfo) Check() error {
    u.CheckBase()
    u.FilterInfo = new(UserQueryFilterInfo)
    if err := ParseMapToStructConvert(u.Filter, u.FilterInfo); err != nil {
        return err
    }
    
    if u.Keywords != "" && u.Mode != "keywords" {
        return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
    }
    return nil
}

func (u *UserQueryReqInfo) QueryFilter(users []*types.User) []*types.User {
    result := make([]*types.User, 0)
    for _, user := range users {
        user.Password = ""
        if u.filter(user) {
            result = append(result, user)
        }
    }
    u.sort(result)
    return result
}

func (u *UserQueryReqInfo) filter(info *types.User) bool {
    if u.FilterInfo != nil && u.FilterInfo.Role != 0 && int(info.Role) != u.FilterInfo.Role {
        return false
    }
    return strings.Contains(strings.ToLower(info.Username), strings.ToLower(u.Keywords)) || strings.Contains(strings.ToLower(info.Email), strings.ToLower(u.Keywords)) || strings.Contains(strings.ToLower(info.Phone), strings.ToLower(u.Keywords))
}

func (u *UserQueryReqInfo) sort(list []*types.User) []*types.User {
    var sortField = []string{"username", "updatedAt", "createdAt"}
    if u.SortInfo == nil || u.SortInfo.Field == "" || slices.Index(sortField, u.SortInfo.Field) == -1 {
        return list
    }
    slices.SortFunc(list, func(a, b *types.User) int {
        switch u.SortInfo.Field {
        case "username":
            return types.QuerySortOrder(u.SortInfo.Order, strings.ToLower(a.Username), strings.ToLower(b.Username))
        case "updatedAt":
            return types.QuerySortOrder(u.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
        case "createdAt":
            return types.QuerySortOrder(u.SortInfo.Order, a.CreatedAt, b.CreatedAt)
        }
        return 1
    })
    return list
}
