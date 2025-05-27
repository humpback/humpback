<script lang="ts" setup>
import { cloneDeep } from "lodash-es"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const stateStore = useStateStore()

const breadcrumbs = computed(() => {
  const list = cloneDeep(route.meta.breadcrumb) || []
  for (const b of list) {
    if (b.href) {
      b.href = b.href?.replace(":groupId", stateStore.getGroup()?.groupId || "")
      b.href = b.href?.replace(":serviceId", stateStore.getService()?.serviceId || "")
    }
    if (b.customName === "service") {
      b.name = stateStore.getService()?.serviceName
    }
    if (b.customName === "group") {
      b.name = stateStore.getGroup()?.groupName
    }
  }
  return list
})

function navigateToRoute(href: string, event: MouseEvent) {
  if (event.ctrlKey || event.metaKey) {
    window.open(href, "_blank")
  } else {
    router.push(href)
  }
}
</script>

<template>
  <el-breadcrumb v-if="breadcrumbs.length > 0" separator="/">
    <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">
      <a v-if="item.isLink" :href="item.href" @click.prevent.stop="navigateToRoute(item.href!, $event)">
        {{ item?.name ? item.name : item?.i18nLabel ? t(item?.i18nLabel) : "" }}
      </a>
      <span v-else>{{ item?.name ? item.name : item?.i18nLabel ? t(item?.i18nLabel) : "" }}</span>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<style lang="scss" scoped></style>
