import { ResponseQuery, ResponseSuccess, TeamInfo } from "@/types"

class TeamService {
  async info(id: string) {
    return await httpClient.get<TeamInfo>(`/webapi/team/${id}/info`).then(res => res.data)
  }

  async list() {
    return await httpClient.get<TeamInfo[]>(`/webapi/team/list`).then(res => res.data)
  }

  async query(data: any) {
    return await httpClient.post<ResponseQuery<TeamInfo>>("/webapi/team/query", data).then(res => res.data)
  }

  async queryByUserId(userId: string) {
    return await httpClient.get<TeamInfo[]>(`/webapi/team/query-by-user/${userId}`).then(res => res.data)
  }

  async create(data: any) {
    return await httpClient.post<string>("/webapi/team", data).then(res => res.data)
  }

  async update(data: any) {
    return await httpClient.put<string>("/webapi/team", data).then(res => res.data)
  }

  async delete(id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/team/${id}`).then(res => res.data)
  }
}

export const teamService = new TeamService()
