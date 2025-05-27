<script lang="ts" setup>
import { cloneDeep, find } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RuleFormatErrEmailOption, RuleFormatErrPhone, RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { TeamInfo, UserInfo } from "@/types"
import { RSAEncrypt } from "utils/rsa.ts"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const userStore = useUserStore()

const isLoading = ref(false)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as UserInfo
})

const teamsOption = ref<TeamInfo[]>([])

const formRef = useTemplateRef<FormInstance>("formRef")
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
    { validator: RuleLimitMax(RuleLength.Phone.Max), trigger: "blur" },
    { validator: RuleFormatErrPhone(), trigger: "blur" }
  ],
  password: [
    { required: true, validator: RulePleaseEnter("label.password"), trigger: "blur" },
    { validator: RuleLimitRange(RuleLength.Password.Min, RuleLength.Password.Max), trigger: "blur" }
  ],
  role: [{ required: true, validator: checkRole, trigger: "change" }],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

function checkRole(rule: any, value: any, callback: any) {
  const options = userStore.isSuperAdmin ? [UserRole.User, UserRole.Admin] : [UserRole.User]
  if (!find(options, x => x === value)) {
    return callback(new Error(t("rules.invalidRole")))
  }
  callback()
}

async function open(info?: UserInfo) {
  dialogInfo.value.info = info ? cloneDeep(info) : NewUserEmptyInfo()
  dialogInfo.value.show = true
  isLoading.value = true
  await Promise.all([dialogInfo.value.info.userId ? getUserInfo() : undefined, getTeams()])
    .catch(() => {
      dialogInfo.value.show = false
    })
    .finally(() => (isLoading.value = false))
}

async function getUserInfo() {
  return await userService.info(dialogInfo.value.info.userId, true).then(info => {
    dialogInfo.value.info = info
  })
}

async function getTeams() {
  return await teamService.list().then(data => {
    teamsOption.value = data
    return data
  })
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  const body: any = {
    username: dialogInfo.value.info.username,
    email: dialogInfo.value.info.email,
    password: RSAEncrypt(dialogInfo.value.info.password),
    description: dialogInfo.value.info.description,
    phone: dialogInfo.value.info.phone,
    role: dialogInfo.value.info.role,
    teams: dialogInfo.value.info.teams
  }
  isAction.value = true
  if (dialogInfo.value.info.userId) {
    body.userId = dialogInfo.value.info.userId
    userService
      .update(body)
      .then(() => {
        ShowSuccessMsg(t("message.saveSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  } else {
    userService
      .create(body)
      .then(() => {
        ShowSuccessMsg(t("message.addSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  }
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="800px">
    <template #header>
      {{ dialogInfo.info.userId ? t("header.editUser") : t("header.addUser") }}
    </template>
    <div class="my-3">
      <el-form ref="formRef" v-loading="isLoading" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-row :gutter="12">
          <el-col>
            <el-form-item :label="t('label.username')" prop="username">
              <v-username-input v-model="dialogInfo.info.username" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('label.role')" prop="role">
              <v-role-select
                v-model="dialogInfo.info.role"
                :clearable="false"
                :only-show-roles="userStore.isSuperAdmin ? [UserRole.User, UserRole.Admin] : [UserRole.User]" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('label.password')" prop="password">
              <v-password-input v-model="dialogInfo.info.password" />
            </el-form-item>
          </el-col>
          <el-col>
            <el-form-item :label="t('label.description')" prop="description">
              <v-description-input v-model="dialogInfo.info.description" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('label.email')" prop="email">
              <v-email-input v-model="dialogInfo.info.email" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('label.phone')" prop="phone">
              <v-phone-input v-model="dialogInfo.info.phone" />
            </el-form-item>
          </el-col>

          <el-col>
            <el-form-item :label="t('label.teams')" prop="teams">
              <v-teams-select v-model="dialogInfo.info.teams" :options="teamsOption" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="isLoading" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
