import { ContainersPerformance } from "@/types"

class GroupContainerService {
  async operate(groupId: string, serviceId: string, data: { nodeId: string; containerId: string; action: "Start" | "Stop" | "Restart" }) {
    return await httpClient.put<any>(`/webapi/group/${groupId}/service/${serviceId}/instance/operate`, data).then(res => res.data)
  }

  async performance(groupId: string, serviceId: string, data: { containers: Array<{ nodeId: string; containerId: string }> }) {
    return await httpClient.post<ContainersPerformance>(`/webapi/group/${groupId}/service/${serviceId}/instance/performances`, data).then(res => res.data)
  }

  async logs(
    groupId: string,
    serviceId: string,
    data: { nodeId: string; containerId: string; line?: number; startAt?: number; endAt?: number; showTimestamp?: boolean }
  ) {
    return await httpClient.post<string[]>(`/webapi/group/${groupId}/service/${serviceId}/instance/logs`, data).then(res => res.data)
  }
}

export const groupContainerService = new GroupContainerService()
