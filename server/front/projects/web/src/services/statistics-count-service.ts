import { StatisticsCountInfo } from "@/types"

class StatisticsCountService {
  async query(data: any) {
    return await httpClient.post<StatisticsCountInfo[]>("/webapi/statistics-count/query", data).then(res => res.data)
  }
}

export const statisticsCountService = new StatisticsCountService()
