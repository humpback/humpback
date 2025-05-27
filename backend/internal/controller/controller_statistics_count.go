package controller

import (
	"fmt"
	"log/slog"
	"time"

	"humpback/api/handle/models"
	"humpback/common/response"
	"humpback/config"
	"humpback/internal/db"
	"humpback/types"
)

var (
	StatisticsCh = make(chan *types.StatisticalCountInfo, 100)
)

type StatisticalCountEvent struct {
	CreateAt int64           `json:"createAt"`
	Type     types.CountType `json:"type"`
	Num      int             `json:"num"`
	UserId   string          `json:"userId"`
}

func ReceiveStatisticsCount(stopCh <-chan struct{}) {
	defer close(StatisticsCh)
	slog.Info("[Statistics Count] Startup receive channel.")
	for {
		select {
		case <-stopCh:
			return
		case info := <-StatisticsCh:
			if err := db.StatisticalCountInsert(info); err != nil {
				slog.Error("[Statistics Count] Insert statistics count failed.", "Type", info.Type, "Id", info.Id, "Error", err)
			}
		}
	}
}

func StatisticsCountGCInterval(stopCh <-chan struct{}) {
	slog.Info("[Statistics Count GC] Startup GC channel.", "Interval", config.DBArgs().StatisticsGCInterval.String())
	ticker := time.NewTicker(config.DBArgs().StatisticsGCInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			slog.Info("[Statistics Count GC] Time is up, start GC...")
			if err := statisticsCountGC(); err != nil {
				slog.Error("[Statistics Count GC] GC failed.", "Error", err)
			} else {
				slog.Info("[Statistics Count GC] GC completed.")
			}
		}
	}
}

func statisticsCountGC() error {
	expiredTimestamp := time.Now().AddDate(0, 0, -1*config.DBArgs().StatisticsRetentionDay).UnixMilli()
	gcList, err := db.StatisticsQueryKeysBeforeTimestamp(expiredTimestamp)
	if err != nil {
		return err
	}
	return db.StatisticsDeleteByKeys(gcList)
}

func InsertStatisticsCount(info *StatisticalCountEvent) {
	StatisticsCh <- &types.StatisticalCountInfo{
		Id:       fmt.Sprintf("%d-%s-%s-%d", info.CreateAt, info.UserId, info.Type, info.Num),
		CreateAt: info.CreateAt,
		Type:     info.Type,
		Num:      info.Num,
		UserId:   info.UserId,
	}
}

func StatisticsCountQuery(info *models.StatisticsCountQueryReqInfo) ([]*types.StatisticalCountInfo, error) {
	list, err := db.StatisticsCountQuery(info)
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return list, nil
}
