# DDD 示例：用户管理 CRUD

基于 DDD（领域驱动设计）分层架构的完整示例，演示如何使用 `go-wind-bootstrap` 声明式引导框架快速搭建一个用户管理 HTTP API。

## 目录结构

```
_examples/ddd/
├── cmd/
│   └── main.go              # 入口：配置加载 → Bootstrap → 依赖注入 → 运行
├── configs/
│   └── config.yaml          # 声明式配置（HTTP / Logger / Database / Cache）
└── internal/
    ├── domain/
    │   └── user.go           # 领域层：User 实体 + UserRepository 接口
    ├── data/
    │   └── user.go           # 数据层：GORM 实现 UserRepository
    ├── service/
    │   └── user.go           # 服务层：业务逻辑（参数校验、分页限制）
    └── server/
        └── http.go           # 接口层：HTTP handler（POST/GET/PUT/DELETE）
```

## 分层说明

| 层 | 包 | 职责 | 依赖方向 |
|---|---|------|---------|
| **domain** | `internal/domain` | 定义实体（Entity）、仓储接口（Repository）。零外部依赖。 | 被所有层引用 |
| **data** | `internal/data` | 实现 domain 接口，与数据库/缓存交互。PO ↔ Entity 转换。 | → domain |
| **service** | `internal/service` | 业务逻辑编排（校验、规则）。只依赖 domain 接口。 | → domain |
| **server** | `internal/server` | HTTP 路由和请求编解码。委托给 service 处理。 | → service, domain |
| **cmd** | `cmd/main` | 启动入口。加载配置、Bootstrap 基础设施、组装依赖。 | → 所有层 |

依赖方向：`cmd → server → service → domain ← data`

## 声明式配置

`configs/config.yaml`：

```yaml
app:
  id: ddd-demo
  name: ddd-example
  version: "1.0.0"
  env: development

server:
  http:
    addr: ":8080"

logger:
  type: zap
  zap:
    level: info
    format: json

database:
  sql:
    driver: go_sqlite        # 纯 Go SQLite，无需 CGO
    source: ":memory:"       # 内存数据库
    migrate: true            # 自动建表

cache:
  local:
    size: 10485760           # 10MB
    default_ttl_seconds: 300
```

`migrate: true` 时，Bootstrap 自动执行 `AutoMigrate`。业务代码只需在 `init()` 中注册 PO 模型：

```go
import gormAdapter "github.com/tx7do/go-wind-bootstrap/database/gorm"

func init() {
    gormAdapter.RegisterMigrateModel(&UserPO{})
}
```

## 启动流程

`cmd/main.go` 执行以下 5 步：

```
1. LoadConfigFromFile  →  从 YAML 加载 Protobuf 配置
2. BootstrapWithContext →  自动创建 HTTP Server / Logger / GORM DB / Cache
3. ctx.Database()       →  从 Bootstrap Context 获取基础设施实例
4. 手动组装 DDD 分层     →  data → service → server
5. ctx.App().Run()      →  启动 HTTP Server，监听信号优雅关闭
```

关键代码：

```go
// 加载配置
cfg, _ := bootstrap.LoadConfigFromFile("configs/config.yaml")

// Bootstrap：自动创建所有基础设施
ctx, _ := bootstrap.BootstrapWithContext(nil, cfg)
defer ctx.Cleanup()

// 获取数据库实例
gormDB := ctx.Database(bootstrap.DatabaseTypeGorm).(*gormCrud.Client).DB

// 组装 DDD 分层
userRepo := data.NewUserRepository(gormDB)
userSvc := service.NewUserService(userRepo)
server.RegisterHTTPRoutes(httpServer, userSvc)

// 运行
ctx.App().Run(signalContext())
```

## 启动

```bash
cd _examples/ddd
go run ./cmd/
```

输出：

```
INFO database initialized dialect=go_sqlite
INFO cache initialized type=local
INFO application started id=ddd-demo env=development
```

## API 测试

```bash
# 创建用户
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'

# 返回：
# {"ID":1,"Name":"Alice","Email":"alice@example.com","CreatedAt":"...","UpdatedAt":"..."}

# 查询单个用户
curl http://localhost:8080/users/1

# 列出用户（支持分页）
curl "http://localhost:8080/users?offset=0&limit=10"

# 更新用户（支持部分更新）
curl -X PUT http://localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Bob"}'

# 删除用户
curl -X DELETE http://localhost:8080/users/1

# 健康检查
curl http://localhost:8080/health
# {"status":"ok"}
```

## 适配层引入

本示例通过空导入（blank import）自动注册 builder：

```go
import (
    _ "github.com/tx7do/go-wind-bootstrap/cache/local"     // FreeCache 内存缓存
    _ "github.com/tx7do/go-wind-bootstrap/database/gorm"   // GORM 数据库
    _ "github.com/tx7do/go-wind-bootstrap/log/zap"         // Zap 日志
    httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http" // HTTP Server
)
```

每个适配层在 `init()` 中调用 `bootstrap.MustRegisterXxxBuilder()` 完成自注册，Bootstrap 阶段根据配置自动查找并调用。

## 扩展指南

### 切换数据库

修改 `config.yaml` 即可切换数据库，无需改代码：

```yaml
# MySQL
database:
  sql:
    driver: mysql
    source: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True"

# PostgreSQL
database:
  sql:
    driver: postgres
    source: "host=127.0.0.1 user=postgres password=secret dbname=test sslmode=disable"
```

### 添加新的业务模块

1. 在 `domain/` 中定义实体和仓储接口
2. 在 `data/` 中实现仓储（注册 PO 模型到迁移表）
3. 在 `service/` 中实现业务逻辑
4. 在 `server/` 中注册 HTTP 路由
5. 在 `cmd/main.go` 中组装依赖
