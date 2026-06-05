<template>
  <div class="shortcut-box">
    <div class="box-title">
      <div>
        <span class="eyebrow">Dashboard KPI</span>
        <strong>数据汇总</strong>
      </div>
      <span class="box-note">订单、金额与通知健康度</span>
    </div>
    <a-grid class="finance-card" :cols="{ xs: 1, sm: 2, lg: 3, xl: 5 }" :col-gap="16" :row-gap="16">
      <a-grid-item v-for="(item, index) in financeData" :key="index">
        <a-card hoverable class="finance-a-card" :class="'animated-fade-up-' + index" :body-style="{ padding: '0' }">
          <div class="finance-card-inner" :style="{ '--accent': item.color }">
            <div class="finance-nav">
              <div class="tag-dot"></div>
              <span class="finance-nav-title">{{ item.title }}</span>
            </div>
            <div class="finance-value">{{ item.value }}</div>
            <div class="finance-sub">
              <span>{{ item.subLabel || "当前范围" }}</span>
              <span>{{ item.subValue }}</span>
            </div>
          </div>
        </a-card>
      </a-grid-item>
    </a-grid>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  homeData: any;
}>();

const formatAmount = (value: any) => Number(value || 0).toFixed(2);

const formatPeriod = (from?: string, to?: string) => {
  if (!from || !to) return "--";
  return `${from.slice(5, 10).replace("-", "/")} - ${to.slice(5, 10).replace("-", "/")}`;
};

const buildFinanceData = (data: any) => {
  const kpi = data?.kpi || {};
  return [
    {
      id: 1,
      title: "订单总数",
      value: kpi.orders_total || 0,
      subLabel: "已支付订单",
      subValue: kpi.orders_success || 0,
      color: "#2563EB"
    },
    {
      id: 2,
      title: "收款金额",
      value: formatAmount(kpi.gmv_paid),
      subLabel: "支付成功率",
      subValue: `${formatAmount(kpi.order_success_rate)}%`,
      color: "#10B981"
    },
    {
      id: 3,
      title: "待付订单",
      value: kpi.orders_pending || 0,
      subLabel: "确认中订单",
      subValue: kpi.orders_confirming || 0,
      color: "#F0B90B"
    },
    {
      id: 4,
      title: "失败订单",
      value: kpi.orders_failed || 0,
      subLabel: "通知失败",
      subValue: kpi.notify_failed || 0,
      color: "#EF4444"
    },
    {
      id: 5,
      title: "统计周期",
      value: formatPeriod(data?.from, data?.to),
      subLabel: "时区",
      subValue: data?.timezone || "--",
      color: "#111827"
    }
  ];
};

const financeData = ref(buildFinanceData(props.homeData));

watch(
  () => props.homeData,
  newData => {
    financeData.value = buildFinanceData(newData);
  },
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.shortcut-box {
  padding: 18px;
  background: rgba(255, 255, 255, 86%);
  border: 1px solid rgba(229, 231, 235, 92%);
  border-radius: 22px;
  box-shadow: 0 16px 42px rgba(17, 24, 39, 5%);
}

.box-title {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
  color: #111827;

  strong {
    display: block;
    margin-top: 4px;
    font-size: 18px;
    font-weight: 900;
  }
}

.eyebrow {
  color: #98a2b3;
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.box-note {
  color: #667085;
  font-size: 12px;
  font-weight: 700;
}

.finance-a-card {
  overflow: hidden;
  background: #ffffff;
  border: 1px solid rgba(229, 231, 235, 88%);
  border-radius: 18px;
  box-shadow: none;
}

.finance-card-inner {
  position: relative;
  min-height: 132px;
  padding: 16px;
  overflow: hidden;

  &::before {
    position: absolute;
    top: -56px;
    right: -46px;
    width: 118px;
    height: 118px;
    content: "";
    background: color-mix(in srgb, var(--accent) 18%, transparent);
    border-radius: 50%;
  }
}

.finance-nav {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #667085;
}

.tag-dot {
  width: 10px;
  height: 10px;
  background: var(--accent);
  border-radius: 999px;
  box-shadow: 0 0 0 5px color-mix(in srgb, var(--accent) 14%, transparent);
}

.finance-nav-title {
  font-size: 13px;
  font-weight: 850;
}

.finance-value {
  position: relative;
  margin-top: 18px;
  color: #111827;
  font-size: 26px;
  font-weight: 950;
  line-height: 32px;
  letter-spacing: -0.04em;
  word-break: break-word;
}

.finance-sub {
  position: relative;
  display: flex;
  gap: 8px;
  align-items: center;
  margin-top: 10px;
  color: #667085;
  font-size: 12px;
  font-weight: 760;
  line-height: 18px;
  white-space: nowrap;

  span:last-child {
    color: #111827;
    font-weight: 900;
  }
}

@media (max-width: 768px) {
  .box-title {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
