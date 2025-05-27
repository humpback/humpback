import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { omitBy } from "lodash-es"

export const sortOptions = ["groupName", "updatedAt", "createdAt"]

export const defaultSort = NewSortInfo("groupName", "asc")
export const defaultPage = NewPageInfo(1, 20)

export class QueryGroupsInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["groupName"], defaultPage, defaultSort, sortOptions)
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
