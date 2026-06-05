<template>
  <a-layout-footer class="footer">
    <div class="footer-title">
      <span class="secondary-info">Self-hosted crypto payment gateway</span>
      <span @click="onPreUSDT" class="primary-copyright">{{ currentYear }} © PreUSDT</span>
      <span class="secondary-info right">Open source · Multi-chain cashier</span>
    </div>
  </a-layout-footer>
</template>

<script setup lang="ts">
const currentYear = ref(new Date().getFullYear());
const onPreUSDT = () => {
  window.open("https://github.com/chunkburst/PreUSDT", "_blank");
};
</script>

<style lang="scss" scoped>
.footer {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  height: 34px;
  padding: 0 18px;
  background: rgba(255, 255, 255, 72%);
  border-top: 1px solid rgba(229, 231, 235, 80%);
}

.footer-title {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  width: 100%;
  align-items: center;
  gap: 12px;
}

.primary-copyright {
  color: #111827;
  font-weight: 850;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    color: #92400e;
    transform: translateY(-1px);
  }
}

.secondary-info {
  color: #98a2b3;
  font-size: 12px;
  font-weight: 700;
}

.right {
  text-align: right;
}
</style>
