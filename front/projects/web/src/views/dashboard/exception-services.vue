<script lang="ts" setup>
import { ResourceExceptionServiceInfo } from "@/types"
import { PageServiceDetail } from "@/models"

const props = defineProps<{ data?: ResourceExceptionServiceInfo[]; enabledServices?: number; isLoading?: boolean }>()
const { t } = useI18n()

const tableList = computed(() => props.data || [])
</script>

<template>
  <v-card>
    <div class="title">{{ t("header.exceptionService") }}</div>
    <div>
      <v-table v-if="tableList.length > 0" v-loading="props.isLoading" :data="tableList" hide-header-bg-color maxHeight="360px" min-height="360px">
        <el-table-column :label="t('label.service')">
          <template #default="scope">
            <v-router-link
              :href="`/ws/group/${scope.row.groupId}/service/${scope.row.serviceId}/${PageServiceDetail.BasicInfo}`"
              :text="scope.row.serviceName" />
          </template>
        </el-table-column>
        <el-table-column :label="t('label.group')">
          <template #default="scope">
            <v-router-link :href="`/ws/group/${scope.row.groupId}/services`" :text="scope.row.groupName" />
          </template>
        </el-table-column>
        <el-table-column :label="t('label.status')" width="120px">
          <template #default="scope">
            <v-service-status-tag :is-enabled="scope.row.isEnabled" :status="scope.row.status" />
            <v-memo v-if="!!scope.row.memo" :icon-size="18" :memo="scope.row.memo" only-icon />
          </template>
        </el-table-column>
      </v-table>
      <div v-else v-loading="props.isLoading" :class="['empty-content', props.enabledServices && 'no-abnormal']">
        <el-empty>
          <template #image>
            <el-icon :size="160">
              <IconMdiCalendarCheckOutline />
            </el-icon>
          </template>
          <template #description>
            <span>{{ props.enabledServices ? t("tips.noExceptionServices") : t("tips.noEnabledServices") }}</span>
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
  font-size: 14px;
  color: gray;
}

.no-abnormal {
  color: #13b913;

  :deep(.el-icon) svg {
    color: #d8f2d8;
  }
}
</style>
