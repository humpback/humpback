import type { RouteRecordRaw } from "vue-router"
import { PageActivity, PageGroupDetail, PageLimitRole, PageServiceDetail, PageUserRelated } from "@/models"
import { find } from "lodash-es"

const administrator = <RouteRecordRaw[]>[
  {
    path: "/ws/registries",
    name: "registries",
    component: () => import("@/views/administration/registries/registries.vue"),
    meta: {
      loginLimit: PageLimitRole.Login,
      onlyAdmin: true
    }
  },
  {
    path: "/ws/nodes",
    name: "nodes",
    component: () => import("@/views/administration/nodes/nodes.vue"),
    meta: {
      loginLimit: PageLimitRole.Login,
      onlyAdmin: true
    }
  },
  {
    path: "/ws/user-related/:mode",
    name: "userRelated",
    component: () => import("@/views/administration/user-related/user-related.vue"),
    beforeEnter: (to, from, next) => {
      return find([PageUserRelated.Users, PageUserRelated.Teams], x => x === to.params.mode) ? next() : next({ name: "404" })
    },
    meta: {
      loginLimit: PageLimitRole.Login,
      onlyAdmin: true,
      webTitle: {
        params: "mode"
      }
    }
  }
]

const serviceManagement = <RouteRecordRaw[]>[
  {
    path: "/ws/groups",
    name: "groups",
    component: () => import("@/views/service-management/groups.vue"),
    meta: {
      loginLimit: PageLimitRole.Login
    }
  },
  {
    path: "/ws/group/:groupId/:mode",
    name: "groupDetail",
    component: () => import("@/views/service-management/group-detail.vue"),
    beforeEnter: (to, from, next) => {
      return find([PageGroupDetail.Services, PageGroupDetail.Nodes], x => x === to.params.mode) ? next() : next({ name: "404" })
    },
    meta: {
      loginLimit: PageLimitRole.Login,
      currentMenu: "groups",
      webTitle: {
        params: "mode"
      },
      routerKeys: ["groupId"],
      breadcrumb: [
        { href: "/ws/groups", i18nLabel: "breadcrumb.groupOverview", isLink: true },
        { isLink: false, customName: "group" }
      ]
    }
  },
  {
    path: "/ws/group/:groupId/service/:serviceId/:mode",
    name: "serviceInfo",
    component: () => import("@/views/service-management/service/service-detail.vue"),
    beforeEnter: (to, from, next) => {
      return find(
        [
          PageServiceDetail.BasicInfo,
          PageServiceDetail.Application,
          PageServiceDetail.Deployment,
          PageServiceDetail.Instances,
          PageServiceDetail.Log,
          PageServiceDetail.Performance,
          PageServiceDetail.Activity
        ],
        x => x === to.params.mode
      )
        ? next()
        : next({ name: "404" })
    },
    meta: {
      loginLimit: PageLimitRole.Login,
      currentMenu: "groups",
      routerKeys: ["groupId", "serviceId"],
      breadcrumb: [
        { href: "/ws/groups", i18nLabel: "breadcrumb.groupOverview", isLink: true },
        { href: `/ws/group/:groupId/${PageGroupDetail.Services}`, isLink: true, customName: "group" },
        { i18nLabel: "breadcrumb.serviceOverview", isLink: false },
        { isLink: false, customName: "service" }
      ]
    }
  }
]

export default <RouteRecordRaw[]>[
  {
    path: "",
    name: "workspace",
    redirect: "/ws/dashboard",
    component: () => import("@/views/layout/layout.vue"),
    children: [
      {
        path: "/ws/my-account",
        name: "myAccount",
        component: () => import("@/views/my-account/my-account.vue"),
        meta: { loginLimit: PageLimitRole.Login }
      },
      {
        path: "/ws/dashboard",
        name: "dashboard",
        component: () => import("@/views/dashboard/dashboard.vue"),
        meta: {
          loginLimit: PageLimitRole.Login
        }
      },
      {
        path: "/ws/configs",
        name: "configs",
        component: () => import("@/views/configs/configs.vue"),
        meta: { loginLimit: PageLimitRole.Login }
      },
      {
        path: "/ws/activities/:mode",
        name: "activities",
        component: () => import("@/views/activity/activities.vue"),
        beforeEnter: (to, from, next) => {
          return find(
            [
              PageActivity.Groups,
              PageActivity.Services,
              PageActivity.Configs,
              PageActivity.Nodes,
              PageActivity.Registries,
              PageActivity.Users,
              PageActivity.Teams
            ],
            x => x === to.params.mode
          )
            ? next()
            : next({ name: "404" })
        },
        meta: {
          loginLimit: PageLimitRole.Login
        }
      },
      ...administrator,
      ...serviceManagement
    ]
  }
]
