<script lang="ts" setup>
import screenfull from "screenfull"
import * as monaco from "monaco-editor"
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker"
import JsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker"
import HtmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker"
import TSWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker"
import YamlWorker from "./yaml.worker.js?worker"
import { configureMonacoYaml } from "monaco-yaml"
import xmlFormat from "xml-formatter"
import * as ini from "ini"

configureMonacoYaml(monaco, {
  enableSchemaRequest: true,
  schemas: [
    {
      fileMatch: ["**/.prettierrc.*"],
      uri: "https://json.schemastore.org/prettierrc.json"
    }
  ]
})

import { CopyToClipboard } from "@/utils/index.ts"

const props = withDefaults(defineProps<{ readOnly?: boolean; language?: string; maxSize?: number }>(), {
  readOnly: false,
  language: "text",
  maxSize: RuleLength.ConfigValue.Max
})

const editContent = defineModel<string>()

const { t } = useI18n()

self.MonacoEnvironment = {
  getWorker(_, label) {
    switch (label) {
      case "json":
        return new JsonWorker()
      case "html":
      case "handlebars":
      case "razor":
        return new HtmlWorker()
      case "typescript":
      case "javascript":
        return new TSWorker()
      case "yaml":
        return new YamlWorker()
      default:
        return new EditorWorker()
    }
  }
}

monaco.languages.registerDocumentFormattingEditProvider("xml", {
  provideDocumentFormattingEdits(model) {
    const content = model.getValue()
    try {
      const formatted = xmlFormat(content, {
        indentation: "  ",
        lineSeparator: "\n"
      })
      return [
        {
          range: model.getFullModelRange(),
          text: formatted
        }
      ]
    } catch {}
  }
})

monaco.languages.registerDocumentFormattingEditProvider("ini", {
  provideDocumentFormattingEdits(model) {
    const content = model.getValue()
    const formatted = ini.stringify(ini.parse(content), { whitespace: true })
    return [
      {
        range: model.getFullModelRange(),
        text: formatted
      }
    ]
  }
})

monaco.languages.typescript.typescriptDefaults.setEagerModelSync(true)

let monacoCtl: monaco.editor.IStandaloneCodeEditor
let listeningContent: monaco.IDisposable
const monacoRef = useTemplateRef("monacoRef")

const boxRef = useTemplateRef<any>("boxRef")

function fullScreen() {
  if (screenfull.isEnabled) {
    screenfull.toggle(boxRef.value)
  }
}

function formatLang(lang: string) {
  if (monacoCtl) {
    monaco.editor.setModelLanguage(monacoCtl.getModel()!, lang)
    setTimeout(() => {
      monacoCtl.getAction("editor.action.formatDocument")!.run()
    }, 100)
  }
}

function getValue() {
  return monacoCtl?.getValue() || ""
}

function copy() {
  CopyToClipboard(getValue())
}

onMounted(() => {
  monacoCtl = monaco.editor.create(monacoRef.value!, {
    value: editContent.value || "",
    language: props.language,
    theme: "vs-dark", //官方自带三种主题vs, hc-black, hc-light, or vs-dark
    selectOnLineNumbers: true, //显示行号
    roundedSelection: false,
    readOnly: props.readOnly, // 只读
    cursorStyle: "line", //光标样式
    automaticLayout: true, //自动布局
    glyphMargin: true, //字形边缘
    useTabStops: false,
    fontSize: 15, //字体大小
    quickSuggestionsDelay: 10, //代码提示延时
    minimap: {
      enabled: true,
      side: "right"
    },
    fontFamily: "var(--hp-font-family)",
    scrollBeyondLastLine: false,
    overviewRulerBorder: false,
    formatOnPaste: true,
    tabSize: 4,
    folding: true,
    foldingHighlight: true,
    foldingStrategy: "indentation",
    showFoldingControls: "always",
    disableLayerHinting: true,
    emptySelectionClipboard: false,
    selectionClipboard: false,
    codeLens: false,
    colorDecorators: true,
    accessibilitySupport: "off",
    lineNumbers: "on",
    lineNumbersMinChars: 1
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
  if (!props.readOnly) {
    listeningContent = monacoCtl.onDidChangeModelContent(() => {
      const currentContent = monacoCtl.getValue()
      if (currentContent.length > props.maxSize) {
        monacoCtl.setValue(editContent.value || "")
      } else {
        editContent.value = currentContent
      }
    })
  }
})

onBeforeUnmount(() => {
  if (listeningContent) {
    listeningContent?.dispose()
  }
  if (monacoCtl) {
    monacoCtl?.dispose()
  }
})

defineExpose({ getValue })
</script>

<template>
  <div ref="boxRef" class="monaco-box">
    <div class="monaco-header">
      <div class="flex-1">
        <slot name="title" />
      </div>
      <el-dropdown :teleported="false" @command="formatLang">
        <el-button link>
          <div class="d-flex">
            {{ t("btn.format") }}
            <el-icon :size="20">
              <IconMdiChevronDown />
            </el-icon>
          </div>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="text">text</el-dropdown-item>
            <el-dropdown-item command="ini">ini</el-dropdown-item>
            <el-dropdown-item command="yaml">yaml</el-dropdown-item>
            <el-dropdown-item command="xml">xml</el-dropdown-item>
            <el-dropdown-item command="javascript">javascript</el-dropdown-item>
            <el-dropdown-item command="json">json</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

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
    <div class="footer-total">
      <el-text>{{ `${editContent?.length || 0} / ${props.maxSize}` }}</el-text>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.monaco-box {
  --hp-monaco-radius: 4px;
  width: 100%;
  height: 100%;
  border-radius: var(--hp-monaco-radius);
  padding: 4px 8px;
  box-sizing: border-box;
  background-color: #f5f5f5;

  .monaco-header {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: left;
    gap: 8px;
  }

  .monaco-content {
    width: 100%;
    height: calc(100% - 72px);
    border: 1px dashed #adadad;
    border-radius: var(--hp-monaco-radius);

    :deep(.monaco-editor) {
      border-radius: var(--hp-monaco-radius);
      //--vscode-editor-background: #acb1b0;
      .overflow-guard {
        border-radius: var(--hp-monaco-radius);

        .margin {
          .margin-view-overlays {
            //background-color: #3b3b3b;
          }
        }
      }
    }
  }

  .footer-total {
    display: flex;
    align-items: center;
    justify-content: right;
    padding-right: 4px;
    height: 24px;
    margin-top: 4px;
  }
}
</style>
