<script lang="ts" setup>
import { BytesToGB } from "@/utils"
import { filter, find, findIndex, map } from "lodash-es"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const stateStore = useStateStore()

const isLoading = ref(false)
const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  keywords: "",
  selectedNodes: [] as string[],
  nodes: [] as NodeInfo[]
})

const validNodes = computed(() =>
  filter(dialogInfo.value.nodes, (x: NodeInfo) => {
    if (find(stateStore.getGroup()?.nodes || [], nodeId => nodeId === x.nodeId)) {
      return false
    }
    return x.isEnable && (IncludesIgnoreCase(x.ipAddress, dialogInfo.value.keywords) || IncludesIgnoreCase(x.name, dialogInfo.value.keywords))
  })
)

function disableSelect(ip: string) {
  return !!find(dialogInfo.value.selectedNodes, i => i === ip)
}

function selectNode(ip: string) {
  if (!find(dialogInfo.value.selectedNodes, x => x === ip)) {
    dialogInfo.value.selectedNodes.push(ip)
  }
}

function removeNode(ip: string) {
  const index = findIndex(dialogInfo.value.selectedNodes, x => x === ip)
  if (index !== -1) {
    dialogInfo.value.selectedNodes.splice(index, 1)
  }
}

function open() {
  dialogInfo.value = {
    show: true,
    keywords: "",
    selectedNodes: [] as string[],
    nodes: [] as NodeInfo[]
  }
  getNodeList()
}

async function getNodeList() {
  isLoading.value = true
  return await nodeService
    .list()
    .then(data => {
      dialogInfo.value.nodes = data
      const selectedNodeList: NodeInfo[] = filter(
        dialogInfo.value.nodes,
        (x: NodeInfo) => x.isEnable && !!find(dialogInfo.value.selectedNodes, ip => ip === x.ipAddress)
      )
      dialogInfo.value.selectedNodes = map(selectedNodeList, (n: NodeInfo) => n.ipAddress)
    })
    .finally(() => (isLoading.value = false))
}

async function save() {
  if (dialogInfo.value.selectedNodes.length === 0) {
    return
  }
  const nodes = filter(
    map(dialogInfo.value.selectedNodes, ip => {
      const info = find(dialogInfo.value.nodes, x => x.ipAddress === ip)
      return info?.nodeId || ""
    }),
    x => !!x
  )
  isAction.value = true
  return await groupService
    .updateNodes(stateStore.getGroup()!.groupId, {
      nodes: nodes,
      isDelete: false
    })
    .then(() => {
      ShowSuccessMsg(t("message.addSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => {
      isAction.value = false
    })
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="800px">
    <template #header>{{ t("header.addNodes") }}</template>

    <div class="d-flex gap-2">
      <v-input v-model="dialogInfo.keywords" :placeholder="t('placeholder.enterIpOrHostname')" class="filter-input">
        <template #prepend>
          <el-icon :size="18">
            <IconMdiSearch />
          </el-icon>
        </template>
      </v-input>
      <el-button plain type="primary" @click="getNodeList()">{{ t("btn.refresh") }}</el-button>
    </div>
    <v-table ref="tableRef" v-loading="isLoading" :data="validNodes" :max-height="400" class="mt-3">
      <el-table-column :label="t('label.ip')" fixed="left" prop="ipAddress" sortable width="160">
        <template #default="scope">
          <el-text type="primary">{{ scope.row.ipAddress }}</el-text>
        </template>
      </el-table-column>
      <el-table-column :label="t('label.hostname')" fixed="left" min-width="140" prop="name" sortable>
        <template #default="scope">{{ scope.row.name || "--" }}</template>
      </el-table-column>
      <el-table-column :label="t('label.cpu')" fixed="left" width="100">
        <template #default="scope">
          <div v-if="scope.row.cpu">{{ `${scope.row.cpu} ${t("label.core")}` }}</div>
          <span v-else>--</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('label.memory')" width="160">
        <template #default="scope">
          <div v-if="scope.row.memoryTotal">
            {{ `${BytesToGB(scope.row.memoryUsed)} ${t("label.gb")} / ${BytesToGB(scope.row.memoryTotal)} ${t("label.gb")}` }}
          </div>
          <span v-else> -- </span>
        </template>
      </el-table-column>
      <el-table-column :label="t('label.status')" width="100">
        <template #default="scope">
          <v-node-status-tag v-if="scope.row.isEnable" :status="scope.row.status" size="small" />
        </template>
      </el-table-column>
      <el-table-column :label="t('label.action')" width="100">
        <template #default="scope">
          <el-button
            :disabled="disableSelect(scope.row.ipAddress)"
            :type="disableSelect(scope.row.ipAddress) ? 'info' : 'primary'"
            link
            @click="selectNode(scope.row.ipAddress)">
            {{ t("btn.select") }}
          </el-button>
        </template>
      </el-table-column>
    </v-table>

    <div class="mt-5 mb-3">
      <div class="mb-3 f-bold">{{ t("label.selectedNodes") }}</div>
      <el-input-tag
        :model-value="dialogInfo.selectedNodes"
        clearable
        tag-effect="dark"
        tag-type="info"
        @clear="dialogInfo.selectedNodes = []"
        @remove-tag="removeNode" />
    </div>
    <v-alert v-if="dialogInfo.selectedNodes.length === 0" type="warning">{{ t("tips.selectLastOneNode") }}</v-alert>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="dialogInfo.selectedNodes.length === 0" :loading="isAction" type="primary" @click="save">
        {{ t("btn.save") }}
      </el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.filter-input {
  :deep(.el-input-group__prepend) {
    padding: 0 10px;
  }
}
</style>
