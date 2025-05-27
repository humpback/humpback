package controller

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/config"
	"humpback/internal/db"
	"humpback/types"
)

var (
	ActivityCh = make(chan *ActivityEvent, 100)

	ActivityTypes = map[string]string{
		"account":    db.ActivityBucketAccount,
		"users":      db.ActivityBucketUsers,
		"teams":      db.ActivityBucketTeams,
		"nodes":      db.ActivityBucketNodes,
		"registries": db.ActivityBucketRegistries,
		"configs":    db.ActivityBucketConfigs,
		"groups":     db.ActivityBucketNodesGroups,
		"services":   db.ActivityBucketServices,
	}
)

type ActivityEvent struct {
	Data   *types.ActivityInfo
	Bucket string
}

func ReceiveActivities(stopCh <-chan struct{}) {
	defer close(ActivityCh)
	slog.Info("[Activity] Startup receive channel.")
	for {
		select {
		case <-stopCh:
			return
		case info := <-ActivityCh:
			if err := db.ActivityUpdate(info.Data, info.Bucket); err != nil {
				slog.Error("[Activity] Insert activity failed.", "Bucket", info.Bucket, "Id", info.Data.ActivityId, "Error", err)
			}
		}
	}
}

func ActivityGCInterval(stopCh <-chan struct{}) {
	slog.Info("[Activity GC] Startup GC channel.", "Interval", config.DBArgs().ActivityGCInterval.String())
	ticker := time.NewTicker(config.DBArgs().ActivityGCInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			slog.Info("[Activity GC] Time is up, start GC...")
			activityGC()
			slog.Info("[Activity GC] GC completed.")
		}
	}
}

func activityGC() {
	var (
		gcList           = make(map[string][]string)
		wg               = &sync.WaitGroup{}
		l                = &sync.Mutex{}
		expiredTimestamp = time.Now().AddDate(0, 0, -1*config.DBArgs().ActivityRetentionDay).UnixMilli()
	)
	wg.Add(len(ActivityTypes))
	for _, bucket := range ActivityTypes {
		go func(bucket string) {
			defer wg.Done()
			keys, err := db.ActivityGetKeysBeforeTimestamp(expiredTimestamp, bucket)
			if err != nil {
				slog.Error("[Activity GC] get expired list failed.", "Bucket", bucket, "Total", len(keys), "Error", err)
			} else {
				if len(keys) > 0 {
					l.Lock()
					gcList[bucket] = keys
					l.Unlock()
				}
				slog.Info("[Activity GC] get expired list completed.", "Bucket", bucket, "Total", len(keys))
			}
		}(bucket)
	}
	wg.Wait()
	wg.Add(len(gcList))
	for bucket, keys := range gcList {
		go func(bucket string, keys []string) {
			defer wg.Done()
			if err := db.ActivityDeleteByKeys(bucket, keys); err != nil {
				slog.Error("[Activity GC] delete failed.", "Bucket", bucket, "Total", len(keys), "Error", err)
			} else {
				slog.Info("[Activity GC] delete completed.", "Bucket", bucket, "Total", len(keys))
			}
		}(bucket, keys)
	}
	wg.Wait()
}

// InsertAccountActivity 构造account的activity，同时写入db
func InsertAccountActivity(oldInfo, currentInfo *types.User, action types.ActivityAction, operateAt int64) {
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}
	description, oldContent, newContent := parseAccountActivityContent(oldInfo, currentInfo, action)
	data := &ActivityEvent{
		Bucket: db.ActivityBucketAccount,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, currentInfo.UserId, action, currentInfo.UserId),
			Action:       action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   currentInfo.UserId,
			Operator:     currentInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   currentInfo.UserId,
			ResourceName: currentInfo.Username,
		},
	}
	ActivityCh <- data
}

