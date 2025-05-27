<script lang="ts" setup>
import { SetWebTitle } from "@/utils"
import { InjectKeyChangeTab, InjectKeyIsLoading, InjectKeyResetLoopSearch, refreshData } from "../common.ts"
import { ServiceInfo } from "@/types"
import VLoading from "@/components/business/v-loading/VLoading.vue"
import { toLower, uniqWith } from "lodash-es"

const { t } = useI18n()
const route = useRoute()
const stateStore = useStateStore()

const groupId = ref(route.params.groupId as string)
const serviceId = ref(route.params.serviceId as string)
const serviceInfo = computed<ServiceInfo | undefined>(() => stateStore.getService(serviceId.value))

const isLoading = inject<any>(InjectKeyIsLoading)
const resetLoopSearch = inject<() => void>(InjectKeyResetLoopSearch)
const menuChange = inject<(v: string, query?: any) => void>(InjectKeyChangeTab)

function portDeduplication(ports: Array<{ bindIP: string; privatePort: number; publicPort: number; type: string }>) {
  return uniqWith(ports, (a, b) => a.privatePort === b.privatePort && a.publicPort === b.publicPort)
}

function routerToLogs(containerId: string) {
  if (menuChange && typeof menuChange === "function") {
    menuChange(PageServiceDetail.Log, { instance: containerId })
  }
}

async function search() {
  isLoading.value = true
  await refreshData(groupId.value, serviceId.value, "instances").finally(() => (isLoading.value = false))
}

async function operateContainer(nodeId: string, containerId: string, action: "Start" | "Stop" | "Restart") {
  await groupContainerService.operate(groupId.value, serviceId.value, { containerId: containerId, nodeId: nodeId, action: action })
  ShowSuccessMsg(t("message.operateSuccess"))
  await search()
  if (resetLoopSearch !== undefined) {
    resetLoopSearch()
  }
}

onMounted(async () => {
  await search()
  SetWebTitle(`${t("webTitle.serviceInfo")} - ${stateStore.getService()?.serviceName}`)
})
</script>

