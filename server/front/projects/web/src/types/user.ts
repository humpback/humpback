import { BaseInfo, NewBaseEmptyInfo } from "#/base.ts"

export interface UserInfo extends BaseInfo {
  userId: string
  username: string
  email: string
  password: string
  description: string
  phone: string
  role: number
  teams: string[]
}

export function NewUserEmptyInfo(): UserInfo {
  return {
    ...NewBaseEmptyInfo(),
    userId: "",
    username: "",
    email: "",
    password: "",
    description: "",
    phone: "",
    role: UserRole.User,
    teams: []
  }
}

export interface TeamInfo extends BaseInfo {
  teamId: string
  name: string
  description: string
  users: string[]
}

export function NewTeamEmptyInfo(): TeamInfo {
  return {
    ...NewBaseEmptyInfo(),
    teamId: "",
    name: "",
    description: "",
    createdAt: 0,
    updatedAt: 0,
    users: []
  }
}
