<script lang="ts" setup>
import { SetWebTitle } from "@/utils"
import { refreshData } from "../common.ts"
import { RuleLength } from "@/models"
import { filter, find, toLower } from "lodash-es"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const stateStore = useStateStore()

const isLoading = ref(false)
const isAction = ref(false)

const groupId = ref(route.params.groupId as string)
const serviceId = ref(route.params.serviceId as string)
const containers = computed<ServiceContainerStatusInfo[]>(() => filter(stateStore.getService(serviceId.value)?.containers || [], x => !!x.containerId))

const searchInfo = ref({
  instance: "",
  line: 1000 as any,
  startAt: 0,
  endAt: 0,
  showTimestamp: false
})

const logs = ref<string[]>([])

function handleLineChange(v: any) {
  if (v) {
    searchInfo.value.line = Number(v).valueOf() > RuleLength?.LogsLine?.Max ? RuleLength?.LogsLine?.Max : Number(v).valueOf()
  } else {
    searchInfo.value.line = undefined
  }
}

function newQuery() {
  return {
    instance: searchInfo.value.instance || undefined,
    line: searchInfo.value.line || undefined,
    startAt: searchInfo.value.startAt || undefined,
    endAt: searchInfo.value.endAt || undefined,
    showTimestamp: searchInfo.value.showTimestamp || undefined
  }
}

function parseQuery() {
  const query = route.query
  searchInfo.value.instance = query?.instance ? (query.instance as string) : ""
  searchInfo.value.line = query?.line ? Number(query.line).valueOf() : 1000
  searchInfo.value.startAt = query?.endAt ? Number(query.endAt).valueOf() : 0
  searchInfo.value.endAt = query?.endAt ? Number(query.endAt).valueOf() : 0
  searchInfo.value.showTimestamp = toLower(query?.showTimestamp as string) === "true"
}

const queryLogs = CreateCancelRequest(async (...args: any[]) => {
  return groupContainerService.logs(args[0], args[1], args[2])
})

async function getLogs() {
  if (isAction.value) {
    return
  }
  const instanceInfo = find(containers.value, x => !!x.containerId && x.containerId === searchInfo.value.instance)
  if (!instanceInfo) {
    if (searchInfo.value.instance) {
      ShowErrMsg(t("err.instanceNotExist"))
    }
    return
  }
  await router.replace({ params: route.params, query: Object.assign({}, newQuery()) as any })
  isAction.value = true
  await queryLogs(groupId.value, serviceId.value, {
    nodeId: instanceInfo.nodeId,
    containerId: searchInfo.value.instance,
    line: searchInfo.value.line || 0,
    startAt: searchInfo.value.startAt || 0,
    endAt: searchInfo.value.endAt || 0,
    showTimestamp: searchInfo.value.showTimestamp
  })
    .then(data => {
      logs.value = data || []
    })
    .finally(() => {
      isAction.value = false
    })
}

async function search() {
  isLoading.value = true
  await refreshData(groupId.value, serviceId.value, "log").finally(() => (isLoading.value = false))
}

onMounted(async () => {
  parseQuery()
  await search()
  SetWebTitle(`${t("webTitle.serviceInfo")} - ${stateStore.getService()?.serviceName}`)
  if (searchInfo.value.instance === "" && containers.value.length > 0) {
    searchInfo.value.instance = containers.value[0].containerId
  }
  await getLogs()
})
</script>

<template>
  <div class="d-flex gap-3">
    <div class="header-icon">
      <el-icon :size="18">
        <IconMdiFileDocumentOutline />
      </el-icon>
    </div>
    <el-text class="f-bold" size="large">{{ t("label.instanceLogs") }}</el-text>
    <v-loading v-if="isLoading" />
  </div>

  <div>
    <div class="mt-3">
      <v-tips>{{ t("tips.logsTips") }}</v-tips>
    </div>
    <div class="d-flex gap-3 flex-wrap mt-5">
      <div class="flex-1" style="min-width: 460px">
        <v-select v-model="searchInfo.instance" :out-label="t('label.instance')" out-label-width="100px" placeholder="" @change="getLogs()">
          <el-option v-for="item in containers" :key="item.containerId" :label="item.containerName" :value="item.containerId" />
        </v-select>
      </div>
      <div class="flex-1" style="min-width: 600px">
        <v-log-search-time-range
          v-model:end-time="searchInfo.endAt"
          v-model:start-time="searchInfo.startAt"
          :out-label="t('label.timeRange')"
          out-label-width="120px"
          show-out-label
          @change="getLogs()" />
      </div>
      <div style="width: 200px">
        <el-input
          :max="RuleLength?.LogsLine?.Max"
          :min="RuleLength?.LogsLine?.Min"
          :model-value="searchInfo.line"
          type="number"
          @change="getLogs()"
          @update:modelValue="handleLineChange">
          <template #prepend>
            <el-text>{{ t("label.lines") }}</el-text>
          </template>
        </el-input>
      </div>

      <el-button :disabled="isAction || !searchInfo.instance" plain type="primary" @click="getLogs()">
        {{ t("btn.refresh") }}
      </el-button>
    </div>

    <div class="log-box">
      <v-log-view v-loading="isAction" :logList="logs">
        <template #header-left>
          <el-switch v-model="searchInfo.showTimestamp" :active-text="t('label.showTimestamp')" @change="getLogs()" />
        </template>
      </v-log-view>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.header-icon {
  color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
  border-radius: 50%;
  padding: 8px;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.log-box {
  height: 600px;
  margin-top: 20px;
}
</style>
