<script lang="ts" setup>
import { ActivityInfo } from "@/types"
import { map, orderBy } from "lodash-es"

const props = defineProps<{ data?: { [key: string]: ActivityInfo[] }; isLoading?: boolean }>()
const { t } = useI18n()

const tableList = computed(() => {
  const list: Array<any> = []
  const activityMap = props.data || {}
  Object.keys(activityMap).forEach((key: string) => {
    const tempData = map(activityMap[key] || [], item => {
      return {
        ...item,
        labelName: key
      }
    })
    list.push(...tempData)
  })
  return orderBy(list, ["operateAt"], ["desc"])
})
</script>

<template>
  <v-card>
    <div class="title">{{ t("header.recentActivities") }}</div>
    <div>
      <v-table v-if="tableList.length > 0" v-loading="props.isLoading" :data="tableList" hide-header-bg-color max-height="360px" minHeight="360px">
        <el-table-column :label="t('label.description')">
          <template #default="scope">
            <span v-if="scope.row.labelName && scope.row.action">
              <span v-if="scope.row.labelName === 'service' && scope.row.instanceName">
                {{ t(`activity.${scope.row.labelName}.${scope.row.action}Instance`, { name: scope.row.instanceName }) }}
              </span>
              <span v-else>{{ t(`activity.${scope.row.labelName}.${scope.row.action}`) }} </span>
            </span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column width="160px">
          <template #default="scope">
            <v-router-link
              v-if="scope.row.labelName === 'service'"
              :href="`/ws/group/${(scope.row.resourceId as string).slice(0, 8)}/service/${scope.row.resourceId}/${PageServiceDetail.BasicInfo}`"
              :text="scope.row.resourceName" />
            <v-router-link v-else-if="scope.row.labelName === 'group'" :href="`/ws/group/${scope.row.resourceId}/services`" :text="scope.row.resourceName" />
            <v-table-column-none v-else :text="scope.row.resourceName" />
          </template>
        </el-table-column>
        <el-table-column :label="t('label.date')" width="160px">
          <template #default="scope">
            <v-date-view :timestamp="scope.row.operateAt" />
          </template>
        </el-table-column>
      </v-table>
      <div v-else v-loading="props.isLoading" class="empty-content">
        <el-empty>
          <template #image>
            <el-icon :size="160">
              <IconMdiAccountFileText />
            </el-icon>
          </template>
          <template #description>
            <span>{{ t("tips.noRecentActivities") }}</span>
          </template>
        </el-empty>
      </div>
    </div>
  </v-card>
</template>

<style lang="scss" scoped>
.title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  padding-left: 4px;
}

.empty-content {
  height: 360px;
  color: gray;
  font-size: 14px;
}
</style>
