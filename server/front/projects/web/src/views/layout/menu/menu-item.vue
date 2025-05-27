<script lang="ts" setup>
import { menuI18nPrefix, MenuInfo } from "./types"

const props = defineProps<{ menuInfo: MenuInfo }>()
const { t } = useI18n()
const router = useRouter()

function navigateToRoute(event: MouseEvent, href: string) {
  if (event.ctrlKey || event.metaKey) {
    window.open(href, "_blank")
  } else {
    router.push(href)
  }
}
</script>

<template>
  <router-link v-slot="{ href }" :to="{ name: props.menuInfo?.routeName || props.menuInfo.name, params: props.menuInfo.params }" custom>
    <div @click="navigateToRoute($event, href)">
      <el-menu-item
        :index="props.menuInfo?.routeName || props.menuInfo.name"
        :route="{ name: props.menuInfo?.routeName || props.menuInfo.name, params: props.menuInfo.params }">
        <template #title>
          {{ t(`${menuI18nPrefix}.${props.menuInfo.name}`) }}
        </template>
        <el-icon v-if="props.menuInfo.icon" :size="20" style="padding-bottom: 2px">
          <component :is="props.menuInfo.icon" />
        </el-icon>
      </el-menu-item>
    </div>
  </router-link>
</template>

<style lang="scss" scoped></style>
