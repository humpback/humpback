package db

/*

Buckets

Users: user_id, user_name, password, email, phone, create_time, update_time, is_admin, groups
Groups: group_id, group_name, create_time, update_time, users, teams, nodes

Registries: registry_id, registry_name, url, user_name, password, isDefault, create_time, update_time, status

Configs: config_id, config_name, create_time, update_time, values(json)

Nodes: node_id, node_name, host_ip, host_port, create_time, update_time, status, infos(json)

templates: template_id, template_name, create_time, update_time, infos(json)

Services: service_id, service_name, create_time, update_time, type, status, containers, infos(json)

serviceId: {groupId}{random-8}
containerName: humpback-{serviceId}-{version-5}-{random-5}

*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"humpback/config"

	bolt "go.etcd.io/bbolt"
)

type dbHelper struct {
	boltDB *bolt.DB
}

const (
	BucketUsers       = "Users"
	BucketTeams       = "Teams"
	BucketSessions    = "Sessions"
	BucketRegistries  = "Registries"
	BucketConfigs     = "Configs"
	BucketNodes       = "Nodes"
	BucketNodesGroups = "NodesGroups"
	BucketServices    = "Services"
	BucketActivities  = "Activities"
	BucketStatistics  = "Statistics"
)

const (
	ActivityBucketAccount     = "ActivityAccounts"
	ActivityBucketUsers       = "ActivityUsers"
	ActivityBucketTeams       = "ActivityTeams"
	ActivityBucketNodes       = "ActivityNodes"
	ActivityBucketRegistries  = "ActivityRegistries"
	ActivityBucketNodesGroups = "ActivityNodesGroups"
	ActivityBucketServices    = "ActivityNodesServices"
	ActivityBucketConfigs     = "ActivityNodesConfigs"
)

var (
	Buckets = []string{
		BucketUsers,
		BucketTeams,
		BucketSessions,
		BucketRegistries,
		BucketConfigs,
		BucketNodes,
		BucketNodesGroups,
		BucketServices,
		BucketActivities,
		BucketStatistics,
	}
	ActivityBuckets = []string{
		ActivityBucketAccount,
		ActivityBucketUsers,
		ActivityBucketTeams,
		ActivityBucketNodes,
		ActivityBucketRegistries,
		ActivityBucketNodesGroups,
		ActivityBucketServices,
		ActivityBucketConfigs,
	}
)

var (
	ErrKeyNotExist    = errors.New("key not found")
	ErrBucketNotExist = errors.New("bucket not exists")
)

var db *dbHelper

func InitDB() error {
	db = &dbHelper{}
	boltDB, err := bolt.Open(config.DBArgs().Root, 0600, nil)
	if err != nil {
		return fmt.Errorf("open db failed: %s", err)
	}
	db.boltDB = boltDB
	return nil
}

func CloseDB() {
	db.boltDB.Close()
}

func EnsureAndInitBuckets() error {
	return db.boltDB.Batch(func(tx *bolt.Tx) error {
		for _, bucket := range Buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return fmt.Errorf("check bucekt %s failed: %s", bucket, err)
			}
		}
		activityBucket := tx.Bucket([]byte(BucketActivities))
		if activityBucket == nil {
			return ErrBucketNotExist
		}
		for _, bucket := range ActivityBuckets {
			if _, err := activityBucket.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return fmt.Errorf("check activity child bucekt %s failed: %s", bucket, err)
			}
		}
		return nil
	})
}

func GetDataById[T any](bucketName string, id string) (*T, error) {
	var result T
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		data := bucket.Get([]byte(id))
		if data == nil {
			return ErrKeyNotExist
		}
		return json.Unmarshal(data, &result)
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetDataByIds[T any](bucketName string, ids []string, ignoreNotExist bool) ([]*T, error) {
	var results = make([]*T, 0)
	if len(ids) == 0 {
		return results, nil
	}
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		for _, id := range ids {
			data := bucket.Get([]byte(id))
			if data == nil {
				if !ignoreNotExist {
					return ErrKeyNotExist
				}
				continue
			}
			result := new(T)
			if err := json.Unmarshal(data, result); err != nil {
				return err
			}
			results = append(results, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetDataAll[T any](bucketName string) ([]*T, error) {
	var results []*T
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		bucket.ForEach(func(k, v []byte) error {
			result := new(T)
			if err := json.Unmarshal(v, &result); err != nil {
				return err
			}
			results = append(results, result)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type CustomFilterData func(key string, value interface{}) bool

func GetDataByQuery[T any](bucketName string, filter CustomFilterData) ([]*T, error) {
	var results = make([]*T, 0)
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		bucket.ForEach(func(k, v []byte) error {
			result := new(T)
			if err := json.Unmarshal(v, &result); err != nil {
				return err
			}
			if filter(string(k), result) {
				results = append(results, result)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetDataByPrefix[T any](bucketName string, prefix string) ([]*T, error) {
	var results = make([]*T, 0)
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		c := bucket.Cursor()
		bp := []byte(prefix)
		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, bp); k, v = c.Next() {
			result := new(T)
			if err := json.Unmarshal(v, result); err != nil {
				return err
			}
			results = append(results, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetDataTotalByPrefix[T any](bucketName string, prefix string) (int, error) {
	result := 0
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		c := bucket.Cursor()
		bp := []byte(prefix)
		for k, _ := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, bp); k, _ = c.Next() {
			result++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}

func SaveData[T any](bucketName string, id string, data T) error {
	return db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		encodedData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to encode data: %s", err)
		}
		return bucket.Put([]byte(id), encodedData)
	})
}

type fn func(tx *bolt.Tx) error

func TransactionUpdates(f fn) error {
	return db.boltDB.Update(f)
}

func TransactionGet(f fn) error {
	return db.boltDB.View(f)
}

func DeleteData(bucketName string, id string) error {
	return db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		return bucket.Delete([]byte(id))
	})
}

func DeleteDataByIds(bucketName string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrBucketNotExist
		}
		for _, id := range ids {
			if err := bucket.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}
