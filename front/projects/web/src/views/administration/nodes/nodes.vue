<script lang="ts" setup>
import { NodeInfo } from "@/types"
import { BytesToGB, TableHeight } from "@/utils"
import { Action, NodeStatus } from "@/models"
import NodeAdd from "./node-add.vue"
import NodeDelete from "./node-delete.vue"
import NodeEnable from "./node-enable.vue"
import NodeEditLabel from "./node-edit-label.vue"
import NodeViewCommand from "./node-view-command.vue"
import { QueryNodesInfo, statusOptions } from "./common.ts"

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const tableHeight = computed(() => TableHeight(286))

const isLoading = ref(false)
const queryInfo = ref<QueryNodesInfo>(new QueryNodesInfo(route.query, []))

const tableList = ref({
  total: 0,
  data: [] as Array<NodeInfo>
})

const addRef = useTemplateRef<InstanceType<typeof NodeAdd>>("addRef")
const editLabelRef = useTemplateRef<InstanceType<typeof NodeEditLabel>>("editLabelRef")
const viewValueRef = useTemplateRef<InstanceType<typeof NodeViewCommand>>("viewValueRef")
const enableRef = useTemplateRef<InstanceType<typeof NodeEnable>>("enableRef")
const deleteRef = useTemplateRef<InstanceType<typeof NodeDelete>>("deleteRef")

async function search() {
  await router.replace(queryInfo.value.urlQuery())
  isLoading.value = true
  return await nodeService
    .query(queryInfo.value.searchParams())
    .then(res => {
      tableList.value.data = res.list
      tableList.value.total = res.total
    })
    .finally(() => (isLoading.value = false))
}

function openAction(action: string, info?: NodeInfo) {
  switch (action) {
    case Action.Add:
      addRef.value?.open()
      break
    case Action.EditLabel:
      editLabelRef.value?.open(info!)
      break
    case Action.View:
      viewValueRef.value?.open(info!)
      break
    case Action.Delete:
      deleteRef.value?.open(info!)
      break
    case Action.Enable:
      enableRef.value?.open(info!)
  }
}

onMounted(() => search())
</script>

<template>
  <div>
    <v-card>
      <v-page-title :title="t('label.nodes')" />

      <v-search
        v-model="queryInfo.keywords"
        :add-label="t('btn.addNodes')"
        :placeholder="t('placeholder.enterIpHostNameOrLabel')"
        @add="openAction(Action.Add)"
        @search="search">
        <template #prefix>
          <div style="width: 220px">
            <v-select
              v-model="queryInfo.filter.status"
              :out-label="t('label.status')"
              :placeholder="t('placeholder.all')"
              clearable
              out-label-width="80px"
              @change="search">
              <el-option v-for="(item, index) in statusOptions" :key="index" :label="t(item.label)" :value="item.value">
                <el-text :type="item.type">{{ t(item.label) }}</el-text>
              </el-option>
            </v-select>
          </div>
        </template>
      </v-search>

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
              <v-router-link
                :href="`/ws/node/${scope.row.nodeId}/detail`"
                :text="scope.row.ipAddress"
                :type="scope.row.isEnable ? 'primary' : 'info'"
                disabled />
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('label.port')" min-width="80" prop="port">
          <template #default="scope">
            <span v-if="scope.row.isEnable && scope.row.status === NodeStatus.Online">{{ scope.row.port }}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('label.hostname')" min-width="180" prop="name" sortable="custom">
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
            <div class="custom-column">
              <v-label-table-view :labels="scope.row.labels" :line="4" />
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('label.action')" align="center" fixed="right" header-align="center" width="130">
          <template #default="scope">
            <el-button link type="primary" @click="openAction(Action.View, scope.row)">
              {{ t("btn.command") }}
            </el-button>
            <el-dropdown class="ml-1" placement="bottom-end" @command="openAction($event, scope.row)">
              <el-button link type="primary">
                <el-icon :size="20">
                  <IconMdiMoreHoriz />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="Action.EditLabel">
                    <el-link :underline="false" type="primary">{{ t("btn.editLabel") }}</el-link>
                  </el-dropdown-item>
                  <el-dropdown-item :command="Action.Enable">
                    <el-link :type="scope.row.isEnable ? 'info' : 'success'" :underline="false">
                      {{ scope.row.isEnable ? t("btn.disable") : t("btn.enable") }}
                    </el-link>
                  </el-dropdown-item>
                  <el-dropdown-item :command="Action.Delete">
                    <el-link :underline="false" type="danger">{{ t("btn.delete") }}</el-link>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </v-table>
    </v-card>

    <node-add ref="addRef" @refresh="search()" />

    <node-view-command ref="viewValueRef" />

    <node-edit-label ref="editLabelRef" @refresh="search()" />

    <node-enable ref="enableRef" @refresh="search()" />

    <node-delete ref="deleteRef" @refresh="search()" />
  </div>
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
      max-width: 600px;
    }
  }

  .status-tag {
    width: 120px;
    text-align: right;
    padding-right: 20px;
    padding-top: 8px;
  }
}
</style>
