import { BaseInfo, NewBaseEmptyInfo } from "#/base.ts"
import { ServiceVolumeType } from "@/models"

export interface ServiceInfo extends BaseInfo {
  serviceId: string
  groupId: string
  serviceName: string
  description: string
  version: string
  action: string
  isEnabled: boolean
  isDelete: boolean
  status: string
  memo: string
  meta?: ServiceMetaDockerInfo
  deployment?: ServiceDeploymentInfo
  containers: ServiceContainerStatusInfo[]
  _expanded?: boolean
}

export interface ServiceMetaDockerInfo {
  image: string
  registryDomain: string
  registryId: string
  alwaysPull: boolean
  command: string
  envConfig: string[]
  env?: string[]
  labels: { [key: string]: string }
  privileged: boolean
  capabilities?: ServiceCapabilitiesInfo
  logConfig?: ServiceLogConfigInfo
  resources?: ServiceResourcesInfo
  volumes: ServiceVolumeInfo[]
  network?: ServiceNetworkInfo
  restartPolicy?: ServiceRestartPolicyInfo
}

export interface ServiceCapabilitiesInfo {
  capAdd: string[]
  capDrop: string[]
}

export interface ServiceLogConfigInfo {
  type: string
  config: { [key: string]: string }
}

export interface ServiceResourcesInfo {
  memory: number
  memoryReservation: number
  maxCpuUsage: number
}

export interface ServiceVolumeInfo {
  type: ServiceVolumeType.VolumeTypeBind | ServiceVolumeType.VolumeTypeVolume
  target: string
  source: string
  "readonly": boolean
}

export interface ServiceNetworkInfo {
  mode: string
  hostname: string
  networkName: string
  useMachineHostname: boolean
  ports: ServicePortInfo[]
}

export interface ServicePortInfo {
  hostPort: number
  containerPort: number
  protocol: string
}

export interface ServiceRestartPolicyInfo {
  mode: string
  maxRetryCount: number
}

export interface ServiceDeploymentInfo {
  type: string
  mode: string
  replicas: number
  placements: ServicePlacementInfo[]
  schedule: ServiceScheduleInfo
  manualExec: boolean
}

export interface ServicePlacementInfo {
  mode: string
  key: string
  value: string
  isEqual: boolean
}

export interface ServiceScheduleInfo {
  timeout: string
  rules: string[]
}

export interface ServiceContainerStatusInfo {
  containerId: string
  containerName: string
  nodeId: string
  ip: string
  state: string
  status: string
  errorMsg: string
  image: string
  command: string
  network: string
  created: number
  started: number
  nextAt: number
  lastHeartbeat: number
  labels: { [key: string]: string }
  env: string[]
  mounts: Array<{ source: string; destination: string }>
  ports: Array<{ bindIP: string; privatePort: number; publicPort: number; type: string }>
  _expanded?: boolean
}

export function NewServiceEmptyInfo(): ServiceInfo {
  return {
    ...NewBaseEmptyInfo(),
    serviceId: "",
    groupId: "",
    serviceName: "",
    description: "",
    version: "",
    action: "",
    isEnabled: false,
    isDelete: false,
    status: "",
    memo: "",
    meta: NewServiceMetaDockerEmptyInfo(),
    deployment: NewServiceDeploymentInfo(),
    containers: []
  }
}

export function NewServiceMetaDockerEmptyInfo(): ServiceMetaDockerInfo {
  return {
    image: "",
    registryDomain: "docker.io",
    registryId: "",
    alwaysPull: false,
    command: "",
    envConfig: [],
    env: [],
    labels: {},
    volumes: [],
    privileged: false,
    capabilities: {
      capAdd: [
        "AUDIT_WRITE",
        "CHOWN",
        "DAC_OVERRIDE",
        "FOWNER",
        "FSETID",
        "KILL",
        "MKNOD",
        "NET_BIND_SERVICE",
        "NET_RAW",
        "SETFCAP",
        "SETGID",
        "SETPCAP",
        "SETUID",
        "SYS_CHROOT"
      ],
      capDrop: [
        "AUDIT_CONTROL",
        "BLOCK_SUSPEND",
        "DAC_READ_SEARCH",
        "IPC_LOCK",
        "IPC_OWNER",
        "LEASE",
        "LINUX_IMMUTABLE",
        "MAC_ADMIN",
        "MAC_OVERRIDE",
        "NET_ADMIN",
        "NET_BROADCAST",
        "SYSLOG",
        "SYS_ADMIN",
        "SYS_BOOT",
        "SYS_MODULE",
        "SYS_NICE",
        "SYS_PACCT",
        "SYS_PTRACE",
        "SYS_RAWIO",
        "SYS_RESOURCE",
        "SYS_TIME",
        "SYS_TTY_CONFIG",
        "WAKE_ALARM"
      ]
    },
    logConfig: {
      type: "json-file",
      config: {
        "max-file": "3",
        "max-size": "10m"
      }
    },
    resources: {
      memory: 0,
      memoryReservation: 0,
      maxCpuUsage: 0
    },
    network: {
      mode: ServiceNetworkMode.NetworkModeHost,
      hostname: "",
      networkName: "",
      useMachineHostname: false,
      ports: []
    },
    restartPolicy: {
      mode: ServiceRestartPolicyMode.RestartPolicyModeNo,
      maxRetryCount: 0
    }
  }
}

export function NewServiceDeploymentInfo(): ServiceDeploymentInfo {
  return {
    type: ServiceDeployType.DeployTypeBackground,
    mode: ServiceDeployMode.DeployModeReplicate,
    replicas: 1,
    placements: [],
    schedule: {
      timeout: "",
      rules: []
    },
    manualExec: false
  }
}

//
// type DeployMode string
// type DeployType string
//
// var (
// 	DeployModeGlobal    DeployMode = "global"
// 	DeployModeReplicate DeployMode = "replicate"
// )
//
// var (
// 	DeployTypeSchedule   DeployType = "schedule"
// 	DeployTypeBackground DeployType = "background"
// )
//
// type PlacementMode string
//
// var (
// 	PlacementModeLabel PlacementMode = "label"
// 	PlacementModeIP    PlacementMode = "ip"
// )
