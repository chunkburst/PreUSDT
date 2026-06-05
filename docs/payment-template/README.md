# 收银台自定义修改指南

## 概述

本文档介绍如何自定义修改 PreUSDT 收银台的网页模板。本指南假设您已具备基本的 Linux/Docker 和前端开发知识。

## 获取默认资源文件

收银台的网页资源文件位于官方仓库：

https://github.com/chunkburst/PreUSDT/tree/main/static/payment

## 目录结构

```bash
./payment
├── assets
│   ├── css          # 样式文件
│   ├── i18n         # 国际化配置
│   ├── img          # 图片资源
│   ├── js           # JavaScript 文件
│   └── locales      # 多语言文件
└── views
    ├── cashier.html
    ├── bsc.bnb.html
    ├── ethereum.eth.html
    ├── index.html
    ├── installed.html
    ├── tron.trx.html
    ├── usdc.aptos.html
    ├── usdc.arbitrum.html
    ├── usdc.base.html
    └── ...
```

**目录说明：**

- **assets** - 静态资源文件（CSS、JavaScript、图片、国际化文件）
- **views** - 收银台 HTML 模板文件，`cashier.html` 为新版开放收银台/多链选择页面，其他 HTML 文件对应单币种收银台页面

### 新版收银台相关资源

新版开放收银台主要由以下文件组成：

- `views/cashier.html` - 收银台页面结构，包含币种/网络选择、二维码、地址和支付说明区域
- `assets/js/cashier.js` - 收银台交互逻辑，包括支付方式选择、切换次数提示、二维码刷新和状态轮询
- `assets/css/cashier.css` - 新版收银台样式，包括初始选择卡片、展开动画、图标和风险提示样式
- `assets/locales/*.json` - 中英文等多语言文案
- `assets/img/*.svg` - 本地币种与网络图标资源

自定义模板时建议保留 `assets/img` 下的本地图标资源，并继续使用相对静态路径引用；不建议在收银台中热链远程图标，避免第三方资源不可用或被替换时影响支付页面可信度。

## 修改步骤

### 1. 修改网页模板

在 `views` 目录下编辑相应的 HTML 文件进行自定义修改。

### 2. 上传资源文件到服务器

#### Linux 直接部署

将修改后的整个 `payment` 目录上传到服务器指定路径：

```bash
# 示例：上传到 /root/test/payment/
scp -r ./payment user@server:/root/test/
```

#### Docker 部署

**选项 A：使用 Volume 挂载（推荐）**

在启动容器时挂载本地目录，避免每次都复制文件：

```bash
docker run -v /root/test/payment:/app/static/payment <image_id>
```

**选项 B：复制到运行中的容器**

```bash
docker cp payment/ <container_id>:/app/static/
```

### 3. 配置静态资源路径

![API设置](./1.png)

1. 登录 PreUSDT 后台管理系统
2. 进入 **系统管理** → **基本设置** → **API 设置**
3. 在**静态资源路径**字段中填入完整目录路径
4. 点击保存

### 4. 重启服务并验证

重启 PreUSDT 服务：

```bash
# Linux 直接部署
systemctl restart preusdt

# Docker 部署
docker restart <container_id>
```

查看服务日志，确认出现资源注册成功提示（如下图所示），则表示配置正确。之后访问当请求到对应的交易类型收银台时便能看到修改效果。

![成功注册](2.png)
