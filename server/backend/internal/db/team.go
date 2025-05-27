package db

import (
    "encoding/json"
    "fmt"
    "strings"
    
    "humpback/types"
    
    bolt "go.etcd.io/bbolt"
)

func TeamsGetAll() ([]*types.Team, error) {
    return GetDataAll[types.Team](BucketTeams)
}

func TeamsGetByIds(ids []string, ignoreNotExist bool) ([]*types.Team, error) {
    return GetDataByIds[types.Team](BucketTeams, ids, ignoreNotExist)
}

func TeamGetById(id string) (*types.Team, error) {
    return GetDataById[types.Team](BucketTeams, id)
}

func TeamsGetByName(name string, isLower bool) ([]*types.Team, error) {
    teams, err := TeamsGetAll()
    if err != nil {
        return nil, err
    }
    var result []*types.Team
    for _, team := range teams {
        if isLower && strings.EqualFold(team.Name, name) {
            result = append(result, team)
        }
        if !isLower && team.Name == name {
            result = append(result, team)
        }
    }
    return result, nil
}

func TeamUpdateAndUsers(teamInfo *types.Team, users []*types.User) (string, error) {
    if err := TransactionUpdates(func(tx *bolt.Tx) error {
        var (
            teamBucket *bolt.Bucket
            userBucket *bolt.Bucket
        )
        teamBucket = tx.Bucket([]byte(BucketTeams))
        if teamBucket == nil {
            return ErrBucketNotExist
        }
        teamData, err := json.Marshal(teamInfo)
        if err != nil {
            return fmt.Errorf("failed to encode team data: %s", err)
        }
        if err = teamBucket.Put([]byte(teamInfo.TeamId), teamData); err != nil {
            return err
        }
        if len(users) > 0 {
            userBucket = tx.Bucket([]byte(BucketUsers))
            if userBucket == nil {
                return ErrBucketNotExist
            }
            for _, user := range users {
                userData, err := json.Marshal(user)
                if err != nil {
                    return fmt.Errorf("failed to encode user data: %s", err)
                }
                if err = userBucket.Put([]byte(user.UserId), userData); err != nil {
                    return err
                }
            }
        }
        return nil
    }); err != nil {
        return "", err
    }
    return teamInfo.TeamId, nil
}

func TeamDeleteAndUsersGroupsUpdate(id string, users []*types.User, groups []*types.NodesGroups) error {
    return TransactionUpdates(func(tx *bolt.Tx) error {
        teamBucket := tx.Bucket([]byte(BucketTeams))
        if teamBucket == nil {
            return ErrBucketNotExist
        }
        
        if err := teamBucket.Delete([]byte(id)); err != nil {
            return err
        }
        
        if len(users) > 0 {
            userBucket := tx.Bucket([]byte(BucketUsers))
            if userBucket == nil {
                return ErrBucketNotExist
            }
            for _, user := range users {
                userData, err := json.Marshal(user)
                if err != nil {
                    return fmt.Errorf("delete team(%s), failed to encode user(%s) data: %s", id, user.UserId, err)
                }
                if err = userBucket.Put([]byte(user.UserId), userData); err != nil {
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
                    return fmt.Errorf("delete team(%s), failed to encode group(%s) data: %s", id, group.GroupId, err)
                }
                if err = groupBucket.Put([]byte(group.GroupId), v); err != nil {
                    return err
                }
            }
        }
        return nil
    })
}
