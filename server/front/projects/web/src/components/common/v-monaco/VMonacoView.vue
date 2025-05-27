<script lang="ts" setup>
import screenfull from "screenfull"
import * as monaco from "monaco-editor"
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker"
import JsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker"
import CssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker"
import HtmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker"
import TSWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker"
import { CopyToClipboard } from "@/utils/index.ts"

const props = withDefaults(defineProps<{ language?: string }>(), {
  language: "text"
})

const { t } = useI18n()

self.MonacoEnvironment = {
  getWorker(_, label) {
    switch (label) {
      case "json":
        return new JsonWorker()
      case "css":
      case "scss":
      case "less":
        return new CssWorker()
      case "html":
      case "handlebars":
      case "razor":
        return new HtmlWorker()
      case "typescript":
      case "javascript":
        return new TSWorker()
      default:
        return new EditorWorker()
    }
  }
}

monaco.languages.typescript.typescriptDefaults.setEagerModelSync(true)

let monacoCtl: monaco.editor.IStandaloneCodeEditor
const monacoRef = useTemplateRef("monacoRef")

const boxRef = useTemplateRef<any>("boxRef")

function fullScreen() {
  if (screenfull.isEnabled) {
    screenfull.toggle(boxRef.value)
  }
}

function setValue(v: string) {
  monacoCtl?.setValue(v)
}

function getValue() {
  return monacoCtl?.getValue() || ""
}

function copy() {
  const data = getValue()
  if (data) {
    CopyToClipboard(data)
  }
}

onMounted(() => {
  monacoCtl = monaco.editor.create(monacoRef.value!, {
    value: "",
    language: props.language,
    theme: "vs-dark", //官方自带三种主题vs, hc-black, hc-light, or vs-dark
    renderValidationDecorations: "off",
    cursorSmoothCaretAnimation: "off",
    cursorStyle: "line",
    cursorBlinking: "solid",
    hideCursorInOverviewRuler: true,
    roundedSelection: false,
    readOnly: true, // 只读
    domReadOnly: true,
    glyphMargin: true, //字形边缘
    fontSize: 14, //字体大小
    renderLineHighlight: "none",
    accessibilitySupport: "off",
    selectionHighlight: false,
    renderLineHighlightOnlyWhenFocus: false,
    minimap: {
      enabled: false,
      side: "right"
    },
    scrollBeyondLastLine: false,
    overviewRulerBorder: false,
    folding: false,
    foldingHighlight: false,
    disableLayerHinting: true,
    selectionClipboard: false,
    codeLens: false,
    colorDecorators: true,
    lineNumbers: "on",
    lineNumbersMinChars: 3
  })

  monaco.editor.defineTheme("customTheme", {
    base: "vs", // 基于黑暗主题
    inherit: true, // 继承基础主题
    rules: [],
    colors: {
      // "editor.lineHighlightBackground": "#e7e5e5" // 光标行背景色
    }
  })

  // 应用自定义主题
  monaco.editor.setTheme("customTheme")

  monacoCtl.onKeyDown(e => {
    if (e.code === "Backspace" || e.code === "Delete") {
      e.preventDefault() // 阻止默认行为
      e.stopPropagation() // 阻止事件冒泡
    }
  })
})

onBeforeUnmount(() => {
  monacoCtl?.dispose()
})

defineExpose({ getValue, setValue })
</script>

<template>
  <div ref="boxRef" class="monaco-box">
    <div class="monaco-header">
      <div class="flex-1">
        <slot name="title" />
      </div>
      <v-tooltip :content="t('btn.copy')" :teleported="false" placement="top-end">
        <el-button link @click="copy()">
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
    <div ref="monacoRef" class="monaco-content" />
  </div>
</template>

<style lang="scss" scoped>
.monaco-box {
  --hp-monaco-radius: 4px;
  width: 100%;
  height: 100%;
  border-radius: var(--hp-monaco-radius);
  padding: 8px 8px;
  box-sizing: border-box;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;

  .monaco-header {
    min-height: 40px;
    display: flex;
    align-items: center;
    justify-content: left;
    gap: 8px;
  }

  .monaco-content {
    flex: 1;
    width: 100%;
    //border: 1px dashed #adadad;
    border-radius: var(--hp-monaco-radius);

    :deep(.monaco-editor) {
      border-radius: var(--hp-monaco-radius);
      //--vscode-editor-background: #acb1b0;

      .cursor {
        display: none !important;
      }

      .current-line,
      .margin-view-overlays .current-line-margin {
        background: none !important;
      }

      .overflow-guard {
        border-radius: var(--hp-monaco-radius);
      }
    }
  }
}
</style>
