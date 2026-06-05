<template>
  <div class="data-box">
    <div class="chart-card sell-histogram">
      <div class="chart-head">
        <div>
          <span class="data-kicker">Payment Trend</span>
          <span class="data-title">收款趋势</span>
        </div>
        <span class="data-subtext">单位: 天</span>
      </div>

      <div class="chart-body">
        <HistogramChart :home-data="props.homeData" />
      </div>
    </div>
    <div class="chart-card daily-analysis">
      <div class="chart-head">
        <div>
          <span class="data-kicker">Token Mix</span>
          <span class="data-title">交易占比</span>
        </div>
      </div>

      <div class="chart-body">
        <AnalysisChart :home-data="props.homeData" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import HistogramChart from "@/views/home/components/histogram-chart.vue";
import AnalysisChart from "@/views/home/components/analysis-chart.vue";

const props = defineProps<{
  homeData: any;
}>();
</script>

<style lang="scss" scoped>
.data-box {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 18px;
  align-items: stretch;
}

.chart-card {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 380px;
  padding: 18px;
  background: rgba(255, 255, 255, 88%);
  border: 1px solid rgba(229, 231, 235, 92%);
  border-radius: 22px;
  box-shadow: 0 16px 42px rgba(17, 24, 39, 5%);
}

.chart-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.chart-body {
  flex: 1;
  min-width: 0;
  min-height: 0;
}

.data-kicker {
  display: block;
  color: #98a2b3;
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.data-title {
  display: block;
  margin-top: 4px;
  color: #111827;
  font-size: 18px;
  font-weight: 900;
}

.data-subtext {
  padding: 6px 10px;
  color: #667085;
  font-size: 12px;
  font-weight: 800;
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 999px;
}

@media screen and (min-width: 1280px) {
  .data-box {
    grid-template-columns: minmax(0, 1.08fr) minmax(0, 0.92fr);
    gap: 16px;
  }
}

@media screen and (min-width: 1600px) {
  .data-box {
    grid-template-columns: minmax(0, 1.16fr) minmax(420px, 0.84fr);
  }
}

@media (max-width: 768px) {
  .chart-card {
    height: 360px;
    padding: 14px;
  }
}
</style>
