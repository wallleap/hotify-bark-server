<a name="unreleased"></a>

## [v0.2.0](https://github.com/wallleap/hotify-bark-server/compare/v0.1.0...v0.2.0)

> 2026-08-07

### Bug Fixes

- only expose device count on /info when Basic Auth is enabled
- surface prominent warning when Basic Auth is disabled

### Build

- add PVC persistence and MySQL credentials via Secret in Helm
- run container as non-root user and fix entrypoint startup

### Documentation

- note sudo chown 1000:1000 for host-mounted data dir
- document log level/format and Prometheus metrics (2.2)
- add analysis notes for 5.7 log redaction and 5.8 replica consistency
- mark 5.4 info device-count leak resolved
- document gotify-max-messages and mark 5.3 done
- mark completed optimization review items as done
- document actual env vars for addr and data in README
- document rate limit, non-root container and Helm in AGENTS.md
- document rate limit options, security guidance and non-root
- add optimization review with additional findings
- remove dockerhub overview workflow and file
- add Docker Hub overview sync and docker run restart policy
- expand API_V2 doc with new fields, batch push, response and auth notes
- fix gotify_url example port in GOTIFY_COMPAT.md

### Features

- wire log level, Prometheus metrics and /metrics endpoint into server
- make gotify max messages configurable
- add per-IP rate limiting for register and mcp endpoints

### Maintenance

- show gotify-max-messages in docker compose examples

### Tests

- add runs-now-apns and database package tests


## [v0.1.0](https://github.com/wallleap/hotify-bark-server/compare/v2.3.5...v0.1.0)

> 2026-08-07

### Bug Fixes

- drop request body from access logs
- mask token query param in access logs
- keep client token out of access logs
- harden apns error handling and basic auth whitelist
- ensure gotify client token is printed on first boot

### Build

- accelerate apk with Alpine mirror (default aliyun)
- **ci:** push tagged images via secrets and add semver docs

### Documentation

- how to view the auto-generated client token via docker logs
- update log policy in AGENTS.md
- record log redaction scope in AGENTS.md
- recommend header-based token and TLS in production
- add quick-start tutorial and drop README TODO list
- mark delete endpoints in README and GOTIFY_COMPAT
- strongly recommend presetting gotify client token
- refresh AGENTS.md with testing and constraints
- add token concept guide (TOKENS.md)

### Features

- add delete message endpoints to gotify compat
- rebrand finb/bark-server fork as hotify-bark-server
- add gotify-compatible monitoring interface for hotify-bridge

### Maintenance

- **deps:** bump golang.org/x/net from 0.54.0 to 0.55.0


## [v2.3.5](https://github.com/wallleap/hotify-bark-server/compare/v2.3.4...v2.3.5)

> 2026-05-20

### Maintenance

- **deps:** bump github.com/buger/jsonparser from 1.1.1 to 1.1.2
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.11 to 2.52.12


## [v2.3.4](https://github.com/wallleap/hotify-bark-server/compare/v2.3.3...v2.3.4)

> 2026-02-24

### Bug Fixes

- increase maximum device token length
- **mcp:** disable SSE streaming to mitigate FIN-WAIT-2 buildup via docker-proxy

### Documentation

- add GitHub Container Registry image usage instructions

### Maintenance

- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.9 to 2.52.11


## [v2.3.3](https://github.com/wallleap/hotify-bark-server/compare/v2.3.2...v2.3.3)

> 2025-12-30


## [v2.3.2](https://github.com/wallleap/hotify-bark-server/compare/v2.3.1...v2.3.2)

> 2025-12-22


## [v2.3.1](https://github.com/wallleap/hotify-bark-server/compare/v2.3.0...v2.3.1)

> 2025-12-22

### Features

- add additional notification options


## [v2.3.0](https://github.com/wallleap/hotify-bark-server/compare/v2.2.9...v2.3.0)

> 2025-12-22


## [v2.2.9](https://github.com/wallleap/hotify-bark-server/compare/v2.2.8...v2.2.9)

> 2025-12-22


## [v2.2.8](https://github.com/wallleap/hotify-bark-server/compare/v2.2.7...v2.2.8)

> 2025-09-28

### Bug Fixes

- return error for invalid device token removal


## [v2.2.7](https://github.com/wallleap/hotify-bark-server/compare/v2.2.6...v2.2.7)

