<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="通知设置">
        <a-form :model="form" :rules="rules" :layout="layoutMode" class="base-setting-form" @submit="onSubmit">
          <a-form-item field="notifier_channel" label="通知渠道">
            <a-select v-model="form.notifier_channel" placeholder="请选择通知渠道" @change="onChannelChange">
              <a-option
                v-for="channel in channelConfigs"
                :key="channel.value"
                :value="channel.value"
                :disabled="channel.disabled"
              >
                {{ channel.label }}
              </a-option>
            </a-select>
          </a-form-item>

          <template v-for="field in currentChannelFields" :key="field.key">
            <a-form-item :field="field.key" :label="field.label" :extra="field.extra">
              <a-textarea
                v-if="field.multiline"
                v-model="form.notifier_params[field.key]"
                :placeholder="field.placeholder"
                allow-clear
                :max-length="4000"
                show-word-limit
                :auto-size="{ minRows: field.minRows || 3, maxRows: field.maxRows || 8 }"
              />
              <a-input
                v-else
                v-model="form.notifier_params[field.key]"
                :placeholder="field.placeholder"
                :type="field.type || 'text'"
                allow-clear
              />
            </a-form-item>
          </template>

          <a-form-item>
            <a-space>
              <a-button type="primary" html-type="submit">保存配置</a-button>
              <a-button v-if="form.notifier_channel !== 'none'" type="outline" @click="onTest" :loading="testLoading">
                推送测试
              </a-button>
              <a-button v-if="form.notifier_channel === 'telegram'" type="outline" @click="onResetTelegramTemplates">
                恢复默认模板
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useDevicesSize } from "@/hooks/useDevicesSize";
import { Message } from "@arco-design/web-vue";
import { notifierAPI, notifierTestAPI } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const { isMobile } = useDevicesSize();
const layoutMode = computed(() => (isMobile.value ? "vertical" : "horizontal"));

interface FieldConfig {
  key: string;
  label: string;
  placeholder: string;
  type?: string;
  required: boolean;
  message?: string;
  validator?: string;
  multiline?: boolean;
  minRows?: number;
  maxRows?: number;
  extra?: string;
}

interface ChannelConfig {
  value: string;
  label: string;
  disabled: boolean;
  fields: FieldConfig[];
}

interface FormData {
  notifier_channel: string;
  notifier_params: Record<string, string>;
}

const telegramTemplateExtra = "支持占位符，如 {{order_id}}、{{money}}、{{fiat}}、{{rate}}、{{amount}}、{{trade_type}}、{{tx_hash}}、{{address}}、{{created_at}}、{{updated_at}} 等。";

const DEFAULT_TELEGRAM_TEMPLATES: Record<string, string> = {
  success_template: `✅ 收款成功
商户订单：{{order_id}}
请求金额：{{money}} {{fiat}}（汇率 {{rate}}）
支付数额：{{amount}} {{trade_type}}
交易哈希：{{tx_hash}}
收款地址：{{address}}
创建时间：{{created_at}}
支付时间：{{updated_at}}`,
  notify_fail_template: `⚠️ 回调失败
商户订单：{{order_id}}
支付数额：{{amount}}
请求金额：{{money}} {{fiat}}（汇率 {{rate}}）
交易类别：{{trade_type}}
确认时间：{{confirmed_at}}
下次回调：{{next_notify_at}}
失败原因：{{reason}}`,
  non_order_transfer_template: `📨 非订单交易
监控方向：{{direction}}
交易数额：{{amount}}
交易类别：{{trade_type}}
交易时间：{{timestamp}}
接收地址：{{recv_address}}
发送地址：{{from_address}}`,
  tron_resource_template: `🔋 资源动态
操作类型：{{action}}
质押数量：{{balance}}
交易时间：{{timestamp}}
操作地址：{{recv_address}}
资源来源：{{from_address}}`,
  welcome_template: `👋 欢迎使用 preusdt，{{desc}}，如果您看到此消息，说明系统已启动成功！
当前版本：{{version}}
开源地址：{{github}}`,
  test_template: `✅ 这是一条测试消息，Telegram 通知配置成功！
当前系统时间：{{now}}`
};

