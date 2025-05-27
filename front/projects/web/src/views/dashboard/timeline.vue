<script lang="ts" setup>
import * as echarts from "echarts/core"
import {
  GridComponent,
  GridComponentOption,
  LegendComponent,
  LegendComponentOption,
  ToolboxComponent,
  ToolboxComponentOption,
  TooltipComponent,
  TooltipComponentOption
} from "echarts/components"
import { BarChart, BarSeriesOption, LineChart, LineSeriesOption } from "echarts/charts"
import { CanvasRenderer } from "echarts/renderers"
import { UniversalTransition } from "echarts/features"
import { map, toLower } from "lodash-es"
import { TimestampToTime } from "@/utils"
import { StatisticsCountInfo } from "@/types"

echarts.use([GridComponent, ToolboxComponent, LegendComponent, BarChart, TooltipComponent, LineChart, CanvasRenderer, UniversalTransition])

type EChartsOption = echarts.ComposeOption<
  GridComponentOption | ToolboxComponentOption | BarSeriesOption | LineSeriesOption | TooltipComponentOption | LegendComponentOption
>

const props = defineProps<{ isLoading?: boolean }>()

const emits = defineEmits<{
  (e: "change"): void
}>()

const { t, locale } = useI18n()

const onlyMe = defineModel<boolean>()

let incrementChart: echarts.ECharts
const incrementRef = useTemplateRef<HTMLDivElement>("incrementRef")

let incrementOptions = ref<EChartsOption>({
  tooltip: {
    trigger: "axis"
  },
  legend: {
    right: "50%",
    selected: {
      "Group": true,
      "Service": true,
      "Node": true,
      "Deploy": true
    }
  },
  grid: {
    left: "1%",
    right: "1%",
    bottom: "3%",
    containLabel: true
  },
  toolbox: {
    show: true,
    feature: {
      dataView: { show: true, readOnly: false },
      magicType: { show: true, type: ["line", "bar", "stack"] },
      restore: { show: true }
    },
    right: "1%"
  },
  calculable: true,
  xAxis: [
    {
      type: "category",
      data: []
    }
  ],
  yAxis: [{ type: "value" }],
  series: [
    {
      name: t("label.group"),
      type: "bar",
      seriesLayoutBy: "column",
      data: []
    },
    {
      name: t("label.service"),
      type: "bar",
      seriesLayoutBy: "column",
      data: []
    },
    {
      name: t("label.node"),
      type: "bar",
      seriesLayoutBy: "column",
      data: []
    },
    {
      name: t("label.deploy"),
      type: "line",
      seriesLayoutBy: "column",
      data: []
    }
  ]
})

function resize() {
  incrementChart?.resize()
}

function refreshData() {
  incrementChart?.setOption(incrementOptions.value)
}

function setData(startAt: number, data: StatisticsCountInfo[]) {
  const tempResultMap: Record<
    number,
    {
      timestamp: number
      groups: number
      services: number
      nodes: number
      deploy: number
    }
  > = {}
  let tempStartAt = new Date(startAt)
  tempStartAt.setHours(0, 0, 0, 0)
  for (let i = 0; i < 30; i++) {
    const timestamp = tempStartAt.getTime()
    tempResultMap[timestamp] = { timestamp: timestamp, groups: 0, services: 0, nodes: 0, deploy: 0 }
    tempStartAt.setDate(tempStartAt.getDate() + 1)
  }
  map(data, item => {
    const date = new Date(item.createAt)
    date.setHours(0, 0, 0, 0)
    const timestamp = date.getTime()
    const tempData = tempResultMap[timestamp]
    if (tempData) {
      switch (toLower(item.type)) {
        case "group":
          tempData.groups += item.num
          break
        case "service":
          tempData.services += item.num
          break
        case "node":
          tempData.nodes += item.num
          break
        case "deploy":
          tempData.deploy += item.num
          break
      }
      tempResultMap[timestamp] = tempData
    }
  })
  const result = {
    timeStr: [] as Array<string>,
    groups: [] as Array<number>,
    services: [] as Array<number>,
    nodes: [] as Array<number>,
    deploy: [] as Array<number>
  }
  Object.values(tempResultMap).forEach(v => {
    result.timeStr.push(TimestampToTime(v.timestamp, 3))
    result.groups.push(v.groups)
    result.services.push(v.services)
    result.nodes.push(v.nodes)
    result.deploy.push(v.deploy)
  })
  incrementOptions.value!.xAxis![0].data = result.timeStr
  incrementOptions.value!.series![0].data = result.groups
  incrementOptions.value!.series![1].data = result.services
  incrementOptions.value!.series![2].data = result.nodes
  incrementOptions.value!.series![3].data = result.deploy
  refreshData()
}

watch(locale, () => {
  refreshData()
})

onMounted(() => {
  incrementChart = echarts.init(incrementRef.value)
  window.addEventListener("resize", resize)
})

onBeforeUnmount(() => {
  incrementChart?.dispose()
  window.removeEventListener("resize", resize)
})

defineExpose({
  setData
})
</script>

<template>
  <v-card>
    <div class="title">
      <span>{{ t("header.newDataOf30Days") }}</span>
      <el-radio-group v-model="onlyMe" size="small" @change="emits('change')">
        <el-radio-button :label="t('label.all')" :value="false" />
        <el-radio-button :label="t('label.self')" :value="true" />
      </el-radio-group>
    </div>
    <div ref="incrementRef" v-loading="props.isLoading" class="chart" />
  </v-card>
</template>

<style lang="scss" scoped>
.title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  padding-left: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart {
  height: 360px;
  width: 100%;
}
</style>
