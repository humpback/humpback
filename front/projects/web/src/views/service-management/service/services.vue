<script lang="ts" setup>
import { ServiceInfo } from "@/types"
import { SetWebTitle, TableHeight } from "@/utils"
import { Action, PageServiceDetail } from "@/models"
import { ActionOptions, QueryServicesInfo, showAction } from "./common.ts"
import { serviceService } from "services/service-service.ts"
import ServiceCreate from "./action/service-create.vue"
import ServiceDelete from "./action/service-delete.vue"
import VServiceStatusTag from "@/components/business/v-service/VServiceStatusTag.vue"
import { toLower } from "lodash-es"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const stateStore = useStateStore()

const tableHeight = computed(() => TableHeight(362))

const groupId = ref(route.params?.groupId as string)

const isLoading = ref(false)
const queryInfo = ref<QueryServicesInfo>(new QueryServicesInfo(route.query))

const createRef = useTemplateRef<InstanceType<typeof ServiceCreate>>("createRef")
const deleteRef = useTemplateRef<InstanceType<typeof ServiceDelete>>("deleteRef")

const tableList = ref({
  total: 0,
  data: [] as Array<ServiceInfo>
})

const setRowClass = ({ row }) => {
  return row.containers.length === 0 || !row.isEnabled ? "hide-expand-icon" : ""
}

function routerToInstanceLog(serviceId: string, containerId: string) {
  router.push({
    name: "serviceInfo",
    params: { groupId: groupId.value, serviceId: serviceId, mode: PageServiceDetail.Log },
    query: { instance: containerId }
  })
}

function serviceInfoIsInComplete(info: ServiceInfo) {
  return !info?.meta || !info?.deployment
}

async function getGroupInfo() {
  return await groupService.info(groupId.value).then(info => {
    stateStore.setGroup(groupId.value, info)
  })
}

async function getServiceTotal() {
  return await serviceService.total(groupId.value).then(total => {
    stateStore.setGroupTotal(groupId.value, total)
  })
}

async function getServices() {
  return await serviceService.query(groupId.value, queryInfo.value.searchParams()).then(info => {
    tableList.value.data = info.list
    tableList.value.total = info.total
  })
}

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  await Promise.all([getGroupInfo(), getServiceTotal(), getServices()]).finally(() => (isLoading.value = false))
}

async function operateService(serviceId: string, action: "Start" | "Stop" | "Restart" | "Enable" | "Disable") {
  await serviceService.operate(groupId.value, { serviceId: serviceId, action: action })
  ShowSuccessMsg(t("message.succeed"))
  await search()
}

async function operateContainer(serviceId: string, nodeId: string, containerId: string, action: "Start" | "Stop" | "Restart") {
  await groupContainerService.operate(groupId.value, serviceId, {
    containerId: containerId,
    nodeId: nodeId,
    action: action
  })
  ShowSuccessMsg(t("message.operateSuccess"))
  await search()
}

function openAction(action: string, info?: ServiceInfo) {
  switch (action) {
    case Action.Add: {
      createRef.value?.open()
      break
    }
    case Action.Delete: {
      deleteRef.value?.open(info!)
    }
  }
}

onMounted(async () => {
  await search()
  SetWebTitle(`${t("webTitle.services")} - ${stateStore.getGroup()?.groupName}`)
})
</script>

