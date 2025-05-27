package controller

import (
	"log/slog"
	"slices"

	"golang.org/x/exp/maps"
	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/pkg/utils"
	"humpback/types"
)

func MeLogin(reqInfo *models.UserLoginReqInfo) (*types.User, string, error) {
	user, err := userCheckNamePsd(reqInfo.Username, reqInfo.Password)
	if err != nil {
		return nil, "", err
	}
	sessionInfo := &types.Session{
		SessionId: utils.NewGuidStr(),
		UserId:    user.UserId,
	}
	if err = SessionUpdate(sessionInfo); err != nil {
		return nil, "", err
	}
	InsertAccountActivity(nil, user, types.ActivityActionLogin, 0)
	return user, sessionInfo.SessionId, nil
}

func userCheckNamePsd(username, password string) (*types.User, error) {
	users, err := Users()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Username == username {
			if user.Password == password {
				return user, nil
			}
			return nil, response.NewBadRequestErr(locales.CodePasswordIsWrong)
		}
	}
	return nil, response.NewBadRequestErr(locales.CodeUserNotExist)
}

func MeLogout(userInfo *types.User, sessionId string) error {
	if err := SessionDelete(sessionId); err != nil {
		return err
	}
	InsertAccountActivity(nil, userInfo, types.ActivityActionLogout, 0)
	return nil
}

func MeUpdate(userInfo *types.User, info *models.MeUpdateReqInfo) error {
	currentUser := info.NewUserInfo(userInfo)
	if err := meUpdate(currentUser); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	InsertAccountActivity(userInfo, currentUser, types.ActivityActionUpdate, currentUser.UpdatedAt)
	return nil
}

func meUpdate(userInfo *types.User) error {
	if err := db.UserUpdate(userInfo.UserId, userInfo); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	return nil
}

func MeChangePassword(userInfo *types.User, reqInfo *models.MeChangePasswordReqInfo) error {
	if userInfo.Password != reqInfo.OldPassword {
		return response.NewBadRequestErr(locales.CodeOldPasswordIsWrong)
	}
	userInfo.Password = reqInfo.NewPassword
	userInfo.UpdatedAt = utils.NewActionTimestamp()
	if err := meUpdate(userInfo); err != nil {
		return err
	}
	if err := SessionDeleteByUserId(userInfo.UserId); err != nil {
		return err
	}
	InsertAccountActivity(nil, userInfo, types.ActivityActionChangePassword, userInfo.UpdatedAt)
	return nil
}

func UserCreate(operator *types.User, reqInfo *models.UserCreateReqInfo) (string, error) {
	err := userCheckNameExist(reqInfo.Username, "")
	if err != nil {
		return "", err
	}
	var (
		teams    = make([]*types.Team, 0)
		userInfo = reqInfo.NewUserInfo()
	)
	if len(reqInfo.Teams) > 0 {
		teams, err = TeamsGetByIds(reqInfo.Teams, false)
		if err != nil {
			return "", err
		}
	}
	for _, team := range teams {
		team.Users = append(team.Users, userInfo.UserId)
	}
	id, err := db.UserAndTeamsUpdate(userInfo, teams)
	if err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertUserActivity(&ActivityUserInfo{
		NewUserInfo:  userInfo,
		NewTeams:     teams,
		Action:       types.ActivityActionAdd,
		OperatorInfo: operator,
		OperateAt:    userInfo.UpdatedAt,
	})
	return id, nil
}

func UserUpdate(reqInfo *models.UserUpdateReqInfo, operator *types.User) (string, error) {
	oldUser, err := userCheckRole(reqInfo, operator)
	if err != nil {
		return "", err
	}
	if err = userCheckNameExist(reqInfo.Username, reqInfo.UserId); err != nil {
		return "", err
	}
	newUser, clearSession := reqInfo.NewUserInfo(oldUser)
	allTeams, updateTeams, err := userUpdateCheckTeams(oldUser.Teams, newUser.Teams, newUser.UserId)
	if err != nil {
		return "", err
	}

	id, err := db.UserAndTeamsUpdate(newUser, updateTeams)
	if err != nil {
		return "", response.NewRespServerErr(err.Error())
	}

	if clearSession {
		if err = SessionDeleteByUserId(newUser.UserId); err != nil {
			return "", err
		}
	}
	InsertUserActivity(&ActivityUserInfo{
		OldUserInfo:  oldUser,
		OldTeams:     allTeams,
		NewUserInfo:  newUser,
		NewTeams:     allTeams,
		Action:       types.ActivityActionUpdate,
		OperatorInfo: operator,
		OperateAt:    newUser.UpdatedAt,
	})
	return id, nil
}

func userCheckRole(reqInfo *models.UserUpdateReqInfo, operator *types.User) (*types.User, error) {
	user, err := User(reqInfo.UserId)
	if err != nil {
		return nil, err
	}
	if types.IsAdmin(user.Role) && !types.IsSuperAdmin(operator.Role) {
		return nil, response.NewNoPermissionErr(locales.CodeNoPermission)
	}
	return user, nil
}

