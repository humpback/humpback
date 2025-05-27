<script lang="ts" setup></script>

<template>
  <v-alert :show-icon="false" class="tips-line">
    <slot name="default" />
  </v-alert>
</template>

<style lang="scss" scoped>
.tips-line {
  --el-color-info: #305d8c;
  background-color: #d6dde7;
  border-left: 4px solid #9cb4c5;
}
</style>
