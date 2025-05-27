<script lang="ts" setup>
import { CopyToClipboard } from "@/utils"

const props = defineProps<{ command: string }>()

const { t } = useI18n()

function copy() {
  CopyToClipboard(props.command)
}
</script>

<template>
  <div class="command-body">
    <v-tooltip :content="t('btn.copy')" :teleported="false" placement="left">
      <el-button class="copy-btn" link plain type="primary" @click="copy()">
        <el-icon :size="16">
          <IconMdiContentCopy />
        </el-icon>
      </el-button>
    </v-tooltip>
    <pre class="command-content">
<code>
{{ props.command }}</code>
    </pre>
  </div>
</template>

<style lang="scss" scoped>
.command-body {
  position: relative;
  border-radius: 8px;
  background-color: #282828;
  padding: 0 16px;

  .copy-btn {
    position: absolute;
    right: 8px;
    top: 8px;
  }

  .command-content {
    margin: 0;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    color: #ffffff;
    word-break: break-word;
  }
}
</style>
