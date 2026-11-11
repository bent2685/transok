---
name: transok-release
description: 发布 Transok 新版本的完整流程：预检发布说明、打 tag 推送、监视 release 流水线、检查草稿 Release 并发布、把发布说明标成已发布并开下一个版本；也用于「开下一个版本」「这一版改发 0.6.0」。只在用户明确要求发版或改版本号时使用；写代码、写发布说明条目、普通推送都不用它。
---

# 发布 Transok

发版的每一步由脚本判断，你负责按顺序执行、在检查点停下、如实汇报。不要凭记忆或 git 历史推断「这个版本发没发」，一律问脚本。

原理：tag 触发 `.github/workflows/release.yml` → 建草稿 Release（正文取自 `.github/release-notes/v<版本>.md`）→ 各平台 `mygo build -upload` 用私钥签名并上传安装包与更新清单 → 人工检查草稿后 Publish → 已安装的应用在后台检查到更新，弹出原生更新窗口（内容同样来自发布说明）。**Publish 之前用户收不到更新。**

## 什么时候用

| 用户说 | 做什么 |
| --- | --- |
| 「发版」「发布 v0.5.3」 | 完整流程，第 1–6 步 |
| 「开下一个版本」「下个版本 0.6.0」 | 只做第 7 步 |
| 「这一版改发 0.6.0」 | 只做「改这一版的版本号」 |

用户只说「改完提交」「推一下 dev」不是发版，不要用这个 skill。

### 直接调用 `/transok-release`

斜杠命令本身就是明确的发版指令，不要再问「要发吗」：

| 命令 | 做什么 |
| --- | --- |
| `/transok-release` | 发布当前 `VERSION`（第 1–6 步），发完开下一个修订版本（第 7 步），推送 main、tag、dev |
| `/transok-release minor` / `major` / `0.6.0` | 同上，但第 7 步下一个版本按参数开 |
| `/transok-release next [patch\|minor\|major\|x.y.z]` | 不发版，只开下一个版本（要求当前已 released） |
| `/transok-release status` | 只汇报版本号与发布状态，什么都不改 |

只在这几种情况停下来问：预检没过、升的不是修订号、出现不可重试的失败。

## 版本号在哪里

单一来源是根目录 `VERSION`；`make version V=...` 会同步到 `mygo.json`（应用版本、更新比较用）和 `backend/consts/default_conf.go`（界面「版本」行）。**永远不要手改其中一处**，用：

```bash
make version V=patch          # patch | minor | major | 具体版本如 0.6.0
node scripts/set-version.mjs --check   # 三处一致才算对
```

## 两条不能破的线

1. **推送 tag 和 Publish 草稿都是对外发布。** 只有用户在当前对话里明确说了发版才执行；以前说过不算。
2. **失败不要自作主张。** 流水线失败先看日志诊断；不改代码、不重新打 tag、不删 tag，汇报后问用户。瞬时网络故障可以 `gh run rerun <id> --failed`。

## 先判断现在在哪一步

```bash
node scripts/release-notes.mjs --status        # unreleased / released
git fetch origin --tags
git ls-remote --tags origin "refs/tags/v$(cat VERSION)"
gh run list --workflow release.yml --limit 3
```

| 看到的 | 从哪一步继续 |
| --- | --- |
| `unreleased`，远端没有 tag | 第 1 步 |
| 远端有 tag，run 在跑或失败 | 第 4 步 |
| run 成功，Release 还是草稿 | 第 5 步 |
| Release 已发布，状态仍是 `unreleased` | 第 6 步 |
| `released` | 第 7 步 |

## 流程

### 1. 预检

```bash
git branch --show-current          # 必须是 dev
git status --short                 # 不属于这次发版的改动不要带上
make check                         # 格式、vet、i18n、版本一致、测试
node scripts/release-notes.mjs --check
```

`--check` 不过就停下：告诉用户发布说明缺什么（没有条目、占位文字、已发布），补好再来。

### 2. 合并 dev 到 main

```bash
git checkout main && git pull --ff-only origin main
git merge --no-ff dev -m "chore: 合并 dev 发布 v<版本>"
```

有冲突就停下汇报，不要自行解决。提交信息不要带任何 `Co-Authored-By` 之类的尾注（见 AGENTS.md）。

### 3. 打 tag 并推送

向用户展示摘要：版本号、自上个 tag 起的提交数（`git log --oneline v<上一版>..HEAD | wc -l`）、发布说明条目数。用户已明确要求发版才继续：

```bash
git tag v<版本>
git push origin main
git push origin v<版本>
```

tag 必须等于 `VERSION`，否则流水线第一步就会拦下。

### 4. 监视流水线

```bash
gh run list --workflow release.yml --limit 1
gh run watch <run id> --exit-status
```

失败按「两条线」处理。

### 5. 检查草稿并发布

```bash
gh release view v<版本> --json isDraft,body,assets --jq '{isDraft, assets: [.assets[].name]}'
```

确认草稿里有：macOS 的 `.dmg` 和 `.tar.gz`、`update-darwin-universal.json`，Windows 的安装包、`.zip`/压缩包、`update-windows-amd64.json`。缺任何一项都不要发布，汇报给用户。齐全后让用户确认（或用户已授权）再发布：

```bash
gh release edit v<版本> --draft=false
```

### 6. 标记已发布

```bash
node scripts/release-notes.mjs --mark-released
```

提交这一处改动（`chore: 标记 v<版本> 已发布`）。

### 7. 开下一个版本

```bash
make version V=patch                    # 默认修订号 +1；用户指定了就按用户的（minor / major / x.y.z）
node scripts/release-notes.mjs --new
```

提交（`chore: 开始 v<新版本>`，包含 `VERSION`、`mygo.json`、`default_conf.go` 和新的发布说明）。把 main 合回 dev（`git checkout dev && git merge main`）。**推送 main、dev 前同样要用户已明确要求推送。**

## 改这一版的版本号

仅限状态为 `unreleased`、tag 还没推时：

```bash
make version V=<新版本>
git mv .github/release-notes/v<旧>.md .github/release-notes/v<新>.md
# 再把文件首行标题改成 # Transok <新版本>
node scripts/release-notes.mjs --check
```

## 汇报

用几句话说清：发的什么版本、Release 链接、各平台产物是否齐全、是否已标记已发布、下一个版本号。失败就给出失败的 job 和日志要点。
