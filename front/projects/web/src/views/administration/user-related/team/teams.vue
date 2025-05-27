<script lang="ts" setup>
import { TeamInfo } from "@/types"
import { TableHeight } from "@/utils"
import { Action } from "@/models"
import TeamEdit from "./team-edit.vue"
import TeamDelete from "./team-delete.vue"
import TeamViewUsers from "./team-view-users.vue"
import { QueryTeamInfo } from "@/views/administration/user-related/team/common.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const tableHeight = computed(() => TableHeight(344))

const isLoading = ref(false)
const queryInfo = ref<QueryTeamInfo>(new QueryTeamInfo(route.query))

const tableList = ref({
  total: 0,
  data: [] as Array<TeamInfo>
})

const teamEditRef = useTemplateRef<InstanceType<typeof TeamEdit>>("teamEditRef")
const teamDeleteRef = useTemplateRef<InstanceType<typeof TeamDelete>>("teamDeleteRef")
const teamViewUsersRef = useTemplateRef<InstanceType<typeof TeamViewUsers>>("teamViewUsersRef")

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  return await teamService
    .query(queryInfo.value.searchParams())
    .then(res => {
      tableList.value.data = res.list
      tableList.value.total = res.total
    })
    .finally(() => (isLoading.value = false))
}

function openAction(action: string, info?: TeamInfo) {
  switch (action) {
    case Action.Add:
    case Action.Edit:
      teamEditRef.value?.open(info)
      break
    case Action.Delete:
      teamDeleteRef.value?.open(info!)
      break
    case Action.View:
      teamViewUsersRef.value?.open(info!)
      break
  }
}

onMounted(() => search())
</script>

<template>
  <v-search v-model="queryInfo.keywords" :add-label="t('btn.addTeam')" :input-label="t('label.name')" @add="openAction(Action.Add)" @search="search" />

  <v-table
    v-loading="isLoading"
    v-model:page-info="queryInfo.pageInfo"
    v-model:sort-info="queryInfo.sortInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :total="tableList.total"
    @page-change="search"
    @sort-change="search">
    <el-table-column :label="t('label.teamName')" fixed="left" min-width="140" prop="name" sortable="custom" />
    <el-table-column :label="t('label.description')" min-width="140" prop="description">
      <template #default="scope">
        <v-table-column-none :text="scope.row.description" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.users')" min-width="120">
      <template #default="scope">
        <el-button v-if="Array.isArray(scope.row.users) && scope.row.users.length > 0" link type="primary" @click="openAction(Action.View, scope.row)">
          {{ t("label.totalUsers", { total: scope.row.users.length }) }}
        </el-button>
        <span v-else>--</span>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.updateDate')" min-width="160" prop="updatedAt" sortable="custom">
      <template #default="scope">
        <v-date-view :timestamp="scope.row.updatedAt" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.createDate')" min-width="160" prop="createdAt" sortable="custom">
      <template #default="scope">
        <v-date-view :timestamp="scope.row.createdAt" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.action')" align="right" fixed="right" header-align="center" width="130">
      <template #default="scope">
        <el-button link type="primary" @click="openAction(Action.Edit, scope.row)">{{ t("btn.edit") }}</el-button>
        <el-button link type="danger" @click="openAction(Action.Delete, scope.row)">{{ t("btn.delete") }}</el-button>
      </template>
    </el-table-column>
  </v-table>

  <team-delete ref="teamDeleteRef" @refresh="search()" />

  <team-edit ref="teamEditRef" @refresh="search()" />

  <team-view-users ref="teamViewUsersRef" />
</template>

<style lang="scss" scoped></style>
