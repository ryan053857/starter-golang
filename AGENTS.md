# GoLang 后端脚手架模板

## 通用约定

- 任何场景下，记住是在任何场景下，都只做单一逻辑，不需要考虑任何防御性逻辑，失败就让它失败，越早失败越容易发现错误，解决错误。

- 包、类型、函数写中文直述型注释：直接说明是什么、做什么、返回什么。结构体默认通过类型注释说明字段含义；较长函数可按步骤酌情加点注释。

- 优先保持函数的完整阅读流，单个函数可到 100–200 行；仅在复用或职责显著独立时抽出函数。标准库与依赖直接调用。

- 当前开发优先考虑代码整洁度和可读性，仅保证最基本的安全考虑，更严格的安全方面问题在正式上线前统一处理。

- 只做集成测试，不做单元测试。

- 数据库统一使用 ent 约定。

- `README.md` 修改和功能增强均需用户明确要求。

- 创建或更新 PR 时，标题与正文须使用中文；为该 PR 提交的 commit message 也使用中文。

## 本地命令

用根目录 `package.json` 的 bun scripts 作为统一入口：

- `bun run db:fresh`：按当前目录路径哈希创建开发库；库已存在则报错。`bun run db:fresh -- --force` 先删后建
- `bun run dev`：`serve`（全部角色）
- `bun run dev:http` / `bun run dev:worker`：仅 HTTP / 仅 Worker
- `bun run format`：格式化（`gofumpt` / `goimports`、`go fix`）
- `bun run check`：`golangci-lint`（`default: all`，过时 API 等）与测试；`git commit` 时由 hook 执行
- `bun run build`：编译到 `build/`

## 进程角色

| 命令     | 角色                         |
| -------- | ---------------------------- |
| `serve`  | 全部角色（HTTP + Worker 等） |
| `http`   | 仅 HTTP                      |
| `worker` | 仅消息消费者                 |

装配与启停在 `internal/app`：按 `Role` 位组装，信号处理与 SIGHUP 热重载（含 HTTP 时）对齐同一套生命周期。当前 demo 使用 Watermill gochannel，发布与订阅必须在同一进程（完整 demo 用 `serve`）；日后换外部 broker 再拆多进程。

## 配置加载

优先级：默认值 < YAML < 环境变量。

- 默认配置为仓库内 `config.yaml`（可用 `--config`）；缺失则失败。
- `APP_ENV=dev` 时读取 `.env`；缺失则失败。
- 示例：`.env.example`；本地 `.env` 不入库。
