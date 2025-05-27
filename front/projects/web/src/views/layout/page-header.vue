<script lang="ts" setup>
import { storageEventBus, StorageEventType } from "@/utils"
import VGlobalSearch from "@/components/business/v-search/VGlobalSearch.vue"

enum Menu {
  Logout = "logout",
  MyAccount = "myAccount",
  Help = "help"
}

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const pageStore = usePageStore()
const userStore = useUserStore()

const userMenu = computed<
  {
    label: string
    value: string
    icon: any
    color: string
    component?: any
    hide?: boolean
  }[]
>(() => {
  return [
    { label: "menu.myAccount", value: Menu.MyAccount, icon: IconMdiUserOutline, color: "green" },
    { label: "menu.logout", value: Menu.Logout, icon: IconMdiLogoutVariant, color: "var(--el-color-info)" }
  ]
})

function handleUserMenuClick(v: string) {
  switch (v) {
    case Menu.Logout:
      userService.logout().finally(() => {
        userStore.clearUserInfo()
        router.push({ name: "login", query: { redirectURL: route.fullPath } })
        storageEventBus.SendMessage(StorageEventType.Logout)
      })
      return
    case Menu.MyAccount:
      router.push({ name: "myAccount" })
  }
}
</script>

<template>
  <div class="header-box">
    <div class="flex-1" style="height: 24px">
      <v-global-search />
    </div>
    <div class="d-flex">
      <div class="d-flex gap-5 mr-5">
        <el-button v-if="!pageStore.isSmallScreen" link> {{ t("btn.help") }}</el-button>
        <Language :iconSize="18" />
        <el-link :underline="false" class="github-icon" href="https://github.com/humpback/humpback" target="_blank" title="GitHub">
          <el-icon :size="20">
            <IconMdiGithub />
          </el-icon>
        </el-link>
        <v-role-admin v-if="!pageStore.isSmallScreen" :role="userStore.userInfo.role" />
      </div>
      <el-dropdown :show-timeout="0" placement="bottom-end" trigger="hover" @command="handleUserMenuClick">
        <div class="user-btn">
          <div class="user-icon">
            <el-icon :size="20">
              <IconMdiUserOutline />
            </el-icon>
          </div>
          <div v-if="userStore.userInfo.username" class="username"> {{ userStore.userInfo.username }}</div>
          <el-icon :size="20">
            <IconMdiChevronDown />
          </el-icon>
        </div>

        <template #dropdown>
          <el-dropdown-menu>
            <template v-for="item in userMenu" :key="item.value">
              <el-dropdown-item v-if="!item.hide" :command="item.value" :style="{ color: item.color }">
                <el-icon :size="20">
                  <component :is="item.icon" />
                </el-icon>
                <component :is="item.component" v-if="item.component"></component>
                <span v-else>
                  {{ t(item.label) }}
                </span>
              </el-dropdown-item>
            </template>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.header-box {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  padding: 0 10px;
  display: flex;
  align-items: center;
  color: inherit;

  .github-icon:hover {
    color: var(--el-link-text-color);
    opacity: 0.7;
  }

  .el-dropdown:focus-visible {
    outline: none;
  }

  .user-btn {
    display: flex;
    align-items: center;
    gap: 6px;

    &:hover {
      cursor: pointer;
      opacity: 0.7;
    }

    &:focus-visible {
      outline: none;
    }

    .user-icon {
      width: 24px;
      height: 24px;
      border-radius: 50%;
      background-color: #eaecf0;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .username {
      max-width: 200px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}
</style>