func parseAccountActivityContent(oldInfo, currentInfo *types.User, action types.ActivityAction) (string, any, any) {
	switch action {
	case types.ActivityActionLogin:
		return "Log in to the system.", nil, nil
	case types.ActivityActionLogout:
		return "Log out of the system.", nil, nil
	case types.ActivityActionUpdate:
		var (
			oldContent map[string]any
			newContent map[string]any
		)
		if oldInfo != nil {
			oldContent = map[string]any{
				"Username":    oldInfo.Username,
				"Description": oldInfo.Description,
				"Email":       oldInfo.Email,
				"Phone":       oldInfo.Phone,
			}
		}
		if currentInfo != nil {
			newContent = map[string]any{
				"Username":    currentInfo.Username,
				"Description": currentInfo.Description,
				"Email":       currentInfo.Email,
				"Phone":       currentInfo.Phone,
			}
		}
		return "Update account information.", oldContent, newContent
	case types.ActivityActionChangePassword:
		return "Change the password.", nil, nil
	}
	return "", nil, nil
}

type ActivityUserInfo struct {
	OldUserInfo  *types.User
	NewUserInfo  *types.User
	OldTeams     []*types.Team
	NewTeams     []*types.Team
	Action       types.ActivityAction
	OperatorInfo *types.User
	OperateAt    int64
}

// InsertUserActivity 构造user的activity，同时写入db
func InsertUserActivity(info *ActivityUserInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseUserActivityContent(info)
		resourceInfo                        = info.NewUserInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}
	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldUserInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketUsers,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.UserId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.UserId,
			ResourceName: resourceInfo.Username,
		},
	}
	ActivityCh <- data
}

func parseUserActivityContent(info *ActivityUserInfo) (string, any, any) {
	var (
		oldTeams    = make([]string, 0)
		newTeams    = make([]string, 0)
		oldUserRole = ""
		newUserRole = ""
		oldContent  map[string]any
		newContent  map[string]any
	)

	if info.OldUserInfo != nil {
		oldUserRole = types.UserRolesMap[info.OldUserInfo.Role]
		if len(info.OldTeams) > 0 {
			for _, teamId := range info.OldUserInfo.Teams {
				index := slices.IndexFunc(info.OldTeams, func(item *types.Team) bool {
					return item.TeamId == teamId
				})
				if index != -1 {
					oldTeams = append(oldTeams, info.OldTeams[index].Name)
				}
			}
		}
		slices.Sort(oldTeams)
		oldContent = map[string]any{
			"Username":    info.OldUserInfo.Username,
			"Role":        oldUserRole,
			"Description": info.OldUserInfo.Description,
			"Email":       info.OldUserInfo.Email,
			"Phone":       info.OldUserInfo.Phone,
			"Teams":       oldTeams,
		}
	}
	if info.NewUserInfo != nil {
		newUserRole = types.UserRolesMap[info.NewUserInfo.Role]
		if len(info.NewTeams) > 0 {
			for _, teamId := range info.NewUserInfo.Teams {
				index := slices.IndexFunc(info.NewTeams, func(item *types.Team) bool {
					return item.TeamId == teamId
				})
				if index != -1 {
					newTeams = append(newTeams, info.NewTeams[index].Name)
				}
			}
		}
		slices.Sort(newTeams)
		newContent = map[string]any{
			"Username":    info.NewUserInfo.Username,
			"Role":        newUserRole,
			"Description": info.NewUserInfo.Description,
			"Email":       info.NewUserInfo.Email,
			"Phone":       info.NewUserInfo.Phone,
			"Teams":       newTeams,
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add user.", nil, newContent
	case types.ActivityActionUpdate:
		return "Update user.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete user.", oldContent, nil
	}
	return "", nil, nil
}

type ActivityTeamInfo struct {
	OldTeamInfo  *types.Team
	NewTeamInfo  *types.Team
	OldUsers     []*types.User
	NewUsers     []*types.User
	Action       types.ActivityAction
	OperatorInfo *types.User
	OperateAt    int64
}

// InsertTeamActivity 构造team的activity，同时写入db
func InsertTeamActivity(info *ActivityTeamInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseTeamActivityContent(info)
		resourceInfo                        = info.NewTeamInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldTeamInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketTeams,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.TeamId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.TeamId,
			ResourceName: resourceInfo.Name,
		},
	}
	ActivityCh <- data
}

