<script lang="ts" setup>
import SubMenu from "./sub-menu.vue"
import MenuItem from "./menu-item.vue"
import { MenuInfo } from "./types"
import { find } from "lodash-es"

const props = defineProps<{ menuInfo: MenuInfo }>()

const userStore = useUserStore()

const hasChild = computed(() => Array.isArray(props.menuInfo.children) && props.menuInfo.children.length > 0)

const menuComponent = computed(() => {
  return hasChild.value ? SubMenu : MenuItem
})

function show() {
  if (Array.isArray(props.menuInfo.rolesLimit) && props.menuInfo.rolesLimit.length > 0) {
    return !!find(props.menuInfo.rolesLimit, x => x === userStore.userRole)
  }
  return true
}
</script>

<template>
  <component :is="menuComponent" v-if="show()" :menuInfo="menuInfo">
    <template v-if="hasChild">
      <menu-container v-for="item in menuInfo.children" :key="item.name" :menuInfo="item" />
    </template>
  </component>
</template>

<style lang="scss" scoped></style>
