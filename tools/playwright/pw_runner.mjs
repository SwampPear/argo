#!/usr/bin/env node
import fs from 'node:fs';
import { chromium } from 'playwright'; // npm i playwright
const println = o => process.stdout.write(JSON.stringify(o) + '\n')
const now = () => new Date().toISOString()

// read a single JSON config from stdin
const cfg = JSON.parse(fs.readFileSync(0, 'utf8') || '{}')
const {
  startUrls = [],
  storageState,
  harPath,
  screenshotDir,
  headless = false,       // <- headful by default
  keepOpen = true         // <- keep the browser open until user closes
} = cfg

const ensureDir = p => { if (p && !fs.existsSync(p)) fs.mkdirSync(p, { recursive: true }) }
ensureDir(screenshotDir)
if (harPath) ensureDir(require('node:path').dirname(harPath))

const log = (m) => println({ type:'log', timestamp: now(), module:'Playwright', ...m })
const obs = (m) => println({ type:'obs', timestamp: now(), ...m })

const run = async () => {
  const browser = await chromium.launch({ headless })
  const context = await browser.newContext({
    storageState: (storageState && fs.existsSync(storageState)) ? storageState : undefined,
    recordHar: harPath ? { path: harPath, mode: 'minimal' } : undefined
  })
  const page = await context.newPage()

  // --- stream events ---
  page.on('console', msg => {
    log({ action:'console', target: page.url(), status:'OK', duration:'0s', confidence:0.4,
          summary: `${msg.type()}: ${msg.text()}` })
  })
  page.on('request', req => {
    log({ action:'request', target:req.url(), status:'OK', duration:'0s', confidence:0.5,
          summary:`${req.method()} ${req.resourceType()}` })
  })
  page.on('response', res => {
    log({ action:'response', target:res.url(), status:String(res.status()), duration:'0s',
          confidence:0.5, summary:`status=${res.status()}` })
  })
  page.on('pageerror', err => {
    log({ action:'pageerror', target: page.url(), status:'Error', duration:'0s',
          confidence:0, summary:String(err) })
  })
  page.on('framenavigated', async f => {
    if (!screenshotDir) return
    const file = `${screenshotDir}/${Date.now()}.png`
    try {
      await page.screenshot({ path:file, fullPage:true })
      obs({ kind:'screenshot', target:f.url(), evidence:[file] })
    } catch {}
  })

  // seed navigations (optional)
  for (const url of startUrls) {
    try {
      const t0 = Date.now()
      log({ action:'navigate', target:url, status:'OK', duration:'0s', confidence:0.4, summary:`open ${url}` })
      await page.goto(url, { waitUntil:'domcontentloaded', timeout:30000 })
      await page.waitForLoadState('networkidle', { timeout:15000 })
      log({ action:'navigate', target:url, status:'OK', duration:`${Date.now()-t0}ms`, confidence:0.6, summary:'loaded' })
    } catch (e) {
      log({ action:'navigate', target:url, status:'Error', duration:'0s', confidence:0, summary:String(e?.message||e) })
    }
  }

  const graceful = async () => { try { await browser.close() } catch {} process.exit(0) }
  process.on('SIGINT', graceful); process.on('SIGTERM', graceful)

  if (keepOpen) {
    log({ action:'ready', target:'-', status:'OK', duration:'0s', confidence:1, summary:'Headful session started' })
    await new Promise(resolve => browser.on('disconnected', resolve))
    await graceful()
  } else {
    await graceful()
  }
}

run().catch(e => {
  log({ action:'boot', target:'-', status:'Error', duration:'0s', confidence:0, summary:String(e?.message||e) })
  process.exit(1)
})
