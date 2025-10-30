import { chromium } from 'playwright'; // npm i -D playwright
const url = process.argv[2] || 'https://example.com'
const userDataDir = process.env.PW_USER_DATA || `${process.env.HOME}/.wails-playwright`
const browser = await chromium.launchPersistentContext(userDataDir, { headless: false })
const page = await browser.newPage()
await page.goto(url)
// keep Node alive so the window stays open; exit when the browser closes
browser.on('close', () => process.exit(0))
process.on('SIGINT', async () => { await browser.close(); process.exit(0) })
process.stdin.resume()