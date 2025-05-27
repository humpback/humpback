<script lang="ts" setup>
import MenuContainer from "./menu-container.vue"
import { MenuInfo, menuList } from "./types"

const emits = defineEmits<{
  (e: "menuClick"): void
}>()

const route = useRoute()
const pageStore = usePageStore()

const activeMenu = ref("dashboard")
const defaultOpenMenus = ref(["workspace"])

const menus = ref<MenuInfo[]>(menuList)

watch(
  () => route.name,
  () => {
    if (route.name) {
      activeMenu.value = route.meta.currentMenu || (route.name as string)
    }
  },
  { immediate: true }
)
</script>

<template>
  <el-menu
    :collapse="pageStore.menuIsCollapse"
    :collapse-transition="false"
    :default-active="activeMenu"
    :default-openeds="defaultOpenMenus"
    class="website-menu"
    popper-class="website-menu-popover"
    text-color=""
    @select="emits('menuClick')">
    <menu-container v-for="item in menus" :key="item.name" :level="0" :menuInfo="item" />
  </el-menu>
</template>

<style lang="scss" scoped>
.website-menu {
  --el-menu-bg-color: var(--hp-aside-bg-color);
  --el-menu-text-color: var(--hp-aside-text-color);
  --el-menu-active-color: var(--hp-aside-active-color);
  --el-menu-base-level-padding: var(--hp-menu-padding);
  --el-menu-sub-item-height: var(--hp-menu-item-height);
  --el-menu-item-height: var(--hp-menu-item-height);
  width: 100%;
  border-right: 0;
  padding: 0 6px;

  :deep(li) {
    margin: 6px 0;
    border-radius: 8px;

    .el-sub-menu__title {
      border-radius: 8px;

      &:hover {
        background-color: var(--hp-aside-active-bg-color);
      }
    }

    &.el-menu-item {
      &:hover {
        background-color: var(--hp-aside-active-bg-color);
      }

      &.is-active {
        background-color: var(--hp-aside-active-bg-color);
      }
    }
  }
}
</style>
