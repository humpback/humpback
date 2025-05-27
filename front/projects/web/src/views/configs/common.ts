import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { cloneDeep, find, omitBy } from "lodash-es"
import { ConfigType } from "@/models"

export const sortOptions = ["configName", "updatedAt", "createdAt"]

export const defaultSort = NewSortInfo("configName", "asc")
export const defaultPage = NewPageInfo(1, 20)
export const defaultFilter = { configType: 0 }

export class QueryConfigsInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["configName"], defaultPage, defaultSort, sortOptions, cloneDeep(defaultFilter))
    const configType = queryInfo["configType"] ? Number(queryInfo["configType"]).valueOf() : defaultFilter.configType
    this.filter.configType =
      !isNaN(configType) && find([0, ConfigType.Volume, ConfigType.Static], x => x === configType) ? configType : defaultFilter.configType
  }

  urlQuery() {
    return {
      query: Object.assign({}, this.filter.configType !== defaultFilter.configType ? { configType: this.filter.configType } : {}, this.getBaseQuery())
    }
  }

  searchParams() {
    return omitBy(this, (value, key) => key.startsWith("_"))
  }
}
