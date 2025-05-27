<script lang="ts" setup>
import { map } from "lodash-es"

const props = withDefaults(defineProps<{ labels?: { [key: string]: string }; line?: number }>(), { line: 4 })
const { t } = useI18n()

const labelMapping = computed(() => {
  if (props.labels) {
    const keys = Object.keys(props.labels).sort()
    return map(keys, key => `${key}:${props.labels![key]}`)
  }
  return []
})
</script>

<template>
  <div class="custom-column">
    <div style="width: 100%">
      <div v-for="(item, index) in labelMapping.slice(0, props.line)" :key="index" class="line">
        <el-text size="small">-- {{ item }}</el-text>
      </div>
      <div v-if="labelMapping.length > props.line">
        <el-popover :width="300" placement="bottom-start" trigger="hover">
          <template #reference>
            <el-button link size="small" type="primary"> {{ t("btn.more") }}</el-button>
          </template>
          <div v-for="(item, index) in labelMapping.slice(props.line)" :key="index" class="line">
            <el-text size="small">-- {{ item }}</el-text>
          </div>
        </el-popover>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.custom-column {
  min-height: 60px;
  display: flex;
  align-items: center;
  justify-content: start;
}

.line {
  font-size: 12px;
  line-height: 18px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
