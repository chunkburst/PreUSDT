<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="首页模板编辑">
        <a-alert type="info" style="margin-bottom: 16px">
          这里编辑的是站点首页 <strong>/</strong> 的展示模板，不是收银台模板。支持安全 HTML/CSS，保存后立即生效；收银台模板请在 API 设置中通过静态资源目录覆盖。
        </a-alert>
        <a-form :model="form" :layout="layoutMode" class="base-setting-form" @submit="onSubmit">
          <a-form-item
            field="homepage_template_html"
            label="首页 HTML"
            extra="仅支持安全 HTML，禁止 script、iframe、内联事件等危险内容。可使用 {{ .title }} 与 {{ .url }}。"
          >
            <a-textarea
              v-model="form.homepage_template_html"
              placeholder="请输入自定义首页 HTML，留空则使用默认首页"
              allow-clear
              :max-length="20000"
              show-word-limit
              :auto-size="{ minRows: 12, maxRows: 24 }"
            />
          </a-form-item>

          <a-form-item
            field="homepage_template_css"
            label="首页 CSS"
            extra="仅支持安全 CSS，禁止 @import、javascript:、expression() 等危险内容。"
          >
            <a-textarea
              v-model="form.homepage_template_css"
              placeholder="请输入自定义首页 CSS，可留空"
              allow-clear
              :max-length="20000"
              show-word-limit
              :auto-size="{ minRows: 8, maxRows: 20 }"
            />
          </a-form-item>

          <a-form-item>
            <a-space>
              <a-button type="primary" html-type="submit">保存配置</a-button>
              <a-button @click="onReset">恢复默认</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { Message } from "@arco-design/web-vue";
import { setsConfAPI } from "@/api/modules/conf/index";
import { useDevicesSize } from "@/hooks/useDevicesSize";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const { isMobile } = useDevicesSize();
const layoutMode = computed(() => (isMobile.value ? "vertical" : "horizontal"));

const form = ref({
  homepage_template_html: "",
  homepage_template_css: ""
});

const onSubmit = async (): Promise<void> => {
  const response = await setsConfAPI([
    { key: "homepage_template_html", value: form.value.homepage_template_html || "" },
    { key: "homepage_template_css", value: form.value.homepage_template_css || "" }
  ]);

  if (response?.code === 200) {
    Message.success(response?.msg || "保存成功");
    emit("refresh");
    return;
  }

  Message.error(response?.msg || "保存失败");
};

const onReset = async () => {
  form.value.homepage_template_html = "";
  form.value.homepage_template_css = "";
  await onSubmit();
};

watch(
  () => data.value,
  () => {
    form.value.homepage_template_html = data.value?.homepage_template_html || "";
    form.value.homepage_template_css = data.value?.homepage_template_css || "";
  },
  { immediate: true }
);
</script>
