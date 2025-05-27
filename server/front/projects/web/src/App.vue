<script lang="ts" setup>
import { StorageEventType, storageEventBus } from "utils/index.ts"
import { GetUILocale } from "@/locales/index.ts"
import VLoadingPage from "@/components/business/v-loading/VLoadingPage.vue"

const { t } = useI18n()
const pageStore = usePageStore()
const userStore = useUserStore()

const handleChannelMessage = (data: any) => {
  switch (data?.type) {
    case StorageEventType.Login:
      if (!data.value || data.value?.userId != userStore.userInfo.userId) {
        pageStore.refreshPage.needRefresh = true
        pageStore.refreshPage.type = StorageEventType.Login
      }
      return
    case StorageEventType.Logout:
      pageStore.refreshPage.needRefresh = true
      pageStore.refreshPage.type = StorageEventType.Logout
  }
}

const handleFocus = () => {
  if (pageStore.refreshPage.needRefresh) {
    if (pageStore.refreshPage.type === StorageEventType.Login) {
      ShowWarningMsg(t("message.loginUserChangeEvent"))
    }
    location.reload()
  }
}

onMounted(() => {
  pageStore.setScreen()
  storageEventBus.SetMessageHandler(handleChannelMessage)
  window.addEventListener("resize", pageStore.setScreen)
  window.addEventListener("focus", handleFocus)
})

onBeforeUnmount(() => {
  storageEventBus.Close()
  window.removeEventListener("resize", pageStore.setScreen)
  window.removeEventListener("focus", handleFocus)
})
</script>

<template>
  <el-config-provider :locale="GetUILocale('elementPlus')">
    <Suspense>
      <router-view />
      <template #fallback>
        <v-loading-page />
      </template>
    </Suspense>
  </el-config-provider>
</template>

<style scoped></style>
