import { NewPageInfo, NewSortInfo, QueryInfo, ServiceInfo } from "@/types"
import { cloneDeep, find, omitBy } from "lodash-es"
import { NodeSwitch, ServiceStatus } from "@/models"
import { groupService } from "services/group-service.ts"

export const InjectKeyIsLoading = Symbol("IsLoading")
export const InjectKeyResetLoopSearch = Symbol("ResetLoopSearch")
export const InjectKeyChangeTab = Symbol("ChangeTab")

export const sortOptions = ["serviceName", "updatedAt", "createdAt"]

export const defaultSort = NewSortInfo("serviceName", "asc")
export const defaultPage = NewPageInfo(1, 20)
export const defaultFilter = { status: "", schedule: "" }

export class QueryServicesInfo extends QueryInfo {
  constructor(queryInfo: any) {
    super(queryInfo, ["keywords"], defaultPage, defaultSort, sortOptions, cloneDeep(defaultFilter))
    const status = find(
      [NodeSwitch.Enabled, NodeSwitch.Disabled, ServiceStatus.ServiceStatusNotReady, ServiceStatus.ServiceStatusRunning, ServiceStatus.ServiceStatusFailed],
      x => x === (queryInfo["status"] as string)
    )
    this.filter.status = status || ""
    const schedule = find(["Yes", "No"], x => x === (queryInfo["schedule"] as string))
    this.filter.schedule = schedule || ""
  }

  urlQuery() {
    return {
      query: Object.assign(
        {},
        {
          status: this.filter.status !== defaultFilter.status ? this.filter.status : undefined,
          schedule: this.filter.schedule !== defaultFilter.schedule ? this.filter.schedule : undefined
        },
        this.getBaseQuery()
      )
    }
  }

  searchParams() {
    return omitBy(this, (value, key) => key.startsWith("_"))
  }
}

export const ActionOptions: Array<{
  action: "Start" | "Stop" | "Restart" | "Enable" | "Disable"
  type: "default" | "info" | "success" | "primary" | "text" | "warning" | "danger"
  i18nLabel: string
  icon: any
}> = [
  { action: "Enable", type: "primary", i18nLabel: "btn.enable", icon: IconMdiPlay },
  { action: "Disable", type: "info", i18nLabel: "btn.disable", icon: IconMdiSquare },
  { action: "Start", type: "success", i18nLabel: "btn.start", icon: IconMdiPlay },
  { action: "Restart", type: "success", i18nLabel: "btn.restart", icon: IconMdiRestart },
  { action: "Stop", type: "primary", i18nLabel: "btn.stop", icon: IconMdiStopCircleOutline }
]

async function getGroupInfo(groupId: string) {
  return await groupService.info(groupId).then(info => {
    useStateStore().setGroup(groupId, info)
  })
}

async function getServiceInfo(groupId: string, serviceId: string) {
  return await serviceService.info(groupId, serviceId).then(info => {
    useStateStore().setService(serviceId, info)
  })
}

async function getGroupNodes(groupId: string) {
  return await groupService.getNodes(groupId).then(nodes => {
    useStateStore().setGroupNodeList(groupId, nodes)
  })
}

export async function refreshData(
  groupId: string,
  serviceId: string,
  mode: "global" | "basic-info" | "application" | "deployment" | "instances" | "log" | "activity",
  init?: boolean
) {
  const taskList: Array<Promise<void> | undefined> = [
    getGroupInfo(groupId),
    getServiceInfo(groupId, serviceId),
    mode === "application" && init ? useRegistryStore().refreshRegistries() : undefined,
    mode === "deployment" && init ? getGroupNodes(groupId) : undefined
  ]
  return await Promise.all(taskList)
}

export function showAction(serviceInfo?: ServiceInfo, action?: "Start" | "Stop" | "Restart" | "Enable" | "Disable") {
  if (!serviceInfo) {
    return false
  }
  switch (action) {
    case "Start":
    case "Stop":
    case "Restart":
    case "Disable":
      {
        if (serviceInfo.isEnabled) {
          return true
        }
      }
      break
    case "Enable": {
      if (!serviceInfo.isEnabled && !!serviceInfo.meta && !!serviceInfo.deployment) {
        return true
      }
    }
  }
  return false
}
