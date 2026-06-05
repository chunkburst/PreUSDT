<template>
  <div class="tabs">
    <a-tabs
      :editable="true"
      :hide-content="true"
      :active-key="currentRoute.path"
      size="medium"
      type="line"
      @tab-click="onTabs"
      @delete="onDelete"
    >
      <a-tab-pane v-for="item of tabsList" :key="item.path" :title="$t(`menu.${item.meta.title}`)" :closable="!item.meta.affix" />
    </a-tabs>
    <div class="tabs-setting">
      <a-space>
        <a-tooltip :content="$t(`system.refresh`)" position="bottom" mini>
          <span ref="refreshRef" id="system-tabs-refresh" :class="rotateOpen && 'refresh'">
            <icon-refresh :size="18" @click="refresh" />
          </span>
        </a-tooltip>
        <a-dropdown trigger="hover" :popup-max-height="false">
          <div class="setting" id="system-tabs-setting"><icon-apps :size="18" /></div>
          <template #content>
            <a-doption @click="closeCurrent">
              <template #icon><icon-close /></template>
              <template #default>{{ $t(`system.close-current`) }}</template>
            </a-doption>
            <a-doption @click="closeSides('left')">
              <template #icon><icon-left /></template>
              <template #default>{{ $t(`system.close-left-side`) }}</template>
            </a-doption>
            <a-doption @click="closeSides('right')">
              <template #icon><icon-right /></template>
              <template #default>{{ $t(`system.close-right-side`) }}</template>
            </a-doption>
            <a-doption @click="closeOther('other')">
              <template #icon><icon-close-circle /></template>
              <template #default>{{ $t(`system.close-other`) }}</template>
            </a-doption>
            <a-doption @click="closeOther('all')">
              <template #icon><icon-folder-delete /></template>
              <template #default>{{ $t(`system.close-all`) }}</template>
            </a-doption>
          </template>
        </a-dropdown>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useThemeConfig } from "@/store/modules/theme-config";
const router = useRouter();
const routerStore = useRouteConfigStore();
const { tabsList, currentRoute } = storeToRefs(routerStore);

const onTabs = (key: string) => {
  router.push(key);
};

const onDelete = (path: string) => {
  routerStore.removeTabsList(path);
  routerStore.removeRouteName(path);
  if (tabsList.value.length == 0) return;
  if (currentRoute.value.path != path) return;
  router.push(tabsList.value.at(-1).path);
};

const rotateOpen = ref(false);
const refresh = () => {
  rotateOpen.value = true;
  setTimeout(() => {
    rotateOpen.value = false;
  }, 500);
  const themeStore = useThemeConfig();
  themeStore.setRefreshPage(false);
  currentRoute.value.meta.keepAlive && routerStore.removeRouteName(currentRoute.value.path);
  nextTick(() => {
    themeStore.setRefreshPage(true);
    currentRoute.value.meta.keepAlive && routerStore.setRoutePaths(currentRoute.value.path);
  });
};

const closeCurrent = () => {
  onDelete(currentRoute.value.path);
};

const closeSides = (type: string) => {
  let currentIndex = tabsList.value.findIndex((item: Menu.MenuOptions) => item.path === currentRoute.value.path);
  let rightList = tabsList.value.filter((item: Menu.MenuOptions, index: number) => {
    if (type == "right") {
      if (index > currentIndex && !item.meta.affix) return item;
    } else {
      if (index < currentIndex && !item.meta.affix) return item;
    }
  });
  let rightPaths = rightList.map((item: Menu.MenuOptions) => item.path);
  tabsList.value = tabsList.value.filter((item: Menu.MenuOptions) => !rightPaths.includes(item.path));
  routerStore.removeRoutePaths(rightPaths);
};

const closeOther = (type: string) => {
  let list = tabsList.value.filter((item: Menu.MenuOptions) => {
    if (type == "other") {
      if (item.path != currentRoute.value.path && !item.meta.affix) {
        return item;
      }
    } else {
      if (!item.meta.affix) {
        return item;
      }
    }
  });
  let rightNames = list.map((item: Menu.MenuOptions) => item.path);
  tabsList.value = tabsList.value.filter((item: Menu.MenuOptions) => !rightNames.includes(item.path));
  routerStore.removeRoutePaths(rightNames);
  if (tabsList.value.length != 0 && !currentRoute.value.meta.affix && type == "all") {
    router.push(tabsList.value.at(-1).path);
  }
};
</script>

<style lang="scss" scoped>
.tabs {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 52px;
  padding: 7px 14px;
  background: rgba(248, 250, 252, 86%);
  border-bottom: 1px solid rgba(229, 231, 235, 88%);
  box-shadow: inset 0 -1px 0 rgba(255, 255, 255, 68%);
  backdrop-filter: blur(14px);

  .tabs-setting {
    flex: 0 0 auto;
    margin: 0 0 0 12px;
    color: #667085;

    .setting,
    #system-tabs-refresh {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 34px;
      height: 34px;
      color: #667085;
      background: #ffffff;
      border: 1px solid rgba(229, 231, 235, 92%);
      border-radius: 12px;
      box-shadow: 0 8px 18px rgba(17, 24, 39, 4%);
      cursor: pointer;
      transition: all 0.18s ease;

      &:hover {
        color: #181a20;
        border-color: rgba(240, 185, 11, 55%);
        box-shadow: 0 10px 22px rgba(240, 185, 11, 12%);
        transform: translateY(-1px);
      }
    }

    .refresh {
      transform: rotate(360deg);
      transition: transform 0.5s;
    }
  }
}

:deep(.arco-tabs) {
  min-width: 0;
  flex: 1;
}

:deep(.arco-tabs-nav-tab-list) {
  gap: 8px;
}

:deep(.arco-tabs-tab) {
  max-width: 176px;
  height: 34px;
  margin: 0;
  padding: 0 12px;
  color: #667085;
  background: rgba(255, 255, 255, 78%);
  border: 1px solid rgba(229, 231, 235, 86%);
  border-radius: 999px;
  box-shadow: 0 8px 18px rgba(17, 24, 39, 3%);
  font-size: 13px;
  font-weight: 760;
  transition: all 0.18s ease;
}

:deep(.arco-tabs-tab:hover) {
  color: #111827;
  border-color: rgba(240, 185, 11, 48%);
  background: #ffffff;
}

:deep(.arco-tabs-tab-active) {
  color: #181a20;
  font-weight: 900;
  background: linear-gradient(135deg, rgba(240, 185, 11, 92%) 0%, rgba(247, 210, 90, 92%) 100%);
  border-color: rgba(240, 185, 11, 86%);
  box-shadow: 0 12px 26px rgba(240, 185, 11, 22%);
}

:deep(.arco-tabs-tab-title) {
  display: inline-flex;
  align-items: center;
  min-width: 0;
  overflow: hidden;
  line-height: 1;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.arco-tabs-tab-close-btn) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-left: 6px;
  color: inherit;
  border-radius: 999px;
  opacity: 0.42;
  transition: all 0.16s ease;
}

:deep(.arco-tabs-tab-close-btn:hover) {
  background: rgba(17, 24, 39, 10%);
  opacity: 1;
}

:deep(.arco-tabs-nav-tab) {
  min-width: 0;

  .arco-tabs-tab-closable {
    svg {
      width: 0.85em;
      transition: all 0.2s;
    }
  }

  &:hover .arco-tabs-tab-title::before {
    background: unset;
  }
}

:deep(.arco-tabs-nav) {
  &::before {
    background: unset;
  }
}

:deep(.arco-tabs-nav-ink) {
  display: none;
}
</style>
