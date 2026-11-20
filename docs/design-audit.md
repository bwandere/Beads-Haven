# Amazing Outfits — Design Audit & Prototype Analysis

**Reference URL:** `https://pixel-perfect-capture-7092.lovable.app/`  
**Breakpoints Audited:** 1440px (Desktop), 768px (Tablet), 390px (Mobile)  
**Inspection Method:** Headless Chromium DOM & computed style extraction, network request capture, and visual screenshot rendering.

---

## 1. Brand Identity & Visual Language

- **Store Name:** Amazing Outfits
- **Tagline:** Premium Sneakers & Football Jerseys in Kenya / "Step into your style"
- **Currency & Market:** Kenyan Shilling (`KSh`, format: `KSh 8,500`), localized for Nairobi & Kenya.
- **Brand Personality:** Bold, energetic, modern urban street culture and athletic footwear.

---

## 2. Color Palette (OKLCH & Exact Hex)

Colors were extracted directly from the `:root` styles of the prototype and verified via Chrome rendering:

| Token | OKLCH Value | Exact Hex | Role & Usage |
|---|---|---|---|
| `--background` | `oklch(97.5% .006 85)` | `#f9f6f2` | Primary page background (warm cream/off-white) |
| `--foreground` | `oklch(20% .03 255)` | `#0c1723` | Primary body text (deep charcoal/navy) |
| `--card` | `oklch(100% 0 0)` | `#ffffff` | Card and component surfaces (pure white) |
| `--card-foreground` | `oklch(20% .03 255)` | `#0c1723` | Card headings and primary text |
| `--primary` | `oklch(55% .21 258)` | `#0069e8` | Vibrant electric royal blue (buttons, active pills, links) |
| `--primary-foreground` | `oklch(98% .01 85)` | `#fbf8f1` | Text on primary elements (crisp light cream) |
| `--primary-glow` | `oklch(68% .18 250)` | `#239cff` | Bright cyan-blue glow highlights |
| `--secondary` | `oklch(93% .012 85)` | `#ebe7df` | Subtle pill/badge background, neutral buttons |
| `--secondary-foreground` | `oklch(20% .03 255)` | `#0c1723` | Text on secondary elements |
| `--muted` | `oklch(94% .008 85)` | `#edebe5` | Muted surfaces and dividers |
| `--muted-foreground` | `oklch(50% .025 255)` | `#5a6472` | Subdued metadata, reviews count, sizes text |
| `--accent` | `oklch(74% .17 60)` | `#f68c11` | Blaze energetic orange (CTAs, Sale badges, highlights, price total) |
| `--accent-foreground` | `oklch(20% .03 255)` | `#0c1723` | Dark text contrast on orange buttons/badges |
| `--ink` | `oklch(17% .035 258)` | `#060f1e` | Deep ink navy-black (Hero background, Header, Footer, dark buttons) |
| `--ink-foreground` | `oklch(96% .015 85)` | `#f6f1e7` | Light text on ink background |
| `--ink-muted` | `oklch(70% .03 255)` | `#92a0b1` | Subdued text on dark ink (Hero subcopy, footer copyright) |
| `--surface` | `oklch(92% .01 250)` | `#e0e5eb` | Image container placeholder backgrounds |
| `--border` / `--input` | `oklch(88% .012 85)` | `#dbd7cf` | Input borders, card dividers |
| `--destructive` | `oklch(58% .22 27)` | `#df2225` | Errors, remove item hover |
| `--hero-particle-orange` | `#f39a2b` | `#f39a2b` | 3D floating octahedron color (1/3 of particles) |
| `--hero-particle-blue` | `#4c8dff` | `#4c8dff` | 3D floating octahedron color (2/3 of particles) |

### Gradients
- **`--gradient-hero` (Hero background):**  
  `radial-gradient(ellipse 70% 60% at 65% 45%, oklch(45% .2 258) 0%, oklch(26% .1 258) 45%, oklch(17% .035 258) 100%)`  
  (Translates from luminous royal blue core behind 3D sneaker to deep ink `#060f1e` edges).
