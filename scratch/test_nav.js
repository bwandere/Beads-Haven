import puppeteer from 'puppeteer-core';
import fs from 'fs';
import path from 'path';

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage', '--window-size=1440,900']
  });

  const page = await browser.newPage();
  console.log('Navigating to prototype...');
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/', { waitUntil: 'domcontentloaded', timeout: 30000 });
  
  // Wait 4 seconds for full render and model to load
  await new Promise(r => setTimeout(r, 4000));

  console.log('Page loaded!');
  
  // Get all navigation links
  const links = await page.evaluate(() => {
    return Array.from(document.querySelectorAll('a')).map(a => ({
      href: a.getAttribute('href'),
      fullHref: a.href,
      text: a.innerText.trim(),
      className: a.className
    }));
  });
  console.log('All links on page:', JSON.stringify(links, null, 2));

  // Get entire inner HTML of body or structure
  const pageStructure = await page.evaluate(() => {
    const header = document.querySelector('header')?.outerHTML;
    const nav = document.querySelector('nav')?.outerHTML;
    const footer = document.querySelector('footer')?.outerHTML;
    const sections = Array.from(document.querySelectorAll('section')).map(s => ({
      className: s.className,
      heading: s.querySelector('h1, h2, h3')?.innerText,
      textPreview: s.innerText.slice(0, 200)
    }));
    return { header, nav, footer, sections };
  });

  fs.writeFileSync('scratch/page-structure.json', JSON.stringify(pageStructure, null, 2));
  console.log('Sections found:', pageStructure.sections);

  // Take screenshot 1440
  await page.setViewport({ width: 1440, height: 900 });
  await page.screenshot({ path: 'scratch/screenshots/home-1440.png', fullPage: true });

  // Take screenshot 768
  await page.setViewport({ width: 768, height: 1024 });
  await page.screenshot({ path: 'scratch/screenshots/home-768.png', fullPage: true });

  // Take screenshot 390
  await page.setViewport({ width: 390, height: 844 });
  await page.screenshot({ path: 'scratch/screenshots/home-390.png', fullPage: true });

  console.log('Screenshots saved!');

  await browser.close();
}

run().catch(console.error);
