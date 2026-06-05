# Linux 手动安装指南

> ⚠️ **前置条件**：本指南假设您具备 Linux 基本操作能力，包括命令行使用、文件系统管理等基础知识。

## 系统需求

| 项目         | 要求               |
|------------|------------------|
| **操作系统**   | Debian 11 或更高版本  |
| **CPU 架构** | amd64（其他架构请自行验证） |

## 安装步骤

### 1. 获取安装包

从 [GitHub Releases](https://github.com/chunkburst/PreUSDT/releases/latest/) 页面下载与系统架构对应的安装包。

```bash
wget -O ./linux-amd64-PreUSDT.tar.gz https://github.com/chunkburst/PreUSDT/releases/latest/download/linux-amd64-PreUSDT.tar.gz
```

### 2. 解压安装包

```bash
tar -zxvf ./linux-amd64-PreUSDT.tar.gz
```

解压后的目录结构：

```
./preusdt
├── preusdt              # 可执行程序文件
└── preusdt.service     # systemd 服务配置文件
```

### 3. 系统集成与自启配置

```bash
# 复制可执行文件至系统路径
mv ./preusdt/preusdt /usr/local/bin/

# 设置可执行权限
chmod +x /usr/local/bin/preusdt

# 复制服务配置文件
mv ./preusdt/preusdt.service /etc/systemd/system/

# 启用开机自启
systemctl enable preusdt.service
```

### 4. 启动服务

```bash
systemctl start preusdt.service
```

### 5. 验证服务状态

```bash
systemctl status preusdt.service
```

✅ **成功指标**：状态显示 `Active: active (running)` 表示服务已正常启动

---

## 常用操作命令

| 操作       | 命令                                  |
|----------|-------------------------------------|
| **查看状态** | `systemctl status preusdt.service`  |
| **查看日志** | `journalctl -u preusdt.service -f`  |
| **重启服务** | `systemctl restart preusdt.service` |
| **停止服务** | `systemctl stop preusdt.service`    |