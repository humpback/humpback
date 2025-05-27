import { GetI18nMessage } from "@/locales"

export async function CopyToClipboard(value: string) {
  await copyText(value)
    .then(() => {
      ShowSuccessMsg(GetI18nMessage("message.copySucceed"))
    })
    .catch(e => {
      ShowErrMsg(GetI18nMessage("message.copyFailed"))
    })
}

async function copyText(text: string): Promise<boolean> {
  if (navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      return fallbackCopyText(text)
    }
  } else {
    return fallbackCopyText(text)
  }
}

function fallbackCopyText(text: string): boolean {
  const textarea = document.createElement("textarea")
  textarea.value = text
  textarea.style.position = "fixed"
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  try {
    return document.execCommand("copy")
  } catch {
    return false
  } finally {
    document.body.removeChild(textarea)
  }
}
