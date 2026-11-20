# Amazing Outfits

An e-commerce website for sneakers and football jerseys, with a 3D animated hero, a product catalogue, cart and checkout, and an admin dashboard for managing products and uploading images.

- **Frontend:** plain HTML, CSS and vanilla JavaScript (ES modules). The 3D hero uses Three.js loaded from a CDN.
- **Backend:** Go (standard library `net/http`) with PostgreSQL via `pgx`.
- **Images:** uploaded by the admin, stored on Cloudinary (or locally in development). The database stores only the image **URL**.

---

## Features

- 3D sneaker hero (Three.js + GLTF model) with a static fallback if WebGL is unavailable
- Shop pages for sneakers and jerseys with filters (category, size, max price), search and sorting
- Product detail pages with image gallery and size selection
- Cart (stored in the browser) and checkout (M-Pesa or pay on delivery)
- Admin page: login, add/delete products, upload product images
- REST API with JWT authentication for admin routes
- **Works without a database:** if the API is unreachable, the frontend falls back to `frontend/products.json`

---

## Project structure

```
Amazing-outfits/
├── backend/cmd/server/main.go    # Server entry point, routes, static file serving
├── internal/
│   ├── db/db.go                  # PostgreSQL connection + auto-creates tables
│   ├── handlers/                 # products, orders, auth, upload
│   ├── middleware/auth.go        # JWT auth, admin guard, CORS
│   └── models/models.go          # Data structures
├── migrations/
│   ├── 001_initial_schema.sql    # Tables and indexes
│   └── 002_seed_data.sql         # 12 sample products, images, sizes, admin user
├── frontend/
│   ├── index.html, shop.html, product.html, cart.html,
│   │   checkout.html, about.html, contact.html, admin.html
│   ├── css/                      # variables, base, components, pages
│   ├── js/                       # api.js, cart.js, ui.js, hero3d.js
│   ├── assets/                   # images and models/sneaker.glb
│   └── products.json             # Offline fallback data
├── docs/                         # design-audit.md, hero-3d-notes.md
├── milestone.md                  # Project milestones
├── .env.example                  # Environment variable template
├── go.mod / go.sum
```

---

## Prerequisites

| Tool | Version | Needed for |
|------|---------|-----------|
| Go | **1.26 or newer** (see `go.mod`) | Running the server |
| PostgreSQL | 13+ (16 recommended) | Real data, orders, admin products |
| A modern browser | WebGL-capable | 3D hero |
| Internet connection | – | Three.js is loaded from `cdn.jsdelivr.net` |

PostgreSQL is optional for a first look; see Option A.

---

## Quick start

### Option A: run without a database (fastest)

Good for viewing the frontend. Products come from `frontend/products.json`, and orders are simulated.

```bash
cd Amazing-outfits
cp .env.example .env
go run ./backend/cmd/server
```

Open <http://localhost:8080>.

The server prints a warning that it could not connect to PostgreSQL. That is expected in this mode. You may also see logged errors for `/api/...` requests; the frontend handles them by using the local JSON data.

> Always run the command from the **project root** (the folder containing `frontend/` and `.env`), because the server looks for `.env` and `frontend/` relative to the current directory.

### Option B: run with PostgreSQL (full functionality)

**1. Start PostgreSQL.** The easiest way is Docker, using settings that match `.env.example`:

```bash
docker run --name amazing-pg \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=amazing_outfits \
  -p 5432:5432 -d postgres:16
```

Or use an existing local installation and create the database:

```bash
createdb -U postgres amazing_outfits
```

**2. Configure environment variables.**

```bash
cp .env.example .env
```