> 2025-09-28


## [v2.2.6](https://github.com/wallleap/hotify-bark-server/compare/v2.2.5...v2.2.6)

> 2025-09-10

### Bug Fixes

- correct error status code for URL path parsing failure

### Maintenance

- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.7 to 2.52.9


## [v2.2.5](https://github.com/wallleap/hotify-bark-server/compare/v2.2.4...v2.2.5)

> 2025-07-18


## [v2.2.4](https://github.com/wallleap/hotify-bark-server/compare/v2.2.3...v2.2.4)

> 2025-07-18


## [v2.2.3](https://github.com/wallleap/hotify-bark-server/compare/v2.2.2...v2.2.3)

> 2025-07-17


## [v2.2.2](https://github.com/wallleap/hotify-bark-server/compare/v2.2.1...v2.2.2)

> 2025-07-17

### Maintenance

- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.6 to 2.52.7


## [v2.2.1](https://github.com/wallleap/hotify-bark-server/compare/v2.2.0...v2.2.1)

> 2025-05-12

### Maintenance

- **deps:** bump golang.org/x/net from 0.36.0 to 0.38.0
- **deps:** bump github.com/golang-jwt/jwt/v4 from 4.5.1 to 4.5.2
- **deps:** bump golang.org/x/net from 0.34.0 to 0.36.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.5 to 2.52.6
- **deps:** bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.5.1
- **deps:** bump golang.org/x/net from 0.31.0 to 0.34.0


## [v2.2.0](https://github.com/wallleap/hotify-bark-server/compare/v2.1.9...v2.2.0)

> 2024-12-31


## [v2.1.9](https://github.com/wallleap/hotify-bark-server/compare/v2.1.8...v2.1.9)

> 2024-12-18


## [v2.1.8](https://github.com/wallleap/hotify-bark-server/compare/v2.1.7...v2.1.8)

> 2024-12-16


## [v2.1.7](https://github.com/wallleap/hotify-bark-server/compare/v2.1.6...v2.1.7)

> 2024-12-12


## [v2.1.6](https://github.com/wallleap/hotify-bark-server/compare/v2.1.5...v2.1.6)

> 2024-12-11

### Maintenance

