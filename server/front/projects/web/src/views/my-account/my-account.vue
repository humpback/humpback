<script lang="ts" setup>
import ChangePassword from "./change-password.vue"
import { RuleFormatErrEmailOption, RuleFormatErrPhone, RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { FormInstance, FormRules } from "element-plus"
import { cloneDeep } from "lodash-es"
import { ActivityInfo } from "@/types"

const { t } = useI18n()
const userStore = useUserStore()

const loading = ref(false)
const isAction = ref(false)
const isLoadingActivities = ref(false)
const userInfo = ref<UserInfo>(NewUserEmptyInfo())

const queryUserActivityInfo = ref<any>({
  filter: {
    "type": "account"
  },
  pageInfo: {
    size: 10,
    index: 1
  }
})
const activities = ref({
  total: 0,
  data: [] as Array<ActivityInfo>
})

const tableRef = useTemplateRef<FormInstance>("tableRef")

const rules = ref<FormRules>({
  username: [
    { required: true, validator: RulePleaseEnter("label.username"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.Username.Min, RuleLength.Username.Max), trigger: "blur" }
  ],
  email: [
    { validator: RuleLimitMax(RuleLength.Email.Max), trigger: "blur" },
    { validator: RuleFormatErrEmailOption(), trigger: "blur" }
  ],
  phone: [
    { validator: RuleLimitMax(RuleLength.Email.Max), trigger: "blur" },
    { validator: RuleFormatErrPhone(), trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

const setRowClass = ({ row }) => {
  return !row.oldContent && !row.newContent ? "hide-expand-icon" : ""
}

async function getUserInfo() {
  loading.value = true
  return await userService
    .getMe()
    .then(data => {
      userInfo.value = data
      userStore.setUserInfo(cloneDeep(data))
    })
    .finally(() => {
      loading.value = false
    })
}

async function getUserActivities() {
  isLoadingActivities.value = true
  return await activityService
    .query(queryUserActivityInfo.value)
    .then(data => {
      activities.value.total = data.total
      activities.value.data = data.list
    })
    .finally(() => (isLoadingActivities.value = false))
}

async function save() {
  if (!(await tableRef.value?.validate())) {
    return
  }
  isAction.value = true
  await userService
    .updateMeInfo({
      username: userInfo.value.username,
      email: userInfo.value.email,
      phone: userInfo.value.phone,
      description: userInfo.value.description
    })
    .finally(() => (isAction.value = false))
  ShowSuccessMsg(t("message.saveSuccess"))
  await Promise.all([getUserInfo(), getUserActivities()])
}

onMounted(async () => {
  await Promise.all([getUserInfo(), getUserActivities()])
})
</script>

<template>
  <div>
    <v-card v-loading="loading">
      <v-page-title :title="t('label.myAccount')" />
      <div>
        <v-role-admin :role="userInfo.role" size="default" />
        <div class="d-flex gap-1 mt-2 pl-1 mb-3">
          <el-text type="info">
            {{ t("label.createDate") }}:
            <v-date-view :format="7" :timestamp="userInfo.createdAt" />
          </el-text>
          <el-divider direction="vertical" />
          <div>
            <change-password />
          </div>
        </div>
      </div>

      <div class="mb-5 mt-2">
        <v-alert> {{ t("tips.usernameChangeTips") }}</v-alert>
      </div>
      <el-form ref="tableRef" :model="userInfo" :rules="rules" label-position="top" label-width="auto">
        <el-row :gutter="12">
          <el-col>
            <el-form-item :label="t('label.username')" prop="username">
              <v-username-input v-model="userInfo.username" />
            </el-form-item>
          </el-col>
          <el-col>
            <el-form-item :label="t('label.description')" prop="description">
              <v-description-input v-model="userInfo.description" />
            </el-form-item>
          </el-col>
          <el-col :md="12" :span="24">
            <el-form-item :label="t('label.email')" prop="email">
              <v-email-input v-model="userInfo.email" />
            </el-form-item>
          </el-col>
          <el-col :md="12" :span="24">
            <el-form-item :label="t('label.phone')" prop="phone">
              <v-phone-input v-model="userInfo.phone" />
            </el-form-item>
          </el-col>
          <el-col>
            <el-form-item>
              <div class="text-align-right w-100">
                <el-button :loading="isAction" type="primary" @click="save()">{{ t("btn.save") }}</el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </v-card>

    <v-card class="mt-5">
      <div class="d-flex gap-2">
        <div class="f-bold">
          {{ t("label.activities") }}
        </div>
        <el-button :disabled="isLoadingActivities" plain size="small" type="success" @click="getUserActivities()"> {{ t("btn.refresh") }}</el-button>
        <v-loading v-if="isLoadingActivities" />
      </div>
      <div class="activity-content">
        <v-table
          v-loading="isLoadingActivities"
          v-model:page-info="queryUserActivityInfo.pageInfo"
          :data="activities.data"
          :row-class-name="setRowClass"
          :total="activities.total"
          hide-header-bg-color
          @page-change="getUserActivities()">
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
          <el-table-column :label="t('label.description')" min-width="200px">
            <template #default="scope">
              <span>{{ scope.row.action ? t(`activity.account.${scope.row.action}`) : "--" }}</span>
            </template>
          </el-table-column>
          <el-table-column width="160px">
            <template #default="scope">
              <v-date-view :timestamp="scope.row.operateAt" />
            </template>
          </el-table-column>
        </v-table>
      </div>
    </v-card>
  </div>
</template>

<style lang="scss" scoped>
.activity-content {
  margin-top: 40px;
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
