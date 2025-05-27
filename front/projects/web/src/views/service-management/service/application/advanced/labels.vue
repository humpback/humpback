<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { GenerateUUID, RulePleaseEnter } from "@/utils"
import { filter, toLower } from "lodash-es"

const props = defineProps<{ hasValid?: boolean }>()
const emits = defineEmits<{
  (e: "check"): void
}>()

const { t } = useI18n()

const labels = defineModel<Array<{ id: string; name: string; value: string }>>()

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  name: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: checkName, trigger: "blur" }
  ],
  value: [{ required: true, validator: RulePleaseEnter("label.value"), trigger: "blur" }]
})

function checkName(rule: any, value: any, callback: any) {
  const name = value as string
  if (filter(labels.value, x => toLower(x.name) === toLower(name)).length > 1) {
    return callback(new Error(`${t("rules.duplicate")} ${t("label.name")}`))
  }
  return callback()
}

function addLabel() {
  labels.value?.push({ id: GenerateUUID(), name: "", value: "" })
}

function removeLabel(index: number) {
  labels.value?.splice(index, 1)
  emits("check")
}

async function validate() {
  return await formRef.value?.validate()
}

onMounted(() => {
  if (props.hasValid) {
    validate()
  }
})

defineExpose({ validate })
</script>

<template>
  <div>
    <el-form ref="formRef" :model="labels" :rules="rules" label-position="top" label-width="auto">
      <div v-for="(volume, index) in labels" :key="volume.id" class="d-flex gap-3">
        <el-form-item :prop="`${index}.name`" :rules="rules.name" class="flex-1">
          <v-input v-model="labels![index].name" @blur="emits('check')">
            <template #prepend>{{ t("label.name") }}</template>
          </v-input>
        </el-form-item>
        <el-form-item :prop="`${index}.value`" :rules="rules.value" class="flex-1">
          <v-input v-model="labels![index].value" @blur="emits('check')">
            <template #prepend>{{ t("label.value") }}</template>
          </v-input>
        </el-form-item>
        <el-form-item>
          <el-button plain style="padding: 4px 12px" text type="danger" @click="removeLabel(index)">
            <el-icon :size="26">
              <IconMdiClose />
            </el-icon>
          </el-button>
        </el-form-item>
      </div>
    </el-form>
    <el-button size="small" type="info" @click="addLabel()">
      <el-icon :size="16">
        <IconMdiAdd />
      </el-icon>
      {{ t("btn.addLabel") }}
    </el-button>
  </div>
</template>

<style lang="scss" scoped>
.tips-line {
  --el-color-info: #305d8c;
  background-color: #d6dde7;
  border-left: 4px solid #9cb4c5;
}
</style>
