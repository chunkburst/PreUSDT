<template>
  <div class="login-form-box">
    <a-form :rules="rules" :model="form" layout="vertical" @submit="onSubmit">
      <a-form-item field="username" label="管理员账号" :hide-asterisk="true">
        <a-input v-model="form.username" allow-clear placeholder="请输入后台管理员账号">
          <template #prefix>
            <icon-user />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item field="password" label="登录密码" :hide-asterisk="true">
        <a-input-password v-model="form.password" allow-clear placeholder="请输入登录密码">
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>
      <a-form-item field="remember">
        <div class="remember">
          <a-checkbox v-model="form.remember">记住本机登录信息</a-checkbox>
          <div class="forgot-password" @click="handleForgotPassword">查看重置教程</div>
        </div>
      </a-form-item>
      <a-form-item>
        <a-button class="login-submit" long type="primary" html-type="submit">进入管理后台</a-button>
      </a-form-item>
    </a-form>
    <div class="security-note">请确认当前访问地址可信，支付密钥、钱包地址和回调配置仅在私有环境中操作。</div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUserInfoStore } from "@/store/modules/user-info";
import { loginAPI } from "@/api/modules/user/index";

let userStores = useUserInfoStore();
const router = useRouter();

const REMEMBER_KEY = "login_remember_info";

const encrypt = (str: string) => {
  return btoa(encodeURIComponent(str));
};

const decrypt = (str: string) => {
  try {
    return decodeURIComponent(atob(str));
  } catch {
    return "";
  }
};

const form = ref({
  username: "",
  password: "",
  verifyCode: null,
  remember: false
});
const rules = ref({
  username: [
    {
      required: true,
      message: "请输入账号"
    }
  ],
  password: [
    {
      required: true,
      message: "请输入密码"
    }
  ]
});

onMounted(() => {
  const savedInfo = localStorage.getItem(REMEMBER_KEY);
  if (savedInfo) {
    try {
      const { username, password, remember } = JSON.parse(savedInfo);
      form.value.username = decrypt(username);
      form.value.password = decrypt(password);
      form.value.remember = remember;
    } catch (error) {
      console.error("读取记住的密码失败:", error);
      localStorage.removeItem(REMEMBER_KEY);
    }
  }
});

const onSubmit = async ({ errors }: any) => {
  if (errors) return;
  onLogin();
};

const onLogin = async () => {
  if (form.value.remember) {
    const rememberInfo = {
      username: encrypt(form.value.username),
      password: encrypt(form.value.password),
      remember: true
    };
    localStorage.setItem(REMEMBER_KEY, JSON.stringify(rememberInfo));
  } else {
    localStorage.removeItem(REMEMBER_KEY);
  }

  let res = await loginAPI(form.value);

  userStores.token = res.data.token;

  await userStores.setAccount();

  arcoMessage("success", "登录成功");
  router.replace("/home");
};

const handleForgotPassword = () => {
  window.open("https://github.com/chunkburst/PreUSDT/blob/main/docs/faq/login-reset.md", "_blank");
};
</script>

<style lang="scss" scoped>
.login-form-box {
  margin-top: 30px;
}

:deep(.arco-form-item-label-col > label) {
  color: #374151;
  font-size: 13px;
  font-weight: 800;
}

:deep(.arco-input-wrapper) {
  min-height: 44px;
  background: #f8fafc;
  border-color: #e5e7eb;
  border-radius: 14px;
  transition: all 0.18s ease;

  &:hover,
  &.arco-input-focus {
    background: #ffffff;
    border-color: #111827;
    box-shadow: 0 0 0 4px rgba(17, 24, 39, 6%);
  }
}

.remember {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 12px;

  .forgot-password {
    color: #92400e;
    font-size: 13px;
    font-weight: 800;
    cursor: pointer;
  }
}

.login-submit {
  height: 44px;
  color: #181a20;
  font-weight: 900;
  background: #f0b90b;
  border: 0;
  border-radius: 14px;
  box-shadow: 0 16px 36px rgba(240, 185, 11, 24%);

  &:hover {
    color: #181a20;
    background: #f6c72d;
  }
}

.security-note {
  margin-top: 12px;
  padding: 12px 14px;
  color: #667085;
  font-size: 12px;
  line-height: 1.65;
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
}
</style>
