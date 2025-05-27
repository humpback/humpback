<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { NewTeamEmptyInfo, TeamInfo, UserInfo } from "@/types"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isLoading = ref(false)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as TeamInfo
})

const userOptions = ref<UserInfo[]>([])

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  name: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.TeamName.Min, RuleLength.TeamName.Max), trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

function open(info?: TeamInfo) {
  dialogInfo.value.info = info ? cloneDeep(info) : NewTeamEmptyInfo()
  dialogInfo.value.show = true
  isLoading.value = true
  getUsers()
    .catch(() => {
      dialogInfo.value.show = false
    })
    .finally(() => (isLoading.value = false))
}

async function getUsers() {
  return await userService.list().then(data => {
    userOptions.value = data
    return data
  })
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }

  const body: any = {
    name: dialogInfo.value.info.name,
    description: dialogInfo.value.info.description,
    users: dialogInfo.value.info.users
  }
  isAction.value = true
  if (dialogInfo.value.info.teamId) {
    body.teamId = dialogInfo.value.info.teamId
    teamService
      .update(body)
      .then(() => {
        ShowSuccessMsg(t("message.saveSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  } else {
    teamService
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
    <template #header>{{ dialogInfo.info.teamId ? t("header.editTeam") : t("header.addTeam") }}</template>
    <div class="my-3">
      <el-form ref="formRef" v-loading="isLoading" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-form-item :label="t('label.teamName')" prop="name">
          <v-input v-model="dialogInfo.info.name" :maxlength="RuleLength.TeamName.Max" clearable show-word-limit />
        </el-form-item>
        <el-form-item :label="t('label.description')" prop="description">
          <v-description-input v-model="dialogInfo.info.description" />
        </el-form-item>
        <el-form-item :label="t('label.users')" prop="teams">
          <v-users-select v-model="dialogInfo.info.users" :options="userOptions" show-footer />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="isLoading" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
