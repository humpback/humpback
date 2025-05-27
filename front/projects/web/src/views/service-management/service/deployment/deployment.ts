import { NewServiceDeploymentInfo, ServiceDeploymentInfo } from "@/types"
import { map, omit, omitBy } from "lodash-es"
import { ServiceDeployMode, ServiceDeployType } from "@/models"

export interface ServiceValidDeploymentInfo extends ServiceDeploymentInfo {
  hasPlacements: boolean
  validPlacements: Array<{ id: string; mode: string; key: string; value: string; isEqual: boolean }>
  enableTimeout: boolean
  enableSchedules: boolean
}

export function NewValidDeploymentInfo(info?: ServiceDeploymentInfo): ServiceValidDeploymentInfo {
  const validData = omitBy(info, (value, key) => value === undefined || value === null)
  const deploymentInfo = Object.assign({}, NewServiceDeploymentInfo(), validData)
  return {
    ...deploymentInfo,
    hasPlacements: deploymentInfo.placements.length > 0,
    enableSchedules: !deploymentInfo.manualExec && deploymentInfo.schedule.rules.length > 0,
    enableTimeout: !!deploymentInfo.schedule.timeout,
    validPlacements: map(deploymentInfo.placements, x => ({
      id: GenerateUUID(),
      mode: x.mode,
      key: x.mode === ServicePlacementMode.PlacementModeIP ? "IP" : x.key,
      value: x.value,
      isEqual: x.isEqual
    }))
  }
}

export function ParseDeploymentInfo(info: ServiceValidDeploymentInfo): ServiceDeploymentInfo {
  return {
    type: info.schedule.rules.length > 0 ? ServiceDeployType.DeployTypeSchedule : ServiceDeployType.DeployTypeBackground,
    mode: info.mode,
    replicas: info.mode === ServiceDeployMode.DeployModeGlobal ? 1 : info.replicas,
    placements: info.hasPlacements && info.validPlacements.length > 0 ? map(info.validPlacements, x => omit(x, ["id"])) : [],
    schedule: {
      timeout: !info.manualExec && info.enableSchedules && info.enableTimeout && info.schedule.rules.length > 0 ? info.schedule.timeout : "",
      rules: !info.manualExec && info.enableSchedules ? info.schedule.rules : []
    },
    manualExec: info.manualExec
  }
}