const TELEGRAM_TEMPLATE_KEYS = Object.keys(DEFAULT_TELEGRAM_TEMPLATES);

const channelConfigs: ChannelConfig[] = [
  {
    value: "none",
    label: "关闭通知",
    disabled: false,
    fields: []
  },
  {
    value: "telegram",
    label: "Telegram",
    disabled: false,
    fields: [
      {
        key: "bot_token",
        label: "Bot Token",
        placeholder: "请输入 Telegram Bot Token",
        required: true,
        message: "Bot Token 不能为空"
      },
      { key: "chat_id", label: "Chat ID", placeholder: "请输入 Telegram Chat ID", required: true, message: "Chat ID 不能为空" },
      { key: "topic_id", label: "Topic ID", placeholder: "请输入 Telegram Topic ID", required: false },
      {
        key: "success_template",
        label: "收款成功模板",
        placeholder: "请输入收款成功通知模板",
        required: false,
        multiline: true,
        minRows: 5,
        maxRows: 10,
        extra: telegramTemplateExtra
      },
      {
        key: "notify_fail_template",
        label: "回调失败模板",
        placeholder: "请输入回调失败通知模板",
        required: false,
        multiline: true,
        minRows: 5,
        maxRows: 10,
        extra: telegramTemplateExtra
      },
      {
        key: "non_order_transfer_template",
        label: "非订单交易模板",
        placeholder: "请输入非订单交易通知模板",
        required: false,
        multiline: true,
        minRows: 4,
        maxRows: 8,
        extra: "支持占位符 {{direction}}、{{amount}}、{{trade_type}}、{{timestamp}}、{{recv_address}}、{{from_address}}。"
      },
      {
        key: "tron_resource_template",
        label: "资源动态模板",
        placeholder: "请输入 Tron 资源动态通知模板",
        required: false,
        multiline: true,
        minRows: 4,
        maxRows: 8,
        extra: "支持占位符 {{action}}、{{balance}}、{{timestamp}}、{{recv_address}}、{{from_address}}。"
      },
      {
        key: "welcome_template",
        label: "启动欢迎模板",
        placeholder: "请输入启动欢迎通知模板",
        required: false,
        multiline: true,
        minRows: 3,
        maxRows: 6,
        extra: "支持占位符 {{desc}}、{{version}}、{{github}}。"
      },
      {
        key: "test_template",
        label: "测试消息模板",
        placeholder: "请输入测试消息模板",
        required: false,
        multiline: true,
        minRows: 2,
        maxRows: 4,
        extra: "支持占位符 {{now}}。"
      }
    ]
  },
  {
    value: "wechat",
    label: "企业微信（开发中）",
    disabled: true,
    fields: [
      {
        key: "webhook_url",
        label: "Webhook URL",
        placeholder: "请输入企业微信 Webhook URL",
        type: "url",
        required: true,
        message: "Webhook URL不能为空",
        validator: "url"
      }
    ]
  },
  {
    value: "email",
    label: "邮箱（开发中）",
    disabled: true,
    fields: [
      {
        key: "email",
        label: "邮箱地址",
        placeholder: "请输入邮箱地址",
        type: "email",
        required: true,
        message: "邮箱地址不能为空",
        validator: "email"
      },
      {
        key: "smtp_server",
        label: "SMTP服务器",
        placeholder: "请输入SMTP服务器地址",
        required: true,
        message: "SMTP服务器不能为空"
      }
    ]
  }
];

const form = ref<FormData>({
  notifier_channel: "telegram",
  notifier_params: {}
});

const testLoading = ref<boolean>(false);

const currentChannelFields = computed<FieldConfig[]>(
  () => channelConfigs.find(config => config.value === form.value.notifier_channel)?.fields || []
);

const currentChannelParamKeys = computed<string[]>(() => currentChannelFields.value.map(field => field.key));

