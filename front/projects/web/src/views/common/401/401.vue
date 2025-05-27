<script lang="ts" setup="">
const { t } = useI18n()
const router = useRouter()
</script>

<template>
  <el-result icon="warning" title="401">
    <template #sub-title>
      <p>{{ t("notify.noPermission") }}</p>
    </template>
    <template #extra>
      <el-button type="primary" @click="router.back()">{{ t("btn.back") }}</el-button>
    </template>
  </el-result>
</template>
