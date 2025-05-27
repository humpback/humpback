<script lang="ts" setup>
import { menuI18nPrefix, MenuInfo } from "./types"

const props = defineProps<{ menuInfo: MenuInfo }>()
const { t } = useI18n()
</script>

<template>
  <el-sub-menu :index="props.menuInfo.name" teleported>
    <template #title>
      <el-icon v-if="props.menuInfo.icon" :size="20" style="padding-bottom: 2px">
        <component :is="props.menuInfo.icon" />
      </el-icon>
      <span> {{ t(`${menuI18nPrefix}.${props.menuInfo.name}`) }}</span>
    </template>
    <slot />
  </el-sub-menu>
</template>

<style lang="scss" scoped></style>
