# Beads Haven — Project Milestones

A Kenyan e-commerce store dedicated to authentic handmade beadwork and African art.

**Brand:** "Beads Haven" · Tagline: "Beads & Art"  
**Palette:** Terracotta Maasai-red `#b3261e`, warm gold `#d9a02b`, deep forest green `#1f6b4f`, warm cream background `#faf5ec`, charcoal `#1a1412`, and ink `#14100e`.  
**Typography:** Playfair Display (headings & brand mark) + DM Sans (body, UI, metadata).

**Architecture:** Plain HTML/CSS/vanilla JS (ES modules) frontend + Go backend (`net/http`) with PostgreSQL (`pgx`) and offline JSON fallback (`frontend/products.json`).

---

## Catalogue & Collection Structure

1. **Bracelets** (`bracelets`)
   - Subcategories: Beaded Bracelets, Cowrie Shell Bracelets
   - Cross-link: Waist Beads
2. **Necklaces** (`necklaces`)
   - Subcategories: Beaded Necklaces, Cowrie Shell Necklaces, Pendant Necklaces
3. **Rings** (`rings`)
   - Subcategories: Beaded Rings (Stack rings, Initial monogram bands)
4. **Bags** (`bags`)
   - Subcategories: Beaded Bags (Structured handbags, evening clutches)
5. **Fashion Accessories** (`fashion-accessories`)
   - Subcategories: Beaded Ties, Beaded Belts, Hair Charms, Waist Beads, Beaded Ankaras
6. **Art & Home Decor** (`art-decor`)
   - Subcategories: Beaded Wall Portraits, Beaded Wall Art, Decorative Beadwork

---

## Milestone 1 — Brand Identity & Token Architecture
- [x] Design original beadwork rosette SVG logo, transparent PNGs, favicon, and app icon
- [x] Establish warm African handmade palette (OKLCH + Hex fallbacks) in `variables.css`
- [x] Integrate Playfair Display and DM Sans typography
- [x] Subtle repeating bead-string motif for section dividers and footer

## Milestone 2 — Hero & Visual Presence
- [x] Decommission 3D sneaker hero and Three.js vendor dependencies
- [x] Build performant, accessible 2D hero with 3 overlapping product cards and CSS float animations
- [x] Ensure full compliance with `prefers-reduced-motion`

## Milestone 3 — Navigation & Multi-Collection Architecture
- [x] Create single source of truth in `frontend/js/collections.js`
- [x] Desktop dropdown mega-menu listing all 6 collections and subcategories
- [x] Mobile drawer accordion for collection browsing
- [x] Consistent, informative footer across all 8 site pages
- [x] Dynamic collection-aware shop filtering and category pills

## Milestone 4 — Image Optimization & Manifest
- [x] Process authentic artisan photography into square 800×800 WebP and 400×400 thumbnails
- [x] Generate on-brand vector and WebP placeholder tiles for unphotographed subcategories
- [x] Archive raw photo assets in `frontend/assets/images/originals/`
- [x] Document image manifest, alt texts, and flagged third-party images in `MANIFEST.md`

## Milestone 5 — Data Model & Go Backend
- [x] Update `type` column semantics to represent the 6 collection slugs
- [x] Validate collection slugs in Go create/update handlers (400 on unknown)
- [x] Transition sizing to collection-specific presets (wrist size, length, ring sizes, dimensions, one size)
- [x] Rewrite seed data (`migrations/002_seed_data.sql`) and `frontend/products.json`
- [x] Idempotent migration `003_convert_to_beads.sql` removing legacy seeds and adding uniqueness constraints
- [x] Valid bcrypt hash for admin credentials (`admin@beadshaven.co.ke` / `admin123`)

## Milestone 6 — Admin Dashboard & Custom Orders
- [x] Collection and subcategory dynamic dropdowns in product creation form
- [x] Contextual option inputs based on single-option vs multi-size collections
- [x] Custom orders inquiry pipeline on contact page and promo banners
