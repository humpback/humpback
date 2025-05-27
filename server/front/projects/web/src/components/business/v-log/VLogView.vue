<script lang="ts" setup>
import screenfull from "screenfull"
import { CopyToClipboard } from "@/utils"

const props = defineProps<{ logList?: string[] }>()

const { t } = useI18n()

const wrapLines = ref(false)
const logViewRef = useTemplateRef<any>("logViewRef")

const logsText = computed(() => (props.logList || []).join(""))

function fullScreen() {
  if (screenfull.isEnabled) {
    screenfull.toggle(logViewRef.value)
  }
}

function copy() {
  if (logsText.value) {
    CopyToClipboard(logsText.value)
  }
}
</script>

<template>
  <div ref="logViewRef" class="log-container">
    <div class="log-header">
      <div class="flex-1 d-flex gap-3">
        <slot name="header-left" />
        <el-switch v-model="wrapLines" :active-text="t('label.wrapLines')" />
      </div>
      <div class="d-flex gap-3">
        <v-tooltip :content="t('btn.copy')" :teleported="false" placement="top-end">
          <el-button :disabled="!logsText" link @click="copy()">
            <el-icon :size="20">
              <IconMdiContentCopy />
            </el-icon>
          </el-button>
        </v-tooltip>
        <v-tooltip :content="t('btn.fullScreen')" :teleported="false" placement="top-end">
          <el-button link style="margin-left: 0" @click="fullScreen">
            <el-icon :size="20">
              <IconMdiArrowExpandAll />
            </el-icon>
          </el-button>
        </v-tooltip>
      </div>
    </div>

    <div :class="['log-content', wrapLines && 'wrap-lines']" class="log-content">
      <pre><code>{{ logsText }}</code></pre>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.log-container {
  background-color: #f5f5f5;
  padding: 8px 4px 4px 4px;
  border-radius: 8px;
  box-sizing: border-box;
  height: 100%;
  display: flex;
  flex-direction: column;

  .log-header {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 12px;
    border-radius: 8px;
    margin-bottom: 4px;
  }

  .log-content {
    flex: 1;
    box-sizing: border-box;
    border: 1px solid #ddd;
    background: black;
    color: #ffffff;
    padding: 4px 16px;
    border-radius: 8px;
    overflow: auto;
    line-height: 18px;

    &.wrap-lines {
      pre code {
        white-space: pre-wrap;
        word-break: break-word;
      }
    }
  }
}
</style>
