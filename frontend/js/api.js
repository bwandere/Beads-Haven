/**
 * Amazing Outfits - API Layer
 * Connects to Go backend with seamless fallback to /products.json
 */

const API_BASE = '/api';

export async function fetchProducts(filters = {}) {
  try {
    const params = new URLSearchParams();
    if (filters.brand) params.append('brand', filters.brand);
    if (filters.category) params.append('category', filters.category);
    if (filters.type) params.append('type', filters.type);
    if (filters.size) params.append('size', filters.size);
    if (filters.price) params.append('price', filters.price);
    if (filters.sort) params.append('sort', filters.sort);
    if (filters.search) params.append('search', filters.search);

    const res = await fetch(`${API_BASE}/products?${params.toString()}`);
    if (res.ok) {
      return await res.json();
    }
    throw new Error(`API error: ${res.status}`);
  } catch (err) {
    // Fallback to local mock data
    console.info('Using local products.json mock fallback:', err.message);
    const mockRes = await fetch('/products.json');
    let products = await mockRes.json();

    // Apply client-side filters on mock data
    if (filters.type) {
      products = products.filter(p => p.type === filters.type);
    }
    if (filters.category && filters.category !== 'All') {
      products = products.filter(p => p.category.toLowerCase() === filters.category.toLowerCase());
    }
    if (filters.brand) {
      products = products.filter(p => p.brand.toLowerCase() === filters.brand.toLowerCase());
    }
    if (filters.size) {
      products = products.filter(p => p.sizes.includes(filters.size));
    }
    if (filters.price) {
      const maxP = parseFloat(filters.price);
      products = products.filter(p => p.price <= maxP);
    }
    if (filters.search) {
      const q = filters.search.toLowerCase();
      products = products.filter(p => p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q));
    }

    if (filters.sort === 'price-asc') {
      products.sort((a, b) => a.price - b.price);
    } else if (filters.sort === 'price-desc') {
      products.sort((a, b) => b.price - a.price);
    } else if (filters.sort === 'rating') {
      products.sort((a, b) => b.rating - a.rating);
    }

    return products;
  }
}

export async function fetchProductById(id) {
  try {
    const res = await fetch(`${API_BASE}/products/${id}`);
    if (res.ok) {
      return await res.json();
    }
    throw new Error(`API error: ${res.status}`);
  } catch (err) {
    const mockRes = await fetch('/products.json');
    const products = await mockRes.json();
    return products.find(p => p.id === id) || null;
  }
}

export async function createOrder(orderPayload) {
  try {
    const res = await fetch(`${API_BASE}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(orderPayload)
    });
    if (res.ok) {
      return await res.json();
    }
    const errData = await res.json().catch(() => ({}));
    throw new Error(errData.error || `Failed to create order: ${res.status}`);
  } catch (err) {
    // For demo/offline mode, return simulated order response
    console.info('Simulated order fallback:', err.message);
    return {
      success: true,
      orderId: 'ORD-' + Math.floor(100000 + Math.random() * 900000),
      message: 'Order placed successfully (demo mode)'
    };
  }
}

export async function fetchOrderById(id) {
  const res = await fetch(`${API_BASE}/orders/${id}`);
  if (!res.ok) throw new Error('Order not found');
  return await res.json();
}

export async function loginUser(email, password) {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Login failed');
  }
  return await res.json();
}

export async function createAdminProduct(productData, token) {
  const headers = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}/admin/products`, {
    method: 'POST',
    headers,
    body: JSON.stringify(productData)
  });
  if (!res.ok) throw new Error('Failed to create product');
  return await res.json();
}

export async function uploadProductImage(productId, file, token) {
  const formData = new FormData();
  formData.append('image', file);

  const headers = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}/admin/products/${productId}/images`, {
    method: 'POST',
    headers,
    body: formData
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Image upload failed');
  }
  return await res.json();
}

export async function deleteProduct(productId, token) {
  const headers = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}/admin/products/${productId}`, {
    method: 'DELETE',
    headers
  });
  if (!res.ok) throw new Error('Failed to delete product');
  return await res.json();
}
