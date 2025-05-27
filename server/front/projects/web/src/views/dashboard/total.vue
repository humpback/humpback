<script lang="ts" setup>
import { DashboardResourceStatisticsInfo } from "@/types"

const props = defineProps<{ info: DashboardResourceStatisticsInfo; isLoading?: boolean }>()
const { t } = useI18n()

const totalOptions = computed(() => [
  { i18nLabel: "label.groups", key: "groups", value: props.info.groups },
  { i18nLabel: "label.services", key: "services", value: props.info.services },
  { i18nLabel: "label.nodes", key: "nodes", value: props.info.nodes },
  { i18nLabel: "label.users", key: "users", value: props.info.users }
])
</script>

<template>
  <el-row :gutter="20">
    <el-col v-for="(item, index) in totalOptions" :key="index" :span="6">
      <v-card>
        <div class="total-title">
          <div>{{ t(item.i18nLabel) }}</div>
        </div>
        <div class="total-content">
          <span v-if="!props.isLoading">{{ item.value }}</span>
          <v-loading v-else />
        </div>
      </v-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
.total-title {
  font-weight: bold;
  font-size: 20px;
}

.total-content {
  margin-top: 20px;
  font-size: 20px;
}
</style>
