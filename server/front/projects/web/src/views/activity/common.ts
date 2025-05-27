import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { cloneDeep, omitBy, toLower } from "lodash-es"

export const sortOptions = ["operateAt"]

export const defaultSort = NewSortInfo("operateAt", "desc")
export const defaultPage = NewPageInfo(1, 20)
export const defaultFilter = { startAt: 0, endAt: 0, groupId: "", action: "", operator: "", type: "" }

export class QueryActivityInfo extends QueryInfo {
  constructor(queryInfo: any, activityType: string, isAdmin?: boolean) {
    super(queryInfo, ["keywords"], defaultPage, defaultSort, sortOptions, cloneDeep(defaultFilter))
    const startAt = queryInfo["startAt"] ? Number(queryInfo["startAt"]).valueOf() : 0
    const endAt = queryInfo["endAt"] ? Number(queryInfo["endAt"]).valueOf() : 0
    this.filter.startAt = !isNaN(startAt) ? startAt : 0
    this.filter.endAt = !isNaN(endAt) ? endAt : 0
    this.filter.action = queryInfo["action"] ? queryInfo["action"] : ""
    this.filter.operator = isAdmin && queryInfo["operator"] ? queryInfo["operator"] : ""
    this.filter.groupId =
      (toLower(activityType) === PageActivity.Groups || toLower(activityType) === PageActivity.Services) && queryInfo["groupId"] ? queryInfo["groupId"] : ""
    this.filter.type = activityType
  }

  urlQuery() {
    return {
      query: Object.assign(
        {},
        {
          startAt: this.filter.startAt || undefined,
          endAt: this.filter.endAt || undefined,
          groupId: this.filter.groupId || undefined,
          action: this.filter.action || undefined,
          operator: this.filter.operator || undefined
        },
        this.getBaseQuery()
      )
    }
  }

  searchParams() {
    return omitBy(this, (value, key) => key.startsWith("_"))
  }
}
