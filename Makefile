# Transok 开发与构建入口。直接运行 `make` 查看所有命令。
#
# 下载页（dash/）构建后嵌进 Go 程序（backend/app/templates/downpage，go:embed），
# 所以 dev / build 都会先构建它；dash 源码没变时会自动跳过。

SHELL := /bin/bash

# 构建目标平台，默认当前机器；mygo 的写法，如 darwin/arm64、windows/amd64。
PLATFORM ?= $(shell go env GOOS)/$(shell go env GOARCH)
# 传给 `mygo dev` 的额外参数，如 DEV_ARGS="-lang ja"。
DEV_ARGS ?=

MYGO      := go tool mygo
DASH_PKG  := transok-download
DASH_OUT  := backend/app/templates/downpage/index.html
DASH_SRC  := $(shell find dash/src dash/public -type f 2>/dev/null) \
             dash/index.html dash/package.json dash/vite.config.ts \
             dash/tailwind.config.js dash/postcss.config.js dash/tsconfig.json
NODE_STAMP := node_modules/.install-stamp

.DEFAULT_GOAL := help
.PHONY: help install dash dash-force dash-dev dev build build-mac build-linux build-win \
        test vet fmt check version clean

help: ## 显示所有命令
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo
	@echo "  常用写法:"
	@echo "    make dev                      桌面应用热重载（先构建下载页）"
	@echo "    make dash-dev                 下载页热更新，代理到应用的 9482 端口（另开终端，应用里先开启分享）"
	@echo "    make build PLATFORM=linux/amd64"
	@echo "    make version V=0.5.3          改版本号并同步到各处"

$(NODE_STAMP): pnpm-lock.yaml package.json dash/package.json
	pnpm install
	@touch $@

install: $(NODE_STAMP) ## 安装前端依赖（依赖文件变了才重装）

$(DASH_OUT): $(DASH_SRC) $(NODE_STAMP)
	pnpm -F $(DASH_PKG) build

dash: $(DASH_OUT) ## 构建下载页（源码没变则跳过）

dash-force: $(NODE_STAMP) ## 强制重新构建下载页
	pnpm -F $(DASH_PKG) build

dash-dev: $(NODE_STAMP) ## 下载页开发服务器 http://localhost:3000（需应用已开启分享）
	pnpm -F $(DASH_PKG) dev

dev: dash ## 桌面应用热重载（Go 代码变化自动重编；r 重启、q 退出）
	$(MYGO) dev $(DEV_ARGS)

build: dash ## 构建当前平台，或 PLATFORM=... 指定（macOS 的 dmg 会追加修复脚本）
	$(MYGO) build -platform $(PLATFORM)
	@case "$(PLATFORM)" in darwin/*) scripts/add-fix-script.sh dist/$(subst /,-,$(PLATFORM))/*.dmg ;; esac

build-mac: dash ## 构建 macOS arm64 与 amd64 两个包（dmg 追加「修复已损坏.command」）
	$(MYGO) build -platform darwin/arm64,darwin/amd64
	@for d in dist/darwin-arm64/*.dmg dist/darwin-amd64/*.dmg; do scripts/add-fix-script.sh $$d; done

build-linux: dash ## 构建 Linux amd64 与 arm64（tar.gz、deb 与 install.sh）
	$(MYGO) build -platform linux/amd64,linux/arm64

build-win: dash ## 构建 Windows amd64 单个 exe（不含安装包，NSIS 在 macOS 上做不了）
	$(MYGO) build -platform windows/amd64 -skip-installer

test: ## 运行 Go 测试
	CGO_ENABLED=0 go test ./...

vet: ## MyGo 的 UI 生命周期检查
	$(MYGO) vet .

fmt: ## 检查 Go 格式（有输出即有文件需要 gofmt）
	@out=$$(gofmt -l backend desktop main.go); \
	if [ -n "$$out" ]; then echo "$$out"; echo "以上文件需要 gofmt -w"; exit 1; fi

check: fmt vet ## 提交前的全套检查：格式、vet、文案、版本一致、测试
	$(MYGO) i18n check
	node scripts/set-version.mjs --check
	node --test scripts/*.test.mjs
	CGO_ENABLED=0 go test ./...

version: ## 改版本号并同步到各处：make version V=0.5.3 或 V=patch|minor|major（不带 V 则只同步）
	node scripts/set-version.mjs $(V)

clean: ## 删除构建产物（dist、.mygo），不动已提交的下载页产物
	rm -rf dist .mygo
