# Hotify-Bark Server

基于 [Bark](https://github.com/Finb/Bark) 服务端（[Finb/bark-server](https://github.com/Finb/bark-server)）的修改版分支：iOS (APNs) 推送服务，扩展了 Gotify 兼容监控接口（供 hotify-bridge 推 HarmonyOS）和 MCP 接口。二进制 / 镜像 / module 均已独立命名，与上游可明确区分。

- iOS 需下载 Bark
- HarmonyOS 需下载 Hotify

## 快速开始

```sh
docker run -dt --name hotify-bark-server --restart unless-stopped \
  -p 18080:8080 \
  -v `pwd`/bark-data:/data \
  -e BARK_SERVER_GOTIFY_CLIENT_TOKEN="your-gotify-client-token" \
  wallleap/hotify-bark-server
```

- 容器内部默认监听 `0.0.0.0:8080`，数据目录 `/data`（建议挂载持久化，需可写）。
- 宿主机映射端口 `18080`

docker-compose 请参考仓库 `deploy/docker-compose.yaml`。

## 推送示例

```sh
# V2 API（JSON）
curl -X "POST" "http://127.0.0.1:18080/push" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d '{"title": "hi", "body": "Test Bark Server", "device_key": "your-device-key"}'

# V1 兼容（URL 路径）
curl "http://127.0.0.1:18080/your-device-key/title/body"
```

## 主要特性

- **V1 / V2 推送 API**：单设备、批量（`device_keys`）、静默推送、归档、自定义字段透传。
- **Gotify 兼容监控**：`/version`、`/message`、`/stream`(WebSocket)，供 hotify-bridge 像监测 Gotify 一样监测 bark。
- **MCP 接口**：`/mcp`、`/mcp/:device_key`，AI 代理可直接发推送。
- **可选 Basic Auth**、MySQL（含 TLS）、gotify 客户端 token、多架构镜像（arm / arm64 / amd64）。

## 环境变量 / 参数

以 `BARK_SERVER_*` 前缀映射 CLI 参数（大小写转换），常用：

| 环境变量 | 对应参数 | 说明 |
|---|---|---|
| `BARK_SERVER_ADDR` | `--addr` | 监听地址，默认 `0.0.0.0:8080` |
| `BARK_SERVER_DATA` | `--data` | 数据目录，默认 `/data` |
| `BARK_SERVER_GOTIFY_CLIENT_TOKEN` | `--gotify-client-token` | Gotify 客户端 token（未设置则首次启动自动生成并打印一次） |
| `BARK_SERVER_BASIC_AUTH_USER` | `--user` | Basic Auth 用户名（可选） |
| `BARK_SERVER_BASIC_AUTH_PASSWORD` | `--password` | Basic Auth 密码（可选） |
