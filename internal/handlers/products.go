package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"beads-haven/internal/httpjson"
	"beads-haven/internal/models"
)

type ProductHandler struct {
	DB *sql.DB
}

var ValidCategories = map[string][]string{
	"bracelets":           {"Beaded Bracelets", "Cowrie Shell Bracelets"},
	"necklaces":           {"Beaded Necklaces", "Cowrie Shell Necklaces", "Pendant Necklaces"},
	"rings":               {"Beaded Rings"},
	"bags":                {"Beaded Bags"},
	"fashion-accessories": {"Beaded Ties", "Beaded Belts", "Hair Charms", "Waist Beads", "Beaded Ankaras"},
	"art-decor":           {"Beaded Wall Portraits", "Beaded Wall Art", "Decorative Beadwork"},
}

func isValidType(t string) bool {
	_, ok := ValidCategories[t]
	return ok
}

func validateCategoryForType(prodType, category string) bool {
	cats, ok := ValidCategories[prodType]
	if !ok {
		return false
	}
	for _, c := range cats {
		if strings.EqualFold(c, category) {
			return true
		}
	}
	return false
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	lastHyphen := true
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastHyphen = false
		} else {
			if !lastHyphen {
				b.WriteRune('-')
				lastHyphen = true
			}
		}
	}
	res := strings.Trim(b.String(), "-")
	if res == "" {
		return "beadwork-piece"
	}
	return res
}