- **deps:** bump golang.org/x/net from 0.30.0 to 0.31.0
- **deps:** bump github.com/sideshow/apns2 from 0.24.0 to 0.25.0
- **deps:** bump golang.org/x/net from 0.29.0 to 0.30.0
- **deps:** bump golang.org/x/net from 0.28.0 to 0.29.0
- **deps:** bump go.etcd.io/bbolt from 1.3.10 to 1.3.11
- **deps:** bump github.com/sideshow/apns2 from 0.23.0 to 0.24.0
- **deps:** bump github.com/urfave/cli/v2 from 2.27.3 to 2.27.4
- **deps:** bump golang.org/x/net from 0.27.0 to 0.28.0
- **deps:** bump github.com/urfave/cli/v2 from 2.27.2 to 2.27.3
- **deps:** bump golang.org/x/net from 0.26.0 to 0.27.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.3 to 2.52.5
- **deps:** bump golang.org/x/net from 0.25.0 to 0.26.0
- **deps:** bump golang.org/x/net from 0.24.0 to 0.25.0
- **deps:** bump go.etcd.io/bbolt from 1.3.9 to 1.3.10
- **deps:** bump github.com/urfave/cli/v2 from 2.27.1 to 2.27.2
- **deps:** bump golang.org/x/net from 0.22.0 to 0.24.0
- **deps:** bump github.com/go-sql-driver/mysql from 1.8.0 to 1.8.1
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.2 to 2.52.3
- **deps:** bump github.com/go-sql-driver/mysql from 1.7.1 to 1.8.0
- **deps:** bump golang.org/x/net from 0.21.0 to 0.22.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.1 to 2.52.2
- **deps:** bump go.etcd.io/bbolt from 1.3.8 to 1.3.9
- **deps:** bump github.com/gofiber/fiber/v2 from 2.52.0 to 2.52.1
- **deps:** bump golang.org/x/net from 0.20.0 to 0.21.0
- **deps:** bump golang.org/x/net from 0.19.0 to 0.20.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.51.0 to 2.52.0
- **deps:** bump github.com/urfave/cli/v2 from 2.27.0 to 2.27.1
- **deps:** bump github.com/urfave/cli/v2 from 2.26.0 to 2.27.0
- **deps:** bump github.com/urfave/cli/v2 from 2.25.7 to 2.26.0
- **deps:** bump golang.org/x/net from 0.18.0 to 0.19.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.50.0 to 2.51.0
- **deps:** bump golang.org/x/net from 0.17.0 to 0.18.0
- **deps:** bump go.etcd.io/bbolt from 1.3.7 to 1.3.8
- **deps:** bump github.com/gofiber/fiber/v2 from 2.49.2 to 2.50.0
- **deps:** bump golang.org/x/net from 0.16.0 to 0.17.0
- **deps:** bump golang.org/x/net from 0.15.0 to 0.16.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.49.1 to 2.49.2
- **deps:** bump github.com/gofiber/fiber/v2 from 2.49.0 to 2.49.1
- **deps:** bump golang.org/x/net from 0.14.0 to 0.15.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.48.0 to 2.49.0
- **deps:** bump golang.org/x/net from 0.13.0 to 0.14.0
- **deps:** bump golang.org/x/net from 0.12.0 to 0.13.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.47.0 to 2.48.0
- **deps:** bump golang.org/x/net from 0.11.0 to 0.12.0
- **deps:** bump github.com/urfave/cli/v2 from 2.25.6 to 2.25.7
- **deps:** bump github.com/gofiber/fiber/v2 from 2.46.0 to 2.47.0
- **deps:** bump golang.org/x/net from 0.10.0 to 0.11.0
- **deps:** bump github.com/urfave/cli/v2 from 2.25.5 to 2.25.6
- **deps:** bump github.com/urfave/cli/v2 from 2.25.4 to 2.25.5
- **deps:** bump github.com/urfave/cli/v2 from 2.25.3 to 2.25.4
- **deps:** bump github.com/gofiber/fiber/v2 from 2.45.0 to 2.46.0
- **deps:** bump golang.org/x/net from 0.9.0 to 0.10.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.44.0 to 2.45.0
- **deps:** bump github.com/urfave/cli/v2 from 2.25.2 to 2.25.3
- **deps:** bump github.com/urfave/cli/v2 from 2.25.0 to 2.25.2
- **deps:** bump github.com/go-sql-driver/mysql from 1.7.0 to 1.7.1
- **deps:** bump github.com/gofiber/fiber/v2 from 2.43.0 to 2.44.0
- **task:** add freebsd support


## [v2.1.5](https://github.com/wallleap/hotify-bark-server/compare/v2.1.4...v2.1.5)

> 2023-04-10

### Bug Fixes

- :bug: fix docker build faild

### Maintenance

- **deps:** bump golang.org/x/net from 0.8.0 to 0.9.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.42.0 to 2.43.0
- **deps:** bump github.com/urfave/cli/v2 from 2.24.4 to 2.25.0
- **deps:** bump golang.org/x/net from 0.7.0 to 0.8.0
- **deps:** bump golang.org/x/net
- **deps:** bump golang.org/x/text from 0.3.7 to 0.3.8
- **deps:** bump github.com/urfave/cli/v2 from 2.24.3 to 2.24.4
- **deps:** bump github.com/gofiber/fiber/v2 from 2.41.0 to 2.42.0
- **deps:** bump github.com/urfave/cli/v2 from 2.24.2 to 2.24.3
- **deps:** bump go.etcd.io/bbolt from 1.3.5 to 1.3.7
- **deps:** bump github.com/urfave/cli/v2 from 2.24.1 to 2.24.2
- **deps:** bump github.com/urfave/cli/v2 from 2.23.7 to 2.24.1
- **deps:** bump github.com/gofiber/fiber/v2 from 2.40.1 to 2.41.0
- **deps:** bump github.com/urfave/cli/v2 from 2.23.6 to 2.23.7
- **deps:** bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0
- **deps:** bump github.com/urfave/cli/v2 from 2.23.5 to 2.23.6
- **deps:** bump github.com/gofiber/fiber/v2 from 2.40.0 to 2.40.1
- **deps:** bump github.com/urfave/cli/v2 from 2.23.4 to 2.23.5
- **deps:** bump github.com/gofiber/fiber/v2 from 2.39.0 to 2.40.0
- **deps:** bump github.com/urfave/cli/v2 from 2.20.3 to 2.23.4
- **deps:** bump github.com/urfave/cli/v2 from 2.20.2 to 2.20.3
- **deps:** bump github.com/gofiber/fiber/v2 from 2.38.1 to 2.39.0
- **deps:** bump github.com/urfave/cli/v2 from 2.19.2 to 2.20.2
- **deps:** bump github.com/urfave/cli/v2 from 2.17.1 to 2.19.2
- **deps:** bump github.com/urfave/cli/v2 from 2.16.3 to 2.17.1
- **deps:** bump github.com/gofiber/fiber/v2 from 2.37.1 to 2.38.1
- **deps:** bump github.com/urfave/cli/v2 from 2.16.2 to 2.16.3
- **deps:** bump github.com/urfave/cli/v2 from 2.14.1 to 2.16.2
- **deps:** bump github.com/gofiber/fiber/v2 from 2.37.0 to 2.37.1
- **deps:** bump github.com/urfave/cli/v2 from 2.14.0 to 2.14.1
- **deps:** bump github.com/urfave/cli/v2 from 2.11.2 to 2.14.0
- **deps:** bump github.com/urfave/cli/v2 from 2.11.1 to 2.11.2
- **deps:** bump github.com/gofiber/fiber/v2 from 2.36.0 to 2.37.0


