import { NewPageInfo, NewSortInfo, QueryInfo } from "@/types"
import { omitBy, toLower } from "lodash-es"

export const sortOptions = ["url", "updatedAt", "createdAt"]

export const defaultSort = NewSortInfo("url", "asc")
export const defaultPage = NewPageInfo(1, 20)

export class QueryRegistryInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["url"], defaultPage, defaultSort, sortOptions, {})
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

export function isDefaultRegistry(url: string) {
  return toLower(url) === "docker.io"
}