func (h *ProductHandler) fetchFullProduct(id string) (*models.Product, error) {
	var p models.Product
	var tagsJSON []byte
	err := h.DB.QueryRow(`
		SELECT id, name, brand, type, category, price, old_price, 
		       rating, reviews, description, tags, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.Brand, &p.Type, &p.Category, &p.Price, &p.OldPrice,
		&p.Rating, &p.Reviews, &p.Description, &tagsJSON, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Tags = []string{}
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &p.Tags)
	}

	// Fetch Images
	p.Images = []string{}
	p.ImageItems = []models.ProductImage{}
	imgRows, err := h.DB.Query("SELECT id, product_id, url, display_order FROM product_images WHERE product_id = $1 ORDER BY display_order ASC, id ASC", p.ID)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img models.ProductImage
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder); err == nil {
				p.ImageItems = append(p.ImageItems, img)
				p.Images = append(p.Images, img.URL)
			}
		}
	}

	// Fetch Sizes
	p.Sizes = []string{}
	sizeRows, err := h.DB.Query("SELECT size FROM product_sizes WHERE product_id = $1 ORDER BY id ASC", p.ID)
	if err == nil {
		defer sizeRows.Close()
		for sizeRows.Next() {
			var s string
			if err := sizeRows.Scan(&s); err == nil {
				p.Sizes = append(p.Sizes, s)
			}
		}
	}

	return &p, nil
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	q := r.URL.Query()
	brand := q.Get("brand")
	size := q.Get("size")
	priceStr := q.Get("price")
	sort := q.Get("sort")
	search := q.Get("search")
	category := q.Get("category")
	prodType := q.Get("type")

	query := `
		SELECT p.id, p.name, p.brand, p.type, p.category, p.price, p.old_price, 
		       p.rating, p.reviews, p.description, p.tags, p.created_at, p.updated_at
		FROM products p
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if brand != "" {
		query += fmt.Sprintf(" AND LOWER(p.brand) = LOWER($%d)", argIdx)
		args = append(args, brand)
		argIdx++
	}

	if prodType != "" {
		query += fmt.Sprintf(" AND p.type = $%d", argIdx)
		args = append(args, prodType)
		argIdx++
	}

	if category != "" && category != "All" {
		query += fmt.Sprintf(" AND LOWER(p.category) = LOWER($%d)", argIdx)
		args = append(args, category)
		argIdx++
	}

	if priceStr != "" {
		if maxP, err := strconv.ParseFloat(priceStr, 64); err == nil && maxP > 0 {
			query += fmt.Sprintf(" AND p.price <= $%d", argIdx)
			args = append(args, maxP)
			argIdx++
		}
	}

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(p.name) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if size != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM product_sizes ps WHERE ps.product_id = p.id AND ps.size = $%d)", argIdx)
		args = append(args, size)
		argIdx++
	}

	switch sort {
	case "price-asc":
		query += " ORDER BY p.price ASC"
	case "price-desc":
		query += " ORDER BY p.price DESC"
	case "rating":
		query += " ORDER BY p.rating DESC"
	default:
		query += " ORDER BY p.created_at ASC"
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Query error: %s", err.Error()))
		return
	}
	defer rows.Close()

	var products []models.Product
	var productIDs []string
	prodMap := make(map[string]*models.Product)

	for rows.Next() {
		var p models.Product
		var tagsJSON []byte
		err := rows.Scan(
			&p.ID, &p.Name, &p.Brand, &p.Type, &p.Category, &p.Price, &p.OldPrice,
			&p.Rating, &p.Reviews, &p.Description, &tagsJSON, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "Error scanning product")
			return
		}
		p.Tags = []string{}
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &p.Tags)
		}
		p.Images = []string{}
		p.ImageItems = []models.ProductImage{}
		p.Sizes = []string{}

		products = append(products, p)
		productIDs = append(productIDs, p.ID)
	}

	for i := range products {
		prodMap[products[i].ID] = &products[i]
	}

	// Batch fetch images and sizes to avoid N+1 queries
	if len(productIDs) > 0 {
		imgRows, err := h.DB.Query(`
			SELECT id, product_id, url, display_order 
			FROM product_images 
			WHERE product_id = ANY($1) 
			ORDER BY display_order ASC, id ASC
		`, productIDs)
		if err == nil {
			defer imgRows.Close()
			for imgRows.Next() {
				var img models.ProductImage
				if err := imgRows.Scan(&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder); err == nil {
					if pr, ok := prodMap[img.ProductID]; ok {
						pr.ImageItems = append(pr.ImageItems, img)
						pr.Images = append(pr.Images, img.URL)
					}
				}
			}
		}

		sizeRows, err := h.DB.Query(`
			SELECT product_id, size 
			FROM product_sizes 
			WHERE product_id = ANY($1) 
			ORDER BY id ASC
		`, productIDs)
		if err == nil {
			defer sizeRows.Close()
			for sizeRows.Next() {
				var pid, sz string
				if err := sizeRows.Scan(&pid, &sz); err == nil {
					if pr, ok := prodMap[pid]; ok {
						pr.Sizes = append(pr.Sizes, sz)
					}
				}
			}
		}
	}

	if products == nil {
		products = []models.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Product ID required")
		return
	}

	p, err := h.fetchFullProduct(id)
	if err == sql.ErrNoRows {
		httpjson.WriteError(w, http.StatusNotFound, "Product not found")
		return
	} else if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Product name is required")
		return
	}

	if p.Price <= 0 {
		httpjson.WriteError(w, http.StatusBadRequest, "Price must be greater than 0")
		return
	}

	if p.Brand == "" {
		p.Brand = "Beads Haven"
	}
	if p.Type == "" {
		p.Type = "bracelets"
	}
	if !isValidType(p.Type) {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid collection type. Must be one of: bracelets, necklaces, rings, bags, fashion-accessories, art-decor")
		return
	}

	p.Category = strings.TrimSpace(p.Category)
	if p.Category == "" || !validateCategoryForType(p.Type, p.Category) {
		httpjson.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Category '%s' is not valid for collection '%s'", p.Category, p.Type))
		return
	}

	// Generate or validate ID
	p.ID = strings.TrimSpace(p.ID)
	if p.ID == "" {
		baseSlug := slugify(p.Name)
		slug := baseSlug
		counter := 1
		for {
			var exists bool
			err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", slug).Scan(&exists)
			if err != nil {
				httpjson.WriteError(w, http.StatusInternalServerError, "Database error checking slug")
				return
			}
			if !exists {
				p.ID = slug
				break
			}
			counter++
			slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		}
	} else {
		p.ID = slugify(p.ID)
		var exists bool
		err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", p.ID).Scan(&exists)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "Database error checking ID")
			return
		}
		if exists {
			httpjson.WriteError(w, http.StatusConflict, fmt.Sprintf("Product with ID '%s' already exists", p.ID))
			return
		}
	}

	if p.Tags == nil {
		p.Tags = []string{"new"}
	}
	tagsJSON, _ := json.Marshal(p.Tags)

	// Transaction: Product + Sizes + Images
	tx, err := h.DB.Begin()
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO products (id, name, brand, type, category, price, old_price, rating, reviews, description, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, p.ID, p.Name, p.Brand, p.Type, p.Category, p.Price, p.OldPrice, p.Rating, p.Reviews, p.Description, tagsJSON)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create product: %s", err.Error()))
		return
	}

	// Insert Sizes
	if len(p.Sizes) == 0 {
		p.Sizes = []string{"One size"}
	}
	for _, s := range p.Sizes {
		s = strings.TrimSpace(s)
		if s != "" {
			_, err = tx.Exec("INSERT INTO product_sizes (product_id, size, stock) VALUES ($1, $2, 10) ON CONFLICT (product_id, size) DO NOTHING", p.ID, s)
			if err != nil {
				httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create size: %s", err.Error()))
				return
			}
		}
	}

	// Insert Images (if any provided)
	for idx, u := range p.Images {
		u = strings.TrimSpace(u)
		if u != "" {
			_, err = tx.Exec("INSERT INTO product_images (product_id, url, display_order) VALUES ($1, $2, $3)", p.ID, u, idx)
			if err != nil {
				httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert image: %s", err.Error()))
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	created, err := h.fetchFullProduct(p.ID)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Product created but failed to load response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Product ID required")
		return
	}

	// Check if product exists
	var existing models.Product
	var exTags []byte
	err := h.DB.QueryRow(`
		SELECT id, name, brand, type, category, price, old_price, rating, reviews, description, tags
		FROM products WHERE id = $1
	`, id).Scan(
		&existing.ID, &existing.Name, &existing.Brand, &existing.Type, &existing.Category,
		&existing.Price, &existing.OldPrice, &existing.Rating, &existing.Reviews, &existing.Description, &exTags,
	)
	if err == sql.ErrNoRows {
		httpjson.WriteError(w, http.StatusNotFound, "Product not found")
		return
	} else if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Product name is required")
		return
	}

	if p.Price <= 0 {
		httpjson.WriteError(w, http.StatusBadRequest, "Price must be greater than 0")
		return
	}

	if p.Brand == "" {
		p.Brand = existing.Brand
	}
	if p.Type == "" {
		p.Type = existing.Type
	}
	if !isValidType(p.Type) {
		httpjson.WriteError(w, http.StatusBadRequest, "Invalid collection type. Must be one of: bracelets, necklaces, rings, bags, fashion-accessories, art-decor")
		return
	}

	p.Category = strings.TrimSpace(p.Category)
	if p.Category == "" {
		p.Category = existing.Category
	}
	if !validateCategoryForType(p.Type, p.Category) {
		httpjson.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Category '%s' is not valid for collection '%s'", p.Category, p.Type))
		return
	}

	if p.Tags == nil {
		p.Tags = []string{}
		if len(exTags) > 0 {
			json.Unmarshal(exTags, &p.Tags)
		}
	}
	tagsJSON, _ := json.Marshal(p.Tags)

	// Transaction to update product and sync sizes preserving stock
	tx, err := h.DB.Begin()
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to begin transaction")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE products
		SET name = $1, brand = $2, type = $3, category = $4, price = $5, old_price = $6,
		    rating = $7, reviews = $8, description = $9, tags = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`, p.Name, p.Brand, p.Type, p.Category, p.Price, p.OldPrice, p.Rating, p.Reviews, p.Description, tagsJSON, id)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update product: %s", err.Error()))
		return
	}

	// Size sync: preserve existing stock values
	if p.Sizes != nil {
		existingSizes := make(map[string]int)
		rows, err := tx.Query("SELECT size, stock FROM product_sizes WHERE product_id = $1", id)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, "Failed to read existing sizes")
			return
		}
		for rows.Next() {
			var sz string
			var st int
			if err := rows.Scan(&sz, &st); err == nil {
				existingSizes[sz] = st
			}
		}
		rows.Close()

		newSizesMap := make(map[string]bool)
		for _, sz := range p.Sizes {
			sz = strings.TrimSpace(sz)
			if sz == "" {
				continue
			}
			newSizesMap[sz] = true
			if _, exists := existingSizes[sz]; !exists {
				// New size - default stock 10
				_, err = tx.Exec("INSERT INTO product_sizes (product_id, size, stock) VALUES ($1, $2, 10)", id, sz)
				if err != nil {
					httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert size: %s", err.Error()))
					return
				}
			}
		}

		// Delete removed sizes
		for oldSz := range existingSizes {
			if !newSizesMap[oldSz] {
				_, err = tx.Exec("DELETE FROM product_sizes WHERE product_id = $1 AND size = $2", id, oldSz)
				if err != nil {
					httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete removed size: %s", err.Error()))
					return
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	updated, err := h.fetchFullProduct(id)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Product updated but failed to reload")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		httpjson.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		httpjson.WriteError(w, http.StatusBadRequest, "Product ID required")
		return
	}

	// Check if product exists and gather image URLs for cleanup
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		httpjson.WriteError(w, http.StatusNotFound, "Product not found")
		return
	}

	imgRows, _ := h.DB.Query("SELECT url FROM product_images WHERE product_id = $1", id)
	var imageURLs []string
	if imgRows != nil {
		for imgRows.Next() {
			var u string
			if err := imgRows.Scan(&u); err == nil {
				imageURLs = append(imageURLs, u)
			}
		}
		imgRows.Close()
	}

	res, err := h.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete product: %s", err.Error()))
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		httpjson.WriteError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Best-effort delete locally stored upload files
	for _, u := range imageURLs {
		if strings.HasPrefix(u, "/assets/uploads/") {
			fn := filepath.Base(u)
			dest := filepath.Join("frontend", "assets", "uploads", fn)
			_ = os.Remove(dest)
		}
	}

	if cloudName, apiKey, apiSecret, ok := cloudinaryCredentials(); ok {
		for _, u := range imageURLs {
			if err := destroyCloudinaryImage(u, cloudName, apiKey, apiSecret); err != nil {
				log.Printf("warning: could not destroy Cloudinary image %s: %v", u, err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
