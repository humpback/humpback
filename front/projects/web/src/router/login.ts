import type { RouteRecordRaw } from "vue-router"
import login from "@/views/login/login.vue"
import { PageLimitRole } from "@/models"

export default <RouteRecordRaw[]>[
  {
    path: "/login",
    name: "login",
    component: login,
    meta: {
      loginLimit: PageLimitRole.Logout
    }
  }
]
