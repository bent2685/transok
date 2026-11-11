#!/usr/bin/env node
// Single source of truth for the project version is the VERSION file at the
// repo root. This script either:
//   - writes a new version to VERSION (when an argument is given), or
//   - syncs the current VERSION value out to mygo.json and default_conf.go.
//
// Usage:
//   node scripts/set-version.mjs 0.4.9   # bump to 0.4.9 and sync
//   node scripts/set-version.mjs patch   # patch | minor | major: bump from the current VERSION
//   node scripts/set-version.mjs         # re-sync current VERSION to consumers
//   node scripts/set-version.mjs --check # verify all locations match VERSION
import { readFileSync, writeFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

const ROOT = join(dirname(fileURLToPath(import.meta.url)), '..')
const VERSION_FILE = join(ROOT, 'VERSION')
const MYGO_JSON = join(ROOT, 'mygo.json')
const GO_CONF = join(ROOT, 'backend/consts/default_conf.go')

const arg = process.argv[2]
let mode = 'sync'
let newVersion = ''
if (arg === '--check') mode = 'check'
else if (arg === '-h' || arg === '--help') {
  console.log(readFileSync(new URL(import.meta.url)).toString().split('\n').slice(1, 10).join('\n'))
  process.exit(0)
} else if (arg) {
  mode = 'set'
  newVersion = arg
}

if (mode === 'set' && /^(patch|minor|major)$/.test(newVersion)) {
  const [maj, min, pat] = readFileSync(VERSION_FILE, 'utf8').trim().split('.').map(Number)
  newVersion = { major: `${maj + 1}.0.0`, minor: `${maj}.${min + 1}.0`, patch: `${maj}.${min}.${pat + 1}` }[newVersion]
}

if (mode === 'set') {
  if (!/^\d+\.\d+\.\d+([.-][0-9A-Za-z.-]+)?$/.test(newVersion)) {
    console.error(`error: '${newVersion}' is not a valid semver (e.g. 0.4.9)`)
    process.exit(1)
  }
  writeFileSync(VERSION_FILE, newVersion + '\n')
}

const version = readFileSync(VERSION_FILE, 'utf8').trim()
if (!version) {
  console.error('error: VERSION file is empty')
  process.exit(1)
}

const mygoRaw = readFileSync(MYGO_JSON, 'utf8')
const mygoData = JSON.parse(mygoRaw)
const mygoCurrent = mygoData.version

const goSrc = readFileSync(GO_CONF, 'utf8')
const goVersionRe = /("version":\s*)"([^"]*)"/
const goMatch = goSrc.match(goVersionRe)
if (!goMatch) {
  console.error(`error: could not locate version field in ${GO_CONF}`)
  process.exit(1)
}
const goCurrent = goMatch[2]

if (mode === 'check') {
  let ok = true
  const drift = (cur) => cur === version ? '' : '   <- DRIFT'
  if (mygoCurrent !== version) ok = false
  if (goCurrent !== version) ok = false
  console.log(`VERSION         ${version}`)
  console.log(`mygo.json       ${mygoCurrent}${drift(mygoCurrent)}`)
  console.log(`default_conf.go ${goCurrent}${drift(goCurrent)}`)
  if (!ok) {
    console.log('\nfix with: pnpm version:set')
    process.exit(2)
  }
  process.exit(0)
}

mygoData.version = version
writeFileSync(MYGO_JSON, JSON.stringify(mygoData, null, 2) + '\n')

writeFileSync(GO_CONF, goSrc.replace(goVersionRe, `$1"${version}"`))

console.log(`synced version ${version} ->`)
console.log(`  ${MYGO_JSON}`)
console.log(`  ${GO_CONF}`)
