<script lang="ts" setup>
import { ServiceInfo } from "@/types"
import { map, toLower } from "lodash-es"

const props = withDefaults(defineProps<{ info: ServiceInfo; size?: "large" | "default" | "small"; effect?: "dark" | "light" | "plain" }>(), {
  size: "small",
  effect: "dark"
})

const { t } = useI18n()

const total = computed(() => {
  const result = {
    running: 0,
    exit: 0,
    failed: 0
  }
  map(props.info.containers, x => {
    switch (toLower(x.state)) {
      case toLower(ContainerStatus.ContainerStatusRunning):
        result.running++
        break
      case toLower(ContainerStatus.ContainerStatusCreated):
      case toLower(ContainerStatus.ContainerStatusPending):
      case toLower(ContainerStatus.ContainerStatusExited):
      case toLower(ContainerStatus.ContainerStatusStarting):
        result.exit++
        break
      case toLower(ContainerStatus.ContainerStatusFailed):
      case toLower(ContainerStatus.ContainerStatusWarning):
        result.failed++
    }
  })
  return result
})
</script>

<template>
  <div class="d-flex gap-1">
    <el-tag :effect="props.effect" :size="props.size" :title="t('label.running')" type="success">{{ total.running }}</el-tag>
    /
    <el-tag
      :effect="props.effect"
      :size="props.size"
      :title="toLower(props.info.deployment?.type) === toLower(ServiceDeployType.DeployTypeBackground) ? t('label.exited') : t('label.idle')"
      type="warning">
      {{ total.exit }}
    </el-tag>
    /
    <el-tag :effect="props.effect" :size="props.size" :title="t('label.failed')" type="danger">{{ total.failed }}</el-tag>
  </div>
</template>

<style lang="scss" scoped></style>
