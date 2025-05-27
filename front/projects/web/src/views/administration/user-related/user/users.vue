<script lang="ts" setup>
import { UserInfo } from "@/types"
import { TableHeight } from "@/utils"
import UserDelete from "./user-delete.vue"
import UserEdit from "./user-edit.vue"
import UserViewTeams from "./user-view-teams.vue"
import { QueryUserInfo } from "./common.ts"

import { Action } from "@/models"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const tableHeight = computed(() => TableHeight(344))

const isLoading = ref(false)
const queryInfo = ref<QueryUserInfo>(new QueryUserInfo(route.query))

const tableList = ref({
  total: 0,
  data: [] as Array<UserInfo>
})

const userEditRef = useTemplateRef<InstanceType<typeof UserEdit>>("userEditRef")
const userDeleteRef = useTemplateRef<InstanceType<typeof UserDelete>>("userDeleteRef")
const userViewTeamsRef = useTemplateRef<InstanceType<typeof UserViewTeams>>("userViewTeamsRef")

function showActionBtn(info: UserInfo) {
  if (info.userId === userStore.userInfo.userId) {
    return false
  }
  if (IsSuperAdmin(info.role)) {
    return false
  }
  return !IsAdmin(info.role) || userStore.isSuperAdmin
}

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  return await userService
    .query(queryInfo.value.searchParams())
    .then(res => {
      tableList.value.data = res.list
      tableList.value.total = res.total
    })
    .finally(() => (isLoading.value = false))
}

function openAction(action: string, info?: UserInfo) {
  switch (action) {
    case Action.Add:
    case Action.Edit:
      userEditRef.value?.open(info)
      break
    case Action.Delete:
      userDeleteRef.value?.open(info!)
      break
    case Action.View:
      userViewTeamsRef.value?.open(info!)
  }
}

onMounted(() => search())
</script>

<template>
  <v-search
    v-model="queryInfo.keywords"
    :add-label="t('btn.addUser')"
    :placeholder="t('placeholder.enterUsernameEmailOrPhone')"
    @add="openAction(Action.Add)"
    @search="search">
    <template #prefix>
      <div style="width: 280px">
        <v-role-query-select v-model="queryInfo.filter.role" :placeholder="t('placeholder.all')" @change="search()" />
      </div>
    </template>
  </v-search>

  <v-table
    v-loading="isLoading"
    v-model:page-info="queryInfo.pageInfo"
    v-model:sort-info="queryInfo.sortInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :total="tableList.total"
    @page-change="search"
    @sort-change="search">
    <el-table-column :label="t('label.username')" fixed="left" min-width="200" prop="username" sortable="custom">
      <template #default="scope">
        <div class="d-flex gap-1">
          <div class="overflow_div flex-1">
            <span>{{ scope.row.username }}</span>
          </div>
          <el-tag v-if="scope.row.userId === userStore.userInfo.userId" class="mr-1" effect="dark" round size="small" type="warning">
            {{ t("label.currentUser") }}
          </el-tag>
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.role')" min-width="160" prop="groupList">
      <template #default="scope">
        <v-role-view :role="scope.row.role" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.description')" min-width="160" prop="description">
      <template #default="scope">
        <v-table-column-none :text="scope.row.description" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.email')" min-width="120" prop="email">
      <template #default="scope">
        <v-table-column-none :text="scope.row.email" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.phone')" min-width="120" prop="phone">
      <template #default="scope">
        <v-table-column-none :text="scope.row.phone" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.teams')" min-width="100">
      <template #default="scope">
        <el-button v-if="Array.isArray(scope.row.teams) && scope.row.teams.length > 0" link type="primary" @click="openAction(Action.View, scope.row)">
          {{ t("label.totalTeams", { total: scope.row.teams.length }) }}
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
        <div v-if="showActionBtn(scope.row)">
          <el-button link type="primary" @click="openAction(Action.Edit, scope.row)">{{ t("btn.edit") }}</el-button>
          <el-button link type="danger" @click="openAction(Action.Delete, scope.row)">{{ t("btn.delete") }}</el-button>
        </div>
      </template>
    </el-table-column>
  </v-table>

  <user-delete ref="userDeleteRef" @refresh="search()" />

  <user-edit ref="userEditRef" @refresh="search()" />

  <user-view-teams ref="userViewTeamsRef" />
</template>

<style lang="scss" scoped></style>
