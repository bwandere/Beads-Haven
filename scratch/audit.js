import puppeteer from 'puppeteer-core';
import fs from 'fs';
import path from 'path';

const SCREENSHOT_DIR = path.resolve('scratch/screenshots');
fs.mkdirSync(SCREENSHOT_DIR, { recursive: true });

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage', '--window-size=1440,900']
  });

  const page = await browser.newPage();

  const networkRequests = [];
  page.on('request', req => {
    networkRequests.push({
      url: req.url(),
      resourceType: req.resourceType(),
      method: req.method()
    });
  });

  page.on('response', async res => {
    const url = res.url();
    if (url.includes('.glb') || url.includes('.gltf') || url.includes('three') || url.includes('spline') || url.includes('model')) {
      console.log('INTERESTING RESPONSE:', url, res.status());
    }
  });

  page.on('console', msg => console.log('PAGE LOG:', msg.text()));

  console.log('Navigating to homepage...');
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/', { waitUntil: 'networkidle2' });

  // Check 3D hero details
  const heroInfo = await page.evaluate(() => {
    const canvases = Array.from(document.querySelectorAll('canvas')).map(c => ({
      width: c.width,
      height: c.height,
      className: c.className,
      id: c.id,
      parentElement: c.parentElement ? c.parentElement.outerHTML.slice(0, 300) : null,
      webgl: !!(c.getContext('webgl') || c.getContext('webgl2'))
    }));

    const heroSection = document.querySelector('section, header, [class*="hero"], [class*="Hero"]');
    
    // Check for Three.js, Spline, etc on window
    const globals = {
      hasTHREE: typeof window.THREE !== 'undefined',
      keys: Object.keys(window).filter(k => /three|spline|model|3d|canvas/i.test(k))
    };

    // Check all elements in hero for transforms, perspective, images
    const heroEl = document.querySelector('[class*="hero"], main > section:first-child, section');
    const heroHTML = heroEl ? heroEl.outerHTML : null;

    // Check computed styles on hero children
    let heroStyles = [];
    if (heroEl) {
      const allEls = [heroEl, ...heroEl.querySelectorAll('*')];
      heroStyles = allEls.map(el => {
        const cs = window.getComputedStyle(el);
        return {
          tag: el.tagName,
          class: el.className,
          transform: cs.transform,
          perspective: cs.perspective,
          transformStyle: cs.transformStyle,
          position: cs.position,
          zIndex: cs.zIndex,
          filter: cs.filter,
          boxShadow: cs.boxShadow
        };
      }).filter(s => s.transform !== 'none' || s.perspective !== 'none' || s.transformStyle === 'preserve-3d');
    }

    return {
      canvases,
      globals,
      heroHTML: heroHTML ? heroHTML.slice(0, 2000) : null,
      heroStyles
    };
  });

  console.log('Hero info:', JSON.stringify(heroInfo, null, 2));

  // Capture at 1440px
  await page.setViewport({ width: 1440, height: 900 });
  await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'home-1440.png'), fullPage: true });

  // Capture at 768px
  await page.setViewport({ width: 768, height: 1024 });
  await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'home-768.png'), fullPage: true });

  // Capture at 390px
  await page.setViewport({ width: 390, height: 844 });
  await page.screenshot({ path: path.join(SCREENSHOT_DIR, 'home-390.png'), fullPage: true });

  // Check all links and routes
  const links = await page.evaluate(() => {
    return Array.from(document.querySelectorAll('a')).map(a => ({
      href: a.href,
      text: a.innerText.trim()
    }));
  });
  console.log('Links on home page:', links);

  // Write requests to file
  fs.writeFileSync('scratch/requests.json', JSON.stringify(networkRequests, null, 2));
  fs.writeFileSync('scratch/hero-info.json', JSON.stringify(heroInfo, null, 2));

  await browser.close();
  console.log('Initial audit step complete.');
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
