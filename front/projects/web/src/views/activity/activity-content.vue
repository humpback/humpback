<script lang="ts" setup>
import { TableHeight } from "@/utils"
import { QueryActivityInfo } from "./common.ts"
import { ActivityInfo } from "@/types"

const props = defineProps<{
  activityType: string
  labelName: "config" | "user" | "group" | "service" | "node" | "registry" | "team"
}>()

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const tableHeight = computed(() => TableHeight(346))

const isLoading = ref(false)
const queryInfo = ref<QueryActivityInfo>(new QueryActivityInfo(route.query, props.activityType, userStore.isAdmin))

const tableList = ref({
  total: 0,
  data: [] as Array<ActivityInfo>
})

const setRowClass = ({ row }) => {
  return !row.oldContent && !row.newContent ? "hide-expand-icon" : ""
}

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  return await activityService
    .query(queryInfo.value.searchParams())
    .then(data => {
      tableList.value.total = data.total
      tableList.value.data = data.list
    })
    .finally(() => (isLoading.value = false))
}

onMounted(() => {
  search()
})
</script>

<template>
  <div class="d-flex gap-3 flex-wrap">
    <v-date-time-range
      v-model:end-time="queryInfo.filter.endAt"
      v-model:start-time="queryInfo.filter.startAt"
      :end-placeholder="t('label.endTime')"
      :out-label="t('label.timeRange')"
      :range-separator="t('label.to')"
      :start-placeholder="t('label.startTime')"
      format="YYYY-MM-DD HH:mm"
      out-label-width="120px"
      time-format="HH:mm" />
    <div v-if="props.activityType === PageActivity.Services || props.activityType === PageActivity.Groups" style="width: 280px">
      <v-group-select v-model="queryInfo.filter.groupId" :placeholder="t('placeholder.all')" show-out-label />
    </div>
    <div style="width: 280px">
      <v-activity-action-select v-model="queryInfo.filter.action" :mode="props.labelName" show-out-label />
    </div>
    <div v-if="userStore.isAdmin" style="width: 280px">
      <v-users-select
        v-model="queryInfo.filter.operator"
        :multiple="false"
        :out-label="t('label.operator')"
        :placeholder="t('placeholder.all')"
        clearable
        show-out-label />
    </div>
    <el-button type="primary" @click="search()">
      <el-icon :size="16">
        <IconMdiFilterVariant />
      </el-icon>
      {{ t("btn.filter") }}
    </el-button>
  </div>

  <v-table
    v-loading="isLoading"
    v-model:page-info="queryInfo.pageInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :row-class-name="setRowClass"
    :total="tableList.total"
    class="mt-5"
    row-key="activityId"
    @page-change="search()">
    <el-table-column align="left" class-name="expand-column" type="expand" width="24">
      <template #default="scope">
        <div class="px-5">
          <v-monaco-diff-editor
            v-if="scope.row.newContent || scope.row.oldContent"
            :new-data="scope.row.newContent ? JSON.stringify(scope.row.newContent, null, 4) : ''"
            :old-data="scope.row.oldContent ? JSON.stringify(scope.row.oldContent, null, 4) : ''"
            language="json" />
        </div>
      </template>
    </el-table-column>

    <el-table-column :label="t(`label.${props.labelName}`)" min-width="140">
      <template #default="scope">
        <span v-if="props.activityType === PageActivity.Groups">
          <v-router-link v-if="scope.row.resourceId" :href="`/ws/group/${scope.row.resourceId}/services`" :text="scope.row.resourceName" />
          <v-table-column-none v-else :text="scope.row.resourceName" />
        </span>
        <span v-else-if="props.activityType === PageActivity.Services">
          <v-router-link
            v-if="scope.row.resourceId"
            :href="`/ws/group/${(scope.row.resourceId as string).slice(0, 8)}/service/${scope.row.resourceId}/${PageServiceDetail.BasicInfo}`"
            :text="scope.row.resourceName" />
          <v-table-column-none v-else :text="scope.row.resourceName" />
        </span>
        <v-table-column-none v-else :text="scope.row.resourceName" />
      </template>
    </el-table-column>
    <el-table-column v-if="props.activityType === PageActivity.Services" :label="t('label.group')" min-width="140">
      <template #default="scope">
        <v-router-link v-if="scope.row.resourceId" :href="`/ws/group/${(scope.row.resourceId as string).slice(0, 8)}/services`" :text="scope.row.groupName" />
        <v-table-column-none v-else :text="scope.row.groupName" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.description')" min-width="200">
      <template #default="scope">
        <span v-if="props.labelName && scope.row.action">
          <span v-if="props.activityType === PageActivity.Services && scope.row.instanceName">
            {{ t(`activity.${props.labelName}.${scope.row.action}Instance`, { name: scope.row.instanceName }) }}
          </span>
          <span v-else>{{ t(`activity.${props.labelName}.${scope.row.action}`) }} </span>
        </span>
        <span v-else>--</span>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.operator')" min-width="160">
      <template #default="scope">
        <v-table-column-none :text="scope.row.operator" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.date')" width="160">
      <template #default="scope">
        <v-date-view :timestamp="scope.row.operateAt" />
      </template>
    </el-table-column>
  </v-table>
</template>

<style lang="scss" scoped>
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
