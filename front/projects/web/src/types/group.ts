import { BaseInfo, NewBaseEmptyInfo } from "#/base.ts"

export interface GroupInfo extends BaseInfo {
  groupId: string
  groupName: string
  description: string
  users: string[]
  teams: string[]
  nodes: string[]
}

export function NewGroupEmptyInfo(): GroupInfo {
  return {
    ...NewBaseEmptyInfo(),
    groupId: "",
    groupName: "",
    description: "",
    users: [] as Array<string>,
    teams: [] as Array<string>,
    nodes: [] as Array<string>
  }
}
