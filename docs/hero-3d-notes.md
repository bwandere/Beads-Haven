# 3D Hero Technique & Engineering Analysis

**Document:** `docs/hero-3d-notes.md`  
**Reference Source Component:** `Hero3DCanvas-CynQfFCb.js`  
**Inspected Asset:** `sneaker.glb` (Binary glTF 2.0, 7.8 MB)

---

## 1. Technique Identified

The prototype uses **Three.js (WebGL2 / WebGL)** rendered on an HTML `<canvas>` element via `@react-three/fiber` and `@react-three/drei`.
- **Canvas Element:** Dynamically mounted in the hero right column (`relative order-1 h-[46vh] min-h-[320px] lg:order-2 lg:h-[72vh]`).
- **Model Format:** Binary glTF (`.glb`).
- **Model URL:** `https://pixel-perfect-capture-7092.lovable.app/__l5e/assets-v1/c2be357d-89cc-4aec-b82d-17d752b8ef0d/sneaker.glb`.
- **Local Download:** Saved in `/assets/models/sneaker.glb`.
- **Attribution:** "Materials Variants Shoe" by Shopify (CC BY 4.0).

---

## 2. Camera Configuration

- **Projection:** Perspective Camera (`THREE.PerspectiveCamera`).
- **Field of View (FOV):** `38°`
- **Initial Camera Position:** `x: 0, y: 0.6, z: 6.5`
- **Look At:** Origin `(0, 0, 0)`
- **Device Pixel Ratio (DPR):**
  - Desktop: `dpr: [1, 2]` (capped at 2 for retina screens)
  - Mobile: `dpr: [1, 1.5]` (capped at 1.5 for performance)
- **Context Config:** `antialias: true`, `alpha: true` (transparent canvas showing the underlying CSS gradient).

---

## 3. Lighting & Environment Setup

The 3D scene uses a multi-light studio setup designed to highlight shoe contours with athletic blue and warm orange rim highlights:

1. **Ambient Light:**
   - Color: Neutral white (`#ffffff`)
   - Intensity: `0.6`
2. **Key Directional Light (Top-Right):**
   - Position: `[4, 6, 4]`
   - Intensity: `2.2`
   - Cast Shadow: `true` (disabled on mobile for performance)
3. **Accent Point Light (Left Rim - Blaze Orange):**
   - Position: `[-4, 1, 2]`
   - Intensity: `30`
   - Color: `#f39a2b` (warm orange rim light)
4. **Accent Point Light (Right Rim - Royal Blue):**
   - Position: `[4, -1, 2]`
   - Intensity: `25`
   - Color: `#3d7bff` (vibrant blue fill light)
5. **Contact Shadow / Soft Ground Shadow:**
   - Position: `[0, -1.7, 0]` (below sole)
   - Opacity: `0.55`
   - Blur: `2.6`
   - Far: `4.0`
   - Scale: `10`
6. **Environment Emissive Panels:**
   - Overhead panel: `position: [0, 5, 0], scale: [10, 10, 1], intensity: 3, rotation.x: Math.PI / 2`
   - Left orange reflector: `position: [-5, 1, -1], color: '#f39a2b', intensity: 2, scale: [20, 1, 1]`
   - Right blue reflector: `position: [5, 1, -1], color: '#4c8dff', intensity: 2, scale: [20, 1, 1]`

---

## 4. Model Scaling & Centering

- The shoe model is centered at `[0, 0, 0]`.
- Normalized bounding-box scale: `scale: 4.2`.
- Floating base movement:
  - Base vertical float oscillation: `speed: 1.6`, `rotationIntensity: 0.25`, `floatIntensity: 0.9`.
  - Sinusoidal floating wave added to vertical position.

---

## 5. Animation Timing & Interactive Physics

The prototype employs smooth exponential damping interpolation for interactive mouse-tracking:

### Sneaker Orientation:
```javascript
// Target calculation (continuous auto-rotation + mouse pointer offset)
const targetY = clock.elapsedTime * 0.25 + pointer.x * 0.5;
const targetX = -pointer.y * 0.25 + 0.15;

// Exponential damping per frame (smoothing rate = 4)
const dt = Math.min(delta, 0.05);
const factor = 1 - Math.exp(-4 * dt);
group.rotation.y += (targetY - group.rotation.y) * factor;
group.rotation.x += (targetX - group.rotation.x) * factor;
```
- **Auto-Rotation:** Continuous slow turn at `0.25 rad/s` on Y axis.
- **Pointer Follow:** Mouse X offsets Y rotation by `±0.5 rad`, Mouse Y tilts X rotation by `±0.25 rad` around base tilt `0.15 rad`.
- **Smoothing:** Smooth spring-like exponential ease (`1 - Math.exp(-4 * dt)`).

---

## 6. Background Floating Geometry (Particles)

Surrounding the sneaker are floating 3D geometric crystals:
- **Geometry:** `THREE.OctahedronGeometry(1, 0)`
- **Count:** `18` on desktop, `8` on mobile.
- **Color Distribution:**
  - Index % 3 == 0: Warm Orange (`#f39a2b`)
  - Others: Royal Blue (`#4c8dff`)
