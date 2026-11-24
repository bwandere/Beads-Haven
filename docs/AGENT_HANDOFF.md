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
- [x] Product Delete returns 404 if missing and cleans local uploads (`internal/handlers/products.go`); Cloudinary destroy still needs verification.
- [x] Image upload: multiple files, size/extension/content sniffing, local fallback only in `frontend/assets/uploads/`, placeholder removal, JSON response (`internal/handlers/upload.go`).
- [x] Image delete and transactional reorder endpoints (`internal/handlers/upload.go`).
- [x] Product JSON keeps `images` and adds `image_items`; public list batches image/size queries (`internal/models/models.go`, `internal/handlers/products.go`).
- [x] JSON errors use one helper with `Content-Type: application/json` (`internal/handlers/common.go`).
- [x] Admin API client uses API-only product fetch, real server errors, auth-failure handling, multi-file upload (`frontend/js/api.js`).
- [x] Admin UI has Edit mode, image manager, escaping, disabled Save state, mobile modal layout, no fake success paths (`frontend/admin.html`).
- [ ] End-to-end verification with a fresh database, curl lifecycle, and browser checks still outstanding.

### Part B — Mobile hero

- [x] Mobile edge-to-edge collage background, dark scrim, subtle image scale animation, reduced-motion support (`frontend/index.html`, `frontend/css/pages.css`).
- [ ] Browser verification at 360×800, 390×844, 430×932, 768×1024, and 1280×800 still outstanding.

### Documentation

- [ ] Update `README.md` for admin endpoints, migrations behaviour, upload storage, and `JWT_SECRET` requirement.

## Verification log

- `go build ./... && go vet ./...`: clean after the pending upload-field fixes.
- Fresh-DB/curl/browser/screenshot verification: not yet done in this session.

## Notes for the next agent

- Continue from the first unchecked item. Keep commits atomic.
- Verify behaviour by actually running the app/curl/browser; do not trust commit messages.
- The uncommitted changes at handoff time avoid duplicate upload `images`/`image` fields and preserve multi-file selection in `admin.html`.