<template>
  <v-search
    v-model="queryInfo.keywords"
    :add-label="t('btn.addService')"
    :placeholder="t('placeholder.enterNameOrImageName')"
    @add="openAction(Action.Add)"
    @search="search">
    <template #prefix>
      <div style="width: 220px">
        <v-service-status-query-select v-model="queryInfo.filter.status" :placeholder="t('placeholder.all')" @change="search" />
      </div>
      <div style="width: 200px">
        <v-service-schedule-query-select v-model="queryInfo.filter.schedule" :placeholder="t('placeholder.all')" @change="search" />
      </div>
    </template>
  </v-search>

  <v-table
    v-loading="isLoading"
    v-model:page-info="queryInfo.pageInfo"
    v-model:sort-info="queryInfo.sortInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :row-class-name="setRowClass"
    :total="tableList.total"
    row-key="serviceId"
    @page-change="search"
    @sort-change="search">
    <el-table-column align="left" class-name="expand-column" type="expand" width="24">
      <template #default="scope">
        <div style="padding: 20px 40px">
          <v-table :data="scope.row.containers" :max-height="500" border minHeight="0">
            <el-table-column :label="t('label.instanceName')" min-width="200">
              <template #default="cscope">
                <v-router-link :href="`/ws/group/${groupId}/service/${scope.row.serviceId}/${PageServiceDetail.Instances}`" :text="cscope.row.containerName" />
              </template>
            </el-table-column>
            <el-table-column :label="t('label.status')" min-width="160">
              <template #default="cscope">
                <div class="d-flex gap-3">
                  <v-container-status :status="cscope.row.state" size="small" />
                  <v-tooltip v-if="scope.row.errorMsg" effect="dark">
                    <template #content>
                      <el-text type="danger">{{ cscope.row.errorMsg }}</el-text>
                    </template>
                    <el-icon :size="20" color="var(--el-color-danger)">
                      <IconMdiWarningCircleOutline />
                    </el-icon>
                  </v-tooltip>
                </div>
              </template>
            </el-table-column>
            <el-table-column :label="t('label.ip')" min-width="160" prop="ip" />
            <el-table-column :label="t('label.createTime')" min-width="140">
              <template #default="cscope">
                <v-date-view :timestamp="cscope.row.created" />
              </template>
            </el-table-column>
            <el-table-column :label="t('label.startTime')" min-width="140">
              <template #default="cscope">
                <v-date-view :timestamp="cscope.row.started" />
              </template>
            </el-table-column>
            <el-table-column :label="t('label.action')" width="180">
              <template #default="cscope">
                <el-button
                  :disabled="!cscope.row.containerId"
                  :title="t('label.restart')"
                  link
                  type="success"
                  @click="operateContainer(scope.row.serviceId, cscope.row.nodeId, cscope.row.containerId, 'Restart')">
                  <el-icon :size="16">
                    <IconMdiRestart />
                  </el-icon>
                </el-button>
                <el-button
                  :disabled="!cscope.row.containerId"
                  :title="t('label.start')"
                  link
                  type="success"
                  @click="operateContainer(scope.row.serviceId, cscope.row.nodeId, cscope.row.containerId, 'Start')">
                  <el-icon :size="16">
                    <IconMdiPlay />
                  </el-icon>
                </el-button>
                <el-button
                  :disabled="!cscope.row.containerId"
                  :title="t('label.stop')"
                  link
                  type="danger"
                  @click="operateContainer(scope.row.serviceId, cscope.row.nodeId, cscope.row.containerId, 'Stop')">
                  <el-icon :size="16">
                    <IconMdiSquare />
                  </el-icon>
                </el-button>
                <el-button
                  :disabled="!cscope.row.containerId"
                  :title="t('label.log')"
                  link
                  type="primary"
                  @click="routerToInstanceLog(scope.row.serviceId, cscope.row.containerId)">
                  <el-icon :size="16">
                    <IconMdiNoteText />
                  </el-icon>
                </el-button>
              </template>
            </el-table-column>
          </v-table>
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.serviceName')" class-name="serviceName-column" min-width="200" prop="serviceName" sortable="custom">
      <template #default="scope">
        <v-router-link :href="`/ws/group/${groupId}/service/${scope.row.serviceId}/${PageServiceDetail.BasicInfo}`" :text="scope.row.serviceName" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.description')" min-width="200" prop="description">
      <template #default="scope">
        <v-table-column-none :text="scope.row.description" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.status')" min-width="130" prop="description">
      <template #default="scope">
        <div class="d-flex gap-1">
          <v-service-status-tag :is-enabled="scope.row.isEnabled" :status="scope.row.status" />
          <v-tooltip v-if="serviceInfoIsInComplete(scope.row)" :content="t('tips.serviceInfoInComplete')" effect="dark" placement="top-start">
            <el-button link type="warning">
              <el-icon :size="18">
                <IconMdiWarningCircleOutline />
              </el-icon>
            </el-button>
          </v-tooltip>
          <v-memo v-else-if="toLower(scope.row.status) === toLower(ServiceStatus.ServiceStatusFailed)" :icon-size="18" :memo="scope.row.memo" only-icon />
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.image')" min-width="200" prop="image">
      <template #default="scope">
        <v-table-column-none :text="[scope.row.meta?.registryDomain, scope.row.meta?.image].join('/')" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.deployMode')" min-width="240">
      <template #default="scope">
        <div v-if="scope.row.deployment" class="d-flex gap-2">
          <div style="min-width: 100px">
            <el-text v-if="scope.row.deployment.mode === ServiceDeployMode.DeployModeGlobal">{{ t("label.global") }}</el-text>
            <el-text v-if="scope.row.deployment.mode === ServiceDeployMode.DeployModeReplicate">{{ t("label.replicated") }}</el-text>
            <el-text v-if="scope.row.deployment.mode === ServiceDeployMode.DeployModeReplicate" type="primary">
              {{ ` (${scope.row.deployment.replicas})` }}
            </el-text>
          </div>
          <v-service-instance-status-statistics v-if="scope.row.isEnabled" :info="scope.row" />
        </div>
        <span v-else>--</span>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.updateDate')" min-width="160" prop="updatedAt" sortable="custom">
      <template #default="scope">
        <v-date-view :timestamp="scope.row.updatedAt" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.action')" :show-overflow-tooltip="false" align="right" fixed="right" header-align="center" width="140">
      <template #default="scope">
        <div class="d-flex gap-3">
          <el-dropdown
            :disabled="serviceInfoIsInComplete(scope.row)"
            placement="bottom-end"
            trigger="click"
            @command="operateService(scope.row.serviceId, $event)">
            <el-link :disabled="serviceInfoIsInComplete(scope.row)" :underline="false" class="d-flex gap-1" link type="primary">
              {{ t("btn.action") }}
              <el-icon>
                <IconMdiChevronDown />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <template v-for="item in ActionOptions" :key="item.action">
                  <el-dropdown-item v-if="showAction(scope.row, item.action)" :command="item.action">
                    <el-button :type="item.type" link size="small">
                      <el-icon :size="14">
                        <component :is="item.icon" />
                      </el-icon>
                      {{ t(item.i18nLabel) }}
                    </el-button>
                  </el-dropdown-item>
                </template>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button link type="danger" @click="openAction(Action.Delete, scope.row)">{{ t("btn.delete") }}</el-button>
        </div>
      </template>
    </el-table-column>
  </v-table>

  <service-delete ref="deleteRef" @refresh="search()" />

  <service-create ref="createRef" />
</template>

<style lang="scss" scoped>
:deep(.expand-column) {
  .cell {
    padding: 0 4px 0 8px;
  }
}

:deep(.hide-expand-icon) {
  .expand-column .cell {
    padding-top: 4px;
    display: none;
  }
}

:deep(.el-table__header) {
  .serviceName-column .cell {
    margin-left: -16px;
  }
}

:deep(.el-table__body) {
  .el-table__row .serviceName-column {
    .cell {
      padding-left: 4px;
    }
  }
}
</style>