- **Material:** `THREE.MeshStandardMaterial({ roughness: 0.3, metalness: 0.4 })`
- **Individual Particle Scales:** Random scale `0.08` to `0.26`.
- **Spread Area:** `x: [-6, 6]`, `y: [-3, 3]`, `z: [-2, -6]`.
- **Animation:**
  - Group parallax follows mouse: `group.position.x = pointer.x * -0.3`, `group.position.y = pointer.y * -0.2`.
  - Group rotation: `group.rotation.z = clock.elapsedTime * 0.02`.
  - Individual independent floating and rotation per particle.

---

## 7. Performance & Quality Plan for Production Rebuild

To fulfill all requirements in plain vanilla JavaScript (without React or bundlers):
1. **Three.js via CDN:** Load `three.module.js` and `GLTFLoader.js` via standard ES module import / importmap.
2. **Responsive Canvas:**
   - Capped `devicePixelRatio` at `2` (desktop) and `1.5` (mobile).
   - Window resize handler updating camera aspect ratio and renderer viewport.
3. **Intersection Observer & Page Visibility:**
   - Pause `requestAnimationFrame` render loop when hero is off-screen.
   - Pause rendering when document tab is hidden (`document.hidden`).
4. **Loading State:**
   - Show loading spinner / skeleton overlay until glb is loaded, then smoothly fade in.
5. **Reduced Motion:**
   - Query `window.matchMedia('(prefers-reduced-motion: reduce)')`.
   - If reduced motion is requested, disable auto-rotation and mouse-tracking, rendering static clean model.
6. **Fallback Mechanism:**
   - If WebGL is unsupported or file load fails, gracefully fall back to high-res static hero sneaker image (`/assets/images/sneaker-1.jpg`) without breaking page layout.
7. **Accessibility & SEO:**
   - All text, headings, badges, and CTAs remain real HTML in the DOM overlay above the WebGL canvas.

---

## 8. Root Cause Analysis & Production Fixes Applied

### Root Causes Identified
1. **Unresolvable Bare Module Specifier (`three`):**  
   `hero3d.js` imported `GLTFLoader.js` from `cdn.jsdelivr.net`. Inside `GLTFLoader.js`, it had `import ... from 'three'`. Because the browser had no `<script type="importmap">`, modern browsers rejected the module resolution with:  
   `TypeError: Failed to resolve module specifier "three". Relative references must start with either "/", "./", or "../".`
2. **Coupled Script Halting:**  
   `index.html` imported `hero3d.js` statically inside the same `<script type="module">` as `initHeader()` and `loadHomePageProducts()`. When module resolution failed for `hero3d.js`, the entire script crashed before running `initHeader()` or loading any products.
3. **No Fail-Safe / Permanent Spinner:**  
   There was no timeout guard, so if model loading or WebGL initialization stalled, the loading spinner spun indefinitely without triggering the fallback.
4. **Sub-optimal Model & Fallback Asset:**  
   The initial 3D model `sneaker.glb` was uncompressed (7.83 MB) with 2048x2048 uncompressed JPEGs. The fallback image was a square non-transparent JPEG (`sneaker-1.jpg`) without proper dimensions or `fetchpriority`.

### Solutions Applied
1. **Self-Hosted Three.js (`vendor/three/`):**  
   - Extracted `three@0.160.0`: `build/three.module.js`, `examples/jsm/loaders/GLTFLoader.js`, and `examples/jsm/utils/BufferGeometryUtils.js`.  
   - Kept relative directory structure so internal imports in GLTFLoader function without alterations.
2. **Import Map in `index.html`:**  
   - Added standard `<script type="importmap">` mapping `"three"` to `/vendor/three/build/three.module.js` and addons to `/vendor/three/examples/jsm/`.
3. **Decoupled Dynamic Hero Loader:**  
   - Header initialization and product loading run independently and immediately.  
   - Hero is loaded via dynamic `import('/js/hero3d.js').then(...).catch(...)`. Any hero failure triggers fallback without affecting the rest of the site.
4. **Fail-Safe Mechanism:**  
   - Added a 12-second timeout guard to switch to fallback if model loading exceeds threshold.  
   - Added `webglcontextlost` event handler to switch to fallback and stop the render loop.
5. **High-Res Transparent Fallback:**  
   - Created `/assets/images/hero-fallback.webp` (3.3 KB) and `.png` (39 KB) with transparent cutout, rendered with `fetchpriority="high"` and explicit dimensions (`width="800" height="600"`).
6. **Model Optimization:**  
   - Preserved original 7.83 MB model as `sneaker.original.glb` outside `frontend/`.  
   - Resized textures to max 1024px, converted to WebP (`quality: 85`), deduplicated, and pruned unreferenced properties via `@gltf-transform/cli`.  
   - Reduced model size from 7.83 MB to 1.23 MB (84% reduction) with zero visible fidelity loss and native browser WebP decoder support (no Draco or KTX2 dependencies required).
7. **Responsive & Motion Verification:**  
   - Verified rendering at 1440px, 768px, and 390px with zero console errors and zero failed requests. Captured screenshots in `docs/screenshots/`.  
   - Confirmed `prefers-reduced-motion` renders static pose and off-screen observer pauses the animation loop.
