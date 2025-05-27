import { GenerateUUID } from "@/utils"

export enum StorageEventType {
  Login = "login",
  Logout = "logout"
}

// 用于登录用户变更的storage事件处理
class StorageEventBus {
  channelKey = "HUMPBACK_EVENT_CHANGE"
  code = GenerateUUID()
  bc: BroadcastChannel = new BroadcastChannel(this.channelKey)

  SetMessageHandler(handleFunc: (data: any) => void) {
    const tempCode = this.code
    this.bc.onmessage = function (e) {
      if (e.data?.code !== tempCode) {
        handleFunc(e.data)
      }
    }
  }

  SendMessage(type: StorageEventType, value?: any) {
    this.bc.postMessage({ code: this.code, type: type, value: value })
  }

  Close() {
    this.bc.close()
  }
}

export const storageEventBus = new StorageEventBus()

// 用于http err的事件处理
type EventFunction = (...args: any[]) => void
type EventName = "API:NO_AUTH" | "API:FAILED" | "API:RESOURCE_NOT_EXIST" | "API:NO_PERMISSION"

class EventEmitter {
  private listeners: Record<EventName, Set<EventFunction>> = {
    "API:NO_AUTH": new Set(),
    "API:FAILED": new Set(),
    "API:RESOURCE_NOT_EXIST": new Set(),
    "API:NO_PERMISSION": new Set()
  }

  on(eventName: EventName, listener: EventFunction) {
    this.listeners[eventName].add(listener)
  }

  emit(eventName: EventName, ...args: any[]) {
    this.listeners[eventName].forEach(listener => listener(...args))
  }
}

export const eventEmitter = new EventEmitter()
