import useClipboard from "vue-clipboard3"
import { GetI18nMessage } from "@/locales"

const { toClipboard } = useClipboard()

export async function CopyToClipboard(value: string) {
  try {
    await toClipboard(value)
    ShowSuccessMsg(GetI18nMessage("message.copySucceed"))
  } catch (e) {
    console.error(e)
    ShowErrMsg(GetI18nMessage("message.copyFailed"))
  }
}
