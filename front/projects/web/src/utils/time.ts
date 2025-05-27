import dayjs from "dayjs"

//防抖函数
export function Debounce<T extends (...args: any[]) => void>(func: T, wait: number): T {
  let timeout: ReturnType<typeof setTimeout>

  return function (...args: Parameters<T>) {
    clearTimeout(timeout)
    timeout = setTimeout(() => func(...args), wait)
  } as T
}

export function TimestampToTime(timestamp?: number, formatType?: number): string {
  if (!timestamp) {
    return "--"
  }
  const date = dayjs(timestamp)
  switch (formatType) {
    case 1:
      return date.format("HH:mm")
    case 2:
      return date.format("HH:mm:ss")
    case 3:
      return date.format("MM-DD")
    case 4:
      return date.format("YYYY-MM-DD")
    case 5:
      return date.format("YYYY-MM-DD HH:mm")
    case 6:
      return date.format("MM/DD HH:mm:ss")
    case 7:
      return date.format("YYYY-MM-DD HH:mm")
    default:
      return date.format("YYYY-MM-DD HH:mm:ss")
  }
}
