<template>
  <div class="analysis-chart" :class="{ empty: isEmpty }">
    <a-empty v-if="isEmpty" description="暂无交易占比数据" />
    <s-echarts v-else :options="option" :theme="theme" :update-options="{ notMerge: true }" />
  </div>
</template>

<script setup lang="ts">
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import { TooltipComponent, LegendComponent } from "echarts/components";
import { useChart } from "@/hooks/useChart";

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent]);

const props = defineProps<{
  homeData: any;
}>();

const colorMap: Record<string, string> = {
  USDT: "#26A17B",
  USDC: "#2775CA",
  TRX: "#EF0027",
  BNB: "#F0B90B",
  ETH: "#627EEA"
};

const chartData = computed(() => {
  const tokenMap = props.homeData?.token_map;
  if (!tokenMap || typeof tokenMap !== "object") return [];

  const totalAmount = Object.values(tokenMap).reduce((sum: number, value: any) => {
    const numValue = Number(value);
    return sum + (Number.isNaN(numValue) ? 0 : numValue);
  }, 0);
  if (totalAmount <= 0) return [];

  return Object.entries(tokenMap)
    .map(([token, amount]) => {
      const numAmount = Number(amount);
      if (Number.isNaN(numAmount) || numAmount <= 0) return null;

      return {
        name: token,
        value: Number(((numAmount / totalAmount) * 100).toFixed(2)),
        amount: Number(numAmount.toFixed(2))
      };
    })
    .filter(Boolean) as Array<{ name: string; value: number; amount: number }>;
});

const isEmpty = computed(() => chartData.value.length === 0);

const { option, theme } = useChart((_, palette = []) => {
  const colors = chartData.value.map(item => colorMap[item.name] || palette[chartData.value.indexOf(item) % palette.length] || "#667085");

  return {
    backgroundColor: "transparent",
    color: colors,
    tooltip: {
      trigger: "item",
      borderWidth: 0,
      formatter: (params: any) => {
        const amount = params.data?.amount ?? 0;
        return `<b>${params.name}</b><br/>${params.marker} 占比 <b style="margin-left: 16px">${params.value}%</b><br/>金额 <b style="margin-left: 16px">${amount}</b>`;
      }
    },
    legend: {
      top: 0,
      left: "center",
      type: "scroll",
      orient: "horizontal",
      icon: "circle",
      itemWidth: 8,
      itemHeight: 8,
      itemGap: 12,
      textStyle: {
        color: "#667085",
        fontSize: 12,
        fontWeight: 700
      },
      data: chartData.value.map(item => item.name)
    },
    series: [
      {
        type: "pie",
        center: ["50%", "56%"],
        padAngle: 2,
        radius: ["42%", "66%"],
        avoidLabelOverlap: true,
        itemStyle: {
          borderColor: "#ffffff",
          borderWidth: 2,
          borderRadius: 6
        },
        label: {
          color: "#344054",
          fontSize: 11,
          fontWeight: 800,
          formatter: "{b}\n{c}%"
        },
        labelLine: {
          length: 8,
          length2: 6
        },
        data: chartData.value,
        emphasis: {
          scale: true,
          scaleSize: 8
        }
      }
    ]
  };
});
</script>

<style lang="scss" scoped>
.analysis-chart {
  width: 100%;
  height: 100%;
  min-height: 300px;

  &.empty {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
