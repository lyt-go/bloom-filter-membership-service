# 布隆过滤器

纯 Go 标准库（`net/http`）实现的后端服务，零第三方依赖，标准分层（cmd/internal/pkg），开箱即跑。

## 运行

```bash
cd origin
go run ./cmd/server
# 默认监听 :8080，可用 PORT/ADDR/MAX_PAGE_SIZE 环境变量覆盖
```

## 业务实体

Filter（过滤器实例）、Strategy（哈希策略）、Element（已登记元素）、Probe（成员检测记录）、Group（过滤器分组）、Tag（过滤器标签）

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`。

| 模块 | 接口 | 说明 |
|------|------|------|
| filters | POST /api/filters | 创建Filter |
| filters | GET /api/filters | 列表查询（分页+筛选） |
| filters | GET /api/filters/{id} | 详情 |
| filters | PUT /api/filters/{id} | 更新 |
| filters | DELETE /api/filters/{id} | 删除 |
| filters | PATCH /api/filters/{id}/status | 状态流转 |
| strategies | POST /api/strategies | 创建Strategy |
| strategies | GET /api/strategies | 列表查询（分页+筛选） |
| strategies | GET /api/strategies/{id} | 详情 |
| strategies | PUT /api/strategies/{id} | 更新 |
| strategies | DELETE /api/strategies/{id} | 删除 |
| strategies | PATCH /api/strategies/{id}/status | 状态流转 |
| elements | POST /api/elements | 创建Element |
| elements | GET /api/elements | 列表查询（分页+筛选） |
| elements | GET /api/elements/{id} | 详情 |
| elements | DELETE /api/elements/{id} | 删除 |
| probes | POST /api/probes | 创建Probe |
| probes | GET /api/probes | 列表查询（分页+筛选） |
| probes | GET /api/probes/{id} | 详情 |
| probes | DELETE /api/probes/{id} | 删除 |
| groups | POST /api/groups | 创建Group |
| groups | GET /api/groups | 列表查询（分页+筛选） |
| groups | GET /api/groups/{id} | 详情 |
| groups | PUT /api/groups/{id} | 更新 |
| groups | DELETE /api/groups/{id} | 删除 |
| tags | POST /api/tags | 创建Tag |
| tags | GET /api/tags | 列表查询（分页+筛选） |
| tags | GET /api/tags/{id} | 详情 |
| tags | PUT /api/tags/{id} | 更新 |
| tags | DELETE /api/tags/{id} | 删除 |
| stats | GET /api/stats/overview | 全局统计总览 |

## 说明

- 数据存储为内存实现，重启即清空。
- 误判率（false_positive_rate）为浮点比例（0~1）；位数组大小/哈希个数为正整数。
