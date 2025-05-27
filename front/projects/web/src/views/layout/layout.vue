<script lang="ts" setup>
import PageHeader from "./page-header.vue"
import PageAside from "./page-aside.vue"
import VLoadingPage from "@/components/business/v-loading/VLoadingPage.vue"
import { map } from "lodash-es"

const route = useRoute()
const pageStore = usePageStore()

const leftWidth = computed(() => (pageStore.menuIsCollapse ? "var(--hp-aside-collapse-width)" : "var(--hp-aside-width)"))

function beforeEnter() {
  document.body.style.overflow = "hidden" // 动画前隐藏滚动条
}

function afterEnter() {
  document.body.style.overflow = "" // 动画后恢复滚动条
}
</script>

<template>
  <div id="page-container" :style="{ paddingLeft: leftWidth }">
    <div id="page-header" :style="{ paddingLeft: leftWidth }">
      <page-header />
    </div>
    <div id="page-aside" :style="{ width: leftWidth }">
      <page-aside />
    </div>
    <div id="page-main">
      <router-view v-slot="{ Component }">
        <transition mode="out-in" name="fade" @before-enter="beforeEnter()" @after-enter="afterEnter()">
          <Suspense>
            <component :is="Component" :key="`${route.name as string}_${map(route.meta.routerKeys, p => route.params[p]).join('_')}`" />
            <template #fallback>
              <v-loading-page />
            </template>
          </Suspense>
        </transition>
      </router-view>
    </div>
  </div>
</template>

<style lang="scss" scoped>
#page-container {
  box-sizing: border-box;
  width: 100%;
  min-width: 320px;
  height: 100%;
  min-height: 500px;

  #page-header {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    width: 100%;
    min-width: 320px;
    height: var(--hp-header-height);
    box-sizing: border-box;
    border-bottom: 1px solid var(--el-border-color);
    //box-shadow: 0 1px 2px var(--el-border-color);
    background-color: var(--hp-header-bg-color);
  }

  #page-aside {
    position: fixed;
    z-index: 1001;
    left: 0;
    top: 0;
    height: 100%;
    background-color: var(--hp-aside-bg-color);
  }

  #page-main {
    --hp-main-pading: 16px;
    box-sizing: border-box;
    padding: calc(var(--hp-header-height) + var(--hp-main-pading)) var(--hp-main-pading) var(--hp-main-pading) var(--hp-main-pading);
    max-width: 100%;
  }
}

/* 定义进入过渡的开始状态 */
.fade-enter-from {
  transform: translateX(30%);
  opacity: 0;
}

/* 定义进入过渡的结束状态 */
.fade-enter-to {
  transform: translateX(0);
  opacity: 1;
}

/* 定义离开过渡的开始状态 */
.fade-leave-from {
  opacity: 1;
}

/* 定义离开过渡的结束状态 */
.fade-leave-to {
  opacity: 0;
}

/* 定义进入过渡的过程 */
.fade-enter-active {
  transition:
    transform 0.2s ease,
    opacity 0.2s ease;
  will-change: transform;
}

/* 定义离开过渡的过程 */
.fade-leave-active {
  transition: opacity 0.2s ease;
  will-change: transform;
}
</style>
