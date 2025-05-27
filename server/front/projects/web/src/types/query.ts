import { ParseNumber } from "@/utils"
import { find, isString } from "lodash-es"

export class QueryInfo {
  keywords: string
  mode: string
  filter: {
    [key: string]: any
  }
  pageInfo?: PageInfo
  sortInfo?: SortInfo

  private _defaultPage?: PageInfo
  private _defaultSort?: SortInfo
  private _sortOptions?: string[]
  private _modeOptions: string[]

  constructor(queryInfo: any, modeOptions: string[], defaultPage?: PageInfo, defaultSort?: SortInfo, sortOptions?: string[], defaultFilter?: any) {
    this._modeOptions = modeOptions
    this._defaultPage = defaultPage
    this._defaultSort = defaultSort

    this.keywords = queryInfo["keywords"] ? queryInfo["keywords"] : ""
    this.mode = find(this._modeOptions, x => x === queryInfo["mode"]) ? queryInfo["mode"] : this._modeOptions[0]
    this.filter = defaultFilter
    this.pageInfo = this._defaultPage
      ? ({
          index: ParseNumber(queryInfo["pageIndex"], this._defaultPage.index)!,
          size: ParseNumber(queryInfo["pageSize"], this._defaultPage.size)!
        } as PageInfo)
      : this._defaultPage
    this.sortInfo =
      this._defaultSort && sortOptions
        ? ({
            field: find(this._sortOptions, x => x === (queryInfo["sortField"] as string)) ? queryInfo["sortField"] : this._defaultSort.field,
            order: ParseQueryOrder(queryInfo["order"], this._defaultSort.order)
          } as SortInfo)
        : this._defaultSort
  }

  getBaseQuery() {
    return {
      keywords: this.keywords !== "" ? this.keywords : undefined,
      mode: this.keywords !== "" && this.mode !== this._modeOptions[0] ? this.mode : undefined,
      pageIndex: this._defaultPage && this.pageInfo!.index === this._defaultPage?.index ? undefined : this.pageInfo!.index,
      pageSize: this._defaultPage && this.pageInfo!.size === this._defaultPage?.size ? undefined : this.pageInfo!.size,
      sortField: this._defaultSort && this.sortInfo!.field === this._defaultSort?.field ? undefined : this.sortInfo!.field,
      order: this._defaultSort && this.sortInfo!.order === this._defaultSort?.order ? undefined : this.sortInfo!.order
    }
  }
}

export interface PageInfo {
  index: number
  size: number
}

export function NewPageInfo(index: number, size: number): PageInfo {
  return {
    index: index,
    size: size
  }
}

export interface SortInfo {
  field: string
  order: "asc" | "desc"
}

export function NewSortInfo(field: string, order: "asc" | "desc"): SortInfo {
  return {
    field: field,
    order: order
  }
}

export interface QueryList<T> {
  total: number
  pageInfo: PageInfo
  objects: T[]
}

function ParseQueryOrder(value: any, defaultValue: "desc" | "asc" = "asc") {
  const result = isString(value) ? (value as string).trim() : defaultValue
  if (["desc", "asc"].indexOf(result) === -1) {
    return defaultValue
  }
  return result
}
