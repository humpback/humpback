package db

import (
	"encoding/json"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
	"humpback/api/handle/models"
	"humpback/types"
)

func StatisticalCountInsert(info *types.StatisticalCountInfo) error {
	return SaveData(BucketStatistics, info.Id, info)
}

func StatisticsCountQuery(info *models.StatisticsCountQueryReqInfo) ([]*types.StatisticalCountInfo, error) {
	var results = make([]*types.StatisticalCountInfo, 0)
	if err := TransactionGet(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketStatistics))
		if bucket == nil {
			return ErrBucketNotExist
		}
		c := bucket.Cursor()
		k, v := c.Last()
		for {
			if k == nil {
				break
			}
			valid, err := info.IsValid(string(k))
			if err != nil {
				return err
			}
			if valid {
				tempValue := new(types.StatisticalCountInfo)
				if err := json.Unmarshal(v, tempValue); err != nil {
					return err
				}
				results = append(results, tempValue)
			}
			k, v = c.Prev()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func StatisticsQueryKeysBeforeTimestamp(timestamp int64) ([]string, error) {
	var results = make([]string, 0)
	if err := TransactionGet(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketStatistics))
		if bucket == nil {
			return ErrBucketNotExist
		}
		c := bucket.Cursor()
		k, _ := c.First()
		for {
			if k == nil {
				break
			}
			keys := strings.Split(string(k), "-")
			t, err := strconv.ParseInt(keys[0], 10, 64)
			if err != nil {
				return err
			}
			if t >= timestamp {
				break
			}
			results = append(results, string(k))
			k, _ = c.Next()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func StatisticsDeleteByKeys(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return TransactionUpdates(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketStatistics))
		if bucket == nil {
			return ErrBucketNotExist
		}
		for _, id := range keys {
			if err := bucket.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}
