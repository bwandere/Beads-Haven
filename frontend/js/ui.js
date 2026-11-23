/**
 * Beads Haven - UI Utilities & Shared Component Logic
 */

import { getCartCount, addToCart, getWishlist, toggleWishlist } from './cart.js';
import { COLLECTIONS, BRAND, isSingleOptionProduct } from './collections.js';

export function formatCurrency(amount) {
  if (amount == null) return 'KSh 0';
  return `KSh ${Number(amount).toLocaleString('en-KE')}`;
}

/**
 * Toast Notification System
 */
let toastContainer = null;

export function showToast(message, type = 'default') {
  if (!toastContainer) {
    toastContainer = document.createElement('div');
    toastContainer.className = 'toast-container';
    document.body.appendChild(toastContainer);
  }

  const toast = document.createElement('div');
  toast.className = `toast ${type === 'error' ? 'toast-error' : ''}`;
  toast.innerHTML = `
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--accent); flex-shrink: 0;">
      <path d="M20 6 9 17l-5-5"/>
    </svg>
    <span>${message}</span>
  `;

  toastContainer.appendChild(toast);

  setTimeout(() => {
    toast.classList.add('toast-leave');
    toast.addEventListener('animationend', () => {
      toast.remove();
    });
  }, 2600);
}

/**
 * Render desktop dropdown & mobile accordion from collections.js
 */
function renderNavCollections() {
  const dropdown = document.querySelector('.shop-dropdown-menu');
  if (dropdown && !dropdown.dataset.rendered) {
    dropdown.dataset.rendered = 'true';
    let gridHtml = '<div class="shop-dropdown-grid">';
    COLLECTIONS.forEach(col => {
      gridHtml += `
        <div class="dropdown-col">
          <a href="/shop.html?type=${col.slug}" class="dropdown-col-title">${col.name}</a>
          ${col.subcategories.map(sub => `
            <a href="/shop.html?type=${col.slug}&category=${encodeURIComponent(sub)}" class="dropdown-sub-link">${sub}</a>
          `).join('')}
          ${col.menuExtra ? col.menuExtra.map(extra => `
            <a href="${extra.url}" class="dropdown-sub-link">${extra.name}</a>
          `).join('') : ''}
        </div>
      `;
    });
    gridHtml += '</div>';
    gridHtml += `
      <div class="dropdown-footer">
        <span style="font-size: 0.8125rem; color: var(--muted-foreground);">Handcrafted with glass seed beads in Nairobi, Kenya</span>
        <a href="/shop.html" class="dropdown-footer-link">
          Browse all beadwork &amp; art
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
        </a>
      </div>
    `;
    dropdown.innerHTML = gridHtml;
  }

  const mobilePanel = document.querySelector('.mobile-accordion-panel');
  if (mobilePanel && !mobilePanel.dataset.rendered) {
    mobilePanel.dataset.rendered = 'true';
    let panelHtml = `
      <div style="padding-bottom: 0.5rem; margin-bottom: 0.5rem; border-bottom: 1px solid rgba(250, 245, 236, 0.1);">
        <a href="/shop.html" class="mobile-col-title" style="color: var(--primary-foreground); font-size: 1.125rem;">Shop All Collections →</a>
      </div>
    `;
    COLLECTIONS.forEach(col => {
      panelHtml += `
        <div class="mobile-col-group">
          <a href="/shop.html?type=${col.slug}" class="mobile-col-title">${col.name}</a>
          ${col.subcategories.map(sub => `
            <a href="/shop.html?type=${col.slug}&category=${encodeURIComponent(sub)}" class="mobile-sub-link">${sub}</a>
          `).join('')}
          ${col.menuExtra ? col.menuExtra.map(extra => `
            <a href="${extra.url}" class="mobile-sub-link">${extra.name}</a>
          `).join('') : ''}
        </div>
      `;
    });
    mobilePanel.innerHTML = panelHtml;
  }
}

/**
 * Common Header Setup
 */
