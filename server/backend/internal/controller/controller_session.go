package controller

import (
	"log/slog"
	"time"

	"humpback/common/response"
	"humpback/config"
	"humpback/internal/db"
	"humpback/pkg/utils"
	"humpback/types"
)

func SessionGCInterval(stopCh <-chan struct{}) {
	slog.Info("[Session GC] Startup GC channel.", "Interval", config.DBArgs().SessionGCInterval.String())
	ticker := time.NewTicker(config.DBArgs().SessionGCInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			slog.Info("[Session GC] Time is up...")
			sessions, err := db.SessionsGetAll()
			if err != nil {
				slog.Error("[Session GC] get all session failed.", "Error", err)
			}
			var (
				nowT  = utils.NewActionTimestamp()
				gcIds = make([]string, 0)
			)
			for _, session := range sessions {
				if session.ExpiredAt < nowT {
					gcIds = append(gcIds, session.SessionId)
				}
			}
			slog.Info("[Session GC] Checked.", "Total", len(gcIds))
			if len(gcIds) > 0 {
				if err = db.SessionGCByIds(gcIds); err != nil {
					slog.Error("[Session GC] GC session failed.", "Total", len(gcIds), "Error", err)
				} else {
					slog.Info("[Session GC] GC session completed.", "Total", len(gcIds))
				}
			}

		case <-stopCh:
			return
		}
	}
}

func SessionGetAndRefresh(sessionId string) (*types.User, bool, error) {
	sessionInfo, expired, err := db.SessionGetById(sessionId)
	if err != nil {
		return nil, expired, response.NewRespServerErr(err.Error())
	}
	if expired {
		return nil, true, nil
	}
	userInfo, err := User(sessionInfo.UserId)
	if err != nil {
		return nil, false, err
	}
	if err = SessionUpdate(sessionInfo); err != nil {
		return nil, false, err
	}
	return userInfo, false, nil
}

func SessionUpdate(sessionInfo *types.Session) error {
	sessionInfo.ExpiredAt = time.Now().Add(config.DBArgs().SessionTimeout).UnixMilli()
	if err := db.SessionUpdate(sessionInfo); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	return nil
}

func SessionDelete(sessionId string) error {
	if err := db.SessionDelete(sessionId); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	return nil
}

func SessionDeleteByUserId(userId string) error {
	if err := db.SessionBatchDeleteByUserId(userId); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	return nil
}
