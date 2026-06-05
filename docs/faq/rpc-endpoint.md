# RPC 节点配置指南

## Tron 节点

### TronGrid Api Key（强烈推荐）

> 系统默认内置的 Tron 公共 RPC 节点为 `grpc.trongrid.io:50051`，虽然目前无明显频率限制，但长期使用后被限流是必然趋势。  
> 强烈建议配置 TronGrid Api Key，以提高 Tron 扫块稳定性，避免因节点频率限制导致订单确认失败等问题；基础计划完全免费，足以满足个人需求，无需额外付费！

#### 获取 Api Key

1. 访问 https://www.trongrid.io/register 使用邮箱注册账号并完成登录。
2. 登录后找到 `API Keys` 选项，点击 `Create API Key`，填写名称后提交即可。

#### 配置 Api Key

你的 Api Key 应类似：`648870c0-xxxx-xxxx-xxxx-c7ac4ec263b0`

拿到之后登录后台，进入 `系统管理` -> `区块节点` -> `Tron 网络`，将 Api Key 填入 `TronGrid Api Key` 输入框，保存即可生效。