const rules = computed(() => {
  const baseRules: Record<string, any[]> = {
    notifier_channel: [{ required: true, message: "请选择通知渠道" }]
  };

  currentChannelFields.value.forEach(field => {
    if (field.required) {
      const fieldPath = `notifier_params.${field.key}`;
      const fieldRules: any[] = [{ required: true, message: field.message }];

      if (field.validator === "email") {
        fieldRules.push({ type: "email", message: "请输入正确的邮箱格式" });
      } else if (field.validator === "url") {
        fieldRules.push({ type: "url", message: "请输入正确的URL格式" });
      }

      baseRules[fieldPath] = fieldRules;
    }
  });

  return baseRules;
});

const fillTelegramTemplates = (params: Record<string, string>): Record<string, string> => {
  const nextParams = { ...params };
  TELEGRAM_TEMPLATE_KEYS.forEach(key => {
    if (!String(nextParams[key] || "").trim()) {
      nextParams[key] = DEFAULT_TELEGRAM_TEMPLATES[key];
    }
  });
  return nextParams;
};

const initParams = (channel = form.value.notifier_channel): Record<string, string> => {
  const params: Record<string, string> = {};
  channelConfigs.forEach(config => {
    config.fields.forEach(field => {
      params[field.key] = "";
    });
  });
  return channel === "telegram" ? fillTelegramTemplates(params) : params;
};

const buildParams = (): Record<string, string> => {
  const filteredParams: Record<string, string> = {};
  currentChannelParamKeys.value.forEach(key => {
    const value = form.value.notifier_params[key];
    if (value !== undefined && value !== null) {
      filteredParams[key] = String(value);
    }
  });
  return filteredParams;
};

const onChannelChange = (): void => {
  form.value.notifier_params = initParams();
};

const onResetTelegramTemplates = (): void => {
  form.value.notifier_params = {
    ...form.value.notifier_params,
    ...DEFAULT_TELEGRAM_TEMPLATES
  };
  Message.success("已恢复 Telegram 默认模板");
};

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit): Promise<void> => {
  if (errors) return;

  try {
    const response = await notifierAPI({
      channel: form.value.notifier_channel,
      params: buildParams()
    });

    if (response?.code === 200) {
      Message.success(response?.msg || "配置成功！");
      emit("refresh");
    } else {
      Message.error(response?.msg || "配置保存失败");
    }
  } catch (error: any) {
    console.error("配置保存失败:", error);
    Message.error(error?.msg || "配置保存失败，请稍后重试");
  }
};

const onTest = async (): Promise<void> => {
  try {
    testLoading.value = true;

    const response = await notifierTestAPI({
      channel: form.value.notifier_channel,
      params: buildParams()
    });

    if (response?.code === 200) {
      Message.success(response?.msg || "推送测试成功！");
    } else {
      Message.error(response?.msg || "推送测试失败");
    }
  } catch (error: any) {
    console.error("推送测试失败:", error);
    Message.error(error?.msg || "推送测试失败，请稍后重试");
  } finally {
    testLoading.value = false;
  }
};

watch(
  () => data.value,
  () => {
    if (data.value) {
      form.value.notifier_channel = String(data.value.notifier_channel || "telegram");

      if (data.value.notifier_params) {
        try {
          const params =
            typeof data.value.notifier_params === "string" ? JSON.parse(data.value.notifier_params) : data.value.notifier_params;

          const parsedParams: Record<string, string> = {};
          Object.keys(params).forEach(key => {
            parsedParams[key] = String(params[key] || "");
          });

          const nextParams = { ...initParams(), ...parsedParams };
          form.value.notifier_params = form.value.notifier_channel === "telegram" ? fillTelegramTemplates(nextParams) : nextParams;
        } catch (e) {
          console.error("解析 notifier_params 失败:", e);
          form.value.notifier_params = initParams();
        }
      } else {
        form.value.notifier_params = initParams();
      }
    }
  },
  { immediate: true }
);

onMounted(() => {
  if (!Object.keys(form.value.notifier_params).length) {
    form.value.notifier_params = initParams();
  }
});
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
