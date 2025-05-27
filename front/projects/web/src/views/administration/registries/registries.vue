<script lang="ts" setup>
import { RegistryInfo } from "@/types"
import { TableHeight } from "@/utils"
import { Action } from "@/models"
import RegistryEdit from "./registry-edit.vue"
import RegistryDelete from "./registry-delete.vue"
import RegistryView from "./registry-view.vue"
import { isDefaultRegistry, QueryRegistryInfo } from "./common.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const tableHeight = computed(() => TableHeight(286))

const isLoading = ref(false)
const queryInfo = ref<QueryRegistryInfo>(new QueryRegistryInfo(route.query))

const tableList = ref({
  total: 0,
  data: [] as Array<RegistryInfo>
})

const registryEditRef = useTemplateRef<InstanceType<typeof RegistryEdit>>("registryEditRef")
const registryDeleteRef = useTemplateRef<InstanceType<typeof RegistryDelete>>("registryDeleteRef")
const registryViewRef = useTemplateRef<InstanceType<typeof RegistryView>>("registryViewRef")

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  return await registryService
    .query(queryInfo.value.searchParams())
    .then(res => {
      tableList.value.data = res.list
      tableList.value.total = res.total
    })
    .finally(() => (isLoading.value = false))
}

function openAction(action: string, info?: RegistryInfo) {
  switch (action) {
    case Action.Add:
    case Action.Edit:
      registryEditRef.value?.open(info)
      break
    case Action.Delete:
      registryDeleteRef.value?.open(info!)
      break
    case Action.View:
      registryViewRef.value?.open(info!)
      break
  }
}

onMounted(() => search())
</script>

<template>
  <div>
    <v-card>
      <v-page-title :title="t('label.registries')" />

      <v-search v-model="queryInfo.keywords" :add-label="t('btn.addRegistry')" :input-label="t('label.url')" @add="openAction(Action.Add)" @search="search" />

      <v-table
        v-loading="isLoading"
        v-model:page-info="queryInfo.pageInfo"
        v-model:sort-info="queryInfo.sortInfo"
        :data="tableList.data"
        :max-height="tableHeight"
        :total="tableList.total"
        @page-change="search"
        @sort-change="search">
        <el-table-column :label="t('label.url')" fixed="left" min-width="200" prop="url" sortable="custom">
          <template #default="scope">
            <div class="d-flex gap-1">
              <div class="overflow_div flex-1">
                <span>{{ scope.row.url }}</span>
              </div>
              <el-tag v-if="scope.row.isDefault" effect="dark" round size="small" type="warning">
                {{ t("label.default") }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('label.authentication')" align="left" min-width="160">
          <template #default="scope">
            <el-button v-if="scope.row?.hasAuth" link type="primary" @click="openAction(Action.View, scope.row)">
              {{ t("btn.view") }}
            </el-button>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('label.updateDate')" min-width="160" prop="updatedAt" sortable="custom">
          <template #default="scope">
            <v-date-view :timestamp="scope.row.updatedAt" />
          </template>
        </el-table-column>
        <el-table-column :label="t('label.createDate')" min-width="160" prop="createdAt" sortable="custom">
          <template #default="scope">
            <v-date-view :timestamp="scope.row.createdAt" />
          </template>
        </el-table-column>

        <el-table-column :label="t('label.action')" align="right" fixed="right" header-align="center" width="130">
          <template #default="scope">
            <el-button link type="primary" @click="openAction(Action.Edit, scope.row)">{{ t("btn.edit") }}</el-button>
            <el-button v-if="!isDefaultRegistry(scope.row.url)" link type="danger" @click="openAction(Action.Delete, scope.row)">
              {{ t("btn.delete") }}
            </el-button>
          </template>
        </el-table-column>
      </v-table>
    </v-card>

    <registry-delete ref="registryDeleteRef" @refresh="search()" />

    <registry-edit ref="registryEditRef" @refresh="search()" />

    <registry-view ref="registryViewRef" />
  </div>
</template>

<style lang="scss" scoped></style>
