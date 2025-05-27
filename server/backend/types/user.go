package types

import (
    "slices"
)

type UserRole int

var (
    UserRoleSuperAdmin UserRole = 1
    UserRoleAdmin      UserRole = 2
    UserRoleUser       UserRole = 3
)

var (
    UserRolesMap = map[UserRole]string{
        UserRoleSuperAdmin: "Super Admin",
        UserRoleAdmin:      "Admin",
        UserRoleUser:       "user",
    }
)

type User struct {
    UserId      string   `json:"userId"`
    Username    string   `json:"username"`
    Description string   `json:"description"`
    Email       string   `json:"email"`
    Password    string   `json:"password"`
    Phone       string   `json:"phone"`
    Role        UserRole `json:"role"`
    CreatedAt   int64    `json:"createdAt"`
    UpdatedAt   int64    `json:"updatedAt"`
    Teams       []string `json:"teams"`
}

func (u *User) InGroup(group *NodesGroups) bool {
    if group == nil {
        return false
    }
    if IsAdmin(u.Role) || IsSuperAdmin(u.Role) {
        return true
    }
    if slices.Index(group.Users, u.UserId) != -1 {
        return true
    }
    groupTeams := make(map[string]bool)
    for _, team := range group.Teams {
        groupTeams[team] = true
    }
    for _, team := range u.Teams {
        if groupTeams[team] {
            return true
        }
    }
    return false
}

type Team struct {
    TeamId      string   `json:"teamId"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    CreatedAt   int64    `json:"createdAt"`
    UpdatedAt   int64    `json:"updatedAt"`
    Users       []string `json:"users"`
}

func IsSuperAdmin(role UserRole) bool {
    return role == UserRoleSuperAdmin
}

func IsAdmin(role UserRole) bool {
    return role == UserRoleAdmin
}

func IsUser(role UserRole) bool {
    return role == UserRoleUser
}
