import puppeteer from 'puppeteer-core';
import fs from 'fs';

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();
  const hexValues = await page.evaluate(() => {
    const tokens = {
      background: 'oklch(97.5% .006 85)',
      foreground: 'oklch(20% .03 255)',
      card: 'oklch(100% 0 0)',
      cardForeground: 'oklch(20% .03 255)',
      primary: 'oklch(55% .21 258)',
      primaryForeground: 'oklch(98% .01 85)',
      primaryGlow: 'oklch(68% .18 250)',
      secondary: 'oklch(93% .012 85)',
      secondaryForeground: 'oklch(20% .03 255)',
      muted: 'oklch(94% .008 85)',
      mutedForeground: 'oklch(50% .025 255)',
      accent: 'oklch(74% .17 60)',
      accentForeground: 'oklch(20% .03 255)',
      ink: 'oklch(17% .035 258)',
      inkForeground: 'oklch(96% .015 85)',
      inkMuted: 'oklch(70% .03 255)',
      surface: 'oklch(92% .01 250)',
      destructive: 'oklch(58% .22 27)',
      border: 'oklch(88% .012 85)',
      input: 'oklch(88% .012 85)',
      ring: 'oklch(55% .21 258)',
      heroMeshOrange: '#f39a2b',
      heroMeshBlue: '#4c8dff'
    };

    const canvas = document.createElement('canvas');
    canvas.width = 1;
    canvas.height = 1;
    const ctx = canvas.getContext('2d', { willReadFrequently: true });

    function toHex(col) {
      ctx.clearRect(0,0,1,1);
      ctx.fillStyle = col;
      ctx.fillRect(0,0,1,1);
      const [r, g, b] = ctx.getImageData(0,0,1,1).data;
      return '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('');
    }

    const hexMap = {};
    for (const [k, v] of Object.entries(tokens)) {
      hexMap[k] = { oklch: v, hex: toHex(v) };
    }
    return hexMap;
  });

  console.log(JSON.stringify(hexValues, null, 2));
  fs.writeFileSync('scratch/hex-tokens.json', JSON.stringify(hexValues, null, 2));
  await browser.close();
}

run().catch(console.error);
