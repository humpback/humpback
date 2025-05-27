import { ResponseQuery, ResponseSuccess } from "@/types"

class ServiceService {
  async info(groupId: string, id: string) {
    return await httpClient.get<ServiceInfo>(`/webapi/group/${groupId}/service/${id}/info`).then(res => res.data)
  }

  async query(groupId: string, data: any) {
    return await httpClient.post<ResponseQuery<ServiceInfo>>(`/webapi/group/${groupId}/service/query`, data).then(res => res.data)
  }

  async total(groupId: string) {
    return await httpClient.get<number>(`/webapi/group/${groupId}/service/total`).then(res => res.data)
  }

  async create(groupId: string, data: any) {
    return await httpClient.post<string>(`/webapi/group/${groupId}/service`, data).then(res => res.data)
  }

  async clone(groupId: string, data: any) {
    return await httpClient.post<string>(`/webapi/group/${groupId}/service/clone`, data).then(res => res.data)
  }

  async update(groupId: string, data: any) {
    return await httpClient.put<string>(`/webapi/group/${groupId}/service`, data).then(res => res.data)
  }

  async operate(groupId: string, data: { serviceId: string; action: "Start" | "Stop" | "Restart" | "Enable" | "Disable" }) {
    return await httpClient.put<string>(`/webapi/group/${groupId}/service/operate`, data).then(res => res.data)
  }

  async delete(groupId: string, id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/group/${groupId}/service/${id}`).then(res => res.data)
  }
}

export const serviceService = new ServiceService()
