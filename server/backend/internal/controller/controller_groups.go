package controller

import (
	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/types"
)

func GroupCreate(operator *types.User, reqInfo *models.GroupCreateReqInfo) (string, error) {
	if err := groupCheckNameExist(reqInfo.GroupName, ""); err != nil {
		return "", err
	}
	newInfo := reqInfo.NewGroupInfo()
	users, err := UsersGetByIds(newInfo.Users, true)
	if err != nil {
		return "", err
	}
	teams, err := TeamsGetByIds(newInfo.Teams, true)
	if err != nil {
		return "", err
	}
	if err := db.GroupUpdate(newInfo); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertGroupActivity(&ActivityGroupInfo{
		NewGroupInfo: newInfo,
		NewUsers:     users,
		NewTeams:     teams,
		Action:       types.ActivityActionAdd,
		OperatorInfo: operator,
		OperateAt:    newInfo.UpdatedAt,
	})
	InsertStatisticsCount(&StatisticalCountEvent{
		CreateAt: newInfo.UpdatedAt,
		Type:     types.CountTypeGroup,
		Num:      1,
		UserId:   operator.UserId,
	})
	return newInfo.GroupId, nil
}

func GroupUpdate(operator *types.User, reqInfo *models.GroupUpdateReqInfo) (string, error) {
	oldInfo, err := Group(operator, reqInfo.GroupId)
	if err != nil {
		return "", err
	}
	if !operator.InGroup(oldInfo) {
		return "", response.NewBadRequestErr(locales.CodeGroupNoPermission)
	}
	if err := groupCheckNameExist(reqInfo.GroupName, reqInfo.GroupId); err != nil {
		return "", err
	}
	newInfo := reqInfo.NewGroupInfo(oldInfo)
	users, err := getGroupUsers(oldInfo.Users, newInfo.Users)
	if err != nil {
		return "", err
	}
	teams, err := getGroupTeams(oldInfo.Teams, newInfo.Teams)
	if err != nil {
		return "", err
	}
	if err = db.GroupUpdate(newInfo); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertGroupActivity(&ActivityGroupInfo{
		OldGroupInfo: oldInfo,
		NewGroupInfo: newInfo,
		OldUsers:     users,
		NewUsers:     users,
		OldTeams:     teams,
		NewTeams:     teams,
		Action:       types.ActivityActionUpdate,
		OperatorInfo: operator,
		OperateAt:    newInfo.UpdatedAt,
	})
	return newInfo.GroupId, err
}

func groupCheckNameExist(name, id string) error {
	sameNames, err := db.GroupsGetByName(name, true)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, info := range sameNames {
		if info.GroupId != id {
			return response.NewBadRequestErr(locales.CodeGroupNameAlreadyExist)
		}
	}
	return nil
}

