<script lang="ts" setup>
const props = defineProps<{ title?: string; showBreadcrumbs?: boolean; breadcrumbParams?: any }>()
</script>

<template>
  <div class="page-title">
    <span v-if="props.title">{{ props.title }} </span>
    <el-button v-else link loading />
  </div>
  <div v-if="props.showBreadcrumbs" class="mt-4 mb-5">
    <v-breadcrumb :params="props.breadcrumbParams" />
  </div>
</template>

<style lang="scss" scoped>
.page-title {
  font-size: 24px;
  font-weight: 600;
  margin: 8px 0 20px 0;
}
</style>
