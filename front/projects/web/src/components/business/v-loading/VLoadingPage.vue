<script lang="ts" setup>
const { t } = useI18n()
</script>

<template>
  <div class="loading-page">
    <div>
      <div class="d-flex gap-3">
        <v-logo direction="horizontal" />
        <span class="neon-text">Humpback.com</span>
      </div>
      <div class="loading-text">
        <el-text>{{ t("tips.loadingHumpback") }}</el-text>
        <v-loading :size="14" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.loading-page {
  height: 100vh;

  width: 100%;

  display: flex;
  align-items: center;
  justify-content: center;
}

.neon-text {
  font-size: 20px;
  font-weight: bold;
  color: var(--el-color-primary);
}

.loading-text {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}
</style>
