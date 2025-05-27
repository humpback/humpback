<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    enabled?: boolean
    effect?: "light" | "dark"
    placement?: "auto" | "auto-start" | "auto-end" | "top" | "bottom" | "right" | "left" | "top-start" | "top-end" | "bottom-start" | "bottom-end"
  }>(),
  { effect: "dark", placement: "top-start" }
)
const { t } = useI18n()
</script>

<template>
  <v-tooltip :content="props.enabled ? t('label.enabled') : t('label.disabled')" :effect="props.effect" :placement="props.placement">
    <div v-if="props.enabled" class="enabled box"></div>
    <div v-else class="disabled box"></div>
  </v-tooltip>
</template>

<style lang="scss" scoped>
.box {
  width: 18px;
  height: 18px;
  border-radius: 50%;
}

.enabled {
  background-color: var(--el-color-success);
}

.disabled {
  background-color: var(--el-color-info);
}
</style>
