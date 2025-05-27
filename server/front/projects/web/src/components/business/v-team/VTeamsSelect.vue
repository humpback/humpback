<script lang="ts" setup>
import { TeamInfo } from "@/types"
import { map } from "lodash-es"

const props = withDefaults(
  defineProps<{
    options?: TeamInfo[]
    filterable?: boolean
    multiple?: boolean
    placeholder?: string
  }>(),
  {
    filterable: true,
    multiple: true,
    placeholder: ""
  }
)

const { t } = useI18n()
const userStore = useUserStore()

const isLoading = ref(false)
const teams = defineModel<string[]>()
const teamList = ref<TeamInfo[]>([])
const selectOptions = computed(() => map(props.options || teamList.value, x => ({ label: x.name, value: x.teamId })))

async function getTeams() {
  isLoading.value = true
  return await teamService
    .list()
    .then(res => {
      teamList.value = res
    })
    .finally(() => (isLoading.value = false))
}

onMounted(async () => {
  if (!props.options) {
    await getTeams()
  }
})
</script>

<template>
  <el-select-v2
    v-model="teams"
    :filterable="props.filterable"
    :loading-text="t('message.loading')"
    :multiple="props.multiple"
    :options="selectOptions"
    :placeholder="props.placeholder">
    <template v-if="isLoading" #prefix>
      <el-button :loading="isLoading" link />
    </template>
    <template v-if="userStore.isAdmin" #footer>
      <div class="text-align-right">
        <el-link href="/ws/user-related/teams" target="_blank" type="primary">
          <strong>{{ t("btn.goToAddTeam") }}</strong>
        </el-link>
      </div>
    </template>
  </el-select-v2>
</template>

<style lang="scss" scoped></style>
