import assert from 'node:assert/strict'
import { test } from 'node:test'
import { parseNotes, problems, toChangelog } from './release-notes.mjs'

const ok = `# Transok 0.5.3
<!-- release-status: unreleased -->

## 新增

- **更新提示**：有新版本时弹出窗口。
`

test('合格的说明没有问题', () => {
  assert.deepEqual(problems('0.5.3', ok), [])
})

test('已发布、无条目、标题不符、未知小节都会被拦下', () => {
  assert.equal(problems('0.5.3', ok.replace('unreleased', 'released 2026-01-01')).length, 1)
  assert.equal(problems('0.5.3', '# Transok 0.5.3\n<!-- release-status: unreleased -->\n').length, 1)
  assert.equal(problems('0.5.4', ok).length, 1)
  assert.equal(problems('0.5.3', ok + '\n## 杂项\n').length, 1)
})

test('没有状态标记按已发布处理', () => {
  assert.equal(parseNotes('# Transok 0.1.0\n\n- a\n').status, '')
})

test('changelog 去掉标题和状态注释，换成版本标题', () => {
  const out = toChangelog('0.5.3', ok)
  assert.ok(out.startsWith('## 0.5.3\n\n### 新增'))
  assert.ok(!out.includes('release-status'))
  assert.ok(!out.includes('# Transok'))
})
