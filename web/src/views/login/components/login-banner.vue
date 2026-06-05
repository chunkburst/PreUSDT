<template>
  <section class="banner-box">
    <div class="banner-brand">
      <div class="brand-mark">P</div>
      <div>
        <div class="brand-title">{{ title }}</div>
        <div class="brand-subtitle">Self-hosted payment gateway</div>
      </div>
    </div>

    <div class="banner-main">
      <div class="status-card">
        <div class="status-top">
          <span>Payment Operations</span>
          <strong>Live</strong>
        </div>
        <div class="status-amount">Multi-chain Cashier</div>
        <div class="status-grid">
          <span>USDT</span>
          <span>USDC</span>
          <span>TRX</span>
          <span>ETH</span>
          <span>BNB</span>
          <span>MQTT</span>
        </div>
      </div>

      <div class="copy-block">
        <p class="eyebrow">PreUSDT Admin</p>
        <h2>把链上收款、订单确认和通知回调放进一个干净后台。</h2>
        <p>管理钱包地址、支付网络、汇率同步、Telegram 模板和首页收银台资源，适合自托管部署和二次开发。</p>
      </div>
    </div>

    <div class="banner-foot">
      <span>TRON</span>
      <span>Ethereum</span>
      <span>BSC</span>
      <span>Polygon</span>
      <span>Base</span>
    </div>
  </section>
</template>
<script setup lang="ts">
const title = import.meta.env.VITE_GLOB_APP_TITLE;
</script>

<style lang="scss" scoped>
.banner-box {
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
  height: 100%;
  padding: 34px;
  overflow: hidden;
  color: #ffffff;
  background:
    radial-gradient(circle at 20% 20%, rgba(240, 185, 11, 28%), transparent 30%),
    linear-gradient(145deg, #111827 0%, #1f2937 58%, #0f172a 100%);
}

.banner-box::after {
  position: absolute;
  right: -90px;
  bottom: -90px;
  width: 260px;
  height: 260px;
  content: "";
  background: rgba(240, 185, 11, 13%);
  border-radius: 50%;
}

.banner-brand,
.banner-main,
.banner-foot {
  position: relative;
  z-index: 1;
}

.banner-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  color: #181a20;
  font-size: 20px;
  font-weight: 950;
  background: #f0b90b;
  border-radius: 16px;
  box-shadow: 0 16px 44px rgba(240, 185, 11, 24%);
}

.brand-title {
  font-size: 18px;
  font-weight: 900;
  letter-spacing: -0.04em;
}

.brand-subtitle {
  margin-top: 2px;
  color: #9ca3af;
  font-size: 12px;
  font-weight: 700;
}

.banner-main {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 28px;
  align-items: center;
}

.status-card {
  padding: 20px;
  background: rgba(255, 255, 255, 8%);
  border: 1px solid rgba(255, 255, 255, 10%);
  border-radius: 24px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 22%);
  backdrop-filter: blur(14px);
}

.status-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #9ca3af;
  font-size: 12px;
  font-weight: 800;
}

.status-top strong {
  color: #22c55e;
}

.status-amount {
  margin: 22px 0;
  font-size: 28px;
  font-weight: 950;
  line-height: 1.05;
  letter-spacing: -0.07em;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.status-grid span {
  padding: 8px 6px;
  color: #e5e7eb;
  font-size: 11px;
  font-weight: 850;
  text-align: center;
  background: rgba(255, 255, 255, 7%);
  border: 1px solid rgba(255, 255, 255, 8%);
  border-radius: 999px;
}

.copy-block .eyebrow {
  margin: 0 0 10px;
  color: #f0b90b;
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.copy-block h2 {
  max-width: 360px;
  margin: 0;
  font-size: 34px;
  font-weight: 950;
  line-height: 1.05;
  letter-spacing: -0.08em;
}

.copy-block p:last-child {
  max-width: 390px;
  margin: 16px 0 0;
  color: #cbd5e1;
  font-size: 14px;
  line-height: 1.75;
}

.banner-foot {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.banner-foot span {
  padding: 8px 10px;
  color: #d1d5db;
  font-size: 12px;
  font-weight: 850;
  background: rgba(255, 255, 255, 7%);
  border: 1px solid rgba(255, 255, 255, 8%);
  border-radius: 999px;
}
</style>