<template>
  <div class="d-flex gap-3">
    <el-text class="f-bold" size="large">{{ t("label.instanceOverview") }}</el-text>
    <div>
      <el-button plain size="small" type="success" @click="search()">{{ t("btn.refresh") }}</el-button>
      <el-button plain size="small" type="primary" @click="menuChange?.(PageServiceDetail.Performance)">
        {{ t("btn.viewMonitor") }}
      </el-button>
    </div>
    <v-loading v-if="isLoading" />
  </div>

  <v-memo v-if="toLower(serviceInfo?.status) === toLower(ServiceStatus.ServiceStatusFailed)" :memo="serviceInfo?.memo" class="mt-5" />

  <v-table :data="serviceInfo?.containers || []" border class="mt-5" row-key="containerName">
    <el-table-column class-name="expand-column" type="expand" width="24">
      <template #default="scope">
        <div class="expand-content">
          <div v-if="scope.row.errorMsg" class="mb-5 d-flex gap-1" style="align-items: start">
            <el-icon :size="16" color="var(--el-color-danger)">
              <IconMdiWarningCircleOutline />
            </el-icon>
            <el-text type="danger">
              {{ scope.row.errorMsg }}
            </el-text>
          </div>

          <el-form label-position="top" label-width="auto">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item>
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiClockTimeFourOutline />
                      </el-icon>
                      {{ t("label.createTime") }}
                    </el-text>
                  </template>
                  <v-date-view :timestamp="scope.row.created" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item>
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiClockTimeFourOutline />
                      </el-icon>
                      {{ t("label.startTime") }}
                    </el-text>
                  </template>
                  <v-date-view :timestamp="scope.row.started" />
                </el-form-item>
              </el-col>
              <el-divider border-style="dashed" />

              <el-col v-if="scope.row.nextAt" :span="24">
                <el-form-item>
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiClockTimeThreeOutline />
                      </el-icon>
                      {{ t("label.nextTime") }}
                    </el-text>
                  </template>
                  <v-date-view :timestamp="scope.row.nextAt" />
                </el-form-item>
                <el-divider border-style="dashed" />
              </el-col>

              <el-col :span="12">
                <el-form-item>
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiAppleKeyboardCommand />
                      </el-icon>
                      {{ t("label.command") }}
                    </el-text>
                  </template>
                  <span>{{ scope.row.command || "--" }}</span>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('label.image')">
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiAlphaCBoxOutline />
                      </el-icon>
                      {{ t("label.image") }}
                    </el-text>
                  </template>
                  <span>{{ scope.row.image || "--" }}</span>
                </el-form-item>
              </el-col>
              <el-divider border-style="dashed" />

              <el-col :span="12">
                <el-form-item :label="t('label.volumes')">
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiFileOutline />
                      </el-icon>
                      {{ t("label.volumes") }}
                    </el-text>
                  </template>
                  <div v-if="scope.row.mounts?.length > 0">
                    <div v-for="(item, index) in scope.row.mounts" :key="index" class="form-line">
                      <div class="line-prefix">-</div>
                      <div>
                        {{ `${item.source}:${item.destination}` }}
                      </div>
                    </div>
                  </div>
                  <el-text v-else>{{ t("tips.noVolumeMappingSetting") }}</el-text>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('label.ports')">
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiAlphaPCircleOutline />
                      </el-icon>
                      {{ t("label.ports") }}
                    </el-text>
                  </template>
                  <div v-if="scope.row.ports?.length > 0">
                    <div v-for="(item, index) in portDeduplication(scope.row.ports)" :key="index" class="form-line">
                      <div class="line-prefix">-</div>
                      <div>
                        {{ item.type }}
                        <el-text type="primary">{{ `${item.privatePort}:${item.publicPort}` }}</el-text>
                      </div>
                    </div>
                  </div>
                  <el-text v-else>{{ t("tips.noPortSetting") }}</el-text>
                </el-form-item>
              </el-col>
              <el-divider border-style="dashed" />
              <el-col :span="12">
                <el-form-item :label="t('label.environments')">
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiMapMarkerPath />
                      </el-icon>
                      {{ t("label.environments") }}
                    </el-text>
                  </template>
                  <div v-if="scope.row.env?.length > 0">
                    <div v-for="(item, index) in scope.row.env" :key="index" class="form-line">
                      <div class="line-prefix">-</div>
                      <div> {{ item }}</div>
                    </div>
                  </div>
                  <el-text v-else>{{ t("tips.noEnvSetting") }}</el-text>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('label.labels')">
                  <template #label>
                    <el-text type="success">
                      <el-icon :size="14">
                        <IconMdiTagTextOutline />
                      </el-icon>
                      {{ t("label.labels") }}
                    </el-text>
                  </template>
                  <div v-if="Object.keys(scope.row.labels || {})?.length > 0">
                    <div v-for="(key, index) in Object.keys(scope.row.labels)" :key="index" class="form-line">
                      <div class="line-prefix">-</div>
                      <div> {{ `${key}:${scope.row.labels[key]}` }}</div>
                    </div>
                  </div>
                  <el-text v-else>{{ t("tips.noLabelSetting") }}</el-text>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.name')" class-name="containerName-column" min-width="300" prop="containerName" sortable />
    <el-table-column :label="t('label.ip')" min-width="160" prop="ip" />
    <el-table-column :label="t('label.status')" min-width="160">
      <template #default="scope">
        <div class="d-flex gap-3">
          <v-container-status :status="scope.row.state" size="small" />
          <v-tooltip v-if="scope.row.errorMsg" effect="dark" max-width="400px" placement="top-start">
            <template #content>
              <el-text type="danger">{{ scope.row.errorMsg }}</el-text>
            </template>
            <el-icon :size="20" color="var(--el-color-danger)">
              <IconMdiWarningCircleOutline />
            </el-icon>
          </v-tooltip>
          <v-loading v-if="serviceInfo?.isEnabled && isLoading" />
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.action')" width="180">
      <template #default="scope">
        <div v-if="scope.row.containerId">
          <el-button :title="t('label.restart')" link type="success" @click="operateContainer(scope.row.nodeId, scope.row.containerId, 'Restart')">
            <el-icon :size="16">
              <IconMdiRestart />
            </el-icon>
          </el-button>
          <el-button :title="t('label.start')" link type="success" @click="operateContainer(scope.row.nodeId, scope.row.containerId, 'Start')">
            <el-icon :size="16">
              <IconMdiPlay />
            </el-icon>
          </el-button>
          <el-button :title="t('label.stop')" link type="danger" @click="operateContainer(scope.row.nodeId, scope.row.containerId, 'Stop')">
            <el-icon :size="16">
              <IconMdiSquare />
            </el-icon>
          </el-button>
          <el-button :title="t('label.log')" link type="primary" @click="routerToLogs(scope.row.containerId)">
            <el-icon :size="16">
              <IconMdiNoteText />
            </el-icon>
          </el-button>
        </div>
      </template>
    </el-table-column>
  </v-table>
</template>

<style lang="scss" scoped>
:deep(.el-table__header) .containerName-column .cell {
  margin-left: -20px;
}

:deep(.el-table__body) .el-table__row {
  .containerName-column .cell {
    padding-left: 4px;
  }
}

:deep(.expand-column) {
  border-right: 0;

  .cell {
    padding: 0 4px 0 8px;
  }
}

.expand-content {
  box-sizing: border-box;
  padding: 20px 40px;

  .el-divider {
    margin-top: 4px;
    margin-bottom: 12px;
  }

  :deep(.el-form-item__label) {
    margin-bottom: 12px;
  }

  :deep(.el-form-item__content) {
    padding-left: 20px;
    font-size: 14px;
    line-height: 22px;
    word-break: break-all;
  }
}

.form-line {
  display: flex;
  align-items: start;
  gap: 4px;

  .line-prefix {
    color: var(--el-color-success);
  }
}
</style>
