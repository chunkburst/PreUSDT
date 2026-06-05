<template>
  <div class="snow-page">
    <div class="home-page">
      <section class="dashboard-hero">
        <div>
          <div class="hero-kicker">PreUSDT Operations</div>
          <h1>支付运营总览</h1>
          <p>聚合订单、链上确认、法币汇率与通知状态，快速判断当前收款网关是否健康。</p>
        </div>
        <div class="hero-metrics">
          <div class="metric-pill">
            <span>当前法币</span>
            <strong>{{ fiat }}</strong>
          </div>
          <div class="metric-pill">
            <span>统计时区</span>
            <strong>{{ timezone }}</strong>
          </div>
        </div>
      </section>

      <div class="dashboard-toolbar">
        <div class="fiat-selector">
          <span class="label">交易法币</span>
          <div class="fiat-options">
            <div
              v-for="item in fiatOptions"
              :key="item.value"
              class="fiat-option"
              :class="{ active: fiat === item.value }"
              @click="handleFiatChange(item.value)"
            >
              <span class="currency-symbol">{{ item.symbol }}</span>
              <span class="currency-name">{{ item.label }}</span>
            </div>
          </div>
        </div>

        <div class="range-actions">
          <a-range-picker
            v-if="range === 'custom'"
            v-model="customDates"
            format="YYYY-MM-DD"
            class="range-picker"
            @change="handleCustomDateChange"
          />
          <a-select v-model="range" class="range-select" @change="handleRangeChange">
            <a-option v-for="item in rangeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
          <a-button class="refresh-button" type="primary" :loading="loading" @click="forceRefresh">
            <template #icon><icon-refresh /></template>
            强制刷新
          </a-button>
        </div>
      </div>

      <Finance :home-data="home" />
      <DataBox :home-data="home" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Finance from "@/views/home/components/finance.vue";
import DataBox from "@/views/home/components/data-box.vue";
import { getDashboardHomeAPI } from "@/api/modules/home/index";

const fiat = ref("CNY");
const range = ref("7d");
const customDates = ref<any[]>([]);
const home = ref<any>(null);
const loading = ref(false);
const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "Asia/Shanghai";
let dashboardRetryTimer: ReturnType<typeof setTimeout> | null = null;

const fiatOptions = ref([
  { value: "CNY", label: "人民币", symbol: "¥" },
  { value: "USD", label: "美元", symbol: "$" },
  { value: "EUR", label: "欧元", symbol: "€" },
  { value: "GBP", label: "英镑", symbol: "£" },
  { value: "JPY", label: "日元", symbol: "¥" }
]);

const rangeOptions = [
  { value: "today", label: "今天" },
  { value: "7d", label: "最近 7 天" },
  { value: "30d", label: "最近 30 天" },
  { value: "custom", label: "自定义" }
];

const clearDashboardRetry = () => {
  if (dashboardRetryTimer) {
    clearTimeout(dashboardRetryTimer);
    dashboardRetryTimer = null;
  }
};

const getDashboardHome = async (force = false, retryCount = 0) => {
  if (range.value === "custom" && (!Array.isArray(customDates.value) || customDates.value.length !== 2)) return;

  if (retryCount === 0) {
    clearDashboardRetry();
  }

  loading.value = true;
  try {
    const params: any = {
      range: range.value,
      tz: timezone,
      fiat: fiat.value,
      force
    };
    if (range.value === "custom") {
      params.from = customDates.value[0];
      params.to = customDates.value[1];
    }

    const data = await getDashboardHomeAPI(params);
    if (!data?.data) {
      throw new Error("仪表盘数据为空");
    }
    home.value = data.data;
  } catch (error) {
    if (retryCount < 3) {
      dashboardRetryTimer = setTimeout(() => {
        getDashboardHome(force, retryCount + 1);
      }, (retryCount + 1) * 1000);
      return;
    }
    console.error("获取首页统计失败:", error);
  } finally {
    loading.value = false;
  }
};

const handleFiatChange = (value: string) => {
  fiat.value = value;
  getDashboardHome();
};

const handleRangeChange = () => {
  if (range.value !== "custom") {
    getDashboardHome();
  }
};

const handleCustomDateChange = () => {
  getDashboardHome();
};

