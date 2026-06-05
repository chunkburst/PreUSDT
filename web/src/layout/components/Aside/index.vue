<template>
  <aside :class="asideDark ? 'aside dark' : 'aside'">
    <Logo />
    <div v-if="!collapsed" class="aside-status">
      <span class="status-dot"></span>
      <span>Self-hosted gateway</span>
    </div>
    <a-layout-sider :collapsed="collapsed" breakpoint="xl" class="layout-side" :width="236">
      <a-scrollbar style="height: 100%; overflow: auto" outer-class="scrollbar"><Menu :route-tree="routeTree" /></a-scrollbar>
    </a-layout-sider>
  </aside>
</template>

<script setup lang="ts">
import Logo from "@/layout/components/Logo/index.vue";
import Menu from "@/layout/components/Menu/index.vue";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useRouteConfigStore } from "@/store/modules/route-config";
const themeStore = useThemeConfig();
const { collapsed, asideDark } = storeToRefs(themeStore);
const routerStore = useRouteConfigStore();
const { routeTree } = storeToRefs(routerStore);
</script>

<style lang="scss" scoped>
.aside {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background:
    radial-gradient(circle at 26px 34px, rgba(240, 185, 11, 26%), transparent 28%),
    radial-gradient(circle at 92% 92%, rgba(22, 93, 255, 18%), transparent 34%),
    linear-gradient(180deg, #111827 0%, #090f1d 100%);
  border-right: 1px solid rgba(255, 255, 255, 10%);
  box-shadow: 18px 0 46px rgba(15, 23, 42, 18%);
}

.aside::before {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: "";
  background-image: linear-gradient(rgba(255, 255, 255, 4%) 1px, transparent 1px);
  background-size: 100% 36px;
  mask-image: linear-gradient(180deg, transparent, #000 14%, #000 82%, transparent);
}

.aside:not(.dark) {
  background:
    radial-gradient(circle at 26px 34px, rgba(240, 185, 11, 14%), transparent 28%),
    radial-gradient(circle at 92% 92%, rgba(22, 93, 255, 8%), transparent 34%),
    rgba(255, 255, 255, 96%);
  border-right-color: rgba(229, 231, 235, 92%);
  box-shadow: 12px 0 34px rgba(17, 24, 39, 5%);
}

.aside-status {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 12px 16px 4px;
  padding: 9px 12px;
  color: rgba(255, 255, 255, 68%);
  font-size: 11px;
  font-weight: 850;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  background: rgba(255, 255, 255, 7%);
  border: 1px solid rgba(255, 255, 255, 9%);
  border-radius: 14px;
}

.status-dot {
  width: 7px;
  height: 7px;
  background: #14f195;
  border-radius: 999px;
  box-shadow: 0 0 0 4px rgba(20, 241, 149, 12%);
}

.aside:not(.dark) .aside-status {
  color: #667085;
  background: #f8fafc;
  border-color: rgba(229, 231, 235, 90%);
}

.layout-side {
  position: relative;
  z-index: 1;
  flex: 1;
  overflow: hidden;

  .scrollbar {
    height: 100%;
  }
}

:deep(.arco-scrollbar-thumb-direction-vertical .arco-scrollbar-thumb-bar) {
  width: 4px;
  margin-left: 8px;
  background: rgba(240, 185, 11, 58%);
  border-radius: 999px;
}

:deep(.arco-layout-sider-light),
:deep(.arco-layout-sider-dark) {
  border-right: 0;
  box-shadow: unset;
}

:deep(.arco-menu) {
  padding: 10px 12px 22px;
  background: transparent;
}

:deep(.arco-menu-item),
:deep(.arco-menu-inline-header) {
  position: relative;
  height: 42px;
  margin: 5px 0;
  color: rgba(255, 255, 255, 70%);
  border: 1px solid transparent;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 820;
  letter-spacing: -0.01em;
  transition: all 0.18s ease;
}

:deep(.arco-menu-item::before),
:deep(.arco-menu-inline-header::before) {
  position: absolute;
  top: 10px;
  bottom: 10px;
  left: -8px;
  width: 3px;
  content: "";
  background: transparent;
  border-radius: 999px;
  transition: all 0.18s ease;
}

:deep(.arco-menu-item:hover),
:deep(.arco-menu-inline-header:hover) {
  color: #ffffff;
  background: rgba(255, 255, 255, 8%) !important;
  border-color: rgba(255, 255, 255, 10%);
  transform: translateX(2px);
}

:deep(.arco-menu-icon) {
  color: rgba(240, 185, 11, 90%);
  font-size: 17px;
}

:deep(.arco-menu-title) {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.arco-menu-selected) {
  color: #181a20 !important;
  background: linear-gradient(135deg, #f0b90b 0%, #f7d25a 100%) !important;
  border-color: rgba(255, 255, 255, 18%);
  box-shadow: 0 14px 30px rgba(240, 185, 11, 24%);
}

:deep(.arco-menu-selected::before) {
  background: #f0b90b;
}

:deep(.arco-menu-selected .arco-menu-icon),
:deep(.arco-menu-selected .arco-icon) {
  color: #181a20 !important;
}

:deep(.arco-menu-inline-content) {
  margin: 2px 0 8px 12px;
  padding-left: 10px;
  border-left: 1px solid rgba(255, 255, 255, 8%);
}

:deep(.arco-menu-inline-header .arco-menu-icon-suffix) {
  color: rgba(255, 255, 255, 42%);
}

.aside:not(.dark) {
  :deep(.arco-menu-item),
  :deep(.arco-menu-inline-header) {
    color: #475467;
  }

  :deep(.arco-menu-item:hover),
  :deep(.arco-menu-inline-header:hover) {
    color: #111827;
    background: #f8fafc !important;
    border-color: rgba(229, 231, 235, 90%);
  }

  :deep(.arco-menu-inline-content) {
    border-left-color: rgba(229, 231, 235, 90%);
  }

  :deep(.arco-menu-inline-header .arco-menu-icon-suffix) {
    color: #98a2b3;
  }
}

:deep(.arco-menu-vertical.arco-menu-collapsed) {
  padding-inline: 8px;

  .arco-menu-item,
  .arco-menu-inline-header {
    transform: none;
  }

  .arco-menu-has-icon {
    justify-content: center;
    padding: 0;
  }

  .arco-menu-icon {
    padding: 10px 0;
    margin-right: 0;
  }

  .arco-menu-title {
    display: none;
  }
}

.arco-layout-sider {
  background: unset;
}
</style>
