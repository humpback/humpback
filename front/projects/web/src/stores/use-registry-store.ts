export default defineStore("registries", () => {
  const registries = ref<RegistryInfo[]>([])

  async function refreshRegistries() {
    return await registryService.list().then(list => {
      registries.value = list
    })
  }

  return {
    registries,
    refreshRegistries
  }
})
