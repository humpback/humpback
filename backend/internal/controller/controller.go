package controller

import (
	"fmt"
	"log/slog"
	"slices"

	"humpback/config"
	"humpback/internal/db"
	"humpback/pkg/utils"
	"humpback/types"
)

func Start(stopCh <-chan struct{}) {
	go SessionGCInterval(stopCh)
	go ActivityGCInterval(stopCh)
	go StatisticsCountGCInterval(stopCh)
	go ReceiveActivities(stopCh)
	go ReceiveStatisticsCount(stopCh)
}

func InitData() error {
	slog.Info("[Init Buckets] Ensure and init all buckets...")
	if err := db.EnsureAndInitBuckets(); err != nil {
		return fmt.Errorf("init Buckets failed: %s", err)
	}
	slog.Info("[Init Buckets] Ensure and init all buckets completed.")

	slog.Info("[Init Super Admin] Ensure and init super admin account...")
	if err := ensureAndInitSuperAdminUser(); err != nil {
		return fmt.Errorf("init super admin acount failed: %s", err)
	}
	slog.Info("[Init Super Admin] Ensure and init super admin account completed.")

	slog.Info("[Init Default Registry] Ensure and init default registry...")
	if err := ensureAndInitRegistry(); err != nil {
		return err
	}
	slog.Info("[Init Default Registry] Ensure and init default registry completed.")
	return nil
}

func ensureAndInitSuperAdminUser() error {
	adminConfig := config.AdminArgs()
	user, err := db.UserGetSuperAdmin()
	if err != nil {
		return err
	}
	if user == nil {
		var (
			t  = utils.NewActionTimestamp()
			id = utils.NewGuidStr()
		)
		if err = db.UserUpdate(id, &types.User{
			UserId:    id,
			Username:  adminConfig.Username,
			Email:     "",
			Password:  adminConfig.Password,
			Phone:     "",
			Role:      types.UserRoleSuperAdmin,
			CreatedAt: t,
			UpdatedAt: t,
			Teams:     nil,
		}); err != nil {
			return fmt.Errorf("create super admin account failed: %s", err)
		}
	}
	return nil
}

func ensureAndInitRegistry() error {
	registries, err := db.RegistryGetAll()
	if err != nil {
		return fmt.Errorf("get registry failed: %s", err)
	}
	if slices.IndexFunc(registries, func(item *types.Registry) bool {
		return item.URL == "docker.io"
	}) == -1 {
		nowT := utils.NewActionTimestamp()
		if err = db.RegistryUpdate([]*types.Registry{
			{
				RegistryId: utils.NewGuidStr(),
				URL:        "docker.io",
				IsDefault:  true,
				Username:   "",
				Password:   "",
				CreatedAt:  nowT,
				UpdatedAt:  nowT,
			},
		}); err != nil {
			return fmt.Errorf("create default registry failed: %s", err)
		}
	}
	return nil
}

func sendServiceEvent(serviceChangeChan chan types.ServiceChangeInfo, serviceId, version, action string) {
	if serviceChangeChan != nil {
		serviceChangeChan <- types.ServiceChangeInfo{
			Action:    action,
			ServiceId: serviceId,
			Version:   version,
		}
	}
}

func sendNodeEvent(nodeChangeChan chan types.NodeSimpleInfo, nodeId, status string) {
	if nodeChangeChan != nil {
		nodeChangeChan <- types.NodeSimpleInfo{
			NodeId: nodeId,
			Status: status,
		}
	}
}