- **`--gradient-promo` (Promo banner):**  
  `linear-gradient(120deg, oklch(55% .21 258) 0%, oklch(42% .2 262) 100%)`  
  (Vibrant royal blue to deep blue-indigo gradient).

---

## 3. Typography & Type Scale

Fonts loaded via Google Fonts CDN:
- **Headings & Display:** `Anton, Impact, sans-serif`  
  Characteristics: uppercase, ultra-bold, condensed letterforms, modern athletic poster style.
- **Body & Interface:** `Barlow, system-ui, -apple-system, sans-serif`  
  Weights used: `400` (regular), `500` (medium), `600` (semibold), `700` (bold).

### Type Scale Breakdown
| Level | Font Family | Size (Mobile -> Desktop) | Weight | Line Height | Tracking | Usage |
|---|---|---|---|---|---|---|
| **Watermark** | `Anton` | `22vw` | 400 | `1` | normal | Background hero text "AMAZING" (opacity 10%) |
| **Hero H1** | `Anton` | `3.75rem (60px)` -> `4.5rem (72px)` -> `6rem (96px)` | 400 | `1.05` | `0.01em` | Homepage hero title ("Step into your style") |
| **H1 (Page)** | `Anton` | `3rem (48px)` -> `3.75rem (60px)` | 400 | `1.1` | normal | "Sneakers", "Your Cart", "Checkout" titles |
| **H2 (Section)**| `Anton` | `2.25rem (36px)` -> `3rem (48px)` | 400 | `1.15` | normal | "Featured Sneakers", "New Arrivals", "Categories" |
| **H3 / Card** | `Anton` | `1.25rem (20px)` -> `1.5rem (24px)` | 400 | `1.2` | normal | Product title in grid ("Volt Runner") |
| **Eyebrow** | `Barlow` | `0.75rem (12px)` | 700 | `1.2` | `0.25em` - `0.3em` | "NEW SEASON · 2026", "HAND-PICKED", "LIMITED OFFER" |
| **Body Large** | `Barlow` | `1.125rem (18px)` | 400 | `1.6` | normal | Hero subtext, story paragraphs |
| **Body Base** | `Barlow` | `1rem (16px)` | 400 / 600 | `1.5` | normal | Product descriptions, labels, form inputs |
| **Body Small** | `Barlow` | `0.875rem (14px)` | 500 / 600 | `1.4` | normal | Metadata, cart item subtitle, table texts |
| **Caption/Tag**| `Barlow` | `0.6875rem (11px)` - `0.75rem (12px)` | 700 | `1.2` | `0.1em` | Badges, Sale pill, sizes summary |

---

## 4. Spacing, Radii, Borders, and Shadows

### Container & Layout Widths
- **Max Container Width:** `max-w-7xl` (`80rem` / `1280px`), centered (`margin: 0 auto`).
- **Horizontal Page Padding:**  
  - Mobile (390px): `px-4` (`16px`)  
  - Tablet/Desktop (768px - 1440px): `px-6` (`24px`)
- **Vertical Section Spacing:** `space-y-24` (`96px`) on desktop, `py-12` to `py-20` for sections.

### Radii System
- Default Token: `--radius: 0.75rem` (`12px`)
- Small: `rounded-md` (`6px`) — size buttons, input boxes, icon buttons
- Standard: `rounded-lg` (`8px`) — thumbnails, badges
- Card: `rounded-2xl` (`16px`) — product cards, cart summary box, category tiles
- Promo Banner: `rounded-3xl` (`24px`) — large promotion container
- Full: `rounded-full` (`9999px`) — filter pills, seasonal tags, round icon buttons

### Borders
- Input & Card: `1px solid var(--border)` (`#dbd7cf`)
- Dark Contrast Borders: `1px solid rgba(246, 241, 231, 0.1)` on header/footer
- Button Outlines: `2px solid` (hero outline, size selector pills)

