<template>
  <div :class="layoutType == 'layoutHead' ? 'logo-head no-border' : 'logo-head'">
    <div class="logo-box" :class="(collapsed || layoutType == 'layoutHead') && 'padding-unset'">
      <div class="logo-mark"><span>P</span></div>
      <div v-if="isTitle" class="logo-copy" :class="isDark ? 'dark' : ''">
        <span class="logo-title">{{ title }}</span>
        <span class="logo-subtitle">payment console</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
const themeStore = useThemeConfig();
const { collapsed, asideDark, layoutType } = storeToRefs(themeStore);

const title = import.meta.env.VITE_GLOB_APP_TITLE;

const isDark = computed(() => asideDark.value && layoutType.value != "layoutHead");

const isTitle = computed(() => !collapsed.value || layoutType.value == "layoutHead");
</script>

<style lang="scss" scoped>
.logo-head {
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  height: 72px;
  border-right: 1px solid rgba(255, 255, 255, 8%);
}

.logo-head::after {
  position: absolute;
  right: 18px;
  bottom: 0;
  left: 18px;
  height: 1px;
  content: "";
  background: rgba(255, 255, 255, 8%);
}

.logo-box {
  display: flex;
  column-gap: 12px;
  align-items: center;
  width: 100%;
  padding: 0 18px;
  overflow: hidden;
}

.padding-unset {
  justify-content: center;
  padding: unset;
}

.logo-mark {
  position: relative;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  color: #181a20;
  font-size: 18px;
  font-weight: 950;
  background: linear-gradient(135deg, #f0b90b 0%, #f9d65d 100%);
  border: 1px solid rgba(255, 255, 255, 42%);
  border-radius: 16px;
  box-shadow: 0 14px 32px rgba(240, 185, 11, 26%);
}

.logo-mark::after {
  position: absolute;
  right: -3px;
  bottom: -3px;
  width: 12px;
  height: 12px;
  content: "";
  background: #14f195;
  border: 3px solid #111827;
  border-radius: 999px;
}

.logo-mark span {
  transform: translateY(-1px);
}

.logo-copy {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.logo-title {
  max-width: 132px;
  overflow: hidden;
  color: #111827;
  font-size: 18px;
  font-weight: 950;
  line-height: 1.05;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.logo-subtitle {
  margin-top: 4px;
  color: #98a2b3;
  font-size: 10px;
  font-weight: 850;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.dark {
  .logo-title {
    color: #ffffff;
  }

  .logo-subtitle {
    color: rgba(255, 255, 255, 48%);
  }
}

.no-border {
  border: unset;

  &::after {
    display: none;
  }
}
</style>
