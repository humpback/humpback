<script lang="ts" setup>
import { GenerateUUID, RulePleaseEnter, SetWebTitle } from "@/utils"
import { PageGroupDetail, RuleLength, ServiceNetworkMode, ServiceNetworkProtocol, ServiceRestartPolicyMode } from "@/models"
import { FormInstance, FormRules } from "element-plus"
import { NewServiceMetaDockerEmptyInfo } from "@/types"
import { cloneDeep, filter, find, findIndex, groupBy, toLower, trim } from "lodash-es"
import VolumesPage from "./advanced/volumes.vue"
import EnvironmentsPage from "./advanced/environments.vue"
import LabelsPage from "./advanced/labels.vue"
import CapabilitiesPage from "./advanced/capabilities.vue"
import ResourcesLogConfigPage from "./advanced/resources-log-config.vue"
import { NewApplicationInfo, ParseMetaInfo, ServiceApplicationInfo } from "./application.ts"
import { refreshData } from "../common.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const registryStore = useRegistryStore()
const stateStore = useStateStore()

const isLoading = ref(false)
const isAction = ref(false)
const isInit = ref(false)

const groupId = ref(route.params.groupId as string)
const serviceId = ref(route.params.serviceId as string)

const advancedOptions = ref<Array<{ label: string; value: string; color?: string }>>([
  { label: "label.volumes", value: "volumes" },
  { label: "label.environments", value: "environments" },
  { label: "label.labels", value: "labels" },
  { label: "label.resourcesAndLogs", value: "resourcesAndLogs" },
  { label: "label.capabilities", value: "capabilities" }
])

const advancedMode = ref<"volumes" | "environments" | "labels" | "resourcesAndLogs" | "capabilities">("volumes")
const metaInfo = ref<ServiceApplicationInfo>(NewApplicationInfo(registryStore.registries, NewServiceMetaDockerEmptyInfo()))

const formRef = useTemplateRef<FormInstance>("formRef")
const volumeRef = useTemplateRef<InstanceType<typeof VolumesPage>>("volumeRef")
const environmentRef = useTemplateRef<InstanceType<typeof EnvironmentsPage>>("environmentRef")
const labelRef = useTemplateRef<InstanceType<typeof LabelsPage>>("labelRef")
const resourcesAndLogsRef = useTemplateRef<InstanceType<typeof ResourcesLogConfigPage>>("resourcesAndLogsRef")

const rules = ref<FormRules>({
  image: [
    { required: true, validator: RulePleaseEnter("label.image"), trigger: "blur" },
    {
      required: true,
      validator: RuleLimitRange(RuleLength.ServiceName.Min, RuleLength.ServiceName.Max),
      trigger: "blur"
    },
    { required: true, validator: checkImage, trigger: "blur" }
  ],
  containerPort: [{ required: true, validator: checkContainerPort, trigger: "blur" }],
  hostPort: [{ required: true, validator: checkHostPort, trigger: "blur" }]
})

const registryDomainExist = computed(() => {
  if (registryStore.registries?.length > 0) {
    return !!find(registryStore.registries, x => x.url === metaInfo.value.registryDomain)
  }
  return true
})

function checkImage(rule: any, value: any, callback: any) {
  const image = value as string
  const list = image.split(":")
  if (list.length < 2 || trim(list[0]) === "" || trim(list[1]) === "") {
    return callback(new Error(t("rules.invalidImage")))
  }
  return callback()
}

function checkContainerPort(rule: any, value: any, callback: any) {
  const containerPort = value ? (value as number) : 0
  if (!containerPort) {
    return callback(new Error(`${t("rules.pleaseEnter")} ${t("label.containerPort")}`))
  }
  if (filter(metaInfo.value.validPorts, x => x.containerPort === containerPort).length > 1) {
    return callback(new Error(`${t("rules.duplicate")} ${t("label.containerPort")}`))
  }
  return callback()
}

