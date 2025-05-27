<script lang="ts" setup>
import { TeamInfo, UserInfo } from "@/types"

const { t } = useI18n()

const dialogInfo = ref({
  show: false,
  info: {} as TeamInfo
})

const isLoading = ref(false)
const list = ref<Array<UserInfo>>([])

function open(info: TeamInfo) {
  list.value = []
  dialogInfo.value.info = info
  dialogInfo.value.show = true
  isLoading.value = true
  userService
    .queryByTeamId(info.teamId)
    .then(data => {
      list.value = data
    })
    .finally(() => (isLoading.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="800px">
    <template #header> {{ t("header.users") }}</template>
    <v-table v-loading="isLoading" :data="list" :max-height="600" :show-header="true" :total="list.length">
      <el-table-column :label="t('label.username')" min-width="140" prop="username" />
      <el-table-column :label="t('label.description')" min-width="140" prop="description">
        <template #default="scope">
          <v-table-column-none :text="scope.row.description" />
        </template>
      </el-table-column>
      <el-table-column :label="t('label.email')" min-width="140" prop="email">
        <template #default="scope">
          <v-table-column-none :text="scope.row.email" />
        </template>
      </el-table-column>
      <el-table-column :label="t('label.phone')" min-width="120" prop="phone">
        <template #default="scope">
          <v-table-column-none :text="scope.row.phone" />
        </template>
      </el-table-column>
      <el-table-column :label="t('label.role')" min-width="160" prop="groupList">
        <template #default="scope">
          <v-role-view :role="scope.row.role" />
        </template>
      </el-table-column>
    </v-table>
    <div class="mt-5 d-flex gap-2">
      <div class="flex-1 pl-2">
        {{ dialogInfo.info.name }}
      </div>
      <el-text v-if="list.length > 0">{{ t("label.totalNumber", { total: list.length }) }}</el-text>
    </div>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
