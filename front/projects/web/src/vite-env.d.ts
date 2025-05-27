/// <reference types="vite/client" />
/// <reference types="element-plus" />

declare const __APP_VERSION_: string

declare module "*.vue" {
  import Vue from "vue"
  export default Vue
}

declare module "v3-easyui"

import "vue-router"
//用于定义router中的meta结构
declare module "vue-router" {
  interface RouteMeta {
    currentMenu?: string
    isNew?: boolean
    onlyAdmin?: boolean
    loginLimit?: number // 0|undefined： 必须登录； -1: 忽略是否登录；1: 必须登出
    routerKeys?: string[]
    breadcrumb?: Array<{
      href?: string
      i18nLabel: string
      name?: string
      isLink?: boolean
      customName?: "service" | "group"
    }>
    webTitle?: {
      params?: string
    }
  }
}
