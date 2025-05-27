<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { PageServiceDetail, RuleLength } from "@/models"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const groupId = ref(route.params.groupId as string)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {
    serviceName: "",
    description: ""
  }
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  serviceName: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.ServiceName.Min, RuleLength.ServiceName.Max), trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

async function open() {
  dialogInfo.value.info = {
    serviceName: "",
    description: ""
  }
  dialogInfo.value.show = true
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }

  isAction.value = true
  serviceService
    .create(groupId.value, {
      serviceName: dialogInfo.value.info.serviceName,
      description: dialogInfo.value.info.description
    })
    .then(serviceId => {
      ShowSuccessMsg(t("message.addSuccess"))
      dialogInfo.value.show = false
      router.push({
        name: "serviceInfo",
        params: { groupId: groupId.value, serviceId: serviceId, mode: PageServiceDetail.BasicInfo }
      })
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="800px">
    <template #header>{{ t("header.addService") }}</template>
    <div class="my-3">
      <el-form ref="formRef" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
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
      <el-button :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
