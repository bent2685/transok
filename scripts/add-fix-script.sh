#!/usr/bin/env bash
# 给 `mygo build` 生成的 dmg 追加「修复已损坏.command」。
#
# 未做 Apple 开发者签名 / 公证的应用，从网上下载后会被 Gatekeeper 拦截，
# 提示"已损坏，无法打开"。这个脚本清除应用的隔离属性来绕过拦截。
# MyGo 的 dmg 配置没有"追加文件"的选项，所以在构建后处理。
#
# 用法: add-fix-script.sh <dmg 路径>
#
# 已用 Developer ID 签名的 dmg 不需要它，会直接跳过；重复运行也是安全的。
set -euo pipefail

DMG="${1:?用法: $0 <dmg 路径>}"
[ -f "$DMG" ] || { echo "[add-fix-script] 找不到 $DMG" >&2; exit 1; }

APP_NAME="Transok"
SCRIPT_NAME="修复已损坏.command"

if codesign -dv "$DMG" 2>&1 | grep -q 'Authority=Developer ID'; then
  echo "[add-fix-script] $(basename "$DMG") 已用 Developer ID 签名，无需修复脚本，跳过"
  exit 0
fi
WAS_SIGNED=0
codesign -dv "$DMG" >/dev/null 2>&1 && WAS_SIGNED=1

WORK="$(mktemp -d -t transok-fixdmg)"
MOUNT=""
cleanup() {
  [ -n "$MOUNT" ] && hdiutil detach "$MOUNT" -quiet -force 2>/dev/null || true
  rm -rf "$WORK"
}
trap cleanup EXIT

attach() { # attach <image> [额外参数]，输出挂载点
  # 过滤 hdiutil attach 的弃用警告（diskutil image 只有较新的 macOS 才有），保留真正的错误。
  hdiutil attach -nobrowse -noverify "${@:2}" "$1" 2> >(grep -v 'is deprecated' >&2) \
    | tail -1 | awk -F'\t' '{print $NF}'
}

# 已经带脚本就不重复处理。
MOUNT="$(attach "$DMG" -readonly)"
if [ -e "$MOUNT/$SCRIPT_NAME" ]; then
  echo "[add-fix-script] $(basename "$DMG") 已包含 ${SCRIPT_NAME}，跳过"
  exit 0
fi
hdiutil detach "$MOUNT" -quiet && MOUNT=""

echo "[add-fix-script] 转为可写镜像..."
hdiutil convert "$DMG" -quiet -format UDRW -o "$WORK/rw.dmg"
# 压缩镜像转回可写后余量很小，留出放脚本的空间。
hdiutil resize -quiet -size "$(( $(stat -f %z "$WORK/rw.dmg") / 1048576 + 5 ))m" "$WORK/rw.dmg"

MOUNT="$(attach "$WORK/rw.dmg")"
cat > "$MOUNT/$SCRIPT_NAME" <<'CMD_EOF'
#!/usr/bin/env bash
# 修复 macOS 对未签名应用报"已损坏，无法打开"的问题。
# 原理：清除 Gatekeeper 加的 com.apple.quarantine 等扩展属性。
set -e

APP_NAME="Transok"
TARGETS=(
  "/Applications/${APP_NAME}.app"
  "$HOME/Applications/${APP_NAME}.app"
)

echo "======================================"
echo "  Transok 修复已损坏应用工具"
echo "======================================"
echo

FOUND=0
for APP in "${TARGETS[@]}"; do
  if [ -d "$APP" ]; then
    echo "[*] 找到应用：$APP"
    echo "[*] 正在清除隔离属性..."
    sudo xattr -cr "$APP"
    echo "[✓] 已修复：$APP"
    FOUND=1
  fi
done

if [ "$FOUND" -eq 0 ]; then
  echo "[!] 未在 /Applications 或 ~/Applications 找到 ${APP_NAME}.app"
  echo "[!] 请先把 ${APP_NAME}.app 拖到 Applications 文件夹后再运行本脚本。"
  echo
  read -n 1 -s -r -p "按任意键关闭窗口..."
  exit 1
fi

echo
echo "完成！现在可以正常打开 ${APP_NAME} 了。"
echo
read -n 1 -s -r -p "按任意键关闭窗口..."
CMD_EOF
chmod +x "$MOUNT/$SCRIPT_NAME"
hdiutil detach "$MOUNT" -quiet && MOUNT=""

echo "[add-fix-script] 压回只读镜像..."
hdiutil convert "$WORK/rw.dmg" -quiet -format UDZO -ov -o "$WORK/out.dmg"
mv -f "$WORK/out.dmg" "$DMG"

# 转换会丢掉 MyGo 加的签名；原来有签名的话补上同样的临时签名。
[ "$WAS_SIGNED" -eq 1 ] && codesign -s - -f "$DMG" >/dev/null 2>&1 || true

echo "[add-fix-script] 已写入 ${SCRIPT_NAME} -> $(basename "$DMG")"
