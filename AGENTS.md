# Development Guide

This is a Go desktop app built with MyGo's native UI.

You **MUST** absolutely obey this guide.

You need to communicate with the user in Chinese.

> 客户端已从 Wails + WebView 换到 MyGo native UI：`frontend/`、`wails.json`、`build/` 里的 Wails 产物均已删除。当前视觉是对迁移前 Web 端的 1:1 还原。

---

## Tech Stack
- Golang 1.27.3+
- MyGo native UI（`github.com/egoist/mygo`，`go tool mygo`）
- gin + 浏览器下载页（给被分享方，与客户端 UI 无关）

---

## MyGo native UI 约定

客户端全部窗口用 MyGo native UI（Go 绘制，无 WebView、无 HTML/CSS/JS）。HTTP 下载页（`backend/app` 的 gin + `templates/downpage`）与 `dash/` 是给被分享方的浏览器页面，不属于客户端 UI。

- Go 版本要求 1.27.3+（`go.mod` 已抬升）。
- CLI 以 tool 形式钉住：`go tool mygo`（v0.8.3），不要另装全局版本。
- 维护护栏技能在 `.agents/skills/mygo-maintenance/`，并用软链挂在 `.claude/skills/`。改动 UI 前先读它。
- 视图是「状态的函数」，每帧重建。**元素（`ui.Element`）只在当前 build pass 内有效**：跨帧只能存 `ui.Handle`、widget state 和业务数据，不能存元素；`*ui.Context` 不进结构体、不进 goroutine。
- 有状态控件先 `c.Key(key)` 再构造；列表先配 `.ItemKey`/`.Grow` 再 `.Rows`。
- 慢 IO / 重计算留在 worker，结果用 `Window.Update` 回投；绘制回调只画，不做业务。
- 验证顺序：`go tool mygo vet .` → `go tool mygo i18n check` → `CGO_ENABLED=0 go test ./...` → `go tool mygo dev` 眼看效果。

---

## UI/UX Design

@DESIGN.md in the packages root directory
**MUST** be followed for all UI/UX design decisions.

---

## Conventions
- Prefer named exports for components
- Use TypeScript strict mode

---

## Code Style
- Use functional components with TypeScript
- Prefer async/await over .then() chains
- Use early returns for cleaner code
- Keep components small and focused

---

## Coding Guidelines (Core Principle: Restraint)

### Do only what is requested

* **Do not proactively add features.** A bug fix does not need a "handy cleanup" of surrounding code; a simple feature does not require extra configuration options.
* **Do not write error handling / fallback / validation for non-existent scenarios.** Only validate at system boundaries (user input, external APIs). Internal code should trust each other.
* **Do not abstract for "potential future needs."** Use one-off code for one-off operations. Three lines of similar code are better than a premature abstraction.
* **Do not add comments / docstrings / type annotations to code you haven't modified.** Only comment where the logic is genuinely non-obvious.
* **Do not perform unrelated renaming / reordering / style unification.** Keep the diff minimal.

### Read before modifying

* **Modifying unread files based on impression is strictly forbidden.** If a user says "modify function X," first `Read` the entire file to understand the context, then `Edit`.
* Before modifying a public type / function, first `Grep` the entire repository for all call sites to confirm the blast radius.
* When changing styles, first check the patterns used in sibling files. **Copy your neighbors; do not invent your own.**

### Avoid destructive "shortcuts"

* When encountering obstacles, **do not resolve them by deleting, bypassing, or disabling checks.** Locate the root cause first.
* Operations like `--no-verify`, `--no-gpg-sign`, `git reset --hard`, `git push --force`, and `rm -rf` are not allowed unless **explicitly** authorized by the user.
* When encountering unfamiliar files / branches / configurations, investigate before handling them — they might be part of what the user is currently working on.

### Tool Preferences

* **Use dedicated tools over bash whenever possible:**
* Reading files → Use `Read`, not `cat` / `head` / `tail`
* Modifying files → Use `Edit`, not `sed` / `awk`
* Finding files → Use `Glob`, not `find` / `ls`
* Finding content → Use `Grep`, not `grep` / `rg`
* Writing files → Use `Write`, not `cat <<EOF` (multiple * **concurrently** --- Execute calls independent message), multiple rather same serially. than the tool within/>

---

## Commit Guidelines

### Conventional Commits Style

Format: `<type>: <subject>`, where the subject is written in Chinese.

| type | Purpose |
| --- | --- |
| `feat` | New features |
| `fix` | Bug fixes |
| `refactor` | Internal refactoring that does not change external behavior |
| `chore` | Miscellaneous (dependencies, configurations, build scripts) |
| `docs` | Documentation (README / AGENTS.md / comments) |
| `test` | Adding tests / modifying tests |
| `style` | Pure formatting adjustments (no semantic changes) |
| `perf` | Performance optimizations |

Examples:

```
feat: xxx
fix: xxx
refactor: xxx
test: xxx

```

### Small Steps, Frequent Commits

**Strict constraint: One commit per task.**

* If a PR contains 3 independent changes → **3 commits**, not 1 massive commit stuffed with everything.
* Refactoring + New Feature → Split them: first commit the refactoring (behavior unchanged), then commit the new feature.
* Bug fix + A typo found along the way → Split them: first commit the fix, then commit the chore.
* Writing code + Writing tests → These can be committed together, but **you must not write code without writing tests**.

The criterion for judgment: **If you would want to precisely roll back this specific change later, give it a separate commit.**

### Commit Message Body

* Explain **why**, not **what was done** (the diff already tells you what was done).
* 1-3 sentences are enough. No need to write a novel.
* Reference associated issues / PRs using the `owner/repo#123` format.

### Prohibited Actions

* **Never `git commit --amend**` unless the user explicitly tells you to "amend". To fix the previous commit, create a new fix commit instead.
* **Never use `--no-verify**`. If a pre-commit hook fails, fix the underlying issue instead of bypassing it.
* **Never `git push --force` to main / master.**
* **Never push when the user hasn't explicitly asked you to.** Local commits are perfectly fine, but pushing introduces side effects.
* **Never commit `.env` files, secret keys, or binaries larger than 1MB.**
* **Never append trailers / sign-off footers to commit messages.** This includes `Co-Authored-By:`, `Signed-off-by:`, `Generated-by:`, `Reviewed-by:` and any similar `Key: Value` footer at the end of the message. The commit body ends with the last sentence — nothing after it.

### Mandatory Checklist When Committing

* Run `git status` first to review the staged changes.
* Run `git diff --staged` to double-check that no unrelated files have slipped in.
* Run `git log -5 --oneline` to check the recent commit message style of the repository, and **copy that style**.
