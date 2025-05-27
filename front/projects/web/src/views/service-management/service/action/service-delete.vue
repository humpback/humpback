<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { ServiceInfo } from "@/types"
import { PageGroupDetail } from "@/models"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const router = useRouter()

const isAction = ref(false)
const isChecked = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as ServiceInfo,
  isJumpToList: false
})

function open(info: ServiceInfo, isJumpToList?: boolean) {
  isChecked.value = false
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.show = true
  dialogInfo.value.isJumpToList = !!isJumpToList
}

async function confirmDelete() {
  if (!isChecked.value) {
    return
  }

  isAction.value = true
  await serviceService.delete(dialogInfo.value.info.groupId, dialogInfo.value.info.serviceId).finally(() => (isAction.value = false))
  ShowSuccessMsg(t("message.deleteSuccess"))

  if (dialogInfo.value.isJumpToList) {
    await router.push({ name: "groupDetail", params: { groupId: dialogInfo.value.info.groupId, mode: PageGroupDetail.Services } })
  } else {
    dialogInfo.value.show = false
    emits("refresh")
  }
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ t("header.deleteService") }}</template>
    <v-alert v-if="dialogInfo.info.containers.length > 0" type="warning">{{ t("tips.deleteServiceTips") }}</v-alert>
    <div class="my-3 f-bold">{{ t("notify.delete") }}</div>
    <v-delete-input-continue v-model="isChecked" :keywords="dialogInfo.info.serviceName" class="mt-5" />
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="!isChecked" :loading="isAction" type="danger" @click="confirmDelete">{{ t("btn.delete") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
