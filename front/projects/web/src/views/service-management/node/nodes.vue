<script lang="ts" setup>
import { NodeInfo } from "@/types"
import { BytesToGB, SetWebTitle, TableHeight } from "@/utils"
import { Action } from "@/models"
import NodeAdd from "./node-add.vue"
import NodeRemove from "./node-remove.vue"
import { QueryGroupNodesInfo } from "./common.ts"
import { serviceService } from "services/service-service.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const stateStore = useStateStore()

const tableHeight = computed(() => TableHeight(362))

const groupId = ref(route.params?.groupId as string)

const isLoading = ref(false)
const queryInfo = ref<QueryGroupNodesInfo>(new QueryGroupNodesInfo(route.query))

const tableList = ref({
  total: 0,
  data: [] as Array<NodeInfo>
})

const addRef = useTemplateRef<InstanceType<typeof NodeAdd>>("addRef")
const deleteRef = useTemplateRef<InstanceType<typeof NodeRemove>>("deleteRef")

async function getGroupInfo() {
  return await groupService.info(groupId.value).then(info => {
    stateStore.setGroup(groupId.value, info)
  })
}

async function getServiceTotal() {
  return await serviceService.total(groupId.value).then(total => {
    stateStore.setGroupTotal(groupId.value, total)
  })
}

async function getGroupNodes() {
  return await groupService.queryNodes(groupId.value, queryInfo.value.searchParams()).then(res => {
    tableList.value.data = res.list
    tableList.value.total = res.total
  })
}

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  await Promise.all([getGroupInfo(), getServiceTotal(), getGroupNodes()]).finally(() => (isLoading.value = false))
}

function openAction(action: string, info?: NodeInfo) {
  switch (action) {
    case Action.Add:
      addRef.value?.open()
      break
    case Action.Delete:
      deleteRef.value?.open(info!)
      break
  }
}

onMounted(async () => {
  await search()
  SetWebTitle(`${t("webTitle.nodes")} - ${stateStore.getGroup()?.groupName}`)
})
</script>

<template>
  <v-search
    v-model="queryInfo.keywords"
    :add-label="t('btn.addNodes')"
    :placeholder="t('placeholder.enterIpHostNameOrLabel')"
    @add="openAction(Action.Add)"
    @search="search" />

  <v-table
    v-loading="isLoading"
    v-model:page-info="queryInfo.pageInfo"
    v-model:sort-info="queryInfo.sortInfo"
    :data="tableList.data"
    :max-height="tableHeight"
    :total="tableList.total"
    @page-change="search"
    @sort-change="search">
    <el-table-column :label="t('label.ip')" fixed="left" min-width="160" prop="ipAddress" sortable="custom">
      <template #default="scope">
        <div class="d-flex gap-2">
          <v-node-enable-tag :enabled="scope.row.isEnable" />
          <v-router-link :text="scope.row.ipAddress" :type="scope.row.isEnable ? 'primary' : 'info'" href="" />
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.hostname')" min-width="200" prop="name" sortable="custom">
      <template #default="scope">
        <v-table-column-none :text="scope.row.name" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.status')" min-width="400">
      <template #default="scope">
        <div class="custom-column">
          <div class="status">
            <div class="status-content">
              <div class="status-cpu">
                <el-text size="small" type="info">
                  <strong>{{ t("label.cpu") }}</strong>
                  <div>{{ scope.row.cpu || "--/--" }} {{ t("label.core") }}</div>
                </el-text>
              </div>
              <div class="status-memory">
                <el-text size="small" type="info">
                  <strong>{{ t("label.memoryUsed") }}</strong>
                  <el-progress v-if="scope.row.memoryTotal" :percentage="scope.row.memoryUsage" />
                  <div v-if="scope.row.memoryTotal">
                    {{ `${BytesToGB(scope.row.memoryUsed)} ${t("label.gb")} / ${BytesToGB(scope.row.memoryTotal)} ${t("label.gb")}` }}
                  </div>
                  <div v-else>
                    {{ `--/-- ${t("label.gb")}` }}
                  </div>
                </el-text>
              </div>
            </div>
            <div class="status-tag">
              <v-node-status-tag v-if="scope.row.isEnable" :status="scope.row.status" />
            </div>
          </div>
        </div>
      </template>
    </el-table-column>
    <el-table-column :label="t('label.labels')" min-width="240">
      <template #default="scope">
        <v-label-table-view :labels="scope.row.labels" :line="4" />
      </template>
    </el-table-column>
    <el-table-column :label="t('label.action')" align="center" fixed="right" header-align="center" width="130">
      <template #default="scope">
        <el-button link type="danger" @click="openAction(Action.Delete, scope.row)">
          {{ t("btn.remove") }}
        </el-button>
      </template>
    </el-table-column>
  </v-table>

  <node-add ref="addRef" @refresh="search()" />

  <node-remove ref="deleteRef" @refresh="search()" />
</template>

<style lang="scss" scoped>
.custom-column {
  min-height: 80px;
  display: flex;
  align-items: center;
}

.status {
  display: flex;
  align-items: start;
  gap: 20px;
  width: 100%;

  .status-content {
    display: flex;
    align-items: start;
    gap: 20px;
    flex: 1;

    .status-cpu {
      flex: 3;
      min-width: 100px;
    }

    .status-memory {
      flex: 7;
      min-width: 180px;
    }
  }

  .status-tag {
    width: 100px;
    text-align: right;
    padding-right: 20px;
    padding-top: 8px;
  }
}
</style>
