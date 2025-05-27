<script lang="ts" setup>
const props = defineProps<{ addLabel?: string; inputLabel?: string; placeholder?: string }>()
const emits = defineEmits<{
  (e: "search"): void
  (e: "add"): void
}>()

const { t } = useI18n()
const slots = useSlots()

const keywords = defineModel<string>()
</script>

<template>
  <el-form @submit.prevent="emits('search')">
    <el-form-item>
      <div class="d-flex w-100 flex-wrap gap-3">
        <slot v-if="!!slots.prefix" name="prefix" />

        <slot v-if="!!slots.input" name="input" />
        <div v-else class="search-input">
          <v-input v-model="keywords" :placeholder="props.placeholder">
            <template #prepend>
              <slot v-if="!!slots.inputPrepend" name="inputPrepend"></slot>
              <span v-else style="color: var(--el-text-color-regular)">
                {{ props.inputLabel || t("label.keywords") }}
              </span>
            </template>
          </v-input>
        </div>

        <div>
          <el-button native-type="submit" type="primary">
            <el-icon :size="16">
              <IconMdiSearch />
            </el-icon>
            {{ t("btn.search") }}
          </el-button>
          <el-button v-if="props.addLabel" plain type="primary" @click="emits('add')">
            <el-icon :size="16">
              <IconMdiAdd />
            </el-icon>
            {{ props.addLabel }}
          </el-button>
        </div>
      </div>
    </el-form-item>
  </el-form>
</template>

<style lang="scss" scoped>
.search-input {
  flex: 1;
  min-width: 300px;
}
</style>