export function initHeader() {
  const badge = document.querySelector('.cart-badge');
  const updateBadge = () => {
    if (!badge) return;
    const count = getCartCount();
    if (count > 0) {
      badge.textContent = count;
      badge.style.display = 'grid';
    } else {
      badge.style.display = 'none';
    }
  };

  updateBadge();
  window.addEventListener('cart-updated', updateBadge);

  // Render navigation dropdown & accordion
  renderNavCollections();

  // Active navigation link highlighting
  const currentPath = window.location.pathname;
  const urlParams = new URLSearchParams(window.location.search);
  const currentType = urlParams.get('type') || '';
  const currentCat = urlParams.get('category') || '';

  // Top-level nav links
  document.querySelectorAll('.nav-link, .nav-mobile-link').forEach(link => {
    const href = link.getAttribute('href');
    if (!href) return;
    if (
      (href === '/' && (currentPath === '/' || currentPath === '/index.html')) ||
      (href === '/about.html' && currentPath.includes('about')) ||
      (href === '/contact.html' && currentPath.includes('contact')) ||
      (href.includes('shop.html') && currentPath.includes('shop'))
    ) {
      link.classList.add('active');
    }
  });

  // Highlight specific collection or subcategory in dropdown
  if (currentType || currentCat) {
    document.querySelectorAll('.dropdown-sub-link, .dropdown-col-title, .mobile-sub-link, .mobile-col-title').forEach(link => {
      const href = link.getAttribute('href');
      if (!href) return;
      if (currentCat && href.includes(`category=${encodeURIComponent(currentCat)}`)) {
        link.classList.add('active');
      } else if (!currentCat && currentType && href.endsWith(`type=${currentType}`)) {
        link.classList.add('active');
      }
    });
  }

  // Mobile menu toggle
  const menuBtn = document.querySelector('.mobile-menu-toggle');
  const mobileNav = document.querySelector('.nav-mobile');
  if (menuBtn && mobileNav) {
    menuBtn.addEventListener('click', () => {
      mobileNav.classList.toggle('open');
      const isOpen = mobileNav.classList.contains('open');
      menuBtn.innerHTML = isOpen
        ? `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>`
        : `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 5h16"/><path d="M4 12h16"/><path d="M4 19h16"/></svg>`;
    });
  }

  // Mobile accordion toggle
  const accordionTrigger = document.querySelector('.mobile-accordion-trigger');
  const accordion = document.querySelector('.mobile-accordion');
  if (accordionTrigger && accordion) {
    accordionTrigger.addEventListener('click', () => {
      const isOpen = accordion.classList.toggle('is-open');
      accordionTrigger.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
    });
  }

  // Header Search button handler
  const searchBtn = document.querySelector('.header-search-btn');
  if (searchBtn) {
    searchBtn.addEventListener('click', () => {
      window.location.href = '/shop.html';
    });
  }
}

/**
 * 3D Card Tilt handler
 */
export function bindCardTilt(card) {
  card.addEventListener('mousemove', e => {
    const rect = card.getBoundingClientRect();
    const x = (e.clientX - rect.left) / rect.width - 0.5;
    const y = (e.clientY - rect.top) / rect.height - 0.5;
    card.style.transform = `perspective(900px) rotateY(${x * 8}deg) rotateX(${-y * 8}deg)`;
  });

  card.addEventListener('mouseleave', () => {
    card.style.transform = '';
  });
}

/**
 * Renders a single product card element matching prototype DOM & behavior
 */
export function renderProductCard(product) {
  const card = document.createElement('div');
  card.className = 'product-card';
  card.dataset.id = product.id;

  const wishlist = getWishlist();
  const isWished = wishlist.includes(product.id);

  const saleBadge = product.oldPrice
    ? `<span class="badge-sale">Sale</span>`
    : '';

  const oldPriceHtml = product.oldPrice
    ? `<span class="product-card-old-price">${formatCurrency(product.oldPrice)}</span>`
    : '';

  const defaultSize = (product.sizes && product.sizes.length)
    ? product.sizes[0]
    : 'One size';

  const singleOpt = isSingleOptionProduct(product.type, product.category);
  const sizesSummary = (!singleOpt && product.sizes && product.sizes.length > 1)
    ? `Options: ${product.sizes.join(' · ')}`
    : '';

  const mainImage = (product.images && product.images[0]) || '/assets/images/placeholders/default.svg';

  const ratingHtml = (product.reviews && product.reviews > 0 && product.rating)
    ? `<span class="product-card-rating">
         <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
           <path d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.123 2.123 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.123 2.123 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.122 2.122 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.122 2.122 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.122 2.122 0 0 0 1.597-1.16z"/>
         </svg>
         ${Number(product.rating).toFixed(1)}
       </span>`
    : '';

  card.innerHTML = `
    <button class="product-card-wishlist ${isWished ? 'active' : ''}" aria-label="Add to wishlist" data-wishlist-id="${product.id}">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5"/>
      </svg>
    </button>
    ${saleBadge}
    <a href="/product.html?id=${product.id}" class="product-card-media">
      <img src="${mainImage}" alt="${product.name}" class="product-card-img" loading="lazy" width="400" height="400">
    </a>
    <div class="product-card-body">
      <div class="product-card-meta">
        <span class="product-card-category">${product.category}</span>
        ${ratingHtml}
      </div>
      <a href="/product.html?id=${product.id}" class="product-card-title">${product.name}</a>
      ${sizesSummary ? `<p class="product-card-sizes">${sizesSummary}</p>` : ''}
      <div class="product-card-footer">
        <div>
          <div class="product-card-price">${formatCurrency(product.price)}</div>
          ${oldPriceHtml}
        </div>
        <button class="btn btn-ink btn-sm add-cart-btn" data-id="${product.id}">Add to cart</button>
      </div>
    </div>
  `;

  // Attach 3D tilt
  bindCardTilt(card);

  // Wishlist toggle
  const wishBtn = card.querySelector('.product-card-wishlist');
  wishBtn.addEventListener('click', e => {
    e.preventDefault();
    const updated = toggleWishlist(product.id);
    wishBtn.classList.toggle('active', updated.includes(product.id));
  });

  // Add to cart click
  const addBtn = card.querySelector('.add-cart-btn');
  addBtn.addEventListener('click', e => {
    e.preventDefault();
    addToCart(product, defaultSize, 1);
    showToast(`${product.name} added to cart`);
  });

  return card;
}
