# 快速上手教程：token / key 填写到哪里

本教程面向从零部署的使用者，把整个链路串起来：**Bark App 注册 → 拿到 device_key → 发 iOS 推送 → hotify-bridge 用 client token 订阅监控（转推 HarmonyOS）**。每一步都明确"哪个凭证填到哪里"。

> 概念（是什么、由谁生成）见 [TOKENS.md](./TOKENS.md)；推送 API 参数见 [API_V2.md](./API_V2.md)；监控接口见 [GOTIFY_COMPAT.md](./GOTIFY_COMPAT.md)。

## 1. 整体流程

```
Bark App ──注册(device_token)──▶ 服务端 ──返回 device_key──▶ Bark App / 你的脚本
                                                              │
你的脚本/服务 ──device_key──▶ POST /push ──▶ APNs ──▶ iOS 推送
                                                              │
                                gotify 兼容接口（client token 认证）
                                                              │
                    hotify-bridge（gotify_token 订阅 /stream）──▶ 华为 Push Kit ──▶ HarmonyOS
```

## 2. device_key 填到哪里

**先拿到它**（二选一）：

- **用 Bark App 自动注册**：打开 Bark，在设置里把服务器指向你的服务端（如 `http://<host>:18080`），App 会自动上报 device_token 并生成/保存一个 key，推送 URL 形如 `http://<host>:18080/<key>`，App 界面里能看到这个 key。
- **手动注册**（自定义 key 或脚本注册）：

  ```sh
  # 注册并自定义 key（推荐可读的 key，如 mydemo）
  curl "http://<host>:18080/register?devicetoken=<device_token>&key=mydemo"
  # 或 POST JSON（body 里 device_key/device_token）
  curl -X POST http://<host>:18080/register \
       -H 'Content-Type: application/json' \
       -d '{"device_key":"mydemo","device_token":"<device_token>"}'
  ```

  不指定 `key` 时服务端用 shortuuid 自动生成一个随机 key，响应里返回（同时带 `key` 和 `device_key` 两个字段，值相同）。

**填到这些地方**：

| 场景 | 填的位置 |
|---|---|
| 推送（推荐） | `POST /push` 请求体 JSON 的 `device_key` 字段 |
| 推送（旧版兼容） | URL 路径：`GET /<device_key>/<title>/<body>` |
| hotify-bridge，Hotify APP | Gotify URL/地址配置：`http://<bark-host>:18080/<device_key>` |
| 批量推送 | `POST /push` 的 `device_keys` 数组（每个元素一个 key） |
| MCP 推送 | `/mcp/<device_key>` 路径，或工具参数 `device_key` |
| iOS 侧 | Bark App 里配置的服务器地址指向你的服务端后，key 由 App 管理 |

## 3. device_token 填到哪里

device_token 由 **iOS 系统**生成、Bark App 注册 APNs 时获得，服务端不生成。

- **只在一个地方用**：注册接口 `POST /register` 或 `GET /register?devicetoken=...` 的 `devicetoken`/`device_token` 参数。
- 用 Bark App 时**不需要手动填**：App 配置好服务器地址后自动上报；只有脚本/自己实现注册时才需要显式传。
- 查看：Bark App 内可见（设置/设备信息），或 App 展示的推送 URL 里的信息。

## 4. client token 填到哪里

client token 是 gotify 兼容监控、操作接口（`/message`、`/stream`）的访问凭证，**强烈推荐预置**。

**服务端侧（预置）**——二选一，效果相同：

```sh
# 环境变量
export BARK_SERVER_GOTIFY_CLIENT_TOKEN='你的token'
# 或启动参数
./hotify-bark-server --gotify-client-token '你的token'
```

不预置则首次启动自动生成并在**日志打印一次**（重启后不再打印，日志会提示持久化位置）。

**消费侧（hotify-bridge）**——填到桥的配置：

```yaml
# bridge_config.yaml
gotify_url: http://<bark-host>:18080 # 或 http://<bark-host>:18080/<device_key>
gotify_token: <上面预置的 client token>
```

或环境变量 `GOTIFY_HTTP_URL` / `GOTIFY_CLIENT_TOKEN`。

**验证监控链路**（用 header 传 token，避免 token 出现在 URL 与日志）：

```sh
curl "http://<host>:18080/version"                             # 无需认证，返回版本
curl -H "X-Gotify-Key: <clientToken>" "http://<host>:18080/message"   # 应返回 {"messages":[...]}
curl -X DELETE -H "X-Gotify-Key: <clientToken>" "http://<host>:18080/message"  # 清空历史（可选）
```

> `?token=` 查询参数也支持（gotify 协议兼容），但会进入 URL/访问日志——生产环境请用 header 并启用 TLS（`--cert`/`--key` 或反向代理），否则 token 在网络层是明文传输。

## 5. 三种凭证一览（填哪里速查）

| 凭证 | 环节 | 填的位置 | 备注 |
|---|---|---|---|
| `device_token` | 注册 | `/register` 的 `devicetoken` 参数 | iOS 系统生成；App 自动上报 |
| `device_key` | 推送 | `/push` 的 `device_key`、URL 路径、MCP | 用户自定义或服务端生成；**即推送凭证，勿公开** |
| `client token` | 监控 | 服务端 `BARK_SERVER_GOTIFY_CLIENT_TOKEN`；桥的 `gotify_token` | 预置只存哈希；自动生成仅打印一次 |

## 6. 常见问题

- **iOS 收不到推送**：确认 `device_key` 与注册时一致、`device_token` 已注册成功（`/register` 返回 200）；推送接口返回 410/`BadDeviceToken` 说明 token 失效需重新注册。
- **桥连不上 `/message`、`/stream`（401）**：`gotify_token` 与 `BARK_SERVER_GOTIFY_CLIENT_TOKEN` 不一致，或用了自动生成 token 但服务重启过（重启后自动 token 不变，若删过 `gotify.db` 才变）。
- **找不到自动生成的 client token**：Docker 部署可在**首次启动**时用 `docker logs -f hotify-bark-server` 看到（该行仅打印一次）；或预置一个（推荐）；或停服删除 `<data>/gotify.db` 重新生成（会同时清空监控历史）。
- **`/messageevil` 之类路径被拒绝（418）**：属正常——Basic Auth 白名单只放行精确路径 `/message` 及其子路径。
