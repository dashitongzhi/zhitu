# 开发协作指南

本文件定义团队的 Git 工作流、分支策略、测试规范与分工原则。所有成员在开始开发前请通读本文件。

## Git 工作流

### 分支策略

| 分支 | 用途 | 保护规则 |
|------|------|----------|
| `main` | 存放经过测试的稳定代码 | 禁止直接 push，只能通过 PR 合并 |
| `feature/*` | 功能开发分支 | 命名格式 `feature/功能名`，如 `feature/resume-preview` |
| `fix/*` | Bug 修复分支 | 命名格式 `fix/问题描述`，如 `fix/login-redirect` |

### 日常开发流程

```
1. git pull origin main              # 每日开工先拉取最新代码
2. git checkout -b feature/xxx       # 创建功能分支
3. ...开发与本地测试...
4. git add . && git commit -m "..."  # 提交代码
5. git push origin feature/xxx       # 推送到远程
6. 在 GitHub 上创建 PR → 通知队友 review
7. 队友 review + 本地测试 → 合并到 main
8. git checkout main && git pull     # 切回主分支并拉取最新
```

### Commit 信息规范

使用 conventional commits 格式：

```
<type>(<scope>): <description>

type:   feat | fix | docs | style | refactor | test | chore
scope:  模块名，如 resume | interview | dashboard | auth
```

示例：

```
feat(resume): 实现简历竞争力雷达的六维评分计算
fix(auth): 修复 GitHub OAuth 回调地址错误
docs(prd): 更新面试官画像库的字段规则
```

## 测试规范

### 本地测试流程

收到队友的 PR 后，按以下步骤进行本地测试：

```
1. git fetch origin pull/<PR编号>/head:pr-<PR编号>   # 拉取 PR 分支到本地
2. git checkout pr-<PR编号>                           # 切换到 PR 分支
3. 安装依赖（如有变更）
4. 运行测试套件
5. 手动验证功能点
6. 在 GitHub 上 review 代码并留下评论
7. 测试通过 → 合并 PR；测试失败 → 评论说明问题
```

### 自动化测试

项目配置了 GitHub Actions CI，每次 PR 提交会自动运行：

- 后端单元测试
- 前端组件测试
- 代码风格检查（lint）

CI 全部通过是合并 PR 的前提条件。

## 分工原则

### 按完整功能划分

每个功能应包含前端 + 后端的完整实现，避免按层级切分（如「你做前端、我做后端」）。这样做的原因：

- 减少跨层沟通成本
- 每个人对完整功能负责，理解更深入
- 避免「等对方接口」的阻塞

### 功能拆分建议

| 功能 | 涉及模块 | 建议负责人 |
|------|----------|------------|
| 登录与认证 | 认证模块 | 后端为主 |
| 简历编辑器 | 简历实验室 | 前端为主 |
| 简历竞争力雷达 | 简历实验室 | 全栈 |
| JD 匹配优化 | 简历实验室 + Agent | 后端为主 |
| 语音面试交互 | 面试训练场 | 前端为主 |
| 面试官画像库 | 面试训练场 | 全栈 |
| 投递看板表格 | 投递看板 | 前端为主 |
| 求职策略洞察 | 投递看板 + Agent | 后端为主 |

### 避免冲突

- **不要同时修改同一文件**。开工前在群里说一声「我要改 xxx 文件」
- **数据库 schema 变更必须提前沟通**。任何涉及表结构、字段增减的改动，先讨论再动手
- **每日 `git pull`**。开工前拉取最新代码，避免基于过期代码开发

## 数据库变更规范

数据库 schema 变更流程：

1. 在群里提出变更需求，说明原因和影响范围
2. 讨论确认后，在对应分支编写 migration 脚本
3. PR 中包含 migration 脚本和对应的 model 更新
4. 合并后所有成员执行 `alembic upgrade head` 同步本地数据库

禁止直接在代码中硬编码 SQL 修改表结构，所有变更通过 migration 脚本管理。
