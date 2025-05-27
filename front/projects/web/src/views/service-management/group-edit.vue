<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { GroupInfo, NewGroupEmptyInfo } from "@/types"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isLoading = ref(false)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as GroupInfo
})

const userOptions = ref<UserInfo[]>([])
const teamOptions = ref<TeamInfo[]>([])
const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  groupName: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.GroupName.Min, RuleLength.GroupName.Max), trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

async function getUsers() {
  return await userService.list().then(data => {
    userOptions.value = data
    return data
  })
}

async function getTeams() {
  return await teamService.list().then(data => {
    teamOptions.value = data
    return data
  })
}

async function open(info?: GroupInfo) {
  dialogInfo.value.info = info ? cloneDeep(info) : NewGroupEmptyInfo()
  dialogInfo.value.show = true
  isLoading.value = true
  await Promise.all([getUsers(), getTeams()])
    .catch(() => {
      dialogInfo.value.show = false
    })
    .finally(() => (isLoading.value = false))
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }

  const body: any = {
    groupName: dialogInfo.value.info.groupName,
    description: dialogInfo.value.info.description,
    users: dialogInfo.value.info.users,
    teams: dialogInfo.value.info.teams
  }
  isAction.value = true
  if (dialogInfo.value.info.groupId) {
    body.groupId = dialogInfo.value.info.groupId
    groupService
      .update(body)
      .then(() => {
        ShowSuccessMsg(t("message.saveSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .catch(err => {
        if (err?.response?.data?.code === "R4Group-006") {
          dialogInfo.value.show = false
          emits("refresh")
        }
      })
      .finally(() => (isAction.value = false))
  } else {
    groupService
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
    <template #header>{{ dialogInfo.info.groupId ? t("header.editGroup") : t("header.addGroup") }}</template>
    <div class="my-3">
      <el-form ref="formRef" v-loading="isLoading" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-form-item :label="t('label.name')" prop="groupName">
          <v-input v-model="dialogInfo.info.groupName" :maxlength="RuleLength.GroupName.Max" clearable show-word-limit />
        </el-form-item>
        <el-form-item :label="t('label.description')" prop="description">
          <v-description-input v-model="dialogInfo.info.description" />
        </el-form-item>
        <el-form-item :label="t('label.users')" prop="groups">
          <v-users-select v-model="dialogInfo.info.users" :options="userOptions" show-footer />
        </el-form-item>
        <el-col>
          <el-form-item :label="t('label.teams')" prop="teams">
            <v-teams-select v-model="dialogInfo.info.teams" :options="teamOptions" />
          </el-form-item>
        </el-col>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="isLoading" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
