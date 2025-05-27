package db

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"humpback/types"
)

func RegistryGetAll() ([]*types.Registry, error) {
	return GetDataAll[types.Registry](BucketRegistries)
}

func RegistryGetById(id string) (*types.Registry, error) {
	return GetDataById[types.Registry](BucketRegistries, id)
}

func RegistryUpdate(List []*types.Registry) error {
	return TransactionUpdates(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketRegistries))
		if bucket == nil {
			return ErrBucketNotExist
		}
		for _, registry := range List {
			registryData, err := json.Marshal(registry)
			if err != nil {
				return fmt.Errorf("failed to encode user data: %s", err)
			}
			if err = bucket.Put([]byte(registry.RegistryId), registryData); err != nil {
				return err
			}
		}
		return nil
	})
}

func RegistryDelete(id string) error {
	return DeleteData(BucketRegistries, id)
}
