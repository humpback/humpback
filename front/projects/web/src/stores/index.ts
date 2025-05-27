const piniaStore = createPinia()

//可以使用piniaStore.use()添加中间件
export default piniaStore

export function disposeStore() {
  useUserStore().clearUserInfo()
  useStateStore().$dispose()
  useRegistryStore().$dispose()
}

export function initStore() {}
