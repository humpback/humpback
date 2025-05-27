<script lang="ts" setup>
import { TabPaneName } from "element-plus"
import TeamPage from "./team/teams.vue"
import UserPage from "./user/users.vue"
import { PageUserRelated } from "@/models"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const activeTab = ref<TabPaneName | undefined>(route.params?.mode as string)

async function changeTab(name: TabPaneName) {
  await router.replace({ params: { mode: name } })
  activeTab.value = name
}

const options = reactive<{ name: string; label: string; component: any }[]>([
  { name: PageUserRelated.Users, label: "header.users", component: shallowRef(UserPage) },
  { name: PageUserRelated.Teams, label: "header.teams", component: shallowRef(TeamPage) }
])
</script>

<template>
  <v-card class="tab-card">
    <v-page-title :title="t('label.userRelated')" />
    <el-tabs :model-value="activeTab" class="tab-box" @update:modelValue="changeTab">
      <template v-for="item in options" :key="item.name">
        <el-tab-pane :label="t(item.label)" :name="item.name">
          <template #label>
            <strong>{{ t(item.label) }}</strong>
          </template>
          <component :is="item.component" v-if="item.name === activeTab" />
        </el-tab-pane>
      </template>
    </el-tabs>
  </v-card>
</template>

<style lang="scss" scoped>
.tab-box {
  :deep(.el-tabs__header.is-top) {
    margin-bottom: 0;
  }

  :deep(.el-tabs__content) {
    padding-top: 20px;
  }
}
</style>
