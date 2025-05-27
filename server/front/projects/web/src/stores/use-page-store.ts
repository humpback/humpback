export default defineStore("page", () => {
  const screenWidth = ref(window.innerWidth)
  const screenHeight = ref(window.innerHeight)
  const menuIsCollapse = ref(false)
  const userSetMenuIsCollapse = ref(false)

  const refreshPage = ref({
    needRefresh: false,
    type: ""
  })

  const isBigScreen = computed(() => {
    return screenWidth.value > 1200
  })

  const isSmallScreen = computed(() => {
    return screenWidth.value < 800
  })

  function setScreen() {
    screenWidth.value = window.innerWidth
    screenHeight.value = window.innerHeight
    if (screenWidth.value < 1200) {
      menuIsCollapse.value = true
    } else if (!userSetMenuIsCollapse.value) {
      menuIsCollapse.value = false
    }
  }

  function changeCollapse() {
    menuIsCollapse.value = !menuIsCollapse.value
    userSetMenuIsCollapse.value = menuIsCollapse.value
  }

  return {
    screenWidth,
    screenHeight,
    isBigScreen,
    isSmallScreen,
    refreshPage,
    menuIsCollapse,
    setScreen,
    changeCollapse
  }
})
