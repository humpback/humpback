import { GroupInfo, ResponseQuery, ResponseSuccess } from "@/types"

class GroupService {
  async info(id: string) {
    return await httpClient.get<GroupInfo>(`/webapi/group/${id}/info`).then(res => res.data)
  }

  async query(data: any) {
    return await httpClient.post<ResponseQuery<GroupInfo>>("/webapi/group/query", data).then(res => res.data)
  }

  async list() {
    return await httpClient.get<GroupInfo[]>(`/webapi/group/list`).then(res => res.data)
  }

  async queryNodes(id: string, data: any) {
    return await httpClient.post<ResponseQuery<NodeInfo>>(`/webapi/group/${id}/node/query`, data).then(res => res.data)
  }

  async getNodes(id: string) {
    return await httpClient.get<NodeInfo[]>(`/webapi/group/${id}/node/list`).then(res => res.data)
  }

  async create(data: any) {
    return await httpClient.post<string>("/webapi/group", data).then(res => res.data)
  }

  async update(data: any) {
    return await httpClient.put<string>("/webapi/group", data).then(res => res.data)
  }

  async updateNodes(groupId: string, data: any) {
    return await httpClient.put<string>(`/webapi/group/${groupId}/node`, data).then(res => res.data)
  }

  async delete(id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/group/${id}`).then(res => res.data)
  }
}

export const groupService = new GroupService()
