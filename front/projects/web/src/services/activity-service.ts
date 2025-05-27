import { ActivityInfo, ResponseQuery } from "@/types"

class ActivityService {
  async query(data: any) {
    return await httpClient.post<ResponseQuery<ActivityInfo>>("/webapi/activity/query", data).then(res => res.data)
  }

  async queryAll(data: any) {
    return await httpClient
      .post<{
        [key: string]: ActivityInfo[]
      }>("/webapi/activity/all/query", data)
      .then(res => res.data)
  }
}

export const activityService = new ActivityService()
