import { NewUserEmptyInfo, UserInfo } from "#/index.ts"
import { GetUserRole, IsAdmin, IsSuperAdmin, IsUser } from "@/utils"

export default defineStore("user", () => {
  const userInfo = ref<UserInfo>(NewUserEmptyInfo())

  const isLogged = computed(() => !!userInfo.value.userId)
  const userRole = computed(() => GetUserRole(userInfo.value.role))
  const isAdmin = computed(() => IsAdmin(userInfo.value.role) || IsSuperAdmin(userInfo.value.role))
  const isSuperAdmin = computed(() => IsSuperAdmin(userInfo.value.role))
  const isUser = computed(() => IsUser(userInfo.value.role))

  const init = async () => {
    return await userService.getMe(true).then(data => {
      userInfo.value = data
    })
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = info
  }

  function clearUserInfo() {
    userInfo.value = NewUserEmptyInfo()
  }

  return {
    userInfo,
    userRole,
    isAdmin,
    isSuperAdmin,
    isUser,
    isLogged,
    init,
    setUserInfo,
    clearUserInfo
  }
})
