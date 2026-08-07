# Hotify-Bark Server

Hotify-Bark Server 是 [Bark](https://github.com/Finb/Bark) 服务端（[Finb/bark-server](https://github.com/Finb/bark-server)）的一个**修改版分支**，扩展了部分功能供 [hotify-bridge](https://github.com/sakura-lolipop/hotify-bridge) 调用，可以在推送消息到 iOS 的同时，依赖 hotify-bridge 推送消息到 HarmonyOS。

> **注意**：本项目基于上游修改，**不会同步回原项目**，也不会再使用上游的构建产物 / 镜像。所有二进制、镜像、module 路径均已独立命名，与上游可明确区分。改动清单见 [DIFFERENCES.md](DIFFERENCES.md)。

- iOS 需下载 Bark
- HarmonyOS 需下载 Hotify

## 与原项目的区别

- 独立的 Go module、二进制名与 Docker 镜像名（`wallleap/hotify-bark-server`）
- 内置 [Gotify 兼容监控接口](./docs/GOTIFY_COMPAT.md)（`/version`、`/message`、`/stream`），供 hotify-bridge 监测 bark 推送
- 内置 [MCP](./docs/MCP.md) 接口（`/mcp`、`/mcp/:device_key`），AI 代理可直接调用推送
- 可选 Basic Auth、MySQL TLS、gotify 客户端 token 等

## 安装

### 部署文件

本项目的部署产物位于 `deploy/` 目录，均已改为独立命名、不依赖原项目：

| 文件 | 说明 |
|---|---|
| `deploy/Dockerfile` | 构建镜像（二进制 `hotify-bark-server`） |
| `deploy/docker-compose.yaml` | Docker Compose 部署（远程镜像 `wallleap/hotify-bark-server`） |
| `deploy/docker-compose.local.yaml` | Docker Compose 部署（本地构建镜像，`bin/up` 默认使用） |
| `deploy/hotify-bark-server.service` | systemd 服务 |
| `deploy/entrypoint.sh` | 容器入口，设置时区 |
| `deploy/helm-chart/` | Kubernetes Helm Chart |

### Docker

```sh
docker run -dt --name hotify-bark-server -p 18080:8080 \
  -v `pwd`/bark-data:/data \
  -e BARK_SERVER_GOTIFY_CLIENT_TOKEN="your-gotify-client-token" \
  -e BARK_SERVER_BASIC_AUTH_USER="admin" \
  -e BARK_SERVER_BASIC_AUTH_PASSWORD="secret" \
  wallleap/hotify-bark-server
```

> 镜像推送到 Docker Hub（用户 `wallleap`）。如需自己的仓库，用 `docker tag wallleap/hotify-bark-server yourname/hotify-bark-server`。

使用 docker-compose：

```sh
# 复制本项目 deploy/docker-compose.yaml 到任意目录
mkdir hotify-bark-server && cd hotify-bark-server
cp /path/to/hotify-bark-server/deploy/docker-compose.yaml .
docker compose up -d
```

`deploy/docker-compose.yaml` 已内置注释的环境变量入口，按需取消注释即可。

> **本地开发**：`bin/up` / `bin/down` 默认使用 `deploy/docker-compose.local.yaml`（本地构建镜像，不拉远程），
> `bin/up` 会自动 `--build`。如需用远程镜像：`COMPOSE_FILE=deploy/docker-compose.yaml bin/up`。

### CI 构建与推送（GitHub Actions）

`.github/workflows/ci.yaml` 会在打上 `v*` 开头的标签时自动构建并推送镜像到 **Docker Hub** 和 **GHCR**：

> **版本规范**：`v*` 开头的 tag 必须符合[语义化版本](https://semver.org/)（`vMAJOR.MINOR.PATCH`，可带 `-prerelease`/`+build`，如 `v1.2.3`、`v2.0.0-rc.1`），推送到 Docker Hub 与 GHCR 的镜像 tag 即该版本号。

- Docker Hub：`<DOCKERHUB_USERNAME>/hotify-bark-server`
- GHCR：`ghcr.io/<GitHub 账号>/hotify-bark-server`

推送镜像需要配置两个 Secrets（仓库 **Settings → Secrets and variables → Actions**）：

| Secret | 说明 |
|---|---|
| `DOCKERHUB_USERNAME` | Docker Hub 用户名，用于登录及镜像命名空间 |
| `DOCKERHUB_TOKEN` | Docker Hub [Access Token](https://hub.docker.com/settings/security)（非登录密码） |

GHCR 推送使用仓库自带 `GITHUB_TOKEN`，需在 **Settings → Actions → General → Workflow permissions** 中开启 **Read and write permissions**（`packages: write` 已在 workflow 内显式声明）。

### systemd

```sh
# 1. 安装二进制
install -m 755 hotify-bark-server /usr/local/bin/hotify-bark-server

# 2. 复制服务文件
cp deploy/hotify-bark-server.service /etc/systemd/system/

# 3. 创建参数环境文件（修改后的参数放在这里）
cat > /etc/hotify-bark-server.env <<'EOF'
BARK_SERVER_GOTIFY_CLIENT_TOKEN=your-gotify-client-token
BARK_SERVER_BASIC_AUTH_USER=admin
BARK_SERVER_BASIC_AUTH_PASSWORD=secret
EOF

# 4. 启动
systemctl daemon-reload
systemctl enable --now hotify-bark-server
```

### 直接运行

1. 自行编译或从 [releases](https://github.com/wallleap/hotify-bark-server/releases) 下载预编译二进制
2. 添加执行权限：`chmod +x hotify-bark-server`
3. 启动（含修改后的参数）：

```sh
./hotify-bark-server --addr 0.0.0.0:8080 --data ./bark-data \
  --gotify-client-token your-gotify-client-token \
  --user admin --password secret
```

4. 测试：`curl localhost:8080/ping`

**注意：服务端默认使用 `/data` 目录存储数据，请确保有写权限，否则用 `--data` 指定目录。**

### 主要参数（相对上游新增/常用）

| 参数 / 环境变量 | 说明 |
|---|---|
| `--gotify-client-token` / `BARK_SERVER_GOTIFY_CLIENT_TOKEN` | Gotify 兼容监控客户端 token，自动生成并持久化。**hotify-bridge 需用它当作 `gotify_token`**（依赖见下文） |
| `--user` / `--password` / `BARK_SERVER_BASIC_AUTH_{USER,PASSWORD}` | 可选 Basic Auth |
| `--dsn` / `BARK_SERVER_DSN` | 改用 MySQL 替代 Bbolt |
| `--mysql-tls` `--mysql-ca` `--mysql-client-cert` `--mysql-client-key` `--mysql-tls-name` `--mysql-tls-skip-verify` | MySQL TLS 相关 |
| `--max-batch-push-count` / `BARK_SERVER_MAX_BATCH_PUSH_COUNT` | 批量推送上限，默认 `-1` 不限 |
| `--max-apns-client-count` | APNs 客户端连接数 |
| `--unix-socket`、`--url-prefix`、`--cert`/`--key` | 监听方式 / 前缀 / TLS |

完整参数见 `./hotify-bark-server --help`。

### 依赖 hotify-bridge（Gotify 兼容监控）

本 fork 的 Gotify 兼容接口供 [hotify-bridge](https://github.com/wallleap/hotify-bark-server) 监测推送。部署时：

1. 设置 `BARK_SERVER_GOTIFY_CLIENT_TOKEN`（**强烈推荐预置**：预置时服务端只存 SHA-256 哈希、不落明文凭证；不设置时自动生成的 token 明文会写进 `<data>/gotify.db`，且仅首次启动在日志打印一次）。
2. 在 hotify-bridge 的 `bridge_config.yaml` 填入：
   ```yaml
   gotify_url: http://<bark-host>:18080
   gotify_token: <上面的 client token>
   ```
   或使用环境变量 `GOTIFY_HTTP_URL` / `GOTIFY_CLIENT_TOKEN`。
3. 数据目录需可写，以便持久化监控消息（`<data>/gotify.db`）。

详见 [GOTIFY_COMPAT.md](./docs/GOTIFY_COMPAT.md)。

### 编译

依赖：

- Golang 1.18+
- Go Mod Enabled（`GO111MODULE=on`）
- Go Mod Proxy Enabled（`GOPROXY=https://goproxy.cn`）
- [go-task](https://taskfile.dev/installation/)

```sh
# 交叉编译全部平台
task

# 编译指定平台（参考 Taskfile.yaml）
task linux_amd64
task linux_amd64_v3
```

也可以使用脚本 `bin/build` 直接编译本机二进制。

发布 Docker 镜像：

```sh
bin/publish              # 构建并推送单架构镜像（宿主机架构）
bin/publish --all        # 多架构构建并推送（linux/amd64, linux/arm64, linux/arm/v7）
PLATFORM=linux/amd64,linux/arm64 bin/publish   # 指定架构
```

### 使用 MySQL 替代 Bbolt

以 `-dsn=user:pass@tcp(mysql_host)/bark` 启动即可使用 MySQL：

```sh
./hotify-bark-server --dsn "user:pass@tcp(mysql_host)/bark"
```

开启 TLS：

```sh
./hotify-bark-server \
  --dsn "user:pass@tcp(mysql_host)/bark" \
  --mysql-tls \
  --mysql-tls-name custom \
  --mysql-ca /etc/ssl/certs/ca.pem \
  --mysql-client-cert /etc/ssl/client-cert.pem \
  --mysql-client-key /etc/ssl/client-key.pem
# 忽略证书校验（仅测试环境）: --mysql-tls-skip-verify
```

## 文档

* [TUTORIAL.md](./docs/TUTORIAL.md) — 快速上手教程：token / key 填写到哪里
* [API_V2.md](./docs/API_V2.md) — 推送 API
* [GOTIFY_COMPAT.md](./docs/GOTIFY_COMPAT.md) — Gotify 兼容监控接口
* [MCP.md](./docs/MCP.md) — MCP 推送接口
* [TOKENS.md](./docs/TOKENS.md) — device_token、device_key、client token 是什么及如何生成
* [DIFFERENCES.md](./docs/DIFFERENCES.md) — 相对上游的改动清单

## License

MIT，见 [LICENSE](LICENSE)。上游版权归原作者（mritd / Finb）所有。