## [v2.1.4](https://github.com/wallleap/hotify-bark-server/compare/v2.1.2-4-gbebf1cc...v2.1.4)

> 2022-08-16

### Maintenance

- **deps:** bump github.com/gofiber/fiber/v2 from 2.35.0 to 2.36.0
- **deps:** bump github.com/urfave/cli/v2 from 2.10.3 to 2.11.1


## [v2.1.2-4-gbebf1cc](https://github.com/wallleap/hotify-bark-server/compare/v2.1.3...v2.1.2-4-gbebf1cc)

> 2022-07-11


## [v2.1.3](https://github.com/wallleap/hotify-bark-server/compare/v2.1.2...v2.1.3)

> 2022-07-11

### Maintenance

- **ci:** add GOAMD64 support


## [v2.1.2](https://github.com/wallleap/hotify-bark-server/compare/v2.1.1...v2.1.2)

> 2022-07-11

### Maintenance

- **task:** add release task


## [v2.1.1](https://github.com/wallleap/hotify-bark-server/compare/v2.1.0...v2.1.1)

> 2022-07-11

### Documentation

- **API_V2.md:** add url field description

### Features

- 优化代码
- 优化 时区设置
- 优化变量名称
- 优化时区设置
- add apns push msg payload size check

### Maintenance

- **ci:** add taskfile
- **deps:** bump github.com/urfave/cli/v2 from 2.4.0 to 2.10.3
- **deps:** bump github.com/gofiber/fiber/v2 from 2.29.0 to 2.35.0
- **deps:** bump github.com/sideshow/apns2 from 0.22.0 to 0.23.0
- **deps:** bump github.com/sideshow/apns2 from 0.20.0 to 0.22.0
- **deps:** fix directory
- **deps:** check Dockerfile deps
- **docker:** fix ci build failed


## [v2.1.0](https://github.com/wallleap/hotify-bark-server/compare/v2.0.3...v2.1.0)

> 2022-03-18

### Maintenance

- **deps:** bump github.com/mritd/logger from 0.0.5 to 0.0.6
- **deps:** bump github.com/gofiber/fiber/v2 from 2.20.0 to 2.29.0
- **deps:** bump github.com/urfave/cli/v2 from 2.3.0 to 2.4.0


## [v2.0.3](https://github.com/wallleap/hotify-bark-server/compare/v2.0.2...v2.0.3)

> 2021-12-26


## [v2.0.2](https://github.com/wallleap/hotify-bark-server/compare/v2.0.1...v2.0.2)

> 2021-12-24

### Bug Fixes

