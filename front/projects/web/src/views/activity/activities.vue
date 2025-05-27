<script lang="ts" setup>
import { PageActivity } from "@/models"
import { TabPaneName } from "element-plus"
import PageActivitiesContent from "./activity-content.vue"
import { find } from "lodash-es"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeTab = ref<TabPaneName | undefined>(route.params?.mode as string)

async function changeTab(name: TabPaneName) {
  await router.replace({ params: { mode: name } })
  activeTab.value = name
}

const options = reactive<
  Array<{
    name: string
    labelName: "config" | "user" | "group" | "service" | "node" | "registry" | "team"
    label: string
    limitAdmin?: boolean
  }>
>([
  { name: PageActivity.Groups, labelName: "group", label: "header.groups" },
  { name: PageActivity.Services, labelName: "service", label: "header.services" },
  { name: PageActivity.Configs, labelName: "config", label: "header.configs" },
  { name: PageActivity.Registries, labelName: "registry", label: "header.registries", limitAdmin: true },
  { name: PageActivity.Nodes, labelName: "node", label: "header.nodes", limitAdmin: true },
  { name: PageActivity.Users, labelName: "user", label: "header.users", limitAdmin: true },
  { name: PageActivity.Teams, labelName: "team", label: "header.teams", limitAdmin: true }
])

onMounted(() => {
  const item = find(options, x => x.name === (route.params.mode as string))
  if (!item || (item.limitAdmin && !userStore.isAdmin)) {
    changeTab(PageActivity.Groups)
  }
})
</script>

<template>
  <v-card class="tab-box">
    <v-page-title :title="t('label.activities')" />
    <el-tabs :model-value="activeTab" class="tab-box" @update:modelValue="changeTab">
      <template v-for="item in options" :key="item.name">
        <el-tab-pane v-if="!item.limitAdmin || userStore.isAdmin" :label="t(item.label)" :name="item.name">
          <template #label>
            <strong>{{ t(item.label) }}</strong>
          </template>
          <page-activities-content v-if="item.name === activeTab" :activity-type="item.name" :label-name="item.labelName" />
        </el-tab-pane>
      </template>
    </el-tabs>
  </v-card>
</template>

<style lang="scss" scoped>
.tab-box {
  min-width: 800px;

  :deep(.el-tabs__header.is-top) {
    margin-bottom: 0;
  }

  :deep(.el-tabs__content) {
    padding-top: 20px;
  }
}
</style>
