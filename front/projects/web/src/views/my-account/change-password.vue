<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"
import { disposeStore } from "@/stores"

const { t } = useI18n()
const router = useRouter()

const show = ref(false)
const loading = ref(false)
const psdInfo = ref<any>()
const tableRef = useTemplateRef<FormInstance>("tableRef")

const rules = ref<FormRules>({
  oldPassword: [
    { required: true, validator: RulePleaseEnter("label.oldPassword"), trigger: "blur" },
    { validator: RuleLimitRange(RuleLength.Password.Min, RuleLength.Password.Max), trigger: "blur" }
  ],
  newPassword: [
    { required: true, validator: RulePleaseEnter("label.newPassword"), trigger: "blur" },
    { validator: RuleLimitRange(RuleLength.Password.Min, RuleLength.Password.Max), trigger: "blur" }
  ],
  confirmTheNewPassword: [
    { required: true, validator: RulePleaseEnter("label.confirmTheNewPassword"), trigger: "blur" },
    { validator: checkConfirmTheNewPassword, trigger: "blur" }
  ]
})

function checkConfirmTheNewPassword(rule: any, value: any, callback: any) {
  if (psdInfo.value.newPassword && value !== psdInfo.value.newPassword) {
    return callback(new Error(t("rules.inconsistentNewPassword")))
  }
  callback()
}

function showDialog() {
  psdInfo.value = {
    oldPassword: "",
    newPassword: "",
    confirmTheNewPassword: ""
  }
  show.value = true
}

async function save() {
  if (!(await tableRef.value?.validate())) {
    return
  }
  loading.value = true
  await userService
    .changePassword({
      oldPassword: RSAEncrypt(psdInfo.value.oldPassword),
      newPassword: RSAEncrypt(psdInfo.value.newPassword)
    })
    .finally(() => {
      loading.value = false
    })
  ShowSuccessMsg(t("message.changePasswordSuccess"))
  show.value = false
  disposeStore()
  await router.push({ name: "login" })
}
</script>

<template>
  <el-button class="focus-outline-none" link type="primary" @click="showDialog()">
    {{ t("btn.changePassword") }}
  </el-button>
  <v-dialog v-model="show" :title="t('header.changePassword')" width="600px">
    <v-alert> {{ t("tips.changePasswordTips") }}</v-alert>
    <el-form ref="tableRef" :model="psdInfo" :rules="rules" class="mt-5" label-position="top" label-width="auto">
      <el-form-item :label="t('label.oldPassword')" prop="oldPassword">
        <v-password-input v-model="psdInfo.oldPassword" />
      </el-form-item>
      <el-form-item :label="t('label.newPassword')" prop="newPassword">
        <v-password-input v-model="psdInfo.newPassword" />
      </el-form-item>
      <el-form-item :label="t('label.confirmTheNewPassword')" prop="confirmTheNewPassword">
        <v-password-input v-model="psdInfo.confirmTheNewPassword" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :loading="loading" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
