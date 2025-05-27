import { ConfigType } from "@/models"
import { BaseInfo, NewBaseEmptyInfo } from "#/base.ts"

export interface ConfigInfo extends BaseInfo {
  configId: string
  configName: string
  description: string
  configValue: string
  configType: number
}

export function NewConfigEmptyInfo(): ConfigInfo {
  return {
    ...NewBaseEmptyInfo(),
    configId: "",
    configName: "",
    description: "",
    configValue: "",
    configType: ConfigType.Static
  }
}

export interface RegistryInfo extends BaseInfo {
  registryId: string
  url: string
  isDefault: boolean
  username: string
  password: string
  hasAuth?: boolean
}

export function NewRegistryEmptyInfo(): RegistryInfo {
  return {
    ...NewBaseEmptyInfo(),
    registryId: "",
    url: "",
    isDefault: false,
    username: "",
    password: ""
  }
}

export interface NodeInfo extends BaseInfo {
  nodeId: string
  name: string
  ipAddress: string
  port: number
  status: string
  isEnable: boolean
  cpuUsage: number
  cpu: number
  memoryUsage: number
  memoryTotal: number
  memoryUsed: number
  labels: { [key: string]: string }
  registerInfo: RegisterInfo
}

export interface RegisterInfo {
  token: string
  expireAt: number
  isRegister: boolean
  registerAt: number
}

export function NewNodeEmptyInfo(): NodeInfo {
  return {
    ...NewBaseEmptyInfo(),
    nodeId: "",
    name: "",
    ipAddress: "",
    port: 0,
    status: "",
    isEnable: false,
    cpuUsage: 0,
    cpu: 0,
    memoryUsage: 0,
    memoryTotal: 0,
    memoryUsed: 0,
    labels: {},
    registerInfo: {
      token: "",
      expireAt: 0,
      isRegister: false,
      registerAt: 0
    }
  }
}