- [#66](https://github.com/wallleap/hotify-bark-server/issues/66)
- **x509:** fix load system cert pool crash on windows

### Features

- 添加并设置时区为中国时区
- optimize the use of alpine apk command
- **version:** Bump v2.0.2

### Maintenance

- **deps:** bump github.com/gofiber/fiber/v2 from 2.10.0 to 2.11.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.9.0 to 2.10.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.16.0 to 2.17.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.15.0 to 2.16.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.14.0 to 2.15.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.13.0 to 2.14.0
- **deps:** bump github.com/json-iterator/go from 1.1.11 to 1.1.12
- **deps:** bump github.com/gofiber/fiber/v2 from 2.12.0 to 2.13.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.17.0 to 2.20.0
- **deps:** bump github.com/gofiber/fiber/v2 from 2.11.0 to 2.12.0
- **deps:** bump github.com/lithammer/shortuuid/v3 from 3.0.6 to 3.0.7
- **deps:** add deps bot
- **deps:** bump github.com/gofiber/fiber/v2 from 2.5.0 to 2.9.0
- **deps:** bump github.com/json-iterator/go from 1.1.9 to 1.1.11
- **depsbot:** fix time
- **fiber:** update fiber to v2.9.0


## [v2.0.1](https://github.com/wallleap/hotify-bark-server/compare/v2.0.0...v2.0.1)

> 2021-03-18

### Bug Fixes

- **v2:** fix device key missing
- **v2:** fix device token missing

### Documentation

- **v2:** update doc

### Features

- **cert:** add GeoTrust Global CA
- **http:** add https support


## [v2.0.0](https://github.com/wallleap/hotify-bark-server/compare/v1.0.2...v2.0.0)

> 2021-02-23

### Bug Fixes

- **register_check:** remove url param
- **url:** fix compat url code

### Documentation

- **README:** update README
- **README:** update README
- **apns:** add comment
- **authkey:** add pem format
- **v2:** update v2 api doc
- **v2:** add v2 api doc

### Features

- **apns:** mv apns to new pkg
- **apns:** add apple CAs
- **apns:** add private key
- **auth:** add basic auth support
- **comapt:** update compat route weight
- **compat:** update compat request
- **graceful-shutdown:** add db close
- **http:** change http framework to fiber
- **http:** remove pre fork support
- **info:** add number of devices
- **log:** update log format
- **log:** update log format
- **register:** add register check
- **register:** remove uuid, fix query parse
- **route:** update register check route
- **router:** update log format
- **v2:** tmp commit

### Maintenance

- **build:** add Apple M1 support
- **buildx:** add buildx support
- **buildx:** add binfmt install
- **ci:** remove linux/386
- **deploy:** update compose version
- **docker:** update base image
- **docker:** update dockerfile
- **docker:** update dockerfile
- **docker:** fix netgo
- **git:** update git ignore
- **make:** add pre release
- **make:** update compile sh
- **mod:** update sum
- **mod:** update go mod
- **mod:** update go version
- **systemd:** add pre-fork option


## [v1.0.2](https://github.com/wallleap/hotify-bark-server/compare/v1.0.1...v1.0.2)

> 2020-11-24

### Documentation

- **README:** add toc
- **README:** update README
- **README:** update README

### Maintenance

- **build:** fix cross compile
- **build:** remove freebsd pre build
- **build:** update build file name
- **docker:** add trimpath
- **docker:** fix dockerfile
- **docker:** fix docker build
- **gox:** remove gox
- **make:** update
- **make:** update make file, update docker base image


## [v1.0.1](https://github.com/wallleap/hotify-bark-server/compare/1.0.0-build2...v1.0.1)

> 2020-01-03

### Bug Fixes

- **server:** fix ping msg

### Maintenance

- **build:** add deploy file to release
- **ci:** update build config, add version cmd
- **deploy:** add cert
- **make:** fix make release
- **make:** fix env
- **release:** fix release files


## [1.0.0-build2](https://github.com/wallleap/hotify-bark-server/compare/1.0.0...1.0.0-build2)

> 2019-12-27

### Documentation

- **readme:** update README.MD

### Maintenance

- **docker:** update base image
- **docker-compose:** update docker-compose
- **gomod:** update go mode dep
- **make:** update makefile


## 1.0.0

> 2019-03-01

### Documentation

- **README:** update docker compose command
- **README:** update README.md
- **README:** update README.md
- **README:** update README.md
- **README:** update README.md
- **README:** remove docker image version tag
- **README:** update README

### Features

- **data:** support custom bark server db path
- **docker:** add docker compose support

### Maintenance

- **ci:** fix docker build cache
- **ci:** fix ca-certificates
- **docker:** fix dockerfile