function checkHostPort(rule: any, value: any, callback: any) {
  const hostPort = value ? (value as number) : 0
  if (hostPort && filter(metaInfo.value.validPorts, x => x.hostPort === hostPort).length > 1) {
    return callback(new Error(`${t("rules.duplicate")} ${t("label.hostPort")}`))
  }
  return callback()
}

function cancel() {
  router.push({ name: "groupDetail", params: { groupId: groupId.value, mode: PageGroupDetail.Services } })
}

function addPort() {
  metaInfo.value.validPorts.push({
    id: GenerateUUID(),
    containerPort: undefined,
    protocol: ServiceNetworkProtocol.NetworkProtocolTCP,
    hostPort: undefined
  })
}

function removePort(index: number) {
  metaInfo.value.validPorts.splice(index, 1)
}

function checkArrayDuplicateKey(
  isFailed: boolean,
  list: any[],
  groupByField: string,
  advancedMode: "volumes" | "environments" | "labels" | "resourcesAndLogs" | "capabilities",
  isAddErrCheck?: boolean
) {
  const obj = groupBy(list, groupByField)
  for (const key in obj) {
    if (obj[key].length > 1) {
      isFailed = true
      break
    }
  }
  const index = findIndex(advancedOptions.value, x => x.value === advancedMode)
  if (index != -1) {
    if (isAddErrCheck) {
      advancedOptions.value[index].color = isFailed ? "validator-error" : undefined
    }
    if (!isAddErrCheck && advancedOptions.value[index].color && !isFailed) {
      advancedOptions.value[index].color = undefined
    }
  }
  return !isFailed
}

function checkVolumes(isAddErrCheck?: boolean) {
  let isFailed = false
  for (const volume of metaInfo.value.validVolumes) {
    if (!volume.target || !volume.source) {
      isFailed = true
      break
    }
  }
  return checkArrayDuplicateKey(isFailed, metaInfo.value.validVolumes, "target", "volumes", isAddErrCheck)
}

function checkEnvironment(isAddErrCheck?: boolean) {
  let isFailed = false
  const envs = cloneDeep(metaInfo.value.validEnv)
  for (const env of envs) {
    env.name = toLower(env.name)
    if (!env.name || !env.value) {
      isFailed = true
      break
    }
  }
  return checkArrayDuplicateKey(isFailed, envs, "name", "environments", isAddErrCheck)
}

function checkLabels(isAddErrCheck?: boolean) {
  let isFailed = false
  const labels = cloneDeep(metaInfo.value.validLabel)
  for (const label of labels) {
    label.name = toLower(label.name)
    if (!label.name || !label.value) {
      isFailed = true
      break
    }
  }
  return checkArrayDuplicateKey(isFailed, labels, "name", "labels", isAddErrCheck)
}

function checkLogConfig(isAddErrCheck?: boolean) {
  let isFailed = false
  const logConfigs = cloneDeep(metaInfo.value.validLogConfig)
  for (const logConfig of logConfigs.options) {
    logConfig.name = toLower(logConfig.name)
    if (!logConfig.name || !logConfig.value) {
      isFailed = true
      break
    }
  }
  return checkArrayDuplicateKey(isFailed, logConfigs.options, "name", "resourcesAndLogs", isAddErrCheck)
}

async function validate() {
  const validList = await Promise.all([
    checkVolumes(true),
    checkEnvironment(true),
    checkLabels(true),
    checkLogConfig(true),
    formRef.value?.validate().catch(() => false),
    volumeRef.value?.validate().catch(() => false),
    environmentRef.value?.validate().catch(() => false),
    labelRef.value?.validate().catch(() => false),
    resourcesAndLogsRef.value?.validate().catch(() => false)
  ])
  return filter(validList, x => typeof x !== "undefined" && !x).length <= 0
}

