package db

import (
    "encoding/json"
    "fmt"
    "strings"
    
    "humpback/types"
    
    bolt "go.etcd.io/bbolt"
)

func UserGetSuperAdmin() (*types.User, error) {
    users, err := UsersGetAll()
    if err != nil {
        return nil, err
    }
    for _, user := range users {
        if types.IsSuperAdmin(user.Role) {
            return user, nil
        }
    }
    return nil, nil
}

func UsersGetTotal() (int, error) {
    total := 0
    if err := TransactionGet(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(BucketUsers))
        if bucket == nil {
            return ErrBucketNotExist
        }
        c := bucket.Cursor()
        for k, _ := c.First(); k != nil; k, _ = c.Next() {
            total++
        }
        return nil
    }); err != nil {
        return 0, err
    }
    return total, nil
}

func UsersGetAll() ([]*types.User, error) {
    return GetDataAll[types.User](BucketUsers)
}

func UserGetById(id string) (*types.User, error) {
    return GetDataById[types.User](BucketUsers, id)
}

func UsersGetByName(name string, isLower bool) ([]*types.User, error) {
    users, err := UsersGetAll()
    if err != nil {
        return nil, err
    }
    var result []*types.User
    for _, user := range users {
        if isLower && strings.EqualFold(user.Username, name) {
            result = append(result, user)
        }
        if !isLower && user.Username == name {
            result = append(result, user)
        }
    }
    return result, nil
}

func UsersGetByIds(ids []string, ignoreNotExist bool) ([]*types.User, error) {
    return GetDataByIds[types.User](BucketUsers, ids, ignoreNotExist)
}

func UserUpdate(id string, data *types.User) error {
    return SaveData(BucketUsers, id, data)
}

func UserAndTeamsUpdate(userInfo *types.User, teams []*types.Team) (string, error) {
    if err := TransactionUpdates(func(tx *bolt.Tx) error {
        var (
            teamBucket *bolt.Bucket
            userBucket *bolt.Bucket
        )
        userBucket = tx.Bucket([]byte(BucketUsers))
        if userBucket == nil {
            return ErrBucketNotExist
        }
        userData, err := json.Marshal(userInfo)
        if err != nil {
            return fmt.Errorf("failed to encode user data: %s", err)
        }
        if err = userBucket.Put([]byte(userInfo.UserId), userData); err != nil {
            return err
        }
        if len(teams) > 0 {
            teamBucket = tx.Bucket([]byte(BucketTeams))
            if teamBucket == nil {
                return ErrBucketNotExist
            }
            for _, team := range teams {
                teamData, err := json.Marshal(team)
                if err != nil {
                    return fmt.Errorf("failed to encode team data: %s", err)
                }
                if err = teamBucket.Put([]byte(team.TeamId), teamData); err != nil {
                    return err
                }
            }
        }
        return nil
    }); err != nil {
        return "", err
    }
    return userInfo.UserId, nil
}

func UserDeleteAndTeamsGroupsUpdate(id string, teams []*types.Team, groups []*types.NodesGroups) error {
    return TransactionUpdates(func(tx *bolt.Tx) error {
        userBucket := tx.Bucket([]byte(BucketUsers))
        if userBucket == nil {
            return ErrBucketNotExist
        }
        if err := userBucket.Delete([]byte(id)); err != nil {
            return err
        }
        if len(teams) > 0 {
            teamBucket := tx.Bucket([]byte(BucketTeams))
            if teamBucket == nil {
                return ErrBucketNotExist
            }
            for _, team := range teams {
                teamData, err := json.Marshal(team)
                if err != nil {
                    return fmt.Errorf("delete user(%s), failed to encode team(%s) data: %s", id, team.TeamId, err)
                }
                if err = teamBucket.Put([]byte(team.TeamId), teamData); err != nil {
                    return err
                }
            }
        }
        if len(groups) > 0 {
            groupBucket := tx.Bucket([]byte(BucketNodesGroups))
            if groupBucket == nil {
                return ErrBucketNotExist
            }
            for _, group := range groups {
                v, err := json.Marshal(group)
                if err != nil {
                    return fmt.Errorf("delete user(%s), failed to encode group(%s) data: %s", id, group.GroupId, err)
                }
                if err = groupBucket.Put([]byte(group.GroupId), v); err != nil {
                    return err
                }
            }
        }
        return nil
    })
}
