import puppeteer from 'puppeteer-core';
import fs from 'fs';

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/', { waitUntil: 'domcontentloaded' });
  await new Promise(r => setTimeout(r, 2000));

  const tokens = await page.evaluate(() => {
    const el = document.documentElement;
    const cs = window.getComputedStyle(el);
    const varNames = [
      '--radius',
      '--background',
      '--foreground',
      '--card',
      '--card-foreground',
      '--popover',
      '--popover-foreground',
      '--primary',
      '--primary-foreground',
      '--primary-glow',
      '--secondary',
      '--secondary-foreground',
      '--muted',
      '--muted-foreground',
      '--accent',
      '--accent-foreground',
      '--ink',
      '--ink-foreground',
      '--ink-muted',
      '--surface',
      '--destructive',
      '--destructive-foreground',
      '--border',
      '--input',
      '--ring',
      '--gradient-hero',
      '--gradient-promo',
      '--shadow-card',
      '--shadow-pop'
    ];

    const result = {};
    for (const v of varNames) {
      result[v] = cs.getPropertyValue(v).trim();
    }

    // Also get RGB/HEX by creating dummy elements
    const rgbMap = {};
    for (const v of varNames) {
      if (result[v] && !result[v].includes('gradient') && !result[v].includes('shadow') && !result[v].includes('rem')) {
        const dummy = document.createElement('div');
        dummy.style.color = `var(${v})`;
        document.body.appendChild(dummy);
        const computedColor = window.getComputedStyle(dummy).color;
        document.body.removeChild(dummy);
        rgbMap[v] = computedColor;
      }
    }

    // Typography
    const bodyCs = window.getComputedStyle(document.body);
    const h1 = document.querySelector('h1');
    const h1Cs = h1 ? window.getComputedStyle(h1) : null;
    const p = document.querySelector('p');
    const pCs = p ? window.getComputedStyle(p) : null;

    return {
      tokens: result,
      rgbMap,
      typography: {
        bodyFontFamily: bodyCs.fontFamily,
        bodyFontSize: bodyCs.fontSize,
        bodyLineHeight: bodyCs.lineHeight,
        h1FontFamily: h1Cs ? h1Cs.fontFamily : null,
        h1FontSize: h1Cs ? h1Cs.fontSize : null,
        h1FontWeight: h1Cs ? h1Cs.fontWeight : null,
        h1LetterSpacing: h1Cs ? h1Cs.letterSpacing : null
      }
    };
  });

  console.log(JSON.stringify(tokens, null, 2));
  fs.writeFileSync('scratch/design-tokens-computed.json', JSON.stringify(tokens, null, 2));
  await browser.close();
}

run().catch(console.error);
