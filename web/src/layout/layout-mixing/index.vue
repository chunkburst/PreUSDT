<template>
  <a-layout class="layout">
    <aside :class="asideDark ? 'aside dark' : 'aside'" v-if="isPc">
      <Logo />
      <div v-if="!collapsed" class="aside-status">
        <span class="status-dot"></span>
        <span>Mixed navigation</span>
      </div>
      <a-layout-sider :collapsed="collapsed" breakpoint="xl" class="layout-side" :width="236">
        <a-scrollbar style="height: 100%; overflow: auto" outer-class="scrollbar"><Menu :route-tree="routeList" /></a-scrollbar>
      </a-layout-sider>
    </aside>
    <a-layout class="layout-right">
      <a-layout-header class="header">
        <div class="header-left">
          <ButtonCollapsed />
          <Breadcrumb v-if="!isPc" />
        </div>

        <div class="layout-head-menu" v-if="isPc">
          <a-menu
            v-if="drawing"
            mode="horizontal"
            :selected-keys="[selectedMenu]"
            @menu-item-click="onMenuItem"
            :popup-max-height="600"
          >
            <template v-for="item in routeTree" :key="item.path">
              <a-menu-item v-if="!item.meta.hide" :key="item.path" :popup-max-height="600">
                <template #icon v-if="item.meta.svgIcon || item.meta.icon">
                  <MenuItemIcon :svg-icon="item.meta.svgIcon" :icon="item.meta.icon" />
                </template>
                <span>{{ $t(`menu.${item.meta.title}`) }}</span>
              </a-menu-item>
            </template>
          </a-menu>
        </div>
        <HeaderRight />
      </a-layout-header>
      <Main />
      <Footer v-if="isFooter" />
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import Logo from "@/layout/components/Logo/index.vue";
import Main from "@/layout/components/Main/index.vue";
import Footer from "@/layout/components/Footer/index.vue";
import Menu from "@/layout/components/Menu/index.vue";
import HeaderRight from "@/layout/components/Header/components/header-right/index.vue";
import MenuItemIcon from "@/layout/components/Menu/menu-item-icon.vue";
import ButtonCollapsed from "@/layout/components/Header/components/button-collapsed/index.vue";
import Breadcrumb from "@/layout/components/Header/components/Breadcrumb/index.vue";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useRoutingMethod } from "@/hooks/useRoutingMethod";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useDevicesSize } from "@/hooks/useDevicesSize";
defineOptions({ name: "LayoutMixing" });
const route = useRoute();
const router = useRouter();
const routerStore = useRouteConfigStore();
const themeStore = useThemeConfig();
const { isFooter, collapsed, asideDark, language } = storeToRefs(themeStore);
const { routeTree } = storeToRefs(routerStore);
const { isPc } = useDevicesSize();

const drawing = ref<boolean>(true);
watch(language, () => {
  drawing.value = false;
  nextTick(() => (drawing.value = true));
});

// 横向菜单点击事件
// 将一级菜单下的children给左侧菜单
// 如果没有children则直接自身菜单给左侧菜单
const routeList = ref<any>([]);
const onMenuItem = (path: string) => {
  const { findLinearArray } = useRoutingMethod();
  const find = findLinearArray(path);
  // 路由存在则存入并跳转，不存在则跳404
  if (find) {
    // 给左侧树赋值
    setAside(find);
    // 若有重定向，则跳转到重定向的路由
    // 如果有子路由则重定向到自己的第一个菜单
    // 如果没有子路由则说明当前父级是一个菜单，直接跳转
    let path = "";
    if (find.redirect) {
      path = find.redirect;
    } else if (find.children && find.children.length > 0) {
      path = find.children[0].path;
    } else {
      path = find.path;
    }
    router.push(path);
  } else {
    router.push("/404");
  }
};

// 给左侧树赋值
const setAsideMenu = (find: Menu.MenuOptions) => {
  // 将父级的chindren给左侧树
  if (find.children && find.children.length > 0) {
    routeList.value = find.children;
  } else {
    // 如果没有则直接将父级给左侧树，做一级兜底
    routeList.value = [find];
  }
};