const forceRefresh = () => {
  getDashboardHome(true);
};

onMounted(() => {
  getDashboardHome();
});

onUnmounted(() => {
  clearDashboardRetry();
});
</script>

<style lang="scss" scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.dashboard-hero {
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  padding: 26px;
  overflow: hidden;
  color: #ffffff;
  background:
    radial-gradient(circle at 12% 20%, rgba(240, 185, 11, 30%), transparent 26%),
    linear-gradient(135deg, #111827 0%, #1f2937 62%, #0f172a 100%);
  border: 1px solid rgba(255, 255, 255, 10%);
  border-radius: 24px;
  box-shadow: 0 22px 60px rgba(17, 24, 39, 16%);

  &::after {
    position: absolute;
    top: -72px;
    right: -44px;
    width: 220px;
    height: 220px;
    content: "";
    background: rgba(240, 185, 11, 14%);
    border: 1px solid rgba(240, 185, 11, 22%);
    border-radius: 50%;
  }

  h1 {
    margin: 8px 0 8px;
    font-size: 30px;
    font-weight: 900;
    line-height: 1.15;
    letter-spacing: -0.04em;
  }

  p {
    max-width: 680px;
    margin: 0;
    color: rgba(255, 255, 255, 72%);
    font-size: 14px;
    line-height: 1.8;
  }
}

.hero-kicker {
  color: #f0b90b;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero-metrics {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.metric-pill {
  min-width: 132px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 10%);
  border: 1px solid rgba(255, 255, 255, 12%);
  border-radius: 16px;
  backdrop-filter: blur(12px);

  span {
    display: block;
    color: rgba(255, 255, 255, 62%);
    font-size: 12px;
    font-weight: 700;
  }

  strong {
    display: block;
    margin-top: 5px;
    color: #ffffff;
    font-size: 14px;
    font-weight: 900;
  }
}

.dashboard-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  padding: 14px;
  background: rgba(255, 255, 255, 86%);
  border: 1px solid rgba(229, 231, 235, 92%);
  border-radius: 18px;
  box-shadow: 0 14px 36px rgba(17, 24, 39, 5%);
}

.fiat-selector {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;

  .label {
    color: #667085;
    font-size: 13px;
    font-weight: 850;
    white-space: nowrap;
  }

  .fiat-options {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .fiat-option {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
    min-width: 70px;
    padding: 8px 12px;
    color: #344054;
    background: #f8fafc;
    border: 1px solid #e5e7eb;
    border-radius: 999px;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      color: #181a20;
      border-color: rgba(240, 185, 11, 60%);
      box-shadow: 0 8px 18px rgba(240, 185, 11, 12%);
      transform: translateY(-1px);
    }

    &.active {
      color: #181a20;
      background: #f0b90b;
      border-color: #f0b90b;
      box-shadow: 0 10px 22px rgba(240, 185, 11, 24%);
    }

    .currency-symbol {
      font-size: 14px;
      font-weight: 950;
    }

    .currency-name {
      font-size: 12px;
      font-weight: 800;
    }
  }
}

.range-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  margin-left: auto;
}

.range-picker {
  width: 250px;
}

.range-select {
  width: 136px;
}

.refresh-button {
  color: #181a20;
  font-weight: 850;
  background: #f0b90b;
  border-color: #f0b90b;
}

@media (max-width: 900px) {
  .dashboard-hero {
    align-items: flex-start;
    flex-direction: column;
  }

  .hero-metrics {
    justify-content: flex-start;
    width: 100%;
  }

  .metric-pill {
    flex: 1;
  }
}

@media (max-width: 768px) {
  .home-page {
    gap: 14px;
  }

  .dashboard-hero {
    padding: 20px;

    h1 {
      font-size: 24px;
    }
  }

  .dashboard-toolbar,
  .fiat-selector,
  .range-actions {
    align-items: stretch;
    width: 100%;
  }

  .fiat-selector {
    flex-direction: column;
    align-items: flex-start;
  }

  .fiat-selector .fiat-options,
  .range-actions {
    width: 100%;
  }

  .fiat-option {
    flex: 1;
  }

  .range-picker,
  .range-select,
  .refresh-button {
    width: 100%;
  }
}
</style>