func parseTeamActivityContent(info *ActivityTeamInfo) (string, any, any) {
	var (
		oldUsers   = make([]string, 0)
		newUsers   = make([]string, 0)
		oldContent map[string]any
		newContent map[string]any
	)

	if info.OldTeamInfo != nil {
		if len(info.OldUsers) > 0 {
			for _, userId := range info.OldTeamInfo.Users {
				index := slices.IndexFunc(info.OldUsers, func(item *types.User) bool {
					return item.UserId == userId
				})
				if index != -1 {
					oldUsers = append(oldUsers, info.OldUsers[index].Username)
				}
			}
		}
		slices.Sort(oldUsers)
		oldContent = map[string]any{
			"Name":        info.OldTeamInfo.Name,
			"Description": info.OldTeamInfo.Description,
			"Users":       oldUsers,
		}
	}
	if info.NewTeamInfo != nil {
		if len(info.NewUsers) > 0 {
			for _, userId := range info.NewTeamInfo.Users {
				index := slices.IndexFunc(info.NewUsers, func(item *types.User) bool {
					return item.UserId == userId
				})
				if index != -1 {
					newUsers = append(newUsers, info.NewUsers[index].Username)
				}
			}
		}
		slices.Sort(newUsers)
		newContent = map[string]any{
			"Name":        info.NewTeamInfo.Name,
			"Description": info.NewTeamInfo.Description,
			"Users":       newUsers,
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add team.", nil, newContent
	case types.ActivityActionUpdate:
		return "Update team.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete team.", oldContent, nil
	}
	return "", nil, nil
}

type ActivityNodeInfo struct {
	OldNodeInfo  *types.Node
	NewNodeInfo  *types.Node
	Action       types.ActivityAction
	OperatorInfo *types.User
	OperateAt    int64
}

// InsertNodeActivity 构造node的activity，同时写入db
func InsertNodeActivity(info *ActivityNodeInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseNodeActivityContent(info)
		resourceInfo                        = info.NewNodeInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldNodeInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketNodes,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.NodeId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.NodeId,
			ResourceName: resourceInfo.IpAddress,
		},
	}
	ActivityCh <- data
}

func parseNodeActivityContent(info *ActivityNodeInfo) (string, any, any) {
	var (
		oldContent map[string]any
		newContent map[string]any
	)

	if info.OldNodeInfo != nil {
		oldContent = map[string]any{
			"IpAddress": info.OldNodeInfo.IpAddress,
			"Labels":    info.OldNodeInfo.Labels,
		}
	}
	if info.NewNodeInfo != nil {
		newContent = map[string]any{
			"IpAddress": info.NewNodeInfo.IpAddress,
			"Labels":    info.NewNodeInfo.Labels,
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add node.", nil, nil
	case types.ActivityActionUpdateLabel:
		return "Update node labels.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete node.", oldContent, nil
	case types.ActivityActionEnable:
		return "Enable node.", nil, nil
	case types.ActivityActionDisable:
		return "Disable node.", nil, nil
	}
	return "", nil, nil
}

type ActivityRegistryInfo struct {
	OldRegistryInfo *types.Registry
	NewRegistryInfo *types.Registry
	Action          types.ActivityAction
	OperatorInfo    *types.User
	OperateAt       int64
}

// InsertRegistryActivity 构造registry的activity，同时写入db
func InsertRegistryActivity(info *ActivityRegistryInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseRegistryActivityContent(info)
		resourceInfo                        = info.NewRegistryInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldRegistryInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketRegistries,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.RegistryId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.RegistryId,
			ResourceName: resourceInfo.URL,
		},
	}
	ActivityCh <- data
}

