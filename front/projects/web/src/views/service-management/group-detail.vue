<script lang="ts" setup>
import PageServices from "@/views/service-management/service/services.vue"
import PageNodes from "@/views/service-management/node/nodes.vue"
import { TabPaneName } from "element-plus"
import { PageGroupDetail } from "@/models"
import VPageTitle from "@/components/business/v-page/VPageTitle.vue"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const stateStore = useStateStore()

const activeTab = ref<TabPaneName | undefined>(route.params?.mode as string)

const groupInfo = computed(() => stateStore.getGroup())

const options = reactive<{ name: string; label: string; component: any }[]>([
  { name: PageGroupDetail.Services, label: "header.services", component: shallowRef(PageServices) },
  { name: PageGroupDetail.Nodes, label: "header.nodes", component: shallowRef(PageNodes) }
])

async function changeTab(name: TabPaneName) {
  await router.replace({ params: Object.assign({}, route.params, { mode: name }) })
  activeTab.value = name
}
</script>

<template>
  <div>
    <v-page-title :title="groupInfo?.groupName" show-breadcrumbs />

    <el-tabs :model-value="activeTab" class="tab-box" type="card" @update:modelValue="changeTab">
      <el-tab-pane v-for="item in options" :key="item.name" :name="item.name">
        <template #label>
          <el-badge :offset="[10, 2]" class="mr-4 mb-2" color="#28c3d7">
            <template #content>{{ groupInfo?.total[item.name] || 0 }}</template>
            <strong>{{ t(item.label) }}</strong>
          </el-badge>
        </template>
        <component :is="item.component" v-if="activeTab === item.name" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style lang="scss" scoped>
.tab-box {
  :deep(.el-tabs__header.is-top) {
    --el-tabs-header-height: 56px;
    margin-bottom: -8px;
    border: none;

    & .el-tabs__nav {
      border: 0;
    }

    & .el-tabs__item {
      border-left: 0;

      &.is-active {
        background-color: #ffffff;
        //border-top: 1px solid var(--el-border-color);
        //border-left: 1px solid var(--el-border-color);
        //border-right: 1px solid var(--el-border-color);
        border-top-right-radius: 8px;
        border-top-left-radius: 8px;
      }
    }
  }

  :deep(.el-tabs__content) {
    background-color: #ffffff;
    padding: 20px;
    //border: 1px solid var(--el-border-color);
    //border-top: none;
    border-radius: 8px;
  }
}
</style>
