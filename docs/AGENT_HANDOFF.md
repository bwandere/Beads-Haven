# Agent Handoff — Beads Haven

Repo branch: `feat/beads-and-art`. Current goal: complete verification/fixes for admin CRUD with real image upload/edit and the mobile hero background.

## Git rules

- Commit every logical change; never push, never open a PR.
- Never touch `scratch/`.
- Never commit `.env`, uploaded images under `frontend/assets/uploads/`, or screenshots.
- Update this handoff file in every commit.

## Task checklist

### Part A — Database + admin CRUD

- [x] Embed `migrations/*.sql` and run pending files transactionally with `schema_migrations` (`migrations/embed.go`, `internal/db/db.go`).
- [x] Fail clearly with `503 {"error":"Database unavailable"}` when DB is unavailable; log loud startup warning (`backend/cmd/server/main.go`).
- [x] Require `JWT_SECRET` outside `DEV=true`; remove demo token/hardcoded login shortcuts (`main.go`, `internal/middleware/auth.go`, `internal/handlers/auth.go`).
- [x] Transactional product Create with validation, unique slug, duplicate-ID 409, full response (`internal/handlers/products.go`).
- [x] Product Update updates all editable fields, syncs sizes preserving stock, returns 404/full product (`internal/handlers/products.go`).
- [x] Product Delete returns 404 if missing and cleans local uploads; Cloudinary images are destroyed by derived public_id when credentials exist (`internal/handlers/products.go`, `internal/handlers/upload.go`).
- [x] Image upload: multiple files, size/extension/content sniffing, local fallback only in `frontend/assets/uploads/`, placeholder removal, JSON response (`internal/handlers/upload.go`).
- [x] Image delete and transactional reorder endpoints (`internal/handlers/upload.go`).
- [x] Product JSON keeps `images` and always includes `image_items` (empty array when none); public list batches image/size queries (`internal/models/models.go`, `internal/handlers/products.go`).
- [x] All API JSON errors use the shared `internal/httpjson.WriteError` helper with `Content-Type: application/json`.
- [x] Admin API client uses API-only product fetch, real server errors, auth-failure handling, multi-file upload (`frontend/js/api.js`).
- [x] Admin UI has Edit mode, image manager, escaping, disabled Save state, mobile modal layout, no fake success paths (`frontend/admin.html`).
- [x] Fresh database migrations, curl CRUD/image lifecycle, 401/409/400/404/503 checks, and browser CRUD/image manager flow verified (see log).

### Part B — Mobile hero

- [x] Mobile edge-to-edge collage background, dark scrim, subtle image scale animation, reduced-motion support (`frontend/index.html`, `frontend/css/pages.css`).
- [x] Browser geometry/screenshot verification at 360×800, 390×844, 430×932, 768×1024, and 1280×800 (see log).

### Documentation

- [ ] Update `README.md` for admin endpoints, migrations behaviour, upload storage, and `JWT_SECRET` requirement.

## Verification log

- `go build ./... && go vet ./...`: clean.
- Fresh DB `beads_haven_verify`: migrations 001–003 applied in order; `schema_migrations` records all three; `/api/products` returns 22 seeded products with images/sizes.
- curl lifecycle: login 200; no token 401; create 201; duplicate id 409; upload 2 images 201; reorder 200; update 200 (price/rating/reviews/description/sizes changed, existing stock preserved); delete image 200; delete product 200; GET/update/delete missing product 404.
- Error paths: demo token 401 JSON; invalid category/price 400 JSON; fake PNG rejected 400 by content sniffing; placeholder row removed when first real image uploaded; DB-down server returns 503 JSON for public and admin APIs with loud startup warning; missing `JWT_SECRET` outside `DEV=true` exits 1.
- Browser admin flow: login, create with image, edit prefill, add/reorder/remove/replace image, delete, reload; no `products.json` requests from admin.
- Browser modal: `max-height: 100dvh`, scrollable body and stacked fields at 360/390/430/768px; desktop keeps 90vh. Simulated image-upload failure keeps the modal in edit mode with "Update Piece" enabled and the real server error in a toast.
- Hero geometry: at 360/390/430/768 CSS px the three tiles cover the hero with zero gaps, `hero-grid` static, text z-index 3, scrim z-index 1, `object-fit: cover`, subtle scale animation; 1280px keeps floating desktop layout with no scrim/scale. Screenshots written outside the repo at `/home/bwandere/hero-*.png`.

## Notes for the next agent

- Continue from the first unchecked item. Keep commits atomic.
- Verify behaviour by actually running the app/curl/browser; do not trust commit messages.
- The uncommitted changes at handoff time avoid duplicate upload `images`/`image` fields and preserve multi-file selection in `admin.html`.
