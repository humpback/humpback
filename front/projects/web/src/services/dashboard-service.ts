import { DashboardResourceStatisticsInfo } from "@/types"

class DashboardService {
  async getResourceStatistics() {
    return await httpClient.get<DashboardResourceStatisticsInfo>("/webapi/dashboard/resource-statistics").then(res => res.data)
  }
}

export const dashboardService = new DashboardService()