Edit `.env` (see [Environment variables](#environment-variables)). At minimum, change `JWT_SECRET`.

**3. Start the server.** On startup it connects to PostgreSQL and creates all tables automatically.

```bash
go run ./backend/cmd/server
```

**4. Load the sample data.** Do this **once**, after the tables exist:

```bash
psql "postgres://postgres:postgres@localhost:5432/amazing_outfits?sslmode=disable" \
  -f migrations/002_seed_data.sql
```

> Do not run the seed file repeatedly. Products and sizes are protected by `ON CONFLICT`, but image rows are not, so running it again duplicates product images.

**5. Open the site:** <http://localhost:8080>

---

## Pages and URLs

| URL | Page |
|-----|------|
| `/` | Home with 3D hero |
| `/shop` | All products |
| `/sneakers`, `/jerseys` | Shop filtered by product type |
| `/product/{id}` | Product details, e.g. `/product/volt-runner` |
| `/cart`, `/checkout` | Cart and checkout |
| `/about`, `/contact` | Info pages |
| `/admin` | Admin dashboard (login required) |

### Admin login

Use the admin account from the seed file:

- Email: `admin@amazingoutfits.co.ke`
- Password: `admin123`

**Change this before going live.** See [Before deploying](#before-deploying-important).

---

## Environment variables

Defined in `.env` (loaded automatically at startup; real environment variables take priority).

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Port the server listens on |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/amazing_outfits?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | built-in dev value | Secret used to sign admin tokens. **Set your own long random value** |
| `CLOUDINARY_CLOUD_NAME` | empty | Cloudinary account name |
| `CLOUDINARY_API_KEY` | empty | Cloudinary API key |
| `CLOUDINARY_API_SECRET` | empty | Cloudinary API secret |

Generate a strong secret with: `openssl rand -hex 32`

---

## Product image uploads

Images are uploaded from the admin page to `POST /api/admin/products/{id}/images`.

- Accepted types: JPG, PNG, WebP. Maximum size: 5 MB.
- **If the three Cloudinary variables are set:** the image is uploaded to Cloudinary and the returned `https://` URL is saved.
- **If they are empty (development):** the image is saved into `frontend/assets/images/` and the path `/assets/images/<file>` is saved.
- In both cases, only the **URL** goes into the `product_images` table. Image files are never stored in PostgreSQL.

To enable Cloudinary: create a free account, copy the cloud name, API key and API secret from the dashboard into `.env`, then restart the server.

---

## API reference

All responses are JSON. Admin routes need the header `Authorization: Bearer <token>`.

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/products` | – | List products. Query: `type`, `category`, `brand`, `size`, `price` (max), `search`, `sort` (`price-asc`, `price-desc`, `rating`) |
| GET | `/api/products/{id}` | – | Single product |
| POST | `/api/admin/products` | Admin | Create product (with sizes and image URLs) |
| PUT | `/api/admin/products/{id}` | Admin | Update product |
| DELETE | `/api/admin/products/{id}` | Admin | Delete product |
| POST | `/api/admin/products/{id}/images` | Admin | Upload image (`multipart/form-data`, field name `image`) |
| POST | `/api/orders` | – | Create an order |
| GET | `/api/orders/{id}` | – | Get an order |
| POST | `/api/auth/register` | – | Register a customer account |
| POST | `/api/auth/login` | – | Log in, returns a JWT (valid 24 hours) |

Example:

```bash
curl "http://localhost:8080/api/products?type=sneaker&sort=price-asc&price=9000"
```

---

## Database tables

`products`, `product_images`, `product_sizes`, `users`, `orders`, `order_items`. The server creates them automatically on startup; `migrations/001_initial_schema.sql` contains the same schema plus indexes, and can be run manually with `psql -f` if you want the indexes.

---

## Building a binary

```bash
go build -o bin/server ./backend/cmd/server
./bin/server
```

Run it from the project root so it can find `frontend/` and `.env`. Note: the `bin/server` file included in the repository is a Linux x86-64 build and will not run on macOS or Windows. Rebuild it on your own machine, or delete it.

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `go: errors parsing go.mod` or "requires go >= 1.26" | Install Go 1.26+ from <https://go.dev/dl/> |
| "Failed to connect to PostgreSQL" | Check PostgreSQL is running and `DATABASE_URL` is correct. The site still works in fallback mode |
| Shop shows products but orders/admin changes are not saved | You are in fallback mode. Connect PostgreSQL |
| Shop is empty with a database connected | The seed file was not run. Run `migrations/002_seed_data.sql` (step 4) |
| Blank hero / no 3D | Check your internet connection (Three.js loads from a CDN) and that WebGL is enabled. A static fallback is shown otherwise |
| Port already in use | Change `PORT` in `.env` |
| Page not found for `/something.html` | Run the server from the project root, not from inside `backend/` |

---

## Before deploying 

This project is a strong development build, but a few things are intentionally loose for demo purposes and **must be fixed before real customers use it**:

1. **Admin bypass tokens.** The backend accepts the literal token `demo-admin-token` on every admin route (`internal/middleware/auth.go`), and the admin page falls back to it if login fails. Anyone who knows it can create or delete products and upload files. Remove both.
2. **Hard-coded admin login.** `internal/handlers/auth.go` accepts `admin@amazingoutfits.co.ke` with `admin123` (or `admin`) without checking the database. Remove this shortcut, and replace the seeded admin password hash with one you generate (the hash in the seed file appears to be a placeholder).
3. **Set a strong `JWT_SECRET`.** Do not rely on the default fallback value in the code or the example value in `.env.example`.
4. **Prices are trusted from the browser.** `POST /api/orders` stores the totals and item prices sent by the client. The server should look up prices from the `products` table and calculate totals itself.
5. **No real payment yet.** M-Pesa is recorded as a payment method only; there is no Daraja (STK Push) integration, so orders are created as `pending`.
6. **Uploads check the file extension only.** Add content-type sniffing (`http.DetectContentType`) before accepting a file.
7. **CORS is open (`*`).** Restrict it to your real domain in production, or remove it since the frontend is served by the same server.
8. **Orders are readable by ID.** `GET /api/orders/{id}` has no authentication, and IDs are short and partly random. Protect it or return less data.
9. **Order IDs can collide.** They use a random 4-digit number per year; switch to a UUID or database sequence.
10. **Use HTTPS** (for example behind Caddy, Nginx or your host's TLS).

### Repository housekeeping

- Add a `.gitignore` that excludes at least: `.env`, `bin/`, `scratch/`, `node_modules/`. **Do not commit `.env`.**
- `scratch/` (about 50 MB) holds audit scripts and screenshots from building the prototype copy. It is not needed to run the site.
- `assets/` at the project root duplicates `frontend/assets/`; the frontend only needs the copy inside `frontend/`.
- Run `go mod tidy` once so dependencies are no longer marked `// indirect`.

---

## Credits

- 3D sneaker model: "Materials Variants Shoe" by Shopify, licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). Keep this attribution (for example in the site footer or About page).
- 3D rendering: [Three.js](https://threejs.org/)
- Design reference: the Lovable prototype this site recreates.

---

## Roadmap

See [`milestone.md`](./milestone.md). Remaining items include: real M-Pesa payments, customer accounts and order history, admin order management, stock decrement on purchase, and deployment.
