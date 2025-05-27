package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
	"humpback/api/handle/models"
	"humpback/common/response"
	"humpback/types"
)

func ActivityUpdate(info *types.ActivityInfo, bucket string) error {
	return TransactionUpdates(func(tx *bolt.Tx) error {
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		childBucket := activityBucket.Bucket([]byte(bucket))
		if childBucket == nil {
			return ErrBucketNotExist
		}
		encodedData, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("failed to encode data: %s", err)
		}
		return childBucket.Put([]byte(info.ActivityId), encodedData)
	})
}

func ActivityQuery(info *models.ActivityQueryReqInfo, bucket string) (*response.QueryResult[types.ActivityInfo], error) {
	var (
		results    = &response.QueryResult[types.ActivityInfo]{Total: 0, List: make([]*types.ActivityInfo, 0)}
		startIndex = (info.PageInfo.Index-1)*info.PageInfo.Size + 1
		endIndex   = info.PageInfo.Index*info.PageInfo.Size + 1
	)
	if err := TransactionGet(func(tx *bolt.Tx) error {
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		childBucket := activityBucket.Bucket([]byte(bucket))
		if childBucket == nil {
			return ErrBucketNotExist
		}
		c := childBucket.Cursor()
		k, v := c.Last()
		for {
			if k == nil {
				break
			}
			keys := strings.Split(string(k), "-")
			isSmallStartAt, validTimeRange, err := info.IsValidTimeRange(keys[0])
			if err != nil {
				return err
			}
			if isSmallStartAt {
				break
			}
			if !validTimeRange {
				continue
			}
			if info.IsValid(keys) {
				results.Total++
				if results.Total >= startIndex && results.Total < endIndex {
					tempValue := new(types.ActivityInfo)
					if err := json.Unmarshal(v, tempValue); err != nil {
						return err
					}
					results.List = append(results.List, tempValue)
				}
			}
			k, v = c.Prev()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func ActivityQueryByStartAtAndUser(info *models.ActivityAllQueryReqInfo, bucket string) ([]*types.ActivityInfo, error) {
	var (
		results = make([]*types.ActivityInfo, 0)
	)
	if err := TransactionGet(func(tx *bolt.Tx) error {
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		childBucket := activityBucket.Bucket([]byte(bucket))
		if childBucket == nil {
			return ErrBucketNotExist
		}
		c := childBucket.Cursor()
		k, v := c.Last()
		for {
			if k == nil {
				break
			}
			keys := strings.Split(string(k), "-")
			valid, err := info.IsValid(keys[0], keys[1])
			if err != nil {
				return err
			}
			if valid {
				tempValue := new(types.ActivityInfo)
				if err := json.Unmarshal(v, tempValue); err != nil {
					return err
				}
				results = append(results, tempValue)
			} else {
				break
			}
			k, v = c.Prev()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func ActivityGetKeysBeforeTimestamp(timestamp int64, bucket string) ([]string, error) {
	var (
		results = make([]string, 0)
	)
	if err := TransactionGet(func(tx *bolt.Tx) error {
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		childBucket := activityBucket.Bucket([]byte(bucket))
		if childBucket == nil {
			return ErrBucketNotExist
		}
		c := childBucket.Cursor()
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

func ActivityDeleteByKeys(bucket string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return TransactionUpdates(func(tx *bolt.Tx) error {
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		childBucket := activityBucket.Bucket([]byte(bucket))
		if childBucket == nil {
			return ErrBucketNotExist
		}
		for _, id := range keys {
			if err := childBucket.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}
