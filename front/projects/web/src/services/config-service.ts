import { ConfigInfo, ResponseQuery, ResponseSuccess } from "@/types"

class ConfigService {
  async info(id: string) {
    return await httpClient.get<ConfigInfo>(`/webapi/config/${id}/info`).then(res => res.data)
  }

  async query(data: any) {
    return await httpClient.post<ResponseQuery<ConfigInfo>>("/webapi/config/query", data).then(res => res.data)
  }

  async create(data: any) {
    return await httpClient.post<string>("/webapi/config", data).then(res => res.data)
  }

  async update(data: any) {
    return await httpClient.put<string>("/webapi/config", data).then(res => res.data)
  }

  async delete(id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/config/${id}`).then(res => res.data)
  }
}

export const configService = new ConfigService()
