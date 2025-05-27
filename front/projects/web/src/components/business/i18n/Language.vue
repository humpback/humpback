<script lang="ts" setup>
import { ChangeLanguage, GetLanguageOptions } from "@/locales"

const props = withDefaults(
  defineProps<{
    trigger?: "hover" | "click" | "contextmenu"
    iconSize?: number
  }>(),
  {
    trigger: "hover",
    iconSize: 16
  }
)

const languageOptions = computed(() => GetLanguageOptions())
</script>

<template>
  <el-dropdown :show-timeout="0" :trigger="props.trigger" @command="ChangeLanguage">
    <div class="lang-content">
      <el-icon :size="props.iconSize">
        <IconMdiLanguage />
      </el-icon>
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <template v-for="item in languageOptions" :key="item.value">
          <el-dropdown-item v-if="!item.disabled" :command="item.value">
            <div> {{ item.name }}</div>
          </el-dropdown-item>
        </template>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style lang="scss" scoped>
.lang-content {
  display: flex;
  align-items: center;
  gap: 4px;

  &:hover {
    cursor: pointer;
    opacity: 0.7;
  }

  &:focus-visible {
    outline: none;
  }
}
</style>
