<script lang="ts" setup>
import { PageServiceDetail } from "@/models"
import { onClickOutside } from "@vueuse/core"
import { ElInput } from "element-plus"
import { trim } from "lodash-es"

const { t } = useI18n()

const name = ref("")
const isLoading = ref(false)
const isSearch = ref(false)
const isFocus = ref(false)

const visible = ref(false)

const contentRef = useTemplateRef<HTMLDivElement>("contentRef")
const inputRef = useTemplateRef<InstanceType<typeof ElInput>>("inputRef")

const result = ref<{
  [key: string]: Array<{ groupId: string; groupName: string; serviceId?: string; serviceName?: string }>
}>({ groups: [], services: [] })

const searchGroupService = CreateCancelRequest(commonService.searchGroupServiceByName)

function focusInput() {
  isFocus.value = true
  if (name.value != "" && isSearch.value) {
    visible.value = true
  }
}

function clickLink(isHref?: boolean) {
  if (!isHref) {
    hidePopover()
  }
}

function hidePopover() {
  visible.value = false
  inputRef.value?.blur()
  isFocus.value = false
}

function search() {
  name.value = trim(name.value)
  if (name.value === "") {
    return
  }
  visible.value = true
  isLoading.value = true
  searchGroupService(name.value)
    .then(data => {
      result.value = data
      isSearch.value = true
    })
    .finally(() => (isLoading.value = false))
}

onClickOutside(contentRef, hidePopover)
</script>

<template>
  <div ref="contentRef" :class="['global-search', isFocus && 'is-focused']">
    <div class="global-search-input">
      <el-input
        id="search-input"
        ref="inputRef"
        v-model="name"
        :placeholder="t('placeholder.searchGroupService')"
        size="small"
        @focus="focusInput()"
        @keydown.enter="search()">
        <template #prefix>
          <el-icon :size="16">
            <IconMdiSearch />
          </el-icon>
        </template>
      </el-input>
    </div>

    <div v-if="visible" class="global-search-body">
      <div class="global-search-content">
        <div class="global-search-arrow" />
        <div v-loading="isLoading" class="global-search-inner">
          <div>
            <el-text class="f-bold"> {{ t("label.groups") }}</el-text>
            <el-divider style="margin: 8px 0 16px 0; border-color: var(--el-color-info-light-9)" />
          </div>

          <div v-if="result.groups && result.groups.length > 0" class="pl-5">
            <div v-for="(item, index) in result.groups" :key="index" class="content">
              <v-router-link :href="`/ws/group/${item.groupId}/services`" :text="item.groupName" show-title @click-route="clickLink">
                <template #prefix-text>-</template>
              </v-router-link>
            </div>
          </div>

          <div v-else class="pl-5 d-flex gap-1">
            <el-text type="warning">
              <el-icon :size="16">
                <IconMdiWarningCircleOutline />
              </el-icon>
            </el-text>
            <el-text size="small" type="warning"> {{ t("tips.noGroupFound") }}</el-text>
          </div>

          <div class="mt-5">
            <el-text class="f-bold"> {{ t("label.services") }}</el-text>
            <el-divider style="margin: 8px 0 16px 0; border-color: var(--el-color-info-light-9)" />
            <div v-if="result.services && result.services.length > 0" class="pl-5">
              <div v-for="(item, index) in result.services" :key="index" class="content">
                <el-text :title="item.groupName" size="small" type="info">{{ item.groupName }}</el-text>
                <br />
                <v-router-link
                  :href="`/ws/group/${item.groupId}/service/${item.serviceId}/${PageServiceDetail.BasicInfo}`"
                  :text="item.serviceName!"
                  show-title
                  @click-route="clickLink">
                  <template #prefix-text>-</template>
                </v-router-link>
              </div>
            </div>
            <div v-else class="pl-5 d-flex gap-1">
              <el-text type="warning">
                <el-icon :size="16">
                  <IconMdiWarningCircleOutline />
                </el-icon>
              </el-text>
              <el-text size="small" type="warning">{{ t("tips.noServiceFound") }}</el-text>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.global-search {
  width: 200px;
  transition: width 0.3s ease;
  position: relative;

  &.is-focused {
    width: 400px;
  }

  .global-search-input {
    width: 100%;
    position: absolute;
    left: 0;
    top: 0;

    :deep(.el-input .el-input__wrapper) {
      border-radius: 16px;
    }
  }

  .global-search-body {
    position: absolute;
    top: 28px;
    left: 0;

    .global-search-content {
      position: relative;

      .global-search-arrow {
        position: absolute;
        top: 0;
        left: 20px;
        width: 10px;
        height: 10px;
        background-color: #ffffff;
        border-top: 1px solid var(--el-border-color);
        border-left: 1px solid var(--el-border-color);
        transform: rotate(45deg);
        z-index: 1001;
      }

      .global-search-inner {
        box-sizing: border-box;
        position: absolute;
        top: 5px;
        left: 0;
        background-color: #ffffff;
        border-radius: 8px;
        border: 1px solid var(--el-border-color);
        box-shadow: var(--el-box-shadow-lighter);
        padding: 16px;
        width: 400px;
        max-height: 600px;
        overflow-y: auto;
        z-index: 1000;

        .content {
          font-size: 14px;
          margin-bottom: 2px;
        }
      }
    }
  }
}
</style>

<style lang="scss"></style>
