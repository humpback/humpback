import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { omitBy } from "lodash-es"

export const sortOptions = ["name", "updatedAt", "createdAt"]

export const defaultSort = NewSortInfo("name", "asc")
export const defaultPage = NewPageInfo(1, 20)

export class QueryTeamInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["name"], defaultPage, defaultSort, sortOptions, {})
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
