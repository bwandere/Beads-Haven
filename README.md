# Beads Haven — Beads & Art

An e-commerce website for authentic handmade African beadwork and art, crafted by an artisan in Kisumu, Kenya. Features a responsive 2D product hero, a 6-collection catalogue with subcategory filters, shopping cart, checkout with M-Pesa and Cash on Delivery, and an admin dashboard for inventory management.

- **Frontend:** Plain HTML, CSS, and vanilla JavaScript (ES modules). Zero heavy client frameworks.
- **Backend:** Go (standard library `net/http`) with PostgreSQL via `pgx`.
- **Catalogue & Collections:** 6 core collections (Bracelets, Necklaces, Rings, Bags, Fashion Accessories, Art & Home Decor) driven by a single source of truth in `frontend/js/collections.js`.
- **Images:** Stored on Cloudinary (or local upload fallback) and served in optimized WebP format.
- **Offline / No-Database Mode:** The public storefront falls back to `frontend/products.json`. Admin operations are disabled and return `503` when the database is unavailable.

---

## Features

- **Responsive 2D Hero:** Floating collage on desktop; edge-to-edge photo background with a readable dark scrim on mobile. Respects `prefers-reduced-motion`.
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
- **Admin Dashboard:** JWT-authenticated product CRUD with an image manager (multiple uploads, set primary, remove), real server errors, and no mock-data fallback.
- **Image Optimization:** 800×800 WebP square canvases, 400×400 thumbnails, and on-brand vector placeholder tiles.

---

## Project Structure

```
Amazing-outfits/
├── backend/cmd/server/main.go    # Go HTTP server, routing, static file server
├── internal/
│   ├── db/db.go                  # PostgreSQL connection pool & migration runner
│   ├── handlers/                 # products, orders, auth, upload handlers
│   ├── httpjson/                 # shared JSON error response helper
│   ├── middleware/auth.go        # JWT authentication, CORS, admin route guards
│   └── models/models.go          # Go domain data structures
├── migrations/
│   ├── embed.go                  # go:embed for all *.sql files
│   ├── 001_initial_schema.sql    # Tables and indexing
│   ├── 002_seed_data.sql         # Seed catalogue (22 beadwork & art pieces) + admin user
│   └── 003_convert_to_beads.sql  # Idempotent conversion, image uniqueness constraint
├── frontend/
│   ├── index.html                # Homepage with responsive collage hero & collection showcases
│   ├── shop.html                 # Dynamic shop with search, pills, price slider, and sort
│   ├── product.html              # Product details, gallery, and option selector
│   ├── cart.html                 # Cart drawer and item management
│   ├── checkout.html             # Customer info, M-Pesa / Cash on Delivery checkout
│   ├── about.html                # Artisan story and brand values
│   ├── contact.html              # Contact details and custom beadwork commission form
│   ├── admin.html                # Admin login, catalogue table, and product/image editor
│   ├── css/                      # variables.css, base.css, components.css, pages.css
│   ├── js/
│   │   ├── collections.js        # Single source of truth for catalogue metadata & options
│   │   ├── api.js                # API client; storefront fallback, strict admin API calls
│   │   ├── cart.js               # Cart and wishlist state management
│   │   └── ui.js                 # Shared header, navigation, toasts, and card rendering
│   ├── assets/uploads/           # Local upload fallback (gitignored except .gitkeep)
│   ├── assets/images/            # Optimized WebP photos, placeholders, and MANIFEST.md
│   └── products.json             # Static offline catalogue mirror (public storefront only)
└── docs/                         # design notes, screenshots, AGENT_HANDOFF.md
```

---

## Quick Start

Copy the example environment file and edit it before starting the server:

```bash
cp .env.example .env
```

`JWT_SECRET` **must** be set outside local development. If it is empty and `DEV=true` is not set, the server fails fast at startup. The fallback secret used in `DEV=true` mode is for local development only.

### Option A: Run without a database

Ideal for exploring the public frontend and reviewing layout and responsiveness. With the database unavailable, all `/api/*` endpoints return `503 {"error":"Database unavailable"}` and the public storefront uses `frontend/products.json`. The admin dashboard shows a visible database-unavailable banner and cannot save changes.

```bash
cd Amazing-outfits
go run ./backend/cmd/server
```

Open <http://localhost:8080>. The server logs a loud startup warning when PostgreSQL is unreachable.

### Option B: Run with PostgreSQL

1. **Start PostgreSQL:**
   ```bash
   docker run --name beads-pg \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=beads_haven \
     -p 5432:5432 -d postgres:16
   ```

2. **Configure `DATABASE_URL` in `.env`:**
   ```env
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/beads_haven?sslmode=disable
   JWT_SECRET=replace-with-a-long-random-value
   DEV=false
   ```

3. **Start the server:**
   ```bash
   go run ./backend/cmd/server
   ```

Migrations are embedded in the binary and applied automatically on startup. The runner creates `schema_migrations`, applies pending `migrations/*.sql` files in filename order, and records each version in the same transaction as its SQL. Applying migrations manually with `psql` is not required; already-applied versions are skipped.

---

## Admin API

All admin routes require `Authorization: Bearer <token>` from `POST /api/auth/login` and an admin user. Errors are JSON objects with an `error` field.

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/api/auth/login` | Exchange email/password for a JWT |
| `GET` | `/api/products` | Public catalogue; also used by the admin table |
| `GET` | `/api/products/{id}` | Product detail with `images`, `image_items`, and `sizes` |
| `POST` | `/api/admin/products` | Create product (name, type, category, positive price required) |
| `PUT` | `/api/admin/products/{id}` | Update all editable fields and sync sizes (preserves existing stock) |
| `DELETE` | `/api/admin/products/{id}` | Delete product, related rows, and stored images |
| `POST` | `/api/admin/products/{id}/images` | Upload one or more images (`images` multipart field; legacy `image` accepted) |
| `DELETE` | `/api/admin/products/{id}/images/{imageId}` | Delete one image and its stored file |
| `PUT` | `/api/admin/products/{id}/images/order` | Renumber order with `{"order":[imageId,...]}`; first is primary |

Product responses keep the storefront-compatible `images: []string` and add `image_items: [{id,url,display_order}]` for the admin image manager.

### Image storage and validation

- Uploads accept JPG, PNG, and WebP up to 5 MB each; both the filename extension and `http.DetectContentType` are checked.
- When all three `CLOUDINARY_*` variables are set, images are uploaded to Cloudinary. Otherwise they are written only to `frontend/assets/uploads/` and served from `/assets/uploads/`.
- That upload directory is gitignored (with a committed `.gitkeep`); uploaded files must not be committed.
- Adding the first real image removes any placeholder rows for that product.
- Deleting a product or image removes local files. If Cloudinary credentials are configured, Cloudinary assets are also destroyed by their public ID on a best-effort basis.

---

## Admin Credentials

- **Email:** `admin@beadshaven.co.ke`
- **Password:** `admin123`
- **URL:** <http://localhost:8080/admin>

The seed user is created by the embedded migrations. Change the password before any non-local deployment.
