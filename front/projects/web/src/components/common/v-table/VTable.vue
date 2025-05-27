<script lang="ts" setup>
import type { Sort, TableInstance, TableProps } from "element-plus"
import { cloneDeep, omit } from "lodash-es"

type Props = Partial<
  TableProps<any> & {
    pageInfo: PageInfo
    total: number
    sortInfo: SortInfo
    pageLayout: string
    minHeight?: string
    hideHeaderBgColor?: boolean
  }
>

const props = withDefaults(defineProps<Props>(), {
  stripe: false,
  fit: true,
  showHeader: true,
  selectOnIndeterminate: true,
  allowDragLastColumn: true,
  showOverflowTooltip: true,
  scrollbarAlwaysOn: true,
  tooltipOptions: () => {
    return { placement: "top-start" }
  },
  minHeight: "396px",
  headerCellClassName: "table-header"
})

const emits = defineEmits<{
  (e: "update:pageInfo", pageInfo: PageInfo): void
  (e: "update:sortInfo", sortInfo: SortInfo): void
  (e: "selection-change", selectedData: any[]): void
  (e: "pageChange"): void
  (e: "sortChange"): void
}>()

const slots = useSlots()
const pageStore = usePageStore()

const pageSizeOptions = [10, 20, 30, 50, 100]

const tableRef = useTemplateRef<TableInstance>("tableRef")

const defaultSort = ref<Sort | undefined>(
  props.sortInfo
    ? {
        prop: props.sortInfo?.field || "",
        order: props.sortInfo?.order === "desc" ? "descending" : "ascending"
      }
    : undefined
)

const tableAttrs = computed(() => {
  const attrs: any = cloneDeep(omit(props, ["pageInfo", "total", "sortInfo", "total", "pageLayout", "minHeight"]))
  attrs.defaultSort = defaultSort.value

  return Object.keys(attrs).reduce((acc, key) => {
    if (typeof attrs[key] !== "undefined") {
      acc[key] = props[key]
    }
    return acc
  }, {})
})

const pageLayout = computed(() => {
  if (props.pageLayout) {
    return props.pageLayout
  }
  return pageStore.isSmallScreen ? "total, prev, pager, next" : "total, sizes, prev, pager, next"
})

function pageChangeEvent(index: number, size: number) {
  if (props.pageInfo) {
    emits("update:pageInfo", { index: index, size: size } as PageInfo)
    emits("pageChange")
  }
}

function sortChangeEvent(sort: any) {
  if (props.sortInfo) {
    emits("update:sortInfo", { field: sort.prop, order: sort.order === "descending" ? "desc" : "asc" } as SortInfo)
    emits("sortChange")
  }
}

function selectionChangeEvent(selectedData: any[]) {
  emits("selection-change", selectedData)
}

function clearSelection() {
  tableRef.value?.clearSelection()
}

function toggleRowExpansion(row: any, expanded?: boolean) {
  tableRef.value?.toggleRowExpansion(row, expanded)
}

defineExpose({ clearSelection, toggleRowExpansion })
</script>

<template>
  <div>
    <el-table
      ref="tableRef"
      :class-name="!props.hideHeaderBgColor ? 'header-bg-color' : undefined"
      :style="{ minHeight: props.minHeight || undefined }"
      v-bind="tableAttrs"
      @select="selectionChangeEvent"
      @sort-change="sortChangeEvent"
      @select-all="selectionChangeEvent">
      <template v-if="!!slots.default" #default>
        <slot name="default" />
      </template>
      <template v-if="!!slots.append" #append>
        <slot name="append" />
      </template>
      <template #empty>
        <slot v-if="!!slots.empty" name="empty" />
        <el-empty v-else :image-size="180" />
      </template>
    </el-table>
    <div v-if="props.pageInfo" class="mt-5 pagination">
      <el-pagination
        :background="true"
        :current-page="props.pageInfo.index"
        :layout="pageLayout"
        :page-size="props.pageInfo.size"
        :page-sizes="pageSizeOptions"
        :pager-count="pageStore.isSmallScreen ? 3 : 5"
        :total="props.total"
        @currentChange="pageChangeEvent($event, props.pageInfo.size)"
        @size-change="pageChangeEvent(props.pageInfo.index, $event)" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.header-bg-color {
  :deep(.table-header) {
    background-color: var(--hp-table-header-bg-color);
  }

  :deep(.el-table-fixed-column--left.table-header) {
    background-color: var(--hp-table-header-bg-color);
  }

  :deep(.el-table-fixed-column--right.table-header) {
    background-color: var(--hp-table-header-bg-color);
  }
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: right;
}
</style>