async function search(init?: boolean) {
  isLoading.value = true
  await refreshData(groupId.value, serviceId.value, "application", init)
    .then(() => {
      metaInfo.value = NewApplicationInfo(registryStore.registries, stateStore.getService(serviceId.value)?.meta)
    })
    .finally(() => (isLoading.value = false))
}

async function save() {
  if (!(await validate()) || !registryDomainExist.value) {
    return
  }
  isAction.value = true
  await serviceService
    .update(groupId.value, {
      type: "application",
      serviceId: serviceId.value,
      data: ParseMetaInfo(metaInfo.value, registryStore.registries)
    })
    .finally(() => (isAction.value = false))
  ShowSuccessMsg(t("message.saveSuccess"))
  await search()
}

onBeforeMount(async () => {
  await search(true).finally(() => (isInit.value = true))
  SetWebTitle(`${t("webTitle.serviceInfo")} - ${stateStore.getService()?.serviceName}`)
})
</script>

<template>
  <el-skeleton v-if="!isInit" :rows="14" animated />
  <div v-else>
    <el-form ref="formRef" v-loading="isLoading" :model="metaInfo" :rules="rules" class="form-box" label-position="top" label-width="auto">
      <el-row :gutter="12">
        <el-col :span="24">
          <el-form-item :label="t('label.image')" prop="image">
            <v-input v-model="metaInfo.image" :maxlength="RuleLength.ImageName?.Max" :placeholder="t('placeholder.egImage')" clearable show-word-limit>
              <template #prepend>
                <el-dropdown placement="bottom-start" trigger="click" @command="metaInfo.registryDomain = $event">
                  <div class="registry-domain">
                    <div style="width: auto">{{ metaInfo.registryDomain }}</div>
                    <el-icon :size="18">
                      <IconMdiChevronDown />
                    </el-icon>
                  </div>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item v-for="item in registryStore.registries" :key="item.registryId" :command="item.url">{{ item.url }}</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </v-input>
            <div v-if="!registryDomainExist" class="mt-3 w-100">
              <v-alert type="error">{{ t("rules.notExistRegistry") }}</v-alert>
            </div>
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item :label="t('label.command')">
            <v-description-input v-model="metaInfo.command" />
          </el-form-item>
        </el-col>

        <el-col :md="metaInfo.restartPolicy!.mode === ServiceRestartPolicyMode.RestartPolicyModeOnFailure ? 12 : 24">
          <el-form-item :label="t('label.restartPolicy')">
            <el-select v-model="metaInfo.restartPolicy!.mode">
              <el-option :label="t('label.no')" :value="ServiceRestartPolicyMode.RestartPolicyModeNo" />
              <el-option :label="t('label.always')" :value="ServiceRestartPolicyMode.RestartPolicyModeAlways" />
              <el-option :label="t('label.onFailure')" :value="ServiceRestartPolicyMode.RestartPolicyModeOnFailure" />
              <el-option :label="t('label.unlessStopped')" :value="ServiceRestartPolicyMode.RestartPolicyModeUnlessStopped" />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col v-if="metaInfo.restartPolicy!.mode === ServiceRestartPolicyMode.RestartPolicyModeOnFailure" :md="12">
          <el-form-item :label="t('label.maxRetryCount')">
            <v-input-number v-model="metaInfo.restartPolicy!.maxRetryCount" :min="0" style="width: 100%" />
          </el-form-item>
        </el-col>

        <el-col
          :span="metaInfo.network!.mode === ServiceNetworkMode.NetworkModeHost ? 24 : metaInfo.network!.mode === ServiceNetworkMode.NetworkModeBridge ? 10 : 6">
          <el-form-item :label="t('label.network')">
            <el-select v-model="metaInfo.network!.mode">
              <el-option :label="t('label.host')" :value="ServiceNetworkMode.NetworkModeHost" />
              <el-option :label="t('label.bridge')" :value="ServiceNetworkMode.NetworkModeBridge" />
              <el-option :label="t('label.custom')" :value="ServiceNetworkMode.NetworkModeCustom" />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col v-if="metaInfo.network!.mode === ServiceNetworkMode.NetworkModeCustom" :span="6">
          <el-form-item :label="t('label.networkName')">
            <v-input v-model="metaInfo.network!.networkName" />
          </el-form-item>
        </el-col>

        <el-col
          v-if="metaInfo.network!.mode !== ServiceNetworkMode.NetworkModeHost"
          :span="metaInfo.network!.mode === ServiceNetworkMode.NetworkModeBridge ? 14 : 12">
          <el-form-item :label="t('label.hostname')" prop="network.hostname">
            <v-input v-model="metaInfo.network!.hostname" :disabled="metaInfo.network!.useMachineHostname">
              <template #prepend>
                <el-checkbox v-model="metaInfo.network!.useMachineHostname">{{ t("label.useMachineHostname") }}</el-checkbox>
              </template>
            </v-input>
          </el-form-item>
        </el-col>

        <el-col v-if="metaInfo.network!.mode !== ServiceNetworkMode.NetworkModeHost">
          <div class="network-box">
            <div class="mb-3">
              <v-tips>{{ t("tips.networkPortTips") }}</v-tips>
            </div>
            <el-row :gutter="12">
              <el-col v-for="(portInfo, index) in metaInfo.validPorts" :key="index" :span="24">
                <div class="d-flex gap-2">
                  <el-form-item :prop="`validPorts.${index}.containerPort`" :rules="rules.containerPort" class="flex-1">
                    <v-input-number
                      v-model="metaInfo.validPorts[index].containerPort"
                      :controls="false"
                      :min="0"
                      :placeholder="t('placeholder.containerPort')"
                      style="width: 100%">
                    </v-input-number>
                  </el-form-item>
                  <el-form-item :prop="`validPorts.${index}.protocol`">
                    <el-select v-model="metaInfo.validPorts[index].protocol" :placeholder="t('placeholder.protocol')" style="width: 200px">
                      <el-option :label="t('label.tcp')" :value="ServiceNetworkProtocol.NetworkProtocolTCP" />
                      <el-option :label="t('label.udp')" :value="ServiceNetworkProtocol.NetworkProtocolUDP" />
                    </el-select>
                  </el-form-item>
                  <el-form-item :prop="`validPorts.${index}.hostPort`" :rules="rules.hostPort" class="flex-1">
                    <v-input-number v-model="metaInfo.validPorts[index].hostPort" :controls="false" :placeholder="t('placeholder.hostPort')" class="flex-1" />
                  </el-form-item>
                  <el-form-item>
                    <el-button plain style="padding: 4px 12px" text type="danger" @click="removePort(index)">
                      <el-icon :size="26">
                        <IconMdiClose />
                      </el-icon>
                    </el-button>
                  </el-form-item>
                </div>
              </el-col>
              <el-col>
                <el-button size="small" type="info" @click="addPort">
                  <el-icon :size="16">
                    <IconMdiAdd />
                  </el-icon>
                  {{ t("btn.addPort") }}
                </el-button>
              </el-col>
            </el-row>
          </div>
        </el-col>

        <el-col>
          <el-form-item class="w-100 mt-5">
            <el-segmented v-model="advancedMode" :options="advancedOptions" block class="advanced-segmented w-100">
              <template #default="{ item }">
                <span :class="(item as any).color">{{ t((item as any).label as string) }}</span>
              </template>
            </el-segmented>
          </el-form-item>
        </el-col>

        <el-col>
          <div class="advanced-box">
            <div class="advanced-content">
              <volumes-page
                v-if="advancedMode === 'volumes'"
                ref="volumeRef"
                v-model="metaInfo.validVolumes"
                :has-valid="!!find(advancedOptions, x => x.value === 'volumes')?.color"
                @check="checkVolumes()" />

              <environments-page
                v-if="advancedMode === 'environments'"
                ref="environmentRef"
                v-model="metaInfo.validEnv"
                :has-valid="!!find(advancedOptions, x => x.value === 'environments')?.color"
                @check="checkEnvironment()" />

              <labels-page
                v-if="advancedMode === 'labels'"
                ref="labelRef"
                v-model="metaInfo.validLabel"
                :has-valid="!!find(advancedOptions, x => x.value === 'labels')?.color"
                @check="checkLabels()" />

              <resources-log-config-page
                v-if="advancedMode === 'resourcesAndLogs'"
                ref="resourcesAndLogsRef"
                v-model:log-config="metaInfo.validLogConfig"
                v-model:resources="metaInfo.resources"
                :has-valid="!!find(advancedOptions, x => x.value === 'resourcesAndLogs')?.color"
                @check="checkLogConfig()" />

              <capabilities-page
                v-if="advancedMode === 'capabilities'"
                v-model:cap-add="metaInfo.capabilities!.capAdd"
                v-model:cap-drop="metaInfo.capabilities!.capDrop" />
            </div>
          </div>
        </el-col>

        <el-col class="mt-5">
          <el-form-item>
            <el-checkbox v-model="metaInfo.alwaysPull">
              <el-text class="f-bold" size="small">{{ t("label.alwaysPull") }}</el-text>
            </el-checkbox>
            <v-tips>{{ t("tips.alwaysPullTips") }}</v-tips>
          </el-form-item>
        </el-col>
        <el-col>
          <el-form-item>
            <el-checkbox v-model="metaInfo.privileged">
              <el-text class="f-bold" size="small">{{ t("label.privilegedMode") }}</el-text>
            </el-checkbox>
            <v-tips>{{ t("tips.privilegedTips") }}</v-tips>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <div class="text-align-right pt-5">
      <el-button @click="cancel()">{{ t("btn.cancel") }}</el-button>
      <el-button :loading="isAction" type="primary" @click="save()">{{ t("btn.save") }}</el-button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.registry-domain {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 80px;

  &:hover {
    cursor: pointer;
    opacity: 0.7;
  }
}

