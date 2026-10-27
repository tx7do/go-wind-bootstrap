# 字段映射约定（protoc-gen-fieldmap 维护手册）

本文是各适配器 `mapping.spec.yaml` / `mapping.gen.go` 的维护约定，对人和
AI 代理同等生效。改映射规则的唯一入口是 spec——`mapping.gen.go` 是生成物，
**任何情况下不手改**。

## 文件布局与命名

```
<domain>/<adapter>/
├── mapping.spec.yaml   # 映射规则（唯一事实来源）
├── mapping.gen.go      # 生成物（protoc-gen-fieldmap 产出）
└── <adapter>.go        # builder：调用生成函数 + 手工段
```

生成函数命名：option 风格 `Map<Name>Options`，struct 风格 `New<Name>Config`，
fill 风格 `Fill<Name>Config`。`<Name>` 与消息名一致（如 `Database.SQL` →
`MapSQLOptions` 按惯例取消息短名）。

## 风格选择

| 上游形态 | style | 判断依据 |
|---|---|---|
| `With*` 选项函数 | `option` | go-crud、wind-plugins 大多数 |
| 配置结构体字面量构造 | `struct` | minio `Config`、go-redis `UniversalOptions` |
| 先构造后逐字段覆盖（如 `api.DefaultConfig()`） | `fill` | registry/config 域 |

## 字段规则（硬约定）

1. **完整清单强制**：`fields` 必须覆盖 proto 消息的全部字段。无映射的字段
   写 `skip: true`（生成器对未提及字段输出告警）。上游加字段时，重新生成
   会立刻点名——把它接进映射或显式 skip，二选一，不允许沉默。
2. **guard 语义与手写一致**：默认推导为字符串判空/数值判零/bool 无条件/
   消息判 nil/集合判 len。手写代码语义不同时必须显式指定
   （`empty`/`zero`/`positive`/`not_nil`/`len`/`none`）。
3. **cast 表**：`int`/`int32`/`int64`/`uint`/`uint32`/`uint64`/`float32`/
   `float64`/`seconds`/`millis`；枚举仅 `int`；变参选项用 `spread: true`。
4. **以下内容永远留手工段**，spec 中对应字段 `skip: true` 并在 builder 里
   注释原因：
   - 复合选项（一个选项收多个字段，如 `WithCredentials(u, p)`、
     `WithBasicAuth(u, p, realm)`）；
   - TLS 装载（`tlsutil` 转换 + 错误处理）；
   - switch/分支/优先级逻辑（如 milvus 的 api_key 优先于用户名密码）；
   - proto 有字段但实现未消费的部分（这是配置面欠账的登记处）。

## 迁移阈值

机械字段 < 3 个、或映射面被控制流支配的适配器**不做生成**，保持手写。
已按此规则跳过的适配器清单见 README「字段映射代码生成」一节，勿硬凑。

## 日常工作流（含上游同步场景）

```bash
# 1. 改 proto（或上游 conf 结构变化）后重新生成——逐适配器或全量：
go run github.com/tx7do/go-wind-toolkit/protoc-gen-fieldmap@latest \
  -I conf/proto -spec database/gorm/mapping.spec.yaml -out .

# 全量（仓库根）：
find . -name 'mapping.spec.yaml' | while read -r f; do args+=( -spec "$f" ); done
go run github.com/tx7do/go-wind-toolkit/protoc-gen-fieldmap@latest \
  -I conf/proto "${args[@]}" -out .

# 2. 处理生成器的输出：
#    - 硬报错（字段不存在/类型不适配）→ 改 spec 对应字段
#    - warning（未提及字段）→ 接进映射或补 skip: true
# 3. 编译验证（必须）：
cd <adapter> && GOWORK=off go build ./... && GOWORK=off go vet ./...
```

CI 的 `fieldmap-drift` job 会做同样的事并比对 `mapping.gen.go`——改了 proto
忘了同步 spec 的 PR 会在那里红。

## AI 代理附加守则

- 上游同步时先跑全量生成看报错/告警，再按报错逐个改 spec；spec 是纯数据，
  按本文约定修改即可，不需要理解生成器实现。
- **不许做**：手改 `mapping.gen.go`；擅自接通 `skip: true` 登记的字段
  （那是行为变更，需人决策）；调整手工段代码的顺序或写法；越过
  `<domain>/<adapter>/` 目录边界。
- 每次产出 PR 时，在描述里列出：改动的 spec、触发原因（上游哪个变更）、
  生成的告警处理结果。
