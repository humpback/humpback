<script lang="ts" setup>
const props = withDefaults(defineProps<{ memo?: string; onlyIcon?: boolean; iconSize?: number }>(), { iconSize: 18 })

const { t } = useI18n()
</script>

<template>
  <v-tooltip v-if="props?.onlyIcon && props?.memo" :content="t(`err.${props?.memo}`)" effect="dark">
    <el-button link type="danger">
      <el-icon :size="props.iconSize">
        <IconMdiWarningCircleOutline />
      </el-icon>
    </el-button>
  </v-tooltip>
  <div v-else-if="props?.memo" :title="t(`err.${props?.memo}`)" class="d-flex gap-1 overflow_div">
    <el-icon :size="props.iconSize" color="var(--el-color-danger)">
      <IconMdiWarningCircleOutline />
    </el-icon>
    <el-text truncated type="danger">
      {{ t(`err.${props?.memo}`) }}
    </el-text>
  </div>
</template>

<style lang="scss" scoped></style>
