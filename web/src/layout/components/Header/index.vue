<template>
  <a-layout-header class="header">
    <HeaderLeft />
    <HeaderRight />
  </a-layout-header>
</template>
<script setup lang="ts">
import HeaderLeft from "@/layout/components/Header/components/header-left/index.vue";
import HeaderRight from "@/layout/components/Header/components/header-right/index.vue";
</script>

<style lang="scss" scoped>
.header {
  position: relative;
  z-index: 2;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  height: 66px;
  padding: 0 20px;
  background:
    radial-gradient(circle at 18% 0%, rgba(240, 185, 11, 10%), transparent 24%),
    rgba(255, 255, 255, 90%);
  border-bottom: 1px solid rgba(229, 231, 235, 86%);
  box-shadow: 0 10px 30px rgba(17, 24, 39, 4%);
  backdrop-filter: blur(14px);
}
</style>
