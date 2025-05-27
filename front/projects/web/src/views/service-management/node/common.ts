import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { cloneDeep, omitBy } from "lodash-es"

export const sortOptions = ["ipAddress", "name"]

export const defaultSort = NewSortInfo("ipAddress", "asc")
export const defaultPage = NewPageInfo(1, 20)
export const defaultFilter = {}

export class QueryGroupNodesInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["keywords"], defaultPage, defaultSort, sortOptions, cloneDeep(defaultFilter))
  }

  urlQuery() {
    return {
      query: Object.assign({}, this.getBaseQuery())
    }
  }

  searchParams() {
    return omitBy(this, (value, key) => key.startsWith("_"))
  }
}
