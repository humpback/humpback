import { includes, isNumber, isString, toLower } from "lodash-es"
import { UserRole } from "@/models"

export function TableHeight(invalidHeight: number, minHeight = 500) {
  const pageStore = usePageStore()
  const height = pageStore.screenHeight - invalidHeight
  if (height < minHeight) {
    return "auto"
  }
  return height
}

export function SetWebTitle(title: string) {
  window.document.title = title
}

export function IncludesIgnoreCase(str?: string, subStr?: string): boolean {
  if (!str) {
    return false
  }
  return includes(toLower(str), toLower(subStr || ""))
}

export function IsSuperAdmin(role: number) {
  return role === UserRole.SuperAdmin
}

export function IsAdmin(role: number) {
  return role === UserRole.Admin
}

export function IsUser(role: number) {
  return role === UserRole.User
}

export function GetUserRole(role: number) {
  switch (role) {
    case UserRole.SuperAdmin:
      return UserRole.SuperAdmin
    case UserRole.Admin:
      return UserRole.Admin
    default:
      return UserRole.User
  }
}

export function ParseNumber(number: any, defaultValue?: number) {
  if (isNumber(number)) {
    return number
  }
  const result = isString(number) ? Number(number) : NaN
  return isNaN(result) ? defaultValue : result
}

export function BytesToGB(bytes: number): number {
  const gb = bytes / 1024 ** 3
  return Number(gb.toFixed(2)).valueOf()
}
