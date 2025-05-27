import "./styles/_all.scss"
import { createApp } from "vue"
import App from "./App.vue"
import i18n from "@/locales"
import router from "./router"
import useUserStore from "@/stores/use-user-store.ts"
import "@vue-js-cron/element-plus/dist/element-plus.css"
import CronElementPlusPlugin from "@vue-js-cron/element-plus"
import { ElButton, ElDropdown, ElDropdownItem, ElDropdownMenu, ElIcon } from "element-plus"
// 注册组件
import { init } from "@/app/app.ts"

const app = createApp(App).use(stores).use(i18n).use(CronElementPlusPlugin)

// cron组件需要指定
app.use(ElDropdown).use(ElDropdownMenu).use(ElDropdownItem).use(ElIcon).use(ElButton)

app.config.errorHandler = (err: any, vm: any, info: any) => {
  if (err.isAxiosError) {
    return
  }
  if (err.toString().includes("Validating error")) {
    return
  }
  console.error(err)
}

init().then(() => {
  useUserStore()
    .init()
    .finally(() => {
      app.use(router).mount("#app")
    })
})
