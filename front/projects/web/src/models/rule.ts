interface RuleLengthLimit {
  Min: number
  Max: number
}

interface RuleLength {
  Username: RuleLengthLimit
  TeamName: RuleLengthLimit
  ConfigName: RuleLengthLimit
  Email: RuleLengthLimit
  Password: RuleLengthLimit
  Phone: RuleLengthLimit
  Description: RuleLengthLimit
  ConfigValue: RuleLengthLimit
  RegistryUrl: RuleLengthLimit
  RegistryUsername: RuleLengthLimit
  RegistryPassword: RuleLengthLimit
  GroupName: RuleLengthLimit
  ServiceName: RuleLengthLimit
  ImageName: RuleLengthLimit
  MemoryLimit: RuleLengthLimit
  MemoryReservation: RuleLengthLimit
  MaxCpuUsage: RuleLengthLimit
  InstanceNum: RuleLengthLimit
  LogsLine: RuleLengthLimit
}

interface RuleFormat {
  Email: string
  Phone: string
  IPAddress: string
}

export let RuleLength = {} as RuleLength
export let RuleFormat = {} as RuleFormat

export function initRule(lengthLimit: RuleLength, formatRule: RuleFormat) {
  RuleLength = lengthLimit
  RuleFormat = formatRule
}
