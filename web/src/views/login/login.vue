<template>
  <div class="login-shell">
    <div class="login-orb login-orb-a"></div>
    <div class="login-orb login-orb-b"></div>
    <div class="login-panel">
      <LoginBanner v-if="isPc" />
      <div class="login-box">
        <div class="login-kicker">PreUSDT Console</div>
        <div class="login-title">管理后台登录</div>
        <div class="login-title-desc">多链收款、钱包地址、回调通知与订单状态的统一控制台。</div>
        <LoginForm />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import LoginBanner from "@/views/login/components/login-banner.vue";
import LoginForm from "@/views/login/components/login-form.vue";
import { useDevicesSize } from "@/hooks/useDevicesSize";
const { isPc } = useDevicesSize();
</script>

<style lang="scss" scoped>
.login-shell {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 32px;
  overflow: hidden;
  background:
    radial-gradient(circle at 18% 16%, rgba(240, 185, 11, 18%), transparent 30%),
    radial-gradient(circle at 88% 20%, rgba(5, 150, 105, 12%), transparent 28%),
    linear-gradient(135deg, #ffffff 0%, #f8fafc 48%, #eef2f7 100%);
}

.login-orb {
  position: absolute;
  border-radius: 999px;
  filter: blur(10px);
  pointer-events: none;
}

.login-orb-a {
  right: 9%;
  bottom: 12%;
  width: 180px;
  height: 180px;
  background: rgba(240, 185, 11, 18%);
}

.login-orb-b {
  top: 10%;
  left: 7%;
  width: 130px;
  height: 130px;
  background: rgba(17, 24, 39, 8%);
}

.login-panel {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr 420px;
  width: min(1080px, 100%);
  min-height: 560px;
  overflow: hidden;
  background: rgba(255, 255, 255, 86%);
  border: 1px solid rgba(229, 231, 235, 88%);
  border-radius: 32px;
  box-shadow: 0 32px 90px rgba(17, 24, 39, 14%);
  backdrop-filter: blur(18px);
}

.login-box {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 56px 44px;
  background: rgba(255, 255, 255, 72%);
}

.login-kicker {
  display: inline-flex;
  width: fit-content;
  margin-bottom: 14px;
  padding: 7px 10px;
  color: #92400e;
  font-size: 11px;
  font-weight: 850;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  background: rgba(255, 251, 235, 92%);
  border: 1px solid rgba(240, 185, 11, 36%);
  border-radius: 999px;
}

.login-title {
  margin-bottom: 10px;
  color: #111827;
  font-size: 32px;
  font-weight: 900;
  line-height: 1;
  letter-spacing: -0.06em;
}

.login-title-desc {
  max-width: 320px;
  color: #667085;
  font-size: 14px;
  line-height: 1.75;
}

@media (max-width: 768px) {
  .login-shell {
    padding: 18px;
  }

  .login-panel {
    display: block;
    min-height: auto;
    border-radius: 26px;
  }

  .login-box {
    padding: 36px 24px;
  }

  .login-title {
    font-size: 28px;
  }
}
</style>
