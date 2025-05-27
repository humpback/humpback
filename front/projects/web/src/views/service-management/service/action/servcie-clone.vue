<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { PageServiceDetail, RuleLength } from "@/models"
import { ServiceInfo } from "@/types"
import { find } from "lodash-es"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const groupId = ref(route.params.groupId as string)
const isLoading = ref(false)
const isAction = ref(false)
const groups = ref<GroupInfo[]>([])
const dialogInfo = ref({
  show: false,
  info: {
    serviceId: "",
    serviceName: "",
    description: "",
    newGroupId: ""
  }
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  serviceName: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.ServiceName.Min, RuleLength.ServiceName.Max), trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }],
  newGroupId: [{ required: true, validator: RulePleaseEnter("label.targetGroup"), trigger: "change" }]
})

async function getGroups() {
  return await groupService.list().then(list => (groups.value = list))
}

async function open(info: ServiceInfo) {
  dialogInfo.value.info = {
    serviceId: info.serviceId,
    serviceName: "",
    description: "",
    newGroupId: info.groupId
  }
  dialogInfo.value.show = true
  isLoading.value = true
  await getGroups()
    .catch(() => (dialogInfo.value.show = false))
    .finally(() => (isLoading.value = false))
  if (!find(groups.value, x => x.groupId === dialogInfo.value.info.newGroupId)) {
    dialogInfo.value.info.newGroupId = ""
  }
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }

  isAction.value = true
  serviceService
    .clone(groupId.value, {
      serviceId: dialogInfo.value.info.serviceId,
      serviceName: dialogInfo.value.info.serviceName,
      description: dialogInfo.value.info.description,
      newGroupId: dialogInfo.value.info.newGroupId
    })
    .then(serviceId => {
      ShowSuccessMsg(t("message.cloneSuccess"))
      dialogInfo.value.show = false
      router
        .push({
          name: "serviceInfo",
          params: { groupId: dialogInfo.value.info.newGroupId, serviceId: serviceId, mode: PageServiceDetail.BasicInfo }
        })
        .then(() => {
          emits("refresh")
        })
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="800px">
    <template #header>{{ t("header.cloneService") }}</template>
    <div class="my-3">
      <el-form ref="formRef" v-loading="isLoading" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-form-item :label="t('label.targetGroup')" prop="newGroupId">
          <v-group-select v-model="dialogInfo.info.newGroupId" :options="groups" />
        </el-form-item>
        <v-alert v-if="dialogInfo.info.newGroupId !== groupId">{{ t("tips.cloneServiceTips") }}</v-alert>
        <el-form-item :label="t('label.name')" prop="serviceName">
          <v-input v-model="dialogInfo.info.serviceName" :maxlength="RuleLength.ServiceName.Max" clearable show-word-limit />
        </el-form-item>
        <el-form-item :label="t('label.description')" prop="description">
          <v-description-input v-model="dialogInfo.info.description" />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="isLoading" :loading="isAction" type="primary" @click="save">{{ t("btn.clone") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
