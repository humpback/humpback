import { ResponseQuery, ResponseSuccess } from "@/types"

class NodeService {
  async info(id: string) {
    return await httpClient.get<NodeInfo>(`/webapi/node/${id}/info`).then(res => res.data)
  }

  async list() {
    return await httpClient.get<NodeInfo[]>(`/webapi/node/list`).then(res => res.data)
  }

  async query(data: any) {
    return await httpClient.post<ResponseQuery<NodeInfo>>("/webapi/node/query", data).then(res => res.data)
  }

  async create(data: any) {
    return await httpClient.post<string>("/webapi/node", data).then(res => res.data)
  }

  async updateLabel(data: any) {
    return await httpClient.put<string>("/webapi/node/labels", data).then(res => res.data)
  }

  async updateSwitch(data: any) {
    return await httpClient.put<string>("/webapi/node/switch", data).then(res => res.data)
  }

  async delete(id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/node/${id}`).then(res => res.data)
  }
}

export const nodeService = new NodeService()
