import puppeteer from 'puppeteer-core';
import fs from 'fs';

async function run() {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/google-chrome',
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();
  await page.setViewport({ width: 1440, height: 900 });

  // 1. Visit product page
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/product/volt-runner', { waitUntil: 'domcontentloaded' });
  await new Promise(r => setTimeout(r, 2000));

  // Click on size 42
  const buttons = await page.$$('button');
  for (const b of buttons) {
    const text = await (await b.getProperty('innerText')).jsonValue();
    if (text.trim() === '42') {
      await b.click();
      console.log('Clicked size 42');
      break;
    }
  }

  // Click ADD TO CART
  for (const b of await page.$$('button')) {
    const text = await (await b.getProperty('innerText')).jsonValue();
    if (text.trim() === 'ADD TO CART') {
      await b.click();
      console.log('Clicked ADD TO CART');
      break;
    }
  }

  await new Promise(r => setTimeout(r, 1000));

  // 2. Go to /cart
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/cart', { waitUntil: 'domcontentloaded' });
  await new Promise(r => setTimeout(r, 2000));
  await page.screenshot({ path: 'scratch/screenshots/cart-with-items-1440.png', fullPage: true });

  const cartHTML = await page.evaluate(() => {
    return document.querySelector('main')?.innerHTML;
  });
  fs.writeFileSync('scratch/cart-populated.html', cartHTML);

  // 3. Go to /checkout
  await page.goto('https://pixel-perfect-capture-7092.lovable.app/checkout', { waitUntil: 'domcontentloaded' });
  await new Promise(r => setTimeout(r, 2000));
  await page.screenshot({ path: 'scratch/screenshots/checkout-with-items-1440.png', fullPage: true });

  const checkoutHTML = await page.evaluate(() => {
    return document.querySelector('main')?.innerHTML;
  });
  fs.writeFileSync('scratch/checkout-populated.html', checkoutHTML);

  console.log('Cart and Checkout populated audit complete!');
  await browser.close();
}

run().catch(console.error);