### Elevation & Shadows
- `--shadow-card`: `0 20px 40px -24px rgba(12, 23, 35, 0.35)` — smooth deep drop shadow under white cards.
- `--shadow-pop`: `0 30px 60px -20px rgba(0, 105, 232, 0.5)` — glowing saturated drop shadow under promo block.
- Product Card Wishlist Button: `0 2px 4px rgba(0,0,0,0.06)`

---

## 5. Components Breakdown and States

### A. Header (Sticky)
- **Structure:**
  - Sticky at `top: 0`, `z-index: 50`.
  - Height: `h-16` (`64px`).
  - Background: `rgba(6, 15, 30, 0.9)` (`bg-ink/90`) with `backdrop-filter: blur(12px)`.
  - Border: `border-b border-ink-foreground/10`.
- **Elements:**
  - Logo: `40px x 40px` rounded image (`logo.png`) + text `Amazing` + orange `Outfits`.
  - Desktop Nav (`hidden md:flex`): `Home`, `Sneakers`, `Jerseys`, `About`, `Contact` (uppercase tracking-wider, active state has `color: var(--accent)`).
  - Search Icon Button: navigates to `/sneakers` (catalog with live search input).
  - Cart Icon Button: shopping bag icon, live count badge pill `rounded-full bg-accent text-[11px] font-bold` displayed if count > 0.
  - Mobile Menu Hamburger Button (`md:hidden`): toggles mobile drawer menu with big `Anton` 3xl navigation links.

### B. Product Card Component
- **Structure:**
  - `overflow-hidden rounded-2xl bg-card shadow-card`
  - 3D Interactive Tilt on desktop pointer move:
    `perspective(900px) rotateY(...) rotateX(...)` with smooth reset on mouseleave.
  - Wishlist Button: top-right circular overlay, heart icon toggling filled orange.
  - "Sale" Badge: top-left `rounded-full bg-accent text-[11px] font-bold` (if `oldPrice` present).
  - Image: square `aspect-square bg-surface` with `group-hover:scale-105` zoom.
  - Card Body:
    - Row: Category name (`uppercase tracking-wider`) + Star icon with rating number (`4.8`).
    - Product Title: `font-display text-xl`, hover state turns royal blue (`text-primary`).
    - Sizes preview: `Sizes: 38 · 39 · 40 · 41 · 42 · 43 · 44 · 45`.
    - Bottom row: Price in `KSh` (plus old price strike-through if on sale) + "Add to cart" button (`variant: ink, size: sm`).
- **States:** Hover transforms card, scales image, changes title color, buttons react to pointer.

### C. Buttons System
- Base: `inline-flex items-center justify-center font-semibold rounded-md transition-all`
- Variants:
  - `accent`: `bg-accent text-accent-foreground uppercase tracking-wider shadow-lg hover:-translate-y-0.5 hover:brightness-105`
  - `ink`: `bg-ink text-ink-foreground uppercase tracking-wider hover:bg-primary`
  - `heroOutline`: `border-2 border-ink-foreground/70 text-ink-foreground uppercase tracking-wider hover:bg-ink-foreground hover:text-ink`
  - `default`: `bg-primary text-primary-foreground shadow hover:bg-primary/90`
- Sizes: `sm` (h-8), `default` (h-10), `lg` (h-12), `xl` (h-14).

### D. Shop Page & Filtering System
- Header with eyebrow and title.
- Search input with search magnifying glass icon: live instant filter.
- Category filter pills: `All`, `Running`, `Lifestyle`, `Basketball`, `Training` (Sneakers) or `All`, `National Teams`, `Clubs` (Jerseys). Active pill is `bg-primary text-primary-foreground`.
- Max Price Slider: range slider from 1,000 to max price, live formatted readout.
- Sort Select Dropdown: `Featured`, `Price: Low to high`, `Price: High to low`, `Top rated`.
- Live count of matching products (`N products`).
- Responsive Grid: 2 columns on mobile/tablet, 4 columns on desktop (`grid-cols-2 lg:grid-cols-4`).

