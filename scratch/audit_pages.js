import puppeteer from 'puppeteer-core';
import fs from 'fs';
import path from 'path';

const routes = [
  '/',
  '/sneakers',
  '/jerseys',
  '/product/volt-runner',
  '/cart',
  '/checkout',
  '/about',
  '/contact'
];

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage', '--window-size=1440,900']
  });

  const page = await browser.newPage();
  const results = {};

  for (const route of routes) {
    const url = `https://pixel-perfect-capture-7092.lovable.app${route}`;
    console.log(`Checking ${route}...`);
    try {
      const resp = await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 20000 });
      await new Promise(r => setTimeout(r, 2000));
      
      const currentUrl = page.url();
      const status = resp ? resp.status() : 'unknown';
      const title = await page.title();
      
      const htmlSnippet = await page.evaluate(() => {
        const main = document.querySelector('main') || document.body;
        return {
          h1: document.querySelector('h1')?.innerText,
          headings: Array.from(document.querySelectorAll('h1, h2, h3')).map(h => h.innerText),
          buttons: Array.from(document.querySelectorAll('button')).map(b => b.innerText.trim()).filter(Boolean),
          inputs: Array.from(document.querySelectorAll('input, select, textarea')).map(i => ({
            name: i.name,
            type: i.type,
            placeholder: i.placeholder
          })),
          textSnippet: main.innerText.slice(0, 500)
        };
      });

      // Screenshots at 1440, 768, 390
      const safeName = route.replace(/\//g, '_') || '_root';
      
      await page.setViewport({ width: 1440, height: 900 });
      await page.screenshot({ path: `scratch/screenshots/${safeName}-1440.png`, fullPage: true });

      await page.setViewport({ width: 768, height: 1024 });
      await page.screenshot({ path: `scratch/screenshots/${safeName}-768.png`, fullPage: true });

      await page.setViewport({ width: 390, height: 844 });
      await page.screenshot({ path: `scratch/screenshots/${safeName}-390.png`, fullPage: true });

      results[route] = {
        currentUrl,
        status,
        title,
        data: htmlSnippet
      };
    } catch (e) {
      console.error(`Error on ${route}:`, e.message);
      results[route] = { error: e.message };
    }
  }

  fs.writeFileSync('scratch/routes-audit.json', JSON.stringify(results, null, 2));
  console.log('All routes audited successfully.');
  await browser.close();
}

run().catch(console.error);
