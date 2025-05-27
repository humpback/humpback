import axios, { type AxiosRequestConfig, type AxiosResponse, type Method } from "axios"
import { find, toLower } from "lodash-es"
import { GetCurrentLocale, GetI18nMessage } from "@/locales"
import { eventEmitter, globalLoading } from "utils/index.ts"

export interface HttpRequestOptions extends AxiosRequestConfig {
  showLoading?: boolean | false
  disableErrMsg?: boolean | false
  loadingMessage?: string
}

const _loadingCount = ref(0)
const loading = computed({
  get() {
    return _loadingCount.value > 0
  },
  set(val: boolean) {
    _loadingCount.value += val ? 1 : -1
    _loadingCount.value = Math.max(0, _loadingCount.value)
  }
})

class HttpClientService {
  public async get<T>(url: string, options?: HttpRequestOptions): Promise<AxiosResponse<T>> {
    return this.request("GET", url, null, options)
  }

  public async put<T>(url: string, data: any, options?: HttpRequestOptions): Promise<AxiosResponse<T>> {
    return this.request("PUT", url, data, options)
  }

  public async post<T>(url: string, data: any, options?: HttpRequestOptions): Promise<AxiosResponse<T>> {
    return this.request("POST", url, data, options)
  }

  public async delete<T>(url: string, options?: HttpRequestOptions): Promise<AxiosResponse<T>> {
    return this.request("DELETE", url, null, options)
  }

  private handleInnerErr: any = (options?: HttpRequestOptions) => (err: any) => {
    console.error(err)
    if (err.status === 401) {
      eventEmitter.emit("API:NO_AUTH", err.response.data)
      throw err
    }
    if (err.status === 403) {
      eventEmitter.emit("API:NO_PERMISSION", err.response.data)
      throw err
    }

    const data = err?.response?.data || ""
    if (typeof data === "object") {
      const resource = ["R4Group-NotExist", "R4Service-NotExist"]
      if (find(resource, x => x === data?.code)) {
        eventEmitter.emit("API:RESOURCE_NOT_EXIST", data)
      } else if (!options?.disableErrMsg) {
        eventEmitter.emit("API:FAILED", data.statusCode === 500 ? undefined : data.errMsg)
      }
      throw err
    }
    if (err.message && toLower(err.name) === "axioserror") {
      eventEmitter.emit("API:FAILED", toLower(err.message).includes("timeout") ? GetI18nMessage("err.timeout") : err.message)
      throw err
    }
    eventEmitter.emit("API:FAILED")
    throw err
  }

  private async request<T>(method: Method, url: string, data: any, options?: HttpRequestOptions): Promise<AxiosResponse<T>> {
    options = options || ({} as HttpRequestOptions)

    options.headers = Object.assign(
      {
        "Accept": "application/json, text/javascript, */*",
        "Accept-Language": GetCurrentLocale(),
        "Content-Language": GetCurrentLocale()
      },
      method === "POST" || method === "PUT" ? { "Content-Type": "application/json" } : {},
      options.headers
    )

    if (options?.showLoading) {
      loading.value = true
      if (loading.value) {
        globalLoading.show(options.loadingMessage)
      }
    }
    return await axios
      .request<T>({
        url: url,
        method: method,
        data: data,
        headers: options.headers,
        responseType: options?.responseType,
        params: Object.assign({}, options?.params, { _t: new Date().valueOf() }),
        timeout: options?.timeout || 60000
      })
      .catch<AxiosResponse<T>>(this.handleInnerErr(options))
      .finally(() => {
        if (options?.showLoading) {
          loading.value = false
          if (!loading.value) {
            globalLoading.close()
          }
        }
      })
  }
}

export const httpClient = new HttpClientService()

const NOOP = () => {}

export function CreateCancelRequest<T extends (...args: any[]) => Promise<any>>(asyncRequest: T): (...args: Parameters<T>) => Promise<Awaited<ReturnType<T>>> {
  let cancel = NOOP
  return (...args: any[]) => {
    return new Promise((resolve, reject) => {
      cancel()
      cancel = () => {
        resolve = reject = NOOP
      }
      asyncRequest(...args).then(
        (res: any) => resolve(res),
        (err: any) => reject(err)
      )
    })
  }
}
