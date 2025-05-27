import { type App } from "vue"
import { createI18n, I18nOptions } from "vue-i18n"
import { find, map } from "lodash-es"

export interface I18nOption {
  value: string
  name: string
  components: any
  disabled: boolean
}

const LAST_SELECTED_LANGUAGE = "HUMPBACK_LANGUAGE"
const defaultLanguage = "en-US"

function getLastLanguage(options: I18nOption[]) {
  let lang = localStorage.getItem(LAST_SELECTED_LANGUAGE)
  if (!lang) {
    lang = navigator.language
  }
  const langInfo = options.find(v => !v.disabled && v.value.toLowerCase() === lang!.trim().toLowerCase())
  return langInfo ? langInfo.value : defaultLanguage
}

function readLanguageFiles() {
  const localeInfo = {
    options: <I18nOption[]>[],
    messages: <any>{}
  }
  const fileModules = import.meta.glob("./**/*.ts", { eager: true })
  Object.keys(fileModules).forEach(item => {
    const content: any = (fileModules[item] as any).default
    const paths = item.split("/")
    let name = paths[1]
    if (paths.length > 2) {
      if (paths[paths.length - 1].toLowerCase() === "index.ts") {
        localeInfo.options.push(Object.assign({ value: name }, content))
        if (!localeInfo.messages[name]) {
          localeInfo.messages[name] = {}
        }
        return
      }
      if (localeInfo.messages[name]) {
        Object.assign(localeInfo.messages[name], content)
      } else {
        localeInfo.messages[name] = content
      }
      return
    }
    name = paths[1].replace(".ts", "")
    const option: any = { value: name }
    Object.assign(option, content)
    delete option.messages
    localeInfo.options.push(option)
    localeInfo.messages[name] = content.messages
  })
  return localeInfo
}

export function ChangeLanguage(language: string) {
  i18n.global.locale.value = language
  localStorage.setItem(LAST_SELECTED_LANGUAGE, language)
}

export function GetUILocale(uiName: string) {
  const languageInfo = find(languageOptions, x => x.value === GetCurrentLocale())
  if (languageInfo && languageInfo.components[uiName]) {
    return languageInfo.components[uiName]
  }
  return undefined
}

export function GetCurrentLanguageName() {
  const languageInfo = find(languageOptions, x => x.value === GetCurrentLocale())
  return languageInfo ? languageInfo.name : "English"
}

export function GetCurrentLocale() {
  return (i18n.global.locale as any).value as string
}

export function GetI18nMessage(key: string, params?: any) {
  return i18n.global.t(key, params)
}

export function GetLanguageOptions() {
  return map(languageOptions, x => {
    return {
      name: x.name,
      value: x.value,
      disabled: x.disabled
    }
  })
}

const localeInfo = readLanguageFiles()
const lastLanguage = getLastLanguage(localeInfo.options)
const languageOptions: I18nOption[] = localeInfo.options
const i18n = createI18n({
  warnHtmlMessage: false,
  legacy: false,
  locale: lastLanguage,
  fallbackLocale: defaultLanguage,
  messages: localeInfo.messages as I18nOptions["messages"]
})

export default {
  install(app: App<any>) {
    app.use(i18n)
  }
}
