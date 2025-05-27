<script lang="ts" setup>
import Menu from "./menu/menu.vue"

const pageStore = usePageStore()

function clickLeftMenuBtn() {
  pageStore.changeCollapse()
}
</script>

<template>
  <div class="aside">
    <div class="aside-header">
      <v-logo v-if="!pageStore.menuIsCollapse" direction="horizontal" showText />
      <el-button link type="info" @click="clickLeftMenuBtn()">
        <el-icon :size="20">
          <IconMdiFormatIndentIncrease v-if="pageStore.menuIsCollapse" />
          <IconMdiFormatIndentDecrease v-else />
        </el-icon>
      </el-button>
    </div>
    <div class="aside-content">
      <el-scrollbar>
        <Menu />
      </el-scrollbar>
    </div>
    <div v-if="!pageStore.menuIsCollapse" class="aside-footer">
      <el-text type="info">@2024 v0.0.1</el-text>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.aside {
  height: 100vh;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;

  .aside-header {
    height: var(--hp-header-height);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 16px 0 16px;
  }

  .aside-content {
    overflow: auto;
    flex: 1;
  }

  .aside-footer {
    height: var(--hp-footer-height);
    display: flex;
    align-items: center;
    padding: 0 16px;
  }
}
</style>