func userCheckNameExist(username, userId string) error {
	users, err := db.UsersGetByName(username, true)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, user := range users {
		if user.UserId != userId {
			return response.NewBadRequestErr(locales.CodeUserNameAlreadyExist)
		}
	}
	return nil
}

func userUpdateCheckTeams(oldTeams, newTeams []string, userId string) ([]*types.Team, []*types.Team, error) {
	if len(oldTeams) == 0 && len(newTeams) == 0 {
		return nil, nil, nil
	}

	var teamIdMap = map[string]int{}
	for _, teamId := range oldTeams {
		teamIdMap[teamId] = -1
	}
	for _, teamId := range newTeams {
		if _, ok := teamIdMap[teamId]; ok {
			teamIdMap[teamId] = 0
		} else {
			teamIdMap[teamId] = 1
		}
	}
	teamList, err := TeamsGetByIds(maps.Keys(teamIdMap), false)
	if err != nil {
		return nil, nil, err
	}
	var (
		resultTeams = make([]*types.Team, 0)
		allTeams    = make([]*types.Team, 0)
	)
	for _, team := range teamList {
		allTeams = append(allTeams, team)
		action, ok := teamIdMap[team.TeamId]
		if !ok && action == 0 {
			continue
		}
		index := slices.Index(team.Users, userId)
		if action == -1 && index != -1 {
			team.Users = append(team.Users[:index], team.Users[index+1:]...)
			resultTeams = append(resultTeams, team)
			continue
		}
		if action == 1 && index == -1 {
			team.Users = append(team.Users, userId)
			resultTeams = append(resultTeams, team)
		}
	}
	return allTeams, resultTeams, nil
}

func User(id string) (*types.User, error) {
	info, err := db.UserGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeUserNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return info, nil
}

func Users() ([]*types.User, error) {
	users, err := db.UsersGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return users, nil
}

func UsersQuery(queryInfo *models.UserQueryReqInfo) (*response.QueryResult[types.User], error) {
	users, err := Users()
	if err != nil {
		return nil, err
	}
	result := queryInfo.QueryFilter(users)
	return response.NewQueryResult[types.User](
		len(result),
		types.QueryPagination[types.User](queryInfo.PageInfo, result),
	), nil
}

func UsersGetByIds(ids []string, ignoreNotExist bool) ([]*types.User, error) {
	users, err := db.UsersGetByIds(ids, ignoreNotExist)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeUserNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return users, nil
}

func UsersGetByTeamId(teamId string) ([]*types.User, error) {
	teamInfo, err := Team(teamId)
	if err != nil {
		return nil, err
	}
	users, err := UsersGetByIds(teamInfo.Users, true)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		user.Password = ""
	}
	return users, nil
}

func UserDelete(id string, operator *types.User) error {
	userInfo, err := db.UserGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil
		}
		return response.NewRespServerErr(err.Error())
	}
	if id == operator.UserId {
		return response.NewBadRequestErr(locales.CodeUserIsOwner)
	}
	if userInfo.Role == types.UserRoleSuperAdmin || (userInfo.Role == types.UserRoleAdmin && operator.Role != types.UserRoleSuperAdmin) {
		return response.NewNoPermissionErr(locales.CodeNoPermission)
	}

	teams, err := userDeleteCheckTeams(userInfo.Teams, id)
	if err != nil {
		return err
	}

	groups, err := db.GroupsGetByUserId(id)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}

	for _, group := range groups {
		userIds := make([]string, 0)
		for _, userId := range group.Users {
			if userId != id {
				userIds = append(userIds, userId)
			}
		}
		group.Users = userIds
	}

	if err = db.UserDeleteAndTeamsGroupsUpdate(id, teams, groups); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if err = db.SessionBatchDeleteByUserId(id); err != nil {
		slog.Warn("[User Delete] Clear session failed.", "UserId", id, "UserName", userInfo.Username, "Error", err.Error())
	}
	InsertUserActivity(&ActivityUserInfo{
		OldUserInfo:  userInfo,
		OldTeams:     teams,
		Action:       types.ActivityActionDelete,
		OperatorInfo: operator,
		OperateAt:    0,
	})
	return nil
}

func userDeleteCheckTeams(teams []string, userId string) ([]*types.Team, error) {
	if len(teams) == 0 {
		return nil, nil
	}
	teamList, err := TeamsGetByIds(teams, true)
	if err != nil {
		return nil, err
	}
	var result []*types.Team
	for _, team := range teamList {
		if index := slices.Index(team.Users, userId); index != -1 {
			team.Users = append(team.Users[:index], team.Users[index+1:]...)
			result = append(result, team)
		}
	}
	return result, nil
}
