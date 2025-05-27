class CommonService {
  async config() {
    return await httpClient.get<{ [key: string]: any }>("/webapi/common/config").then(res => res.data)
  }

  async searchGroupServiceByName(name: string) {
    return await httpClient
      .get<{
        [key: string]: Array<{ groupId: string; groupName: string; serviceId?: string; serviceName?: string }>
      }>(`/webapi/common/search/group-service/by-name?name=${name}`)
      .then(res => res.data)
  }
}

export const commonService = new CommonService()
