export interface ContainersPerformance {
  statsAt: number
  containers: Array<{
    containerId: string
    nodeId: string
    isSuccess: boolean
    error: string
    stats?: ContainerStats
  }>
}

export interface ContainerStats {
  cpuPercent: number
  memoryUsed: number
  memoryLimit: number
  ioRead: number
  ioWrite: number
  statsAt: number
  networks: ContainerNetworkStats[]
}

export interface ContainerNetworkStats {
  name: string
  rxBytes: number
  txBytes: number
}
