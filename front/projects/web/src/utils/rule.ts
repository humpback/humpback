import { GetI18nMessage } from "@/locales"
import { RuleFormat } from "@/models"

export function IsEmpty(value: any): boolean {
  return !value || !(value as string).trim()
}

export function IsValidEmail(email: string): boolean {
  return new RegExp(RuleFormat.Email).test(email)
}

// --------------------通用规则定义-----------------------

export function RuleCannotBeEmpty(rule: any, value: any, callback: any) {
  return IsEmpty(value) ? callback(new Error(GetI18nMessage("rules.cannotBeEmpty"))) : callback()
}

export function RulePleaseEnter(fieldI18nName: string) {
  return (rule: any, value: any, callback: any) => {
    return IsEmpty(value) ? callback(new Error(`${GetI18nMessage("rules.pleaseEnter")} ${GetI18nMessage(fieldI18nName)}`)) : callback()
  }
}

export function RuleIsRequired(fieldI18nName: string) {
  return (rule: any, value: any, callback: any) => {
    return IsEmpty(value) ? callback(new Error(`${GetI18nMessage(fieldI18nName)} ${GetI18nMessage("rules.isRequired")}`)) : callback()
  }
}

export function RuleLimitMax(max: number) {
  return (rule: any, value: any, callback: any) => {
    return (value as string).length > max ? callback(new Error(GetI18nMessage("rules.limitLengthMax", { max: max }))) : callback()
  }
}

export function RuleLimitRange(min: number, max: number) {
  return (rule: any, value: any, callback: any) => {
    const len = (value as string).length
    if (len < min || len > max) {
      return callback(new Error(GetI18nMessage("rules.limitLengthRange", { min: min, max: max })))
    }
    callback()
  }
}

// --------------------特定规则定义-----------------------

export function RuleFormatErrPhone(isRequired?: boolean) {
  return (rule: any, value: any, callback: any) => {
    if (!isRequired && IsEmpty(value)) {
      return callback()
    }
    return new RegExp(RuleFormat.Phone).test(value) ? callback() : callback(new Error(GetI18nMessage("rules.formatErrPhone")))
  }
}

export function RuleFormatErrEmailOption(isRequired?: boolean) {
  return (rule: any, value: any, callback: any) => {
    if (!isRequired && IsEmpty(value)) {
      return callback()
    }
    return IsValidEmail(value) ? callback() : callback(new Error(GetI18nMessage("rules.formatErrEmail")))
  }
}
