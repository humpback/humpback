<script lang="ts" setup>
import type { FormInstance, FormRules } from "element-plus"
import { StorageEventType, RulePleaseEnter, storageEventBus } from "@/utils"
import { RSAEncrypt } from "utils/rsa.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loginStep = ref(0)

const formRef = useTemplateRef<FormInstance>("formRef")
const formData = reactive({
  username: "",
  password: ""
})

const formRules = reactive<FormRules>({
  username: [{ required: true, validator: RulePleaseEnter("placeholder.username"), trigger: "blur" }],
  password: [{ required: true, validator: RulePleaseEnter("placeholder.password"), trigger: "blur" }]
})

async function login() {
  if (loginStep.value !== 0 || !(await formRef.value?.validate())) {
    return
  }
  const body = {
    username: RSAEncrypt(formData.username),
    password: RSAEncrypt(formData.password)
  }
  loginStep.value = 1
  userService
    .login(body)
    .then(data => {
      loginStep.value = 2
      userStore.setUserInfo(data)
      storageEventBus.SendMessage(StorageEventType.Login, data)
      ShowSuccessMsg(t("message.loginSuccess"))
      if (route.query?.redirectUrl) {
        router.push(route.query.redirectUrl as string)
        return
      }
      router.push({ name: "workspace" })
    })
    .catch(() => {
      loginStep.value = 0
    })
}
</script>

<template>
  <div class="pub-page">
    <el-card class="v-content">
      <template #header>
        <div class="my-3">
          <v-logo />
        </div>
      </template>
      <el-form ref="formRef" :model="formData" :rules="formRules" @submit.prevent="login()">
        <el-form-item prop="username">
          <v-username-input v-model="formData.username" :clearable="false" :placeholder="t('placeholder.username')" :show-word-limit="false" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <v-password-input v-model="formData.password" :placeholder="t('placeholder.password')" size="large" />
        </el-form-item>
        <el-form-item>
          <el-button v-if="loginStep === 0" class="w-100 mt-3" native-type="submit" size="large" type="primary">
            <el-icon class="mr-1">
              <icon-mdi-login-variant />
            </el-icon>
            {{ t("btn.login") }}
          </el-button>
          <el-button v-if="loginStep === 1" class="w-100 mt-3" disabled size="large" type="primary">
            {{ t("btn.loggingIn") }}
          </el-button>
          <el-button v-if="loginStep === 2" class="w-100 mt-3" disabled size="large" type="primary">
            {{ t("btn.redirecting") }}
          </el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="text-align-right">
          <Language />
        </div>
      </template>
    </el-card>
    <div class="copy-right">
      <el-text size="small" type="info">{{ t("tips.allRightsReserved") }}</el-text>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.pub-page {
  box-sizing: border-box;
  min-height: calc(100vh);
  width: 100%;
  padding-top: 20vh;
  //background-image: linear-gradient(-20deg, #fffeff 50%, #d7fffe 100%);
  //background-color: #dcd9d4;
  background-image:
    linear-gradient(to bottom, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.5) 100%),
    radial-gradient(at 50% 0%, rgba(255, 255, 255, 0.1) 0%, rgba(0, 0, 0, 0.5) 50%);
  background-blend-mode: soft-light, screen;

  :deep(.el-card__footer) {
    border: none;
  }
}

:deep(.v-content) {
  box-sizing: border-box;
  box-shadow: unset;
  border: none;
  margin: 0 auto 48px auto;
  padding: 24px;
  width: 400px;
  border-radius: 16px;
  //background-color: #ffffff;
  background-image: linear-gradient(to top, #c4c5c7 0%, #dcdddf 52%, #ebebeb 100%);

  .el-card__header {
    border: none;
    padding: 0 0 30px 0;
  }

  .el-card__body {
    padding: 0;
  }

  .el-card__footer {
    padding: 30px 0 0 0;
  }
}

.copy-right {
  margin: 0 auto;
  text-align: center;
  font-size: 12px;
  font-family: Poppins, "Noto Sans SC", "Microsoft YaHei", "Ping Fang SC", Helvetica, Arial, sans-serif;
}
</style>
