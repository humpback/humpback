<script lang="ts" setup>
import { TableHeight } from "@/utils"
import { ActivityInfo } from "@/types"
import VLoading from "@/components/business/v-loading/VLoading.vue"
import { refreshData } from "@/views/service-management/service/common.ts"

const { t } = useI18n()
const route = useRoute()

const tableHeight = computed(() => TableHeight(362))

const groupId = ref(route.params.groupId as string)
const serviceId = ref(route.params.serviceId as string)

const isLoading = ref(false)
const isLoadingActivity = ref(false)
const queryInfo = ref<any>({
  filter: {
    "type": "services",
    "serviceId": route.params.serviceId as string
  },
  pageInfo: {
    size: 20,
    index: 1
  }
})

const tableList = ref({
  total: 0,
  data: [] as Array<ActivityInfo>
})

const setRowClass = ({ row }) => {
  return !row.oldContent && !row.newContent ? "hide-expand-icon" : ""
}

async function search() {
  isLoading.value = true
  await refreshData(groupId.value, serviceId.value, "activity").finally(() => (isLoading.value = false))
}

async function getActivities() {
  isLoadingActivity.value = true
  return await activityService
    .query(queryInfo.value)
    .then(data => {
      tableList.value.total = data.total
      tableList.value.data = data.list
    })
    .finally(() => (isLoadingActivity.value = false))
}

onBeforeMount(async () => {
  await search()
  await getActivities()
})
</script>

<template>
  <div class="d-flex gap-3">
    <div class="header-icon">
      <el-icon :size="18">
        <IconMdiAccountFileText />
      </el-icon>
    </div>
    <div class="d-flex gap-1">
      <el-text class="f-bold" size="large">{{ t("label.activities") }}</el-text>
    </div>
    <el-button plain size="small" type="success" @click="getActivities()">{{ t("btn.refresh") }}</el-button>
    <v-loading v-if="isLoading" />
  </div>

  <v-table
    v-loading="isLoadingActivity"
    v-model:page-info="queryInfo.pageInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :row-class-name="setRowClass"
    :total="tableList.total"
    class="mt-5"
    row-key="activityId"
    @page-change="getActivities()">
    <el-table-column align="left" class-name="expand-column" type="expand" width="24">
      <template #default="scope">
        <div class="px-5">
          <v-monaco-diff-editor
            v-if="scope.row.oldContent || scope.row.newContent"
            :new-data="scope.row.newContent ? JSON.stringify(scope.row.newContent, null, 4) : ''"
            :old-data="scope.row.oldContent ? JSON.stringify(scope.row.oldContent, null, 4) : ''"
            language="json" />
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.description')" min-width="200">
      <template #default="scope">
        <span v-if="scope.row.action">
          {{ t(`activity.service.${scope.row.action}`, { name: scope.row.instanceName }) }}
        </span>
        <span v-else>--</span>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.operator')" min-width="140" prop="operator" />
    <el-table-column :label="t('label.date')" width="160">
      <template #default="scope">
        <v-date-view :timestamp="scope.row.operateAt" />
      </template>
    </el-table-column>
  </v-table>
</template>

<style lang="scss" scoped>
.header-icon {
  color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
  border-radius: 50%;
  padding: 8px;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.expand-column) {
  .cell {
    padding: 0 4px 0 8px;
  }
}

:deep(.hide-expand-icon) {
  .expand-column .cell {
    padding-top: 4px;
    display: none;
  }
}
</style>
