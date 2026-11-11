#!/usr/bin/env node
// 发布说明（.github/release-notes/v<版本>.md）的维护脚本。版本号以根目录 VERSION 为准。
//
//   --status          查看当前版本的发布状态
//   --new             为当前版本新建 unreleased 说明文件
//   --check           打 tag 前的预检：文件存在、unreleased、有条目、无占位文字
//   --mark-released   发版完成后标为已发布（可带日期 YYYY-MM-DD，默认今天）
//   --changelog [out] 生成 MyGo 要的 changelog（## <版本> + 正文），默认 .mygo/changelog.md
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const ROOT = join(dirname(fileURLToPath(import.meta.url)), '..')
const NOTES_DIR = join(ROOT, '.github/release-notes')
const STATUS_RE = /^<!--\s*release-status:\s*(unreleased|released(?:\s+\d{4}-\d{2}-\d{2})?)\s*-->$/m
const SECTIONS = ['新增', '改进', '修复', '其他']

export function parseNotes(text) {
  const lines = text.replace(/\r\n/g, '\n').split('\n')
  const status = (STATUS_RE.exec(text)?.[1] ?? '').split(/\s+/)[0]
  const body = lines.filter((l, i) => i > 0 && !STATUS_RE.test(l)).join('\n').trim()
  const items = lines.filter((l) => /^\s*[-*]\s+\S/.test(l))
  return { title: lines[0] ?? '', status, body, items }
}

export function problems(version, text) {
  const n = parseNotes(text)
  const out = []
  if (n.title.trim() !== `# Transok ${version}`) out.push(`标题应为「# Transok ${version}」，现在是「${n.title}」`)
  if (n.status !== 'unreleased') out.push(`发布状态应为 unreleased，现在是「${n.status || '没有标记'}」`)
  if (n.items.length === 0) out.push('没有任何条目')
  if (/TODO|待补|占位/.test(n.body)) out.push('含有 TODO / 待补 / 占位文字')
  for (const m of n.body.matchAll(/^## (.+)$/gm)) {
    if (!SECTIONS.includes(m[1].trim())) out.push(`未知小节「${m[1]}」，只能用：${SECTIONS.join(' / ')}`)
  }
  return out
}

export function toChangelog(version, text) {
  // MyGo 以下一个「## 」作为该版本一节的结尾，所以把小节降一级。
  const body = parseNotes(text).body.replace(/^## /gm, '### ')
  return `## ${version}\n\n${body}\n`
}

const template = (v) => `# Transok ${v}
<!-- release-status: unreleased -->

## 新增

## 改进

## 修复
`

function main() {
  const version = readFileSync(join(ROOT, 'VERSION'), 'utf8').trim()
  const file = join(NOTES_DIR, `v${version}.md`)
  const [flag, arg] = process.argv.slice(2)
  const fail = (msg) => { console.error(`error: ${msg}`); process.exit(1) }
  const read = () => existsSync(file) ? readFileSync(file, 'utf8') : fail(`找不到 ${file}，先运行 --new`)

  switch (flag) {
    case '--status': {
      const n = parseNotes(read())
      console.log(`v${version}: ${n.status || '没有标记（按已发布处理）'}`)
      return
    }
    case '--new': {
      if (existsSync(file)) fail(`${file} 已存在`)
      writeFileSync(file, template(version))
      console.log(`已创建 ${file}`)
      return
    }
    case '--check': {
      const list = problems(version, read())
      if (list.length) fail(`v${version} 的发布说明没通过：\n  - ${list.join('\n  - ')}`)
      console.log(`v${version} 的发布说明通过（${parseNotes(read()).items.length} 条）`)
      return
    }
    case '--mark-released': {
      const text = read()
      if (parseNotes(text).status !== 'unreleased') fail('当前版本不是 unreleased')
      const date = arg ?? new Date().toISOString().slice(0, 10)
      writeFileSync(file, text.replace(STATUS_RE, `<!-- release-status: released ${date} -->`))
      console.log(`v${version} 已标为 released ${date}`)
      return
    }
    case '--changelog': {
      const out = arg ?? join(ROOT, '.mygo/changelog.md')
      mkdirSync(dirname(out), { recursive: true })
      writeFileSync(out, toChangelog(version, read()))
      console.log(`已生成 ${out}`)
      return
    }
    default:
      fail('用法：--status | --new | --check | --mark-released [日期] | --changelog [输出路径]')
  }
}

if (process.argv[1] === fileURLToPath(import.meta.url)) main()
