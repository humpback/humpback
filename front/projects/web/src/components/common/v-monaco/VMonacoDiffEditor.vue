<script lang="ts" setup>
import * as monaco from "monaco-editor"
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker"
import JsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker"

self.MonacoEnvironment = {
  getWorker(_, label) {
    switch (label) {
      case "json":
        return new JsonWorker()
      default:
        return new EditorWorker()
    }
  }
}

const props = defineProps<{ oldData?: string; newData?: string; language?: string }>()

const { t } = useI18n()

let diffRef = useTemplateRef<HTMLDivElement>("diffRef")
let diffEditor: monaco.editor.IStandaloneDiffEditor

const resize = () => {
  if (diffEditor) {
    diffEditor.layout()
  }
}

onMounted(() => {
  diffEditor = monaco.editor.createDiffEditor(diffRef.value!, {
    enableSplitViewResizing: true,
    readOnly: true,
    theme: "vs-dark",
    lineNumbers: "on",
    minimap: {
      maxColumn: 200,
      side: "right",
      size: "proportional",
      renderCharacters: true
    }
  })
  diffEditor.setModel({
    original: monaco.editor.createModel("", props.language ? props.language : "text"),
    modified: monaco.editor.createModel("", props.language ? props.language : "text")
  })
  diffEditor.getOriginalEditor().setValue(props.oldData || "")
  diffEditor.getModifiedEditor().setValue(props.newData || "")
  window.addEventListener("resize", resize)
})

onBeforeUnmount(() => {
  diffEditor?.dispose()
  window.removeEventListener("resize", resize)
})
</script>

<template>
  <div class="w-100">
    <el-row class="mb-1">
      <el-col :span="12">
        <slot name="oldTitle">
          <el-text>{{ t("label.old") }}</el-text>
        </slot>
      </el-col>
      <el-col :span="12">
        <slot name="newTitle">
          <el-text>{{ t("label.new") }}</el-text>
        </slot>
      </el-col>
    </el-row>
    <div :id="new Date().getTime().toString() + '_monaco'" ref="diffRef" class="hp-monaco-diff"></div>
  </div>
</template>

<style scoped>
:deep(.hp-monaco-diff) {
  width: 100%;
  height: 400px;

  .diffOverview {
    background-color: #1e1e1e;
  }
}
</style>
