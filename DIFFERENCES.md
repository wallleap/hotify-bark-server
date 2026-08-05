# 相对上游的改动

本仓库是 [Finb/bark-server](https://github.com/Finb/bark-server) 的修改版分支，**不会同步回原项目**。本文件记录相对上游 master 的改动，便于审计与升级时对照。

上游基线：`3df8990`（Finb/bark-server master）。

## 新增功能

| 功能 | 说明 | 文档 |
|---|---|---|
| Gotify 兼容监控 | `GET /version`、`GET /message`、`GET /stream`(WebSocket)，让 hotify-bridge 能像监测 Gotify 一样监测 bark | [GOTIFY_COMPAT.md](docs/GOTIFY_COMPAT.md) |
| MCP 推送 | `POST /mcp`、`POST /mcp/:device_key`，AI 代理可通过 Model Context Protocol 直接发推送 | [MCP.md](docs/MCP.md) |
| Basic Auth | 可选 `--user/--password`，`/ping` `/register` `/healthz` `/version` `/message` `/stream` 白名单 | |
| MySQL TLS | `--mysql-tls` 及配套 `mysql-ca`/`mysql-client-cert`/`mysql-client-key`/`mysql-tls-name`/`mysql-tls-skip-verify` | |
| Gotify 客户端 token | `--gotify-client-token`，SHA-256 哈希持久化，自动生成并打印一次 | |

## 命名变更

- Go module：`github.com/finb/bark-server/v2` → `github.com/wallleap/hotify-bark-server`
- 二进制：`bark-server` → `hotify-bark-server`
- Docker 镜像：`finab/bark-server` → `wallleap/hotify-bark-server`
- HTTP ServerHeader：`Bark` → `Hotify-Bark`
- CLI 名称 / 日志 / MCP 服务名均改为 Hotify-Bark 前缀

## 其它调整

- `deploy/` 下部署产物（Dockerfile、docker-compose.yaml、systemd unit、helm chart）已改用新二进制与镜像名，不依赖原项目
- `.github/workflows/ci.yaml` 镜像标签改为自有镜像名
- README 重写为本项目自有内容，不再指向原项目部署方式
