# Beads Haven — Beads & Art

An e-commerce website for authentic handmade African beadwork and art, crafted by an artisan in Kisumu, Kenya. Features a responsive 2D product hero, a 6-collection catalogue with subcategory filters, shopping cart, checkout with M-Pesa and Cash on Delivery, and an admin dashboard for inventory management.

- **Frontend:** Plain HTML, CSS, and vanilla JavaScript (ES modules). Zero heavy client frameworks.
- **Backend:** Go (standard library `net/http`) with PostgreSQL via `pgx`.
- **Catalogue & Collections:** 6 core collections (Bracelets, Necklaces, Rings, Bags, Fashion Accessories, Art & Home Decor) driven by a single source of truth in `frontend/js/collections.js`.
- **Images:** Stored on Cloudinary (or local upload fallback) and served in optimized WebP format.
- **Offline / No-Database Mode:** Operates fully offline using `frontend/products.json` if PostgreSQL is not running.

---

## Features

- **Performant 2D Hero:** CSS-animated overlapping collage of artisan bags, bracelets, and rings with `prefers-reduced-motion` support.
- **Multi-Collection Catalogue:** 6 primary collections with subcategory filtering:
  - **Bracelets:** Beaded Bracelets, Cowrie Shell Bracelets (with quick link to Waist Beads)
  - **Necklaces:** Beaded Necklaces, Cowrie Shell Necklaces, Pendant Necklaces
  - **Rings:** Beaded Rings (Stack rings, Initial monogram bands)
  - **Bags:** Beaded Bags (Structured handbags, clutches, totes)
  - **Fashion Accessories:** Beaded Ties, Beaded Belts, Hair Charms, Waist Beads, Beaded Ankaras
  - **Art & Home Decor:** Beaded Wall Portraits, Beaded Wall Art, Decorative Beadwork
- **Collection-Aware Options:** Contextual option selectors (e.g. wrist size for bracelets, necklace lengths, ring sizes, waist measurements, dimensions for art, and hidden pickers for "One size" pieces).
- **Navigation:** Accessible desktop dropdown mega-menu and collapsible mobile drawer accordion.
- **Cart & Wishlist:** Persisted in browser localStorage; supports M-Pesa phone number validation and Pay on Delivery.
- **Admin Dashboard:** Secure JWT authentication, inventory table, and collection-dependent product creation form.
- **Image Optimization:** 800×800 WebP square canvases, 400×400 thumbnails, and on-brand vector placeholder tiles.

---

## Project Structure

```
Amazing-outfits/
├── backend/cmd/server/main.go    # Go HTTP server, routing, static file server
├── internal/
│   ├── db/db.go                  # PostgreSQL connection pool & auto-schema setup
│   ├── handlers/                 # products, orders, auth, upload handlers
│   ├── middleware/auth.go        # JWT authentication, CORS, admin route guards
│   └── models/models.go          # Go domain data structures
├── migrations/
│   ├── 001_initial_schema.sql    # Tables and indexing
│   ├── 002_seed_data.sql         # Seed catalogue (21+ beadwork & art pieces) + admin user
│   └── 003_convert_to_beads.sql  # Idempotent migration cleaning legacy data
├── frontend/
│   ├── index.html                # Homepage with 2D collage hero & collection showcases
│   ├── shop.html                 # Dynamic shop with search, pills, price slider, and sort
│   ├── product.html              # Product details, gallery, and option selector
│   ├── cart.html                 # Cart drawer and item management
│   ├── checkout.html             # Customer info, M-Pesa / Cash on Delivery checkout
│   ├── about.html                # Artisan story and brand values
│   ├── contact.html              # Contact details and custom beadwork commission form
│   ├── admin.html                # Admin login and catalogue management
│   ├── css/                      # variables.css, base.css, components.css, pages.css
│   ├── js/
│   │   ├── collections.js        # Single source of truth for catalogue metadata & options
│   │   ├── api.js                # API client with fallback to products.json
│   │   ├── cart.js               # Cart and wishlist state management
│   │   └── ui.js                 # Shared header, navigation, toasts, and card rendering
│   ├── assets/images/            # Optimized WebP photos, placeholders, and MANIFEST.md
│   └── products.json             # Static offline catalogue mirror
└── docs/                         # design-audit.md, screenshots
```

---

## Quick Start

### Option A: Run without a database (Fastest)

Ideal for exploring the frontend and reviewing layout and responsiveness. The frontend automatically falls back to `frontend/products.json`.

```bash
cd Amazing-outfits
go run ./backend/cmd/server
```

Open <http://localhost:8080> in your browser.

The server logs a warning that PostgreSQL is not connected and operates in mock mode.

### Option B: Run with PostgreSQL

1. **Start PostgreSQL:**
   ```bash
   docker run --name beads-pg \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=beads_haven \
     -p 5432:5432 -d postgres:16
   ```

2. **Apply migrations:**
   ```bash
   psql -U postgres -d beads_haven -f migrations/001_initial_schema.sql
   psql -U postgres -d beads_haven -f migrations/002_seed_data.sql
   psql -U postgres -d beads_haven -f migrations/003_convert_to_beads.sql
   ```

3. **Start the server:**
   ```bash
   go run ./backend/cmd/server
   ```

---

## Admin Credentials

- **Email:** `admin@beadshaven.co.ke`
- **Password:** `admin123`
- **URL:** <http://localhost:8080/admin>
