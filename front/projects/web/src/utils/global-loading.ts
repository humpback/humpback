import { ElLoading } from "element-plus"
import { GetI18nMessage } from "@/locales"

class GlobalLoading {
  private loading: any

  public show(message?: string) {
    this.loading = ElLoading.service({ lock: true, text: message || "" })
  }

  public showLoading() {
    this.show(GetI18nMessage("message.loading"))
  }

  public close() {
    if (this.loading) {
      this.loading.close()
    }
    this.loading = null
  }
}

export const globalLoading = new GlobalLoading()
