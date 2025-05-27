import type { RouteRecordRaw } from "vue-router"
import { createRouter, createWebHistory, NavigationGuardNext, RouteLocationNormalized } from "vue-router"
import { configure, done, start } from "nprogress"
import { GetI18nMessage } from "@/locales"
import "nprogress/nprogress.css"
import { eventEmitter, SetWebTitle } from "@/utils"
import { disposeStore } from "@/stores"
import { SessionStorageCurrentGroupId } from "@/models"

configure({
  easing: "ease", // 动画方式
  speed: 800, // 递增进度条的速度
  showSpinner: false, // 是否显示加载ico
  trickleSpeed: 100, // 自动递增间隔
  minimum: 0.3 // 初始化时的最小百分比
})

const routes: RouteRecordRaw[] = [
  {
    path: "/:pathMatch(.*)*",
    name: "404",
    component: () => import("@/views/common/404/404.vue"),
    meta: {}
  },
  {
    path: "/401",
    name: "401",
    component: () => import("@/views/common/401/401.vue"),
    meta: {}
  }
]

const modules = import.meta.glob("./**/*.ts", { eager: true })

// 获取当前目录下所有ts文件中的router规则，并合并到routes中。
Object.keys(modules).forEach(item => {
  const modulesRoutes = (modules[item] as any).default || null
  if (!modulesRoutes) {
    return
  }
  if (Array.isArray(modulesRoutes)) {
    routes.push(...modulesRoutes)
  } else {
    routes.push(modulesRoutes)
  }
})

const router = createRouter({
  history: createWebHistory("/"),
  routes
})

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  if (to.name !== from.name) {
    start()
  }
  const userStore = useUserStore()
  if (!userStore.isLogged) {
    if (to.meta?.loginLimit === PageLimitRole.Login) {
      next({ name: "login" })
      return
    }
    next()
    return
  }
  if (to.meta?.loginLimit === PageLimitRole.Logout) {
    next({ name: "dashboard" })
    return
  }
  if (to.meta?.onlyAdmin && !userStore.isAdmin) {
    next({ name: "401" })
    return
  }
  next()
})

router.afterEach((to: RouteLocationNormalized, from: RouteLocationNormalized) => {
  if (to.fullPath !== from.fullPath) {
    let titleKey = to.meta?.webTitle?.params ? (to.params[to.meta?.webTitle?.params] as string) : (to.name as string)
    SetWebTitle(`${GetI18nMessage("webTitle." + titleKey)}`)
  }

  if (to.name !== from.name) {
    done()
  }
})

export default router

eventEmitter.on("API:NO_AUTH", (...args: any[]) => {
  const data: any = Array.isArray(args) && args.length > 0 ? args[0] : {}
  switch (data?.code) {
    case "R40101": {
      ShowErrMsg(data.errMsg)
      disposeStore()
      if (router) {
        router.push({ name: "login" })
      } else {
        window.location.href = "/login"
      }
      return
    }
    case "R40102": {
      disposeStore()
    }
  }
})

eventEmitter.on("API:RESOURCE_NOT_EXIST", (...args: any[]) => {
  const data: any = Array.isArray(args) && args.length > 0 ? args[0] : {}
  switch (data?.code) {
    case "R4Group-NotExist": {
      ShowErrMsg(data.errMsg)
      if (router) {
        router.push({ name: "groups" })
      } else {
        window.location.href = "/ws/groups"
      }
      return
    }
    case "R4Service-NotExist": {
      ShowErrMsg(data.errMsg)
      const groupId = sessionStorage.getItem(SessionStorageCurrentGroupId)
      if (groupId) {
        if (router) {
          router.push({ name: "groupDetail", params: { groupId: groupId, mode: PageGroupDetail.Services } })
        } else {
          window.location.href = `/ws/group/${groupId}/${PageGroupDetail.Services}`
        }
      } else {
        if (router) {
          router.push({ name: "groups" })
        } else {
          window.location.href = "/ws/groups"
        }
      }
    }
  }
})
