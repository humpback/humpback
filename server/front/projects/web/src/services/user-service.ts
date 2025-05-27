import { ResponseQuery, ResponseSuccess, UserInfo } from "#/index.ts"

class UserService {
  async login(data: any) {
    return await httpClient.post<UserInfo>("/webapi/user/login", data).then(res => res.data)
  }

  async logout() {
    return await httpClient.post("/webapi/user/logout", null, { disableErrMsg: true }).then(res => res.data)
  }

  async changePassword(data: any) {
    return await httpClient.put("/webapi/user/me/change-psd", data).then(res => res.data)
  }

  async getMe(startup?: boolean) {
    return await httpClient.get<UserInfo>("/webapi/user/me", { params: startup ? { startup: "true" } : undefined }).then(res => res.data)
  }

  async updateMeInfo(data: any) {
    return await httpClient.put("/webapi/user/me", data).then(res => res.data)
  }

  async info(id: string, includePassword?: boolean) {
    let url = includePassword ? `/webapi/user/${id}/info?p=true` : `/webapi/user/${id}/info`
    return await httpClient.get<UserInfo>(url).then(res => res.data)
  }

  async list() {
    return await httpClient.get<UserInfo[]>(`/webapi/user/list`).then(res => res.data)
  }

  async query(data: any) {
    return await httpClient.post<ResponseQuery<UserInfo>>("/webapi/user/query", data).then(res => res.data)
  }

  async queryByTeamId(teamId: string) {
    return await httpClient.get<UserInfo[]>(`/webapi/user/query-by-team/${teamId}`).then(res => res.data)
  }

  async create(data: any) {
    return await httpClient.post<string>("/webapi/user", data).then(res => res.data)
  }

  async update(data: any) {
    return await httpClient.put<string>("/webapi/user", data).then(res => res.data)
  }

  async delete(id: string) {
    return await httpClient.delete<ResponseSuccess>(`/webapi/user/${id}`).then(res => res.data)
  }
}

export const userService = new UserService()
