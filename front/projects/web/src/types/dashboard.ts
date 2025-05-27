import { NodeInfo } from "#/setting.ts"
import { ServiceInfo } from "#/service.ts"

export interface DashboardResourceStatisticsInfo {
  services: number
  nodes: number
  users: number
  groups: number
  ownGroups: number
  ownServices: number
  enableOwnServices: number
  enableOwnNodes: number
  exceptionServices: ResourceExceptionServiceInfo[]
  abnormalNodes: NodeInfo[]
}

export interface ResourceExceptionServiceInfo extends ServiceInfo {
  groupName: string
}

export function NewDashboardResourceStatisticsInfo(): DashboardResourceStatisticsInfo {
  return {
    services: 0,
    nodes: 0,
    users: 0,
    groups: 0,
    ownGroups: 0,
    ownServices: 0,
    enableOwnServices: 0,
    enableOwnNodes: 0,
    exceptionServices: [],
    abnormalNodes: []
  }
}
