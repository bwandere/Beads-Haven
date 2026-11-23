# Product Imagery Manifest — Beads Haven (Beads & Art)

This manifest documents all imagery in `frontend/assets/images/`, separating authentic product photography from SVG/WebP placeholders, and flagging items containing third-party branding or copyrighted characters.

---

## 1. Real Product Photographs (Authentic Beadwork)

All original raw photos are archived in `frontend/assets/images/originals/` and converted to optimized WebP (`800×800` square canvas) + `400×400` thumbnail (`-thumb.webp`).

| File | Collection / Category | Product Usage | Dimensions | Alt Text |
| :--- | :--- | :--- | :--- | :--- |
| `bag1.webp` | Bags · Beaded Bags | Handcrafted Royal Pearl Beaded Handbag | 800×800 | Structured white and cream beaded handbag with pearl bead handles |
| `bag2.webp` | Bags · Beaded Bags | Noir Elegance Beaded Evening Bag | 800×800 | Black glass seed bead evening bag with beaded loop handle |
| `bag3.webp` | Bags · Beaded Bags | Sunset Glow Patterned Bead Bag | 800×800 | Colorful handcrafted beaded handbag with geometric African motifs |
| `bag4.webp` | Bags · Beaded Bags | Amber Crystal Beaded Handbag | 800×800 | Translucent amber crystal bead handbag with secure top flap |
| `bracelet1.webp` | Bracelets · Beaded Bracelets | Maasai Heritage Beaded Cuff | 800×800 | Traditional Kenyan black beaded cuff bracelet with red and green geometric pattern |
| `ring1.webp` | Rings · Beaded Rings | Vibrant Stack Beaded Rings (Set of 3) | 800×800 | Set of three stacked handmade glass bead rings on fingers |
| `ring2.webp` | Rings · Beaded Rings | Custom Monogram Initial Beaded Ring | 800×800 | Black woven seed bead ring featuring a custom white initial letter |

---

## 2. Flagged Third-Party Brand & Character Images

> [!WARNING]
> The following photos feature third-party trademarks, telecommunications logos (e.g., Airtel), or animated characters (e.g., SpongeBob). 
> **Per store policy, these images are strictly excluded from the hero collage, promo banners, collection cards, and featured carousels.** They are seeded as generic custom products. The store owner should review them to decide whether to retain or replace them with original artisanal pieces.

| File | Subcategory | Seeded Product Name | Flagged Content | Note |
| :--- | :--- | :--- | :--- | :--- |
| `bracelet3.webp` | Beaded Bracelets | Custom Pixel-Pattern Cuff | Airtel brand logo & brand colors | Excluded from hero/featured. Owner review advised. |
| `bracelet4.webp` | Beaded Bracelets | Pop-Art Beaded Wristband | Cartoon character (SpongeBob motif) | Tall image (720×1280) square-padded. Excluded from hero. |
| `bracelet5.webp` | Beaded Bracelets | Custom Branded Accent Cuff | Telecommunications brand pattern | Seeded as generic custom band. Excluded from hero. |
| `3dpotrait.webp` | Art & Home Decor | Custom Beaded Portrait Medallion | Pop culture character portrait | Framed bead portrait. Seeded with generic title. |

---

## 3. On-Brand Placeholder Tiles (`frontend/assets/images/placeholders/`)

Categories and subcategories without physical photography have been provided with elegant, on-brand SVG tiles and matching WebP equivalents (800×800 full and 400×400 thumbnail). Products using these have `"placeholder": true` in `products.json`.

> [!NOTE]
> **Action item for store owner**: Supply real artisanal photography for the items listed below to replace these placeholders.

| Placeholder File | Target Subcategory | Status | Owner Photo Needed |
| :--- | :--- | :--- | :--- |
| `cowrie-shell-bracelets.svg` | Cowrie Shell Bracelets | Vector & WebP Placeholder | Real photo of cowrie shell bracelet on wrist |
| `beaded-necklaces.svg` | Beaded Necklaces | Vector & WebP Placeholder | Real photo of Maasai multi-strand collar |
| `cowrie-shell-necklaces.svg` | Cowrie Shell Necklaces | Vector & WebP Placeholder | Real photo of cowrie choker / necklace |
| `pendant-necklaces.svg` | Pendant Necklaces | Vector & WebP Placeholder | Real photo of beaded pendant necklace |
| `beaded-ties.svg` | Beaded Ties | Vector & WebP Placeholder | Real photo of handwoven beaded tie |
| `beaded-belts.svg` | Beaded Belts | Vector & WebP Placeholder | Real photo of wide beadwork belt |
| `hair-charms.svg` | Hair Charms | Vector & WebP Placeholder | Real photo of beaded hair cuffs / braid charms |
| `waist-beads.svg` | Waist Beads | Vector & WebP Placeholder | Real photo of beaded waist strands |
| `beaded-ankaras.svg` | Beaded Ankaras | Vector & WebP Placeholder | Real photo of Ankara fabric piece with bead trim |
| `beaded-wall-art.svg` | Beaded Wall Art | Vector & WebP Placeholder | Real photo of framed African beaded wall art |
| `decorative-beadwork.svg` | Decorative Beadwork | Vector & WebP Placeholder | Real photo of decorative sculptural beadwork |
| `necklaces.svg` | Necklaces Collection Card | Collection Card Asset | Primary collection hero photo |
| `fashion-accessories.svg` | Fashion Accessories Card | Collection Card Asset | Primary collection hero photo |
| `art-decor.svg` | Art & Home Decor Card | Collection Card Asset | Primary collection hero photo |
| `default.svg` | Neutral System Fallback | UI Fallback Asset | Permanent system placeholder |

---

## 4. Brand & Icon Assets

| Asset | Dimensions | Purpose |
| :--- | :--- | :--- |
| `logo.svg` | 512×512 (Vector) | Primary brand logo (light backgrounds) |
| `logo-dark.svg` | 512×512 (Vector) | Brand logo variant (dark background / footer) |
| `logo.png` | 512×512 (PNG) | Transparent high-res brand emblem |
| `logo-dark.png` | 512×512 (PNG) | Transparent dark-mode brand emblem |
| `favicon.png` | 64×64 (PNG) | Browser tab icon (concentric bead rosette) |
| `app-icon.png` | 512×512 (PNG) | PWA / mobile app icon on warm cream squircle |
