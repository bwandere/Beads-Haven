/**
 * Amazing Outfits - 3D Hero Canvas Implementation
 * Three.js + GLTFLoader
 */

import * as THREE from 'three';
import { GLTFLoader } from '/vendor/three/examples/jsm/loaders/GLTFLoader.js';

export function initHero3D() {
  const container = document.getElementById('hero-stage');
  const canvas = document.getElementById('hero-canvas');
  const loaderEl = document.getElementById('hero-loader');

  if (!container || !canvas) return;

  function enableFallback() {
    clearTimeout(loadTimeout);
    container.classList.add('is-fallback');
    if (loaderEl) {
      loaderEl.classList.add('hidden');
      loaderEl.style.display = 'none';
    }
  }

  // Check WebGL availability
  function isWebGLAvailable() {
    try {
      const c = document.createElement('canvas');
      return !!(window.WebGLRenderingContext && (c.getContext('webgl2') || c.getContext('webgl') || c.getContext('experimental-webgl')));
    } catch (e) {
      return false;
    }
  }

  if (!isWebGLAvailable()) {
    enableFallback();
    return;
  }

  // Fail-safe timeout: 12 seconds max before showing static fallback
  let modelLoaded = false;
  const loadTimeout = setTimeout(() => {
    if (!modelLoaded) {
      console.warn('Model loading timed out after 12s, enabling static fallback');
      enableFallback();
    }
  }, 12000);

  // Handle WebGL context lost
  canvas.addEventListener('webglcontextlost', (e) => {
    e.preventDefault();
    console.warn('WebGL context lost, falling back to static image');
    stopLoop();
    enableFallback();
  }, false);

  const motionMediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
  let isReducedMotion = motionMediaQuery.matches;
  motionMediaQuery.addEventListener?.('change', (e) => {
    isReducedMotion = e.matches;
  });

  const isMobile = window.innerWidth < 768;

  // Scene & Camera
  const scene = new THREE.Scene();

  const width = container.clientWidth || 600;
  const height = container.clientHeight || 500;
  const camera = new THREE.PerspectiveCamera(38, width / height, 0.1, 100);
  camera.position.set(0, 0.6, 6.5);
  camera.lookAt(0, 0, 0);

  // Renderer
  let renderer;
  try {
    renderer = new THREE.WebGLRenderer({
      canvas: canvas,
      antialias: true,
      alpha: true,
      powerPreference: 'high-performance'
    });
  } catch (err) {
    console.warn('WebGL init failed, falling back to static image', err);
    enableFallback();
    return;
  }

  const maxDpr = isMobile ? 1.5 : 2;
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, maxDpr));
  renderer.setSize(width, height, false);
  renderer.shadowMap.enabled = !isMobile;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;

  // Lighting
  const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
  scene.add(ambientLight);

  const dirLight = new THREE.DirectionalLight(0xffffff, 2.2);
  dirLight.position.set(4, 6, 4);
  if (!isMobile) {
    dirLight.castShadow = true;
    dirLight.shadow.mapSize.width = 1024;
    dirLight.shadow.mapSize.height = 1024;
    dirLight.shadow.camera.near = 0.5;
    dirLight.shadow.camera.far = 15;
  }
  scene.add(dirLight);

  // Blaze Orange point light
  const orangeLight = new THREE.PointLight(0xf39a2b, 30, 20);
  orangeLight.position.set(-4, 1, 2);
  scene.add(orangeLight);

  // Royal Blue point light
  const blueLight = new THREE.PointLight(0x3d7bff, 25, 20);
  blueLight.position.set(4, -1, 2);
  scene.add(blueLight);

  // Ground Contact Shadow plane
  const shadowGeo = new THREE.PlaneGeometry(8, 8);
  const canvasShadow = document.createElement('canvas');
  canvasShadow.width = 128;
  canvasShadow.height = 128;
  const ctx = canvasShadow.getContext('2d');
  const gradient = ctx.createRadialGradient(64, 64, 0, 64, 64, 64);
  gradient.addColorStop(0, 'rgba(6, 15, 30, 0.6)');
  gradient.addColorStop(0.5, 'rgba(6, 15, 30, 0.25)');
  gradient.addColorStop(1, 'rgba(6, 15, 30, 0)');
  ctx.fillStyle = gradient;
  ctx.fillRect(0, 0, 128, 128);

  const shadowTexture = new THREE.CanvasTexture(canvasShadow);
  const shadowMat = new THREE.MeshBasicMaterial({
    map: shadowTexture,
    transparent: true,
    opacity: 0.55,
    depthWrite: false
  });
  const shadowMesh = new THREE.Mesh(shadowGeo, shadowMat);
  shadowMesh.rotation.x = -Math.PI / 2;
  shadowMesh.position.set(0, -1.7, 0);
  scene.add(shadowMesh);

  // Floating Octahedron Crystals
  const particleGroup = new THREE.Group();
  scene.add(particleGroup);

  const particleCount = isMobile ? 8 : 18;
  const octaGeometry = new THREE.OctahedronGeometry(1, 0);

  const orangeMat = new THREE.MeshStandardMaterial({
    color: 0xf39a2b,
    roughness: 0.3,
    metalness: 0.4
  });
  const blueMat = new THREE.MeshStandardMaterial({
    color: 0x4c8dff,
    roughness: 0.3,
    metalness: 0.4
  });

  const particles = [];
  for (let i = 0; i < particleCount; i++) {
    const isOrange = i % 3 === 0;
    const mesh = new THREE.Mesh(octaGeometry, isOrange ? orangeMat : blueMat);
    const s = 0.08 + Math.random() * 0.18;
    mesh.scale.set(s, s, s);

    const basePos = new THREE.Vector3(
      (Math.random() - 0.5) * 12,
      (Math.random() - 0.5) * 6,
      -2 - Math.random() * 4
    );
    mesh.position.copy(basePos);

    particleGroup.add(mesh);
    particles.push({
      mesh,
      speed: 1 + (i % 3) * 0.5,
      rotSpeed: 0.8 + Math.random(),
      basePos
    });
  }

  // Model Setup
  const sneakerGroup = new THREE.Group();
  scene.add(sneakerGroup);

  const loader = new GLTFLoader();
  const modelUrl = '/assets/models/sneaker.glb';

  loader.load(
    modelUrl,
    gltf => {
      modelLoaded = true;
      clearTimeout(loadTimeout);

      const model = gltf.scene;

      // Auto-center & compute bounding box
      const box = new THREE.Box3().setFromObject(model);
      const size = new THREE.Vector3();
      box.getSize(size);
      const center = new THREE.Vector3();
      box.getCenter(center);

      model.position.x += -center.x;
      model.position.y += -center.y;
      model.position.z += -center.z;

      // Scale model
      const maxDim = Math.max(size.x, size.y, size.z);
      const targetScale = 4.2 / maxDim;
      model.scale.set(targetScale, targetScale, targetScale);

      model.traverse(child => {
        if (child.isMesh) {
          child.castShadow = !isMobile;
          child.receiveShadow = !isMobile;
        }
      });

      sneakerGroup.add(model);

      if (loaderEl) {
        loaderEl.classList.add('hidden');
        loaderEl.style.display = 'none';
      }
    },
    undefined,
    err => {
      console.warn('Failed to load sneaker.glb model, enabling static fallback', err);
      enableFallback();
    }
  );

  // Mouse & Touch Tracking
  const pointer = { x: 0, y: 0 };
  const targetPointer = { x: 0, y: 0 };

  function onPointerMove(e) {
    const rect = container.getBoundingClientRect();
    if (rect.width === 0 || rect.height === 0) return;
    const clientX = e.clientX ?? (e.touches && e.touches[0]?.clientX) ?? 0;
    const clientY = e.clientY ?? (e.touches && e.touches[0]?.clientY) ?? 0;

    targetPointer.x = ((clientX - rect.left) / rect.width) * 2 - 1;
    targetPointer.y = -(((clientY - rect.top) / rect.height) * 2 - 1);
  }

  window.addEventListener('mousemove', onPointerMove, { passive: true });
  window.addEventListener('touchmove', onPointerMove, { passive: true });

  // Responsive Resizing
  function resize() {
    const w = container.clientWidth;
    const h = container.clientHeight;
    if (w === 0 || h === 0) return;

    camera.aspect = w / h;
    camera.updateProjectionMatrix();

    renderer.setSize(w, h, false);
  }

  window.addEventListener('resize', resize);

  // Animation Loop & Performance Controls
  let isVisible = true;
  let isTabActive = !document.hidden;
  let animationFrameId = null;
  const clock = new THREE.Clock();

  // IntersectionObserver to pause rendering when off-screen
  const observer = new IntersectionObserver(
    entries => {
      const entry = entries[0];
      isVisible = entry.isIntersecting;
      if (isVisible && isTabActive) {
        startLoop();
      } else {
        stopLoop();
      }
    },
    { threshold: 0.05 }
  );
  observer.observe(container);

  // Page Visibility API
  document.addEventListener('visibilitychange', () => {
    isTabActive = !document.hidden;
    if (isVisible && isTabActive) {
      startLoop();
    } else {
      stopLoop();
    }
  });

  function startLoop() {
    if (!animationFrameId) {
      clock.start();
      animate();
    }
  }

  function stopLoop() {
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId);
      animationFrameId = null;
      clock.stop();
    }
  }

  function animate() {
    animationFrameId = requestAnimationFrame(animate);

    const delta = clock.getDelta();
    const elapsedTime = clock.getElapsedTime();
    const dt = Math.min(delta, 0.05);

    // Smooth pointer interpolation
    pointer.x += (targetPointer.x - pointer.x) * 0.1;
    pointer.y += (targetPointer.y - pointer.y) * 0.1;

    if (!isReducedMotion) {
      // Sneaker Interactive Physics & Float
      const targetY = elapsedTime * 0.25 + pointer.x * 0.5;
      const targetX = -pointer.y * 0.25 + 0.15;
      const factor = 1 - Math.exp(-4 * dt);

      sneakerGroup.rotation.y += (targetY - sneakerGroup.rotation.y) * factor;
      sneakerGroup.rotation.x += (targetX - sneakerGroup.rotation.x) * factor;

      // Vertical floating wave
      sneakerGroup.position.y = Math.sin(elapsedTime * 1.6) * 0.12;

      // Crystals floating & group parallax
      particleGroup.rotation.z = elapsedTime * 0.02;
      particleGroup.position.x = pointer.x * -0.3;
      particleGroup.position.y = pointer.y * -0.2;

      particles.forEach((p, idx) => {
        p.mesh.rotation.x += 0.01 * p.rotSpeed;
        p.mesh.rotation.y += 0.015 * p.rotSpeed;
        p.mesh.position.y = p.basePos.y + Math.sin(elapsedTime * p.speed + idx) * 0.2;
      });
    } else {
      // Reduced motion: static clean view
      sneakerGroup.rotation.y = 0.4;
      sneakerGroup.rotation.x = 0.1;
      sneakerGroup.position.y = 0;
    }

    renderer.render(scene, camera);
  }

  startLoop();
}