func getGroupUsers(oldUsers, newUsers []string) ([]*types.User, error) {
	var (
		userIdsMap = make(map[string]bool)
		userIds    = make([]string, 0)
	)
	for _, id := range oldUsers {
		userIds = append(userIds, id)
		userIdsMap[id] = true
	}
	for _, id := range newUsers {
		if !userIdsMap[id] {
			userIds = append(userIds, id)
			userIdsMap[id] = true
		}
	}
	users, err := UsersGetByIds(userIds, true)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func getGroupTeams(oldTeams, newTeams []string) ([]*types.Team, error) {
	var (
		teamIdsMap = make(map[string]bool)
		teamIds    = make([]string, 0)
	)
	for _, id := range oldTeams {
		teamIds = append(teamIds, id)
		teamIdsMap[id] = true
	}
	for _, id := range newTeams {
		if !teamIdsMap[id] {
			teamIds = append(teamIds, id)
			teamIdsMap[id] = true
		}
	}
	teams, err := TeamsGetByIds(teamIds, true)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GroupUpdateNodes(operator *types.User, nodeCh chan types.NodeSimpleInfo, groupInfo *types.NodesGroups, info *models.GroupUpdateNodesReqInfo) (string, error) {
	var (
		newGroup = info.NewGroupInfo(groupInfo)
		action   = types.ActivityActionRemoveNode
		nodes    = make([]*types.Node, 0)
		err      error
	)
	if info.IsDelete {
		nodes, err = NodesGetByIds(groupInfo.Nodes, true)
		if err != nil {
			return "", err
		}
	} else {
		action = types.ActivityActionAddNode
		nodes, err = NodesGetByIds(newGroup.Nodes, false)
		if err != nil {
			return "", err
		}
	}
	if err = db.GroupUpdate(newGroup); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	for _, nodeId := range info.Nodes {
		sendNodeEvent(nodeCh, nodeId, "")
	}
	InsertGroupActivity(&ActivityGroupInfo{
		OldGroupInfo: groupInfo,
		NewGroupInfo: newGroup,
		OldNodes:     nodes,
		NewNodes:     nodes,
		Action:       action,
		OperatorInfo: operator,
		OperateAt:    newGroup.UpdatedAt,
	})
	return info.GroupId, nil
}

func Group(userInfo *types.User, id string) (*types.NodesGroups, error) {
	info, err := db.GroupGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeGroupNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	if !userInfo.InGroup(info) {
		return nil, response.NewBadRequestErr(locales.CodeGroupNoPermission)
	}
	return info, nil
}

func Groups(userInfo *types.User) ([]*types.NodesGroups, error) {
	groups, err := db.GroupsGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	result := make([]*types.NodesGroups, 0)
	for _, group := range groups {
		if userInfo.InGroup(group) {
			result = append(result, group)
		}
	}
	return result, nil
}

func GroupQuery(queryInfo *models.GroupQueryReqInfo) (*response.QueryResult[types.NodesGroups], error) {
	groups, err := db.GroupsGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	result := queryInfo.QueryFilter(groups)
	return response.NewQueryResult[types.NodesGroups](
		len(result),
		types.QueryPagination[types.NodesGroups](queryInfo.PageInfo, result),
	), nil
}

func GroupNodes(groupId string, userInfo *types.User) ([]*types.Node, error) {
	groupInfo, err := Group(userInfo, groupId)
	if err != nil {
		return nil, err
	}
	nodes, err := NodesGetByIds(groupInfo.Nodes, true)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func GroupDelete(operator *types.User, svcChan chan types.ServiceChangeInfo, id string) error {
	groupInfo, err := db.GroupGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil
		}
		return response.NewRespServerErr(err.Error())
	}
	users, err := UsersGetByIds(groupInfo.Users, true)
	if err != nil {
		return err
	}
	teams, err := TeamsGetByIds(groupInfo.Teams, true)
	if err != nil {
		return err
	}
	nodes, err := NodesGetByIds(groupInfo.Nodes, true)
	if err != nil {
		return err
	}
	services, err := db.ServicesGetValidByPrefix(id)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if err = db.GroupDeleteAndServiceSoftDelete(id, services); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, service := range services {
		if service.IsEnabled {
			sendServiceEvent(svcChan, service.ServiceId, service.Version, types.ServiceActionDelete)
		}
	}
	InsertGroupActivity(&ActivityGroupInfo{
		OldGroupInfo: groupInfo,
		OldUsers:     users,
		OldTeams:     teams,
		OldNodes:     nodes,
		Action:       types.ActivityActionDelete,
		OperatorInfo: operator,
		OperateAt:    0,
	})
	return nil
}

func GroupNodesQuery(groupId string, userInfo *types.User, queryInfo *models.GroupQueryNodesReqInfo) (*response.QueryResult[types.Node], error) {
	groupInfo, err := Group(userInfo, groupId)
	if err != nil {
		return nil, err
	}
	if !userInfo.InGroup(groupInfo) {
		return nil, response.NewBadRequestErr(locales.CodeGroupNoPermission)
	}
	nodes, err := NodesGetByIds(groupInfo.Nodes, true)
	if err != nil {
		return nil, err
	}
	result := queryInfo.QueryFilter(nodes)
	return response.NewQueryResult[types.Node](
		len(result),
		types.QueryPagination[types.Node](queryInfo.PageInfo, result),
	), nil
}
