package controller

import (
	"strings"

	"humpback/types"
)

func SearchGroupAndServcieByName(userInfo *types.User, name string) (map[string][]*types.SearchGroupServiceSimpleInfo, error) {
	groups, err := Groups(userInfo)
	if err != nil {
		return nil, err
	}
	services, err := Services()
	if err != nil {
		return nil, err
	}
	var (
		searchName    = strings.ToLower(name)
		groupsMap     = make(map[string]*types.NodesGroups)
		matchGroups   = make([]*types.SearchGroupServiceSimpleInfo, 0)
		matchServices = make([]*types.SearchGroupServiceSimpleInfo, 0)
	)
	for _, group := range groups {
		groupsMap[group.GroupId] = group
		if strings.Contains(strings.ToLower(group.GroupName), searchName) {
			matchGroups = append(matchGroups, &types.SearchGroupServiceSimpleInfo{
				GroupId:   group.GroupId,
				GroupName: group.GroupName,
			})
		}
	}
	for _, service := range services {
		groupInfo := groupsMap[service.GroupId]
		if groupInfo != nil && strings.Contains(strings.ToLower(service.ServiceName), searchName) {
			matchServices = append(matchServices, &types.SearchGroupServiceSimpleInfo{
				GroupId:     groupInfo.GroupId,
				GroupName:   groupInfo.GroupName,
				ServiceName: service.ServiceName,
				ServiceId:   service.ServiceId,
			})
		}
	}
	return map[string][]*types.SearchGroupServiceSimpleInfo{
		"groups":   matchGroups,
		"services": matchServices,
	}, nil
}
