export interface ActivityInfo {
  activityId: string
  action: string
  description: string
  oldContent: any
  newContent: any
  operator: string
  operatorId: string
  operateAt: number
  resourceId: string
  resourceName: string
  instanceName?: string
}
