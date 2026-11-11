# 发布说明（Release Notes）

每个版本一个文件，命名为 `v<major>.<minor>.<patch>.md`，例如 `v0.5.3.md`。
这里写的内容有两个去处，同源：

- GitHub Release 页面的正文（`release.yml` 发布时直接用这个文件）。
- 应用内「发现新版本」窗口里展示的更新日志（构建时由 `scripts/release-notes.mjs --changelog` 生成 MyGo 要的 changelog）。

## 什么时候写

新增功能、修复缺陷、调整行为之后，把条目补进**当前开发中版本**的文件里（版本号以仓库根目录的 `VERSION` 为准）。
文件不存在就用 `node scripts/release-notes.mjs --new` 新建。发布流水线会校验该文件存在且有条目，否则发版失败。

## 怎么写

- 标题用 `# Transok <版本号>`，下一行是发布状态标记（见下）。
- 分 `## 新增` / `## 改进` / `## 修复` / `## 其他` 四节，没有内容的小节直接省略。
- 面向用户描述影响：说清楚用户会看到什么变化，而不是改了哪个函数。行为有取舍时把代价一并写出来。
- 关联 PR 或 issue 时统一用 `owner/repo#123` 格式，例如 `bent2685/transok#32`。
- 已发布版本的文件不再修改；唯一的例外是发版完成后把状态标记改成 `released`。

## 发布状态标记

标题下一行是一条 HTML 注释，在 Release 页面和更新窗口都不显示：

```
# Transok 0.5.3
<!-- release-status: unreleased -->
```

- `unreleased`：开发中，可以往里写条目，也只有这个状态能打 tag 发布。
- `released YYYY-MM-DD`：已发布，不能再写条目，也不能再发一次。

标记由脚本维护，不要手改：

| 命令 | 作用 |
| --- | --- |
| `node scripts/release-notes.mjs --status` | 查看当前版本的状态 |
| `node scripts/release-notes.mjs --new` | 为当前 `VERSION` 新建带 `unreleased` 的说明文件 |
| `node scripts/release-notes.mjs --check` | 打 tag 前的预检：文件存在、`unreleased`、有条目、无占位文字 |
| `node scripts/release-notes.mjs --mark-released` | 发版完成后标为已发布（默认今天） |

完整流程见 `transok-release` skill。