.form-box {
  :deep(.el-form-item__label) {
    font-weight: 600;
    font-size: 12px;
    margin-bottom: 4px;
  }
}

.network-box {
  border: 1px solid var(--el-border-color);
  padding: 16px;
  border-radius: 4px;
  box-sizing: border-box;
  background-color: #ecf0f5;
}

.advanced-segmented {
  --hp-active-segmented-bg-color: #e3e3e4;
  padding: 0;
  border: 1px solid var(--el-border-color);
  border-right: none;

  :deep(.el-segmented__group) {
    .el-segmented__item-selected {
      display: none !important;
    }

    .el-segmented__item {
      border-right: 1px solid var(--el-border-color);
      background-color: #ffffff;

      &:has(.validator-error) {
        background-color: var(--el-color-danger);
        color: #ffffff;
      }

      &:not(:has(.validator-error)):hover {
        background-color: #f5f7fa;
      }

      &.is-selected:not(:has(.validator-error)) {
        background-color: var(--hp-active-segmented-bg-color);
        color: #2f4050ff;
      }
    }
  }

  .el-radio-button {
    flex: 1;

    .span {
      flex: 1;
    }
  }
}

.advanced-box {
  border: 1px solid var(--el-border-color);
  border-top: none;
  margin-top: -18px;
  padding: 20px;
  border-bottom-right-radius: 4px;
  border-bottom-left-radius: 4px;
  box-sizing: border-box;

  .advanced-content {
    border: 1px solid var(--el-border-color);
    padding: 16px;
    border-radius: 4px;
    box-sizing: border-box;
    background-color: #ecf0f5;
  }

  .volume-type {
    :deep(.el-radio-button):not(.is-active) {
      & .el-radio-button__inner:hover {
        color: var(--el-text-color-regular);
      }
    }
  }
}
</style>
