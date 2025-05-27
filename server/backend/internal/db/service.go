package db

import (
    "encoding/json"
    
    bolt "go.etcd.io/bbolt"
    "humpback/types"
)

func ServicesGetAll() ([]*types.Service, error) {
    return GetDataAll[types.Service](BucketServices)
}

func ServicesGetValidByPrefix(prefix string) ([]*types.Service, error) {
    services, err := GetDataByPrefix[types.Service](BucketServices, prefix)
    if err != nil {
        return nil, err
    }
    result := make([]*types.Service, 0)
    for _, service := range services {
        if !service.IsDelete {
            result = append(result, service)
        }
    }
    return result, nil
}

func ServicesGetTotalAndExceptionByGroups(groups map[string]string) (int, int, int, []*types.ResourceExceptionServiceInfo, error) {
    var (
        total             = 0
        ownTotal          = 0
        ownEnableTotal    = 0
        exceptionServices = make([]*types.ResourceExceptionServiceInfo, 0)
    )
    if err := TransactionGet(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(BucketServices))
        if bucket == nil {
            return ErrBucketNotExist
        }
        c := bucket.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            tempService := new(types.Service)
            if err := json.Unmarshal(v, tempService); err != nil {
                return err
            }
            if !tempService.IsDelete {
                total++
            }
            if groupName, ok := groups[tempService.GroupId]; ok {
                if !tempService.IsDelete {
                    ownTotal++
                    if tempService.IsEnabled {
                        ownEnableTotal++
                        if tempService.Status == types.ServiceStatusFailed {
                            exceptionServices = append(exceptionServices, &types.ResourceExceptionServiceInfo{
                                GroupName: groupName,
                                Service:   tempService,
                            })
                        }
                    }
                }
            }
        }
        return nil
    }); err != nil {
        return 0, 0, 0, nil, err
    }
    return total, ownTotal, ownEnableTotal, exceptionServices, nil
}

func ServiceGetTotalByPrefix(prefix string) (int, error) {
    return GetDataTotalByPrefix[types.Service](BucketServices, prefix)
}

func ServiceGetById(serviceId string) (*types.Service, error) {
    return GetDataById[types.Service](BucketServices, serviceId)
}

func ServiceUpdate(data *types.Service) error {
    return SaveData(BucketServices, data.ServiceId, data)
}

func ServiceDelete(id string) error {
    return DeleteData(BucketServices, id)
}
