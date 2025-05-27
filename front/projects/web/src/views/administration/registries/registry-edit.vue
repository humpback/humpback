<script lang="ts" setup>
import { NewRegistryEmptyInfo, RegistryInfo } from "@/types"
import { cloneDeep } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { RSAEncrypt } from "utils/rsa.ts"
import { isDefaultRegistry } from "@/views/administration/registries/common.ts"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const registryStore = useRegistryStore()

const isLoading = ref(false)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as RegistryInfo
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  url: [
    { required: true, validator: RulePleaseEnter("label.url"), trigger: "blur" },
    { required: true, validator: RuleLimitMax(RuleLength.RegistryUrl.Max), trigger: "blur" }
  ],
  username: [{ validator: RuleLimitMax(RuleLength.RegistryUsername.Max), trigger: "blur" }],
  password: [{ validator: RuleLimitMax(RuleLength.RegistryPassword.Max), trigger: "blur" }]
})

function clearTail() {
  dialogInfo.value.info.url = dialogInfo.value.info.url.replace(/\/+/g, "/").replace(/^\/|\/$/g, "")
}

function open(info?: RegistryInfo) {
  dialogInfo.value.info = info ? cloneDeep(info) : NewRegistryEmptyInfo()
  dialogInfo.value.show = true
  if (info) {
    isLoading.value = true
    registryService
      .info(info.registryId, true)
      .then(data => {
        dialogInfo.value.info = data
      })
      .catch(() => (dialogInfo.value.show = false))
      .finally(() => (isLoading.value = false))
  }
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  const body: any = {
    url: dialogInfo.value.info.url,
    isDefault: dialogInfo.value.info.isDefault,
    username: dialogInfo.value.info.username ? RSAEncrypt(dialogInfo.value.info.username) : "",
    password: dialogInfo.value.info.password ? RSAEncrypt(dialogInfo.value.info.password) : ""
  }
  isAction.value = true
  if (dialogInfo.value.info.registryId) {
    body.registryId = dialogInfo.value.info.registryId
    registryService
      .update(body)
      .then(() => {
        registryStore.refreshRegistries()
        ShowSuccessMsg(t("message.saveSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  } else {
    registryService
      .create(body)
      .then(() => {
        registryStore.refreshRegistries()
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
    <template #header>{{ dialogInfo.info.registryId ? t("header.editRegistry") : t("header.addRegistry") }}</template>
    <v-alert>{{ t("tips.registryAuthTips") }}</v-alert>
    <div v-loading="isLoading" class="my-3">
      <el-form ref="formRef" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-form-item :label="t('label.url')" prop="url">
          <div class="d-flex gap-3 w-100">
            <v-input
              v-model="dialogInfo.info.url"
              :disabled="isDefaultRegistry(dialogInfo.info.url)"
              :maxlength="RuleLength.RegistryUrl.Max"
              clearable
              show-word-limit
              @blur="clearTail()" />
            <el-checkbox v-model="dialogInfo.info.isDefault" :label="t('label.isDefault')" border class="default-class" />
          </div>
        </el-form-item>
        <el-form-item :label="t('label.username')" prop="username">
          <v-input
            v-model="dialogInfo.info.username"
            :disabled="isDefaultRegistry(dialogInfo.info.url)"
            :maxlength="RuleLength.RegistryUsername.Max"
            clearable
            show-word-limit />
        </el-form-item>
        <el-form-item :label="t('label.password')" prop="password">
          <v-password-input
            v-model="dialogInfo.info.password"
            :disabled="isDefaultRegistry(dialogInfo.info.url)"
            :maxlength="RuleLength.RegistryPassword.Max"
            :minlength="0" />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="isLoading" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.edit {
  width: 100%;
  height: 500px;
}

.default-class {
  --el-color-primary: var(--el-color-warning);
  --el-checkbox-input-border-color-hover: var(--el-color-warning);
  --el-checkbox-checked-text-color: var(--el-color-warning);
  --el-checkbox-checked-bg-color: var(--el-color-warning);
  --el-checkbox-checked-input-border-color: var(--el-color-warning);
}
</style>