func parseRegistryActivityContent(info *ActivityRegistryInfo) (string, any, any) {
	var (
		oldContent map[string]any
		newContent map[string]any
	)

	if info.OldRegistryInfo != nil {
		oldContent = map[string]any{
			"Url":       info.OldRegistryInfo.URL,
			"Username":  info.OldRegistryInfo.Username,
			"IsDefault": info.OldRegistryInfo.IsDefault,
		}
	}
	if info.NewRegistryInfo != nil {
		newContent = map[string]any{
			"Url":       info.NewRegistryInfo.URL,
			"Username":  info.NewRegistryInfo.Username,
			"IsDefault": info.NewRegistryInfo.IsDefault,
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add registry.", nil, newContent
	case types.ActivityActionUpdate:
		return "Update registry.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete registry.", oldContent, nil
	}
	return "", nil, nil
}

type ActivityConfigInfo struct {
	OldConfigInfo *types.Config
	NewConfigInfo *types.Config
	Action        types.ActivityAction
	OperatorInfo  *types.User
	OperateAt     int64
}

// InsertConfigActivity 构造config的activity，同时写入db
func InsertConfigActivity(info *ActivityConfigInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseConfigActivityContent(info)
		resourceInfo                        = info.NewConfigInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldConfigInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketConfigs,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.ConfigId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.ConfigId,
			ResourceName: resourceInfo.ConfigName,
		},
	}
	ActivityCh <- data
}

func parseConfigActivityContent(info *ActivityConfigInfo) (string, any, any) {
	var (
		oldContent map[string]any
		newContent map[string]any
	)

	if info.OldConfigInfo != nil {
		oldContent = map[string]any{
			"ConfigName":  info.OldConfigInfo.ConfigName,
			"Description": info.OldConfigInfo.Description,
			"ConfigType":  types.ConfigTypesMap[info.OldConfigInfo.ConfigType],
			"ConfigValue": info.OldConfigInfo.ConfigValue,
		}
	}
	if info.NewConfigInfo != nil {
		newContent = map[string]any{
			"ConfigName":  info.NewConfigInfo.ConfigName,
			"Description": info.NewConfigInfo.Description,
			"ConfigType":  types.ConfigTypesMap[info.NewConfigInfo.ConfigType],
			"ConfigValue": info.NewConfigInfo.ConfigValue,
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add config.", nil, newContent
	case types.ActivityActionUpdate:
		return "Update config.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete config.", oldContent, nil
	}
	return "", nil, nil
}

type ActivityGroupInfo struct {
	OldGroupInfo *types.NodesGroups
	NewGroupInfo *types.NodesGroups
	OldUsers     []*types.User
	NewUsers     []*types.User
	OldTeams     []*types.Team
	NewTeams     []*types.Team
	OldNodes     []*types.Node
	NewNodes     []*types.Node
	Action       types.ActivityAction
	OperatorInfo *types.User
	OperateAt    int64
}

// InsertGroupActivity 构造group的activity，同时写入db
func InsertGroupActivity(info *ActivityGroupInfo) {
	var (
		operateAt                           = info.OperateAt
		description, oldContent, newContent = parseGroupActivityContent(info)
		resourceInfo                        = info.NewGroupInfo
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldGroupInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketNodesGroups,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.GroupId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.GroupId,
			ResourceName: resourceInfo.GroupName,
		},
	}
	ActivityCh <- data
}

