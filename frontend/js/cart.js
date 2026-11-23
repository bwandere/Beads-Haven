/**
 * Beads Haven - Cart Management (localStorage)
 */

const CART_STORAGE_KEY = 'beads_haven_cart';
const WISHLIST_STORAGE_KEY = 'beads_haven_wishlist';

export function getCart() {
  try {
    const raw = localStorage.getItem(CART_STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch (e) {
    console.error('Failed to read cart from localStorage', e);
    return [];
  }
}

export function saveCart(lines) {
  try {
    localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(lines));
    window.dispatchEvent(new CustomEvent('cart-updated', { detail: lines }));
  } catch (e) {
    console.error('Failed to save cart to localStorage', e);
  }
}

export function addToCart(product, size, qty = 1) {
  const cart = getCart();
  const existingIndex = cart.findIndex(item => item.product.id === product.id && item.size === size);

  if (existingIndex >= 0) {
    cart[existingIndex].qty += qty;
  } else {
    cart.push({ product, size, qty });
  }

  saveCart(cart);
  return cart;
}

export function updateQuantity(productId, size, qty) {
  const cart = getCart();
  const item = cart.find(i => i.product.id === productId && i.size === size);
  if (item) {
    item.qty = Math.max(1, qty);
    saveCart(cart);
  }
  return cart;
}

export function removeFromCart(productId, size) {
  let cart = getCart();
  cart = cart.filter(i => !(i.product.id === productId && i.size === size));
  saveCart(cart);
  return cart;
}

export function clearCart() {
  saveCart([]);
}

export function getCartCount() {
  const cart = getCart();
  return cart.reduce((total, item) => total + item.qty, 0);
}

export function getCartSubtotal() {
  const cart = getCart();
  return cart.reduce((sum, item) => sum + item.qty * item.product.price, 0);
}

export function getDeliveryFee() {
  const cart = getCart();
  return cart.length > 0 ? 300 : 0;
}

export function getCartTotal() {
  return getCartSubtotal() + getDeliveryFee();
}

/* Wishlist Helpers */
export function getWishlist() {
  try {
    const raw = localStorage.getItem(WISHLIST_STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch (e) {
    return [];
  }
}

export function toggleWishlist(productId) {
  let list = getWishlist();
  if (list.includes(productId)) {
    list = list.filter(id => id !== productId);
  } else {
    list.push(productId);
  }
  localStorage.setItem(WISHLIST_STORAGE_KEY, JSON.stringify(list));
  window.dispatchEvent(new CustomEvent('wishlist-updated', { detail: list }));
  return list;
}