### E. Product Detail Page
- Two-column layout on desktop (`grid lg:grid-cols-2 gap-10`):
  - Left: Main square product image + 4-thumbnail gallery selector below it with active border highlight.
  - Right:
    - Category tag (uppercase tracking-wide royal blue).
    - H1 product title (`Anton` 5xl - 6xl).
    - Star rating + review count.
    - Large price + strike-through old price.
    - Description paragraph.
    - Size selection buttons: interactive grid of size buttons (`38` - `45` or `S` - `XXL`), active button turns dark ink `bg-ink text-ink-foreground`.
    - Quantity selector: `[-] count [+]` stepper.
    - Action buttons: "Add to cart" (`variant: ink, size: xl`) and "Buy now" (`variant: accent, size: xl`).
    - Delivery note: Truck icon + `Same-day delivery in Nairobi · Countrywide in 2–3 days`.
- Bottom Section: "You may also like" / "Related products" grid of 4 items.

### F. Cart Page
- Empty State: Centered card with "Your cart is empty." and "Shop sneakers" accent CTA button.
- Populated State:
  - Left column: List of cart item cards.
    - Thumbnail image.
    - Name, size, unit price.
    - Quantity stepper (`-`, qty, `+`).
    - Remove button (trash icon, turns red on hover).
    - Line item total price.
  - Right column: `Order summary` dark card (`bg-ink text-ink-foreground`).
    - Itemized breakdown (`qty × name (size)`).
    - Subtotal, Delivery (KSh 300 flat), and Total (accent orange).
    - "Checkout" full-width CTA button (`variant: accent, size: xl`).

### G. Checkout Page
- Form sections in white cards (`rounded-2xl bg-card p-6 shadow-card`):
  1. Customer Details: Full name, Phone, Email, Delivery location.
  2. Payment Method Selector:
     - 2 selectable cards: "M-Pesa (Get an STK push on your phone)" and "Pay on delivery (Nairobi only)". Active card highlighted with blue border and tint.
     - When M-Pesa is selected: input field for "M-Pesa number".
  3. "Place order" button (`variant: accent, size: xl`).
- Right column: Dark `Order summary` sidebar with real-time totals.

### H. About & Contact Pages
- **About:** Brand story, Nairobi roots, 3 value cards (Authentic Quality, Fast Delivery, Community First).
- **Contact:** Contact info cards (Phone, WhatsApp, Email, Physical Location in Nairobi CBD) + interactive message form.

### I. Footer
- Dark background (`bg-ink text-ink-foreground`).
- 3 Columns:
  1. Brand logo + title + mission statement.
  2. "Shop" links (Sneakers, Jerseys, Cart).
  3. "Company" links (About, Contact).
- Bottom copyright bar: `© 2026 Amazing Outfits · 3D shoe model “Materials Variants Shoe” by Shopify, CC BY 4.0`.

---

## 6. Complete Asset Inventory

1. **3D Model:**
   - File: `sneaker.glb` (7.8 MB binary glTF format).
   - Local path: `/assets/models/sneaker.glb`.
   - Attribution: "Materials Variants Shoe" by Shopify (CC BY 4.0).
2. **Brand Imagery:**
   - Logo: `/assets/images/logo.png`
   - Favicon: `/assets/images/favicon.png`
3. **Product Imagery:**
   - `sneaker-1.jpg` (Volt Runner / Volt Runner Pro)
   - `sneaker-2.jpg` (Court Classic Low / Street Court)
   - `sneaker-3.jpg` (Night High / Blackout Mid)
   - `sneaker-4.jpg` (Retro Wave / Dune Trainer)
   - `jersey-1.jpg` (Royal Blues Home / Jerseys category)
   - `jersey-2.jpg` (Red Stripes Classic)
   - `jersey-3.jpg` (Harambee Stars Home 26 / Lions Away Kit)

---

## 7. Unknowns & Non-Guessed Items

- Backend DB schema and admin authorization was not present in the frontend prototype (lovable demo purely client-side React with in-memory state).
- The prototype checkout was demo-only (`toast.success("Order placed (demo only)")`).
- Per specification, we implement full Go backend with standard `net/http` router, PostgreSQL persistence, and Cloudinary admin image uploads.