func parseGroupActivityContent(info *ActivityGroupInfo) (string, any, any) {
	var (
		oldContent map[string]any
		newContent map[string]any
		oldUsers   = make([]string, 0)
		newUsers   = make([]string, 0)
		oldTeams   = make([]string, 0)
		newTeams   = make([]string, 0)
		oldNodes   = make([]string, 0)
		newNodes   = make([]string, 0)
	)

	if info.OldGroupInfo != nil {
		if len(info.OldUsers) > 0 {
			for _, userId := range info.OldGroupInfo.Users {
				for _, user := range info.OldUsers {
					if user.UserId == userId {
						oldUsers = append(oldUsers, user.Username)
						break
					}
				}
			}
		}
		if len(info.OldTeams) > 0 {
			for _, teamId := range info.OldGroupInfo.Teams {
				for _, team := range info.OldTeams {
					if team.TeamId == teamId {
						oldTeams = append(oldTeams, team.Name)
						break
					}
				}
			}
		}
		if len(info.OldNodes) > 0 {
			for _, nodeId := range info.OldGroupInfo.Nodes {
				for _, node := range info.OldNodes {
					if node.NodeId == nodeId {
						oldNodes = append(oldNodes, node.IpAddress)
					}
				}
			}
		}

		slices.Sort(oldUsers)
		slices.Sort(oldTeams)
		slices.Sort(oldNodes)
		oldContent = map[string]any{
			"GroupName":   info.OldGroupInfo.GroupName,
			"Description": info.OldGroupInfo.Description,
			"Users":       oldUsers,
			"Teams":       oldTeams,
		}
		if info.Action == types.ActivityActionDelete {
			oldContent["Nodes"] = oldNodes
		}
		if info.Action == types.ActivityActionRemoveNode || info.Action == types.ActivityActionAddNode {
			oldContent = map[string]any{
				"Nodes": oldNodes,
			}
		}
	}
	if info.NewGroupInfo != nil {
		if len(info.NewUsers) > 0 {
			for _, userId := range info.NewGroupInfo.Users {
				for _, user := range info.NewUsers {
					if user.UserId == userId {
						newUsers = append(newUsers, user.Username)
						break
					}
				}
			}
		}
		if len(info.NewTeams) > 0 {
			for _, teamId := range info.NewGroupInfo.Teams {
				for _, team := range info.NewTeams {
					if team.TeamId == teamId {
						newTeams = append(newTeams, team.Name)
						break
					}
				}
			}
		}
		if len(info.NewNodes) > 0 {
			for _, nodeId := range info.NewGroupInfo.Nodes {
				for _, node := range info.NewNodes {
					if node.NodeId == nodeId {
						newNodes = append(newNodes, node.IpAddress)
					}
				}
			}
		}
		slices.Sort(newUsers)
		slices.Sort(newTeams)
		slices.Sort(newNodes)
		newContent = map[string]any{
			"GroupName":   info.NewGroupInfo.GroupName,
			"Description": info.NewGroupInfo.Description,
			"Users":       newUsers,
			"Teams":       newTeams,
		}
		if info.Action == types.ActivityActionRemoveNode || info.Action == types.ActivityActionAddNode {
			newContent = map[string]any{
				"Nodes": newNodes,
			}
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add group.", nil, newContent
	case types.ActivityActionUpdate:
		return "Update group.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete group.", oldContent, nil
	case types.ActivityActionRemoveNode:
		return "Remove group node.", oldContent, newContent
	case types.ActivityActionAddNode:
		return "Add group node.", oldContent, newContent
	}
	return "", nil, nil
}

type ActivityServiceInfo struct {
	OldServiceInfo *types.Service
	NewServiceInfo *types.Service
	Action         types.ActivityAction
	InstanceName   string
	OperatorInfo   *types.User
	OperateAt      int64
	GroupName      string
}

// InsertServiceActivity 构造service的activity，同时写入db
func InsertServiceActivity(info *ActivityServiceInfo) {
	var (
		operateAt                           = info.OperateAt
		resourceInfo                        = info.NewServiceInfo
		description, oldContent, newContent = parseServiceActivityContent(info)
	)
	if operateAt == 0 {
		operateAt = time.Now().UnixMilli()
	}

	if info.Action == types.ActivityActionDelete {
		resourceInfo = info.OldServiceInfo
	}

	data := &ActivityEvent{
		Bucket: db.ActivityBucketServices,
		Data: &types.ActivityInfo{
			ActivityId:   fmt.Sprintf("%d-%s-%s-%s", operateAt, info.OperatorInfo.UserId, info.Action, resourceInfo.ServiceId),
			Action:       info.Action,
			Description:  description,
			OldContent:   oldContent,
			NewContent:   newContent,
			OperatorId:   info.OperatorInfo.UserId,
			Operator:     info.OperatorInfo.Username,
			OperateAt:    operateAt,
			ResourceId:   resourceInfo.ServiceId,
			ResourceName: resourceInfo.ServiceName,
			InstanceName: info.InstanceName,
			GroupName:    info.GroupName,
		},
	}
	ActivityCh <- data
}

func parseServiceActivityContent(info *ActivityServiceInfo) (string, any, any) {
	var (
		oldContent map[string]any
		newContent map[string]any
	)

	if info.OldServiceInfo != nil {
		switch info.Action {
		case types.ActivityActionUpdateBasic:
			oldContent = map[string]any{
				"ServiceName": info.OldServiceInfo.ServiceName,
				"Description": info.OldServiceInfo.Description,
				"IsEnabled":   info.OldServiceInfo.IsEnabled,
			}
		case types.ActivityActionUpdateApplication:
			oldContent = parseServiceMetaActivityContent(info.OldServiceInfo)
		case types.ActivityActionUpdateDeployment:
			oldContent = parseServiceDeploymentActivityContent(info.OldServiceInfo)
		case types.ActivityActionDelete:
			oldContent = map[string]any{
				"ServiceName": info.OldServiceInfo.ServiceName,
				"Description": info.OldServiceInfo.Description,
				"IsEnabled":   info.OldServiceInfo.IsEnabled,
				"Application": parseServiceMetaActivityContent(info.OldServiceInfo),
				"Deployment":  parseServiceDeploymentActivityContent(info.OldServiceInfo),
			}
		}
	}
	if info.NewServiceInfo != nil {
		switch info.Action {
		case types.ActivityActionUpdateBasic:
			newContent = map[string]any{
				"ServiceName": info.NewServiceInfo.ServiceName,
				"Description": info.NewServiceInfo.Description,
				"IsEnabled":   info.NewServiceInfo.IsEnabled,
			}
		case types.ActivityActionUpdateApplication:
			newContent = parseServiceMetaActivityContent(info.NewServiceInfo)
		case types.ActivityActionUpdateDeployment:
			newContent = parseServiceDeploymentActivityContent(info.NewServiceInfo)
		case types.ActivityActionDelete:
			newContent = map[string]any{
				"ServiceName": info.NewServiceInfo.ServiceName,
				"Description": info.NewServiceInfo.Description,
				"IsEnabled":   info.NewServiceInfo.IsEnabled,
				"Application": parseServiceMetaActivityContent(info.NewServiceInfo),
				"Deployment":  parseServiceDeploymentActivityContent(info.NewServiceInfo),
			}
		}
	}

	switch info.Action {
	case types.ActivityActionAdd:
		return "Add service.", nil, newContent
	case types.ActivityActionUpdateBasic:
		return "Update service basic information.", oldContent, newContent
	case types.ActivityActionUpdateApplication:
		return "Update service application.", oldContent, newContent
	case types.ActivityActionUpdateDeployment:
		return "Update service deployment.", oldContent, newContent
	case types.ActivityActionDelete:
		return "Delete service.", oldContent, nil
	case types.ActivityActionEnable:
		return "Enable service.", nil, nil
	case types.ActivityActionDisable:
		return "Disable service.", nil, nil
	case types.ActivityActionStart:
		if info.InstanceName != "" {
			return fmt.Sprintf("Start servcie instance %s.", info.InstanceName), nil, nil
		}
		return "Start service.", nil, nil
	case types.ActivityActionStop:
		if info.InstanceName != "" {
			return fmt.Sprintf("Stop servcie instance %s.", info.InstanceName), nil, nil
		}
		return "Stop service.", nil, nil
	case types.ActivityActionRestart:
		if info.InstanceName != "" {
			return fmt.Sprintf("Restart servcie instance %s.", info.InstanceName), nil, nil
		}
		return "Restart service.", nil, nil
	}
	return "", nil, nil
}

func parseServiceMetaActivityContent(serviceInfo *types.Service) map[string]any {
	var content map[string]any
	if serviceInfo != nil && serviceInfo.Meta != nil {
		meta := serviceInfo.Meta
		content = map[string]any{
			"Image":        strings.Join([]string{meta.RegistryDomain, meta.Image}, "/"),
			"AlwaysPull":   meta.AlwaysPull,
			"Command":      meta.Command,
			"Environments": meta.EnvConfig,
			"Labels":       meta.Labels,
			"Privileged":   meta.Privileged,
		}
		if meta.Capabilities != nil {
			content["Capabilities"] = meta.Capabilities
		}
		if meta.LogConfig != nil {
			content["LogConfig"] = meta.LogConfig
		}
		if meta.Resources != nil {
			content["Resources"] = map[string]any{
				"MemoryLimit":       meta.Resources.Memory,
				"MemoryReservation": meta.Resources.MemoryReservation,
				"MaximumCpuUsage":   meta.Resources.MaxCpuUsage,
			}
		}
		if meta.Network != nil {
			ports := make([]map[string]any, 0)
			for _, port := range meta.Network.Ports {
				ports = append(ports, map[string]any{
					"Protocol":      port.Protocol,
					"HostPort":      port.HostPort,
					"ContainerPort": port.ContainerPort,
				})
			}
			content["Network"] = map[string]any{
				"Mode":               meta.Network.Mode,
				"Hostname":           meta.Network.Hostname,
				"NetworkName":        meta.Network.NetworkName,
				"UseMachineHostname": meta.Network.UseMachineHostname,
				"Ports":              ports,
			}
		}
		if meta.RestartPolicy != nil {
			content["RestartPolicy"] = map[string]any{
				"Mode":          meta.RestartPolicy.Mode,
				"MaxRetryCount": meta.RestartPolicy.MaxRetryCount,
			}
		}
		if len(meta.Volumes) > 0 {
			volumes := make([]map[string]any, 0)
			for _, volume := range meta.Volumes {
				volumes = append(volumes, map[string]any{
					"Type":          volume.Type,
					"ContainerPath": volume.Target,
					"HostPath":      volume.Source,
					"ReadOnly":      volume.Readonly,
				})
			}
			content["Volumes"] = volumes
		}
	}
	return content
}

func parseServiceDeploymentActivityContent(serviceInfo *types.Service) map[string]any {
	var content map[string]any
	if serviceInfo != nil && serviceInfo.Deployment != nil {
		deployment := serviceInfo.Deployment
		content = map[string]any{
			"DeployMode":  deployment.Mode,
			"InstanceNum": deployment.Replicas,
		}
		if len(deployment.Placements) > 0 {
			placements := make([]map[string]any, 0)
			for _, placement := range deployment.Placements {
				placements = append(placements, map[string]any{
					"Mode":    placement.Mode,
					"Key":     placement.Key,
					"Value":   placement.Value,
					"IsEqual": placement.IsEqual,
				})
			}
			content["PlacementConstraints"] = placements
		}
		if deployment.Schedule != nil && len(deployment.Schedule.Rules) > 0 {
			content["Schedules"] = map[string]any{
				"Timeout":         deployment.Schedule.Timeout,
				"Cron":            deployment.Schedule.Rules,
				"ManualExecution": false,
			}
		} else if deployment.ManualExec {
			content["Schedules"] = map[string]any{
				"ManualExecution": true,
			}
		}
	}
	return content
}

func ActivityQuery(info *models.ActivityQueryReqInfo) (*response.QueryResult[types.ActivityInfo], error) {
	bucket, ok := ActivityTypes[info.FilterInfo.Type]
	if !ok {
		return nil, response.NewBadRequestErr(locales.CodeActivityTypeInvalid)
	}
	result, err := db.ActivityQuery(info, bucket)
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return result, nil
}

func ActivityAllQuery(info *models.ActivityAllQueryReqInfo) (map[string][]*types.ActivityInfo, error) {
	var (
		buckets = map[string]string{
			"user":     db.ActivityBucketUsers,
			"team":     db.ActivityBucketTeams,
			"node":     db.ActivityBucketNodes,
			"registry": db.ActivityBucketRegistries,
			"config":   db.ActivityBucketConfigs,
			"group":    db.ActivityBucketNodesGroups,
			"service":  db.ActivityBucketServices,
		}
		results = make(map[string][]*types.ActivityInfo)
		wg      = &sync.WaitGroup{}
		l       = &sync.Mutex{}
	)
	if types.IsUser(info.UserInfo.Role) {
		buckets = map[string]string{
			"config":  db.ActivityBucketConfigs,
			"group":   db.ActivityBucketNodesGroups,
			"service": db.ActivityBucketServices,
		}
	}
	wg.Add(len(buckets))
	for key, bucket := range buckets {
		go func(name string, bucket string) {
			defer wg.Done()
			list, err := db.ActivityQueryByStartAtAndUser(info, bucket)
			if err != nil {
				slog.Error("[ActivityQueryByStartAtAndUser] List activities failed.", "Bucket", bucket, "Error", err.Error())
				return
			}
			l.Lock()
			results[name] = list
			l.Unlock()
		}(key, bucket)
	}
	wg.Wait()
	return results, nil
}