const setAside = debounce(setAsideMenu, 150);

const selectedMenu = computed(() => {
  const { getAllParentRoute } = useRoutingMethod();
  const find = getAllParentRoute(route.matched.at(-1).path);
  return find?.[0]?.path || "";
});

watch(
  selectedMenu,
  () => {
    const { getAllParentRoute } = useRoutingMethod();
    const find = getAllParentRoute(route.matched.at(-1).path);
    if (find?.[0]) setAside(find[0]);
  },
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.layout {
  height: 100vh;
}

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

:deep(.arco-menu-vertical) {
  padding: 10px 12px 22px;
  background: transparent;
}

:deep(.arco-menu-vertical .arco-menu-item),
:deep(.arco-menu-vertical .arco-menu-inline-header) {
  height: 42px;
  margin: 5px 0;
  color: rgba(255, 255, 255, 70%);
  border: 1px solid transparent;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 820;
}

:deep(.arco-menu-vertical .arco-menu-item:hover),
:deep(.arco-menu-vertical .arco-menu-inline-header:hover) {
  color: #ffffff;
  background: rgba(255, 255, 255, 8%) !important;
  border-color: rgba(255, 255, 255, 10%);
}

:deep(.arco-menu-vertical .arco-menu-icon) {
  color: rgba(240, 185, 11, 90%);
}

:deep(.arco-menu-vertical .arco-menu-selected) {
  color: #181a20 !important;
  background: linear-gradient(135deg, #f0b90b 0%, #f7d25a 100%) !important;
  box-shadow: 0 14px 30px rgba(240, 185, 11, 24%);
}

:deep(.arco-menu-vertical .arco-menu-selected .arco-menu-icon),
:deep(.arco-menu-vertical .arco-menu-selected .arco-icon) {
  color: #181a20 !important;
}

.aside:not(.dark) {
  :deep(.arco-menu-vertical .arco-menu-item),
  :deep(.arco-menu-vertical .arco-menu-inline-header) {
    color: #475467;
  }

  :deep(.arco-menu-vertical .arco-menu-item:hover),
  :deep(.arco-menu-vertical .arco-menu-inline-header:hover) {
    color: #111827;
    background: #f8fafc !important;
    border-color: rgba(229, 231, 235, 90%);
  }
}

.arco-layout-sider {
  background: unset;
}

.layout-right {
  display: grid;
  grid-template-rows: auto 1fr auto;
  height: 100%;

  .header {
    position: relative;
    z-index: 2;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 66px;
    padding: 0 $padding;
    overflow: hidden;
    background:
      radial-gradient(circle at 18% 0%, rgba(240, 185, 11, 10%), transparent 24%),
      rgba(255, 255, 255, 90%);
    border-bottom: 1px solid rgba(229, 231, 235, 86%);
    box-shadow: 0 10px 30px rgba(17, 24, 39, 4%);
    backdrop-filter: blur(14px);

    .header-left {
      display: flex;
      align-items: center;
      gap: 10px;
    }
  }

  .layout-head-menu {
    display: flex;
    flex: 1;
    min-width: 0;
    padding: 0 18px;
    overflow: hidden;
  }
}

:deep(.arco-menu-pop) {
  white-space: nowrap;
}

:deep(.arco-menu-horizontal) {
  flex: 1;
  overflow: hidden;
  background: transparent;

  .arco-menu-inner {
    padding-left: 0;

    .arco-menu-overflow-wrap {
      white-space: nowrap;
    }
  }

  .arco-menu-item {
    height: 38px;
    margin: 0 3px;
    padding: 0 14px;
    color: #667085;
    border-radius: 999px;
    font-weight: 800;
  }

  .arco-menu-selected {
    color: #181a20 !important;
    background: linear-gradient(135deg, rgba(240, 185, 11, 92%) 0%, rgba(247, 210, 90, 92%) 100%) !important;
  }
}
</style>
