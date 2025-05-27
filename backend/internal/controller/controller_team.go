package controller

import (
	"slices"

	"golang.org/x/exp/maps"

	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/types"
)

func TeamCreate(operator *types.User, reqInfo *models.TeamCreateReqInfo) (string, error) {
	err := teamCheckNameExist(reqInfo.Name, "")
	if err != nil {
		return "", err
	}
	var (
		users = make([]*types.User, 0)
		team  = reqInfo.NewTeamInfo()
	)
	if len(reqInfo.Users) > 0 {
		users, err = UsersGetByIds(reqInfo.Users, false)
		if err != nil {
			return "", err
		}
	}
	for _, user := range users {
		user.Teams = append(user.Teams, team.TeamId)
	}
	id, err := db.TeamUpdateAndUsers(team, users)
	if err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertTeamActivity(&ActivityTeamInfo{
		NewTeamInfo:  team,
		NewUsers:     users,
		Action:       types.ActivityActionAdd,
		OperatorInfo: operator,
		OperateAt:    team.UpdatedAt,
	})
	return id, nil
}

func TeamUpdate(operator *types.User, reqInfo *models.TeamUpdateReqInfo) (string, error) {
	if err := teamCheckNameExist(reqInfo.Name, reqInfo.TeamId); err != nil {
		return "", err
	}

	oldTeam, err := Team(reqInfo.TeamId)
	if err != nil {
		return "", err
	}

	newTeam := reqInfo.NewTeamInfo(oldTeam)
	allUsers, updateUsers, err := teamUpdateCheckUsers(oldTeam.Users, newTeam.Users, newTeam.TeamId)
	if err != nil {
		return "", err
	}

	id, err := db.TeamUpdateAndUsers(newTeam, updateUsers)
	if err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertTeamActivity(&ActivityTeamInfo{
		OldTeamInfo:  oldTeam,
		NewTeamInfo:  newTeam,
		OldUsers:     allUsers,
		NewUsers:     allUsers,
		Action:       types.ActivityActionUpdate,
		OperatorInfo: operator,
		OperateAt:    newTeam.UpdatedAt,
	})
	return id, nil
}

func teamCheckNameExist(name, id string) error {
	teams, err := db.TeamsGetByName(name, true)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, team := range teams {
		if team.TeamId != id {
			return response.NewBadRequestErr(locales.CodeTeamAlreadyExist)
		}
	}
	return nil
}

func teamUpdateCheckUsers(oldUsers, newUsers []string, teamId string) ([]*types.User, []*types.User, error) {
	if len(oldUsers) == 0 && len(newUsers) == 0 {
		return nil, nil, nil
	}

	var userIdMap = map[string]int{}
	for _, userId := range oldUsers {
		userIdMap[userId] = -1
	}
	for _, userId := range newUsers {
		if _, ok := userIdMap[userId]; ok {
			userIdMap[userId] = 0
		} else {
			userIdMap[userId] = 1
		}
	}
	userList, err := UsersGetByIds(maps.Keys(userIdMap), false)
	if err != nil {
		return nil, nil, err
	}
	var (
		resultUsers = make([]*types.User, 0)
		allUsers    = make([]*types.User, 0)
	)
	for _, user := range userList {
		allUsers = append(allUsers, user)
		action, ok := userIdMap[user.UserId]
		if !ok && action == 0 {
			continue
		}
		index := slices.Index(user.Teams, teamId)
		if action == -1 && index != -1 {
			user.Teams = append(user.Teams[:index], user.Teams[index+1:]...)
			resultUsers = append(resultUsers, user)
			continue
		}
		if action == 1 && index == -1 {
			user.Teams = append(user.Teams, teamId)
			resultUsers = append(resultUsers, user)
		}
	}
	return allUsers, resultUsers, nil
}

func Team(id string) (*types.Team, error) {
	info, err := db.TeamGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeTeamNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return info, nil
}

func TeamsGetByUserId(userId string) ([]*types.Team, error) {
	userInfo, err := User(userId)
	if err != nil {
		return nil, err
	}
	teams, err := TeamsGetByIds(userInfo.Teams, true)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func TeamsGetByIds(ids []string, ignoreNotExist bool) ([]*types.Team, error) {
	teams, err := db.TeamsGetByIds(ids, ignoreNotExist)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeTeamNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return teams, nil
}

func Teams() ([]*types.Team, error) {
	teams, err := db.TeamsGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return teams, nil
}

func TeamQuery(queryInfo *models.TeamQueryReqInfo) (*response.QueryResult[types.Team], error) {
	users, err := Teams()
	if err != nil {
		return nil, err
	}
	result := queryInfo.QueryFilter(users)
	return response.NewQueryResult[types.Team](
		len(result),
		types.QueryPagination[types.Team](queryInfo.PageInfo, result),
	), nil
}

func TeamDelete(operator *types.User, id string) error {
	info, err := db.TeamGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil
		}
		return response.NewRespServerErr(err.Error())
	}

	users, err := teamDeleteCheckUsers(info.Users, id)
	if err != nil {
		return err
	}

	groups, err := db.GroupsGetByTeamId(id)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}

	for _, group := range groups {
		teamIds := make([]string, 0)
		for _, teamId := range group.Teams {
			if teamId != id {
				teamIds = append(teamIds, teamId)
			}
		}
		group.Teams = teamIds
	}

	if err = db.TeamDeleteAndUsersGroupsUpdate(id, users, groups); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	InsertTeamActivity(&ActivityTeamInfo{
		OldTeamInfo:  info,
		OldUsers:     users,
		Action:       types.ActivityActionDelete,
		OperatorInfo: operator,
		OperateAt:    0,
	})
	return nil
}

func teamDeleteCheckUsers(users []string, teamId string) ([]*types.User, error) {
	if len(users) == 0 {
		return nil, nil
	}
	userList, err := UsersGetByIds(users, true)
	if err != nil {
		return nil, err
	}
	var result []*types.User
	for _, user := range userList {
		if index := slices.Index(user.Teams, teamId); index != -1 {
			user.Teams = append(user.Teams[:index], user.Teams[index+1:]...)
			result = append(result, user)
		}
	}
	return result, nil
}
