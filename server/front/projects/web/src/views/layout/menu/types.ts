import { shallowRef } from "vue"
import { PageActivity, PageUserRelated } from "@/models"

export interface MenuInfo {
  icon: any
  name: string
  routeName?: string
  params?: any
  query?: any
  path?: string
  children?: MenuInfo[]
  rolesLimit?: number[]
}

export const menuI18nPrefix = "menu"

export const menuList: MenuInfo[] = [
  {
    icon: shallowRef(IconMdiViewDashboard),
    name: "dashboard"
  },
  {
    icon: shallowRef(IconMdiCompany),
    name: "serviceManagement",
    routeName: "groups"
  },
  {
    icon: shallowRef(IconMdiTextBoxOutline),
    name: "configs"
  },
  {
    icon: shallowRef(IconMdiAccountFileText),
    name: "activities",
    params: { mode: PageActivity.Groups }
  },
  {
    icon: shallowRef(IconMdiCogOutline),
    name: "administration",
    rolesLimit: [UserRole.SuperAdmin, UserRole.Admin],
    children: [
      {
        icon: shallowRef(IconMdiAlphaCBoxOutline),
        name: "registries"
      },
      {
        icon: shallowRef(IconMdiServer),
        name: "nodes"
      },
      {
        icon: shallowRef(IconMdiAccount),
        name: "userRelated",
        params: { mode: PageUserRelated.Users }
      }
    ]
  }
]
