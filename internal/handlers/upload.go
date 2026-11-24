package handlers

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"beads-haven/internal/models"
)

type UploadHandler struct {
	DB *sql.DB
}

type CloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Error     *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	productID := r.PathValue("id")
	if productID == "" {
		writeJSONError(w, http.StatusBadRequest, "Product ID required")
		return
	}

	// Verify product exists
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil || !exists {
		writeJSONError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Allow up to 32 MB multipart payload for multi-file upload
	maxMultipart := int64(32 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxMultipart)
	if err := r.ParseMultipartForm(maxMultipart); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Request size exceeds limit or malformed form")
		return
	}

	var fileHeaders []*multipart.FileHeader
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs, ok := r.MultipartForm.File["images"]; ok {
			fileHeaders = append(fileHeaders, fhs...)
		}
		if fhs, ok := r.MultipartForm.File["image"]; ok {
			fileHeaders = append(fileHeaders, fhs...)
		}
	}

	if len(fileHeaders) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Image file is required")
		return
	}

	// Validate all files before saving
	for _, fh := range fileHeaders {
		if fh.Size > 5*1024*1024 {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("File '%s' exceeds 5 MB limit", fh.Filename))
			return
		}

		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid file format for '%s'. Only JPG, PNG, and WebP are allowed", fh.Filename))
			return
		}

		f, err := fh.Open()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read file '%s'", fh.Filename))
			return
		}
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		f.Close()

		sniffType := http.DetectContentType(buf[:n])
		if sniffType != "image/jpeg" && sniffType != "image/png" && sniffType != "image/webp" {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("File '%s' failed content verification (detected: %s)", fh.Filename, sniffType))
			return
		}
	}

	// Remove placeholder images for this product automatically when real image is added
	_, _ = h.DB.Exec(`
		DELETE FROM product_images 
		WHERE product_id = $1 AND (url LIKE '%/placeholders/%' OR url LIKE '%default.svg')
	`, productID)

	var maxOrder int
	h.DB.QueryRow("SELECT COALESCE(MAX(display_order), -1) FROM product_images WHERE product_id = $1", productID).Scan(&maxOrder)

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	useCloudinary := cloudName != "" && apiKey != "" && apiSecret != ""

	var createdImages []models.ProductImage

	for i, fh := range fileHeaders {
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		f, err := fh.Open()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to process image file")
			return
		}

		var imageURL string
		if useCloudinary {
			cloudURL, err := uploadToCloudinary(f, fh.Filename, cloudName, apiKey, apiSecret)
			f.Close()
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Cloudinary upload failed: %s", err.Error()))
				return
			}
			imageURL = cloudURL
		} else {
			filename := fmt.Sprintf("%s-%d-%d%s", productID, time.Now().UnixNano(), i, ext)
			destPath := filepath.Join("frontend", "assets", "uploads", filename)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				f.Close()
				writeJSONError(w, http.StatusInternalServerError, "Failed to create uploads directory")
				return
			}

			out, err := os.Create(destPath)
			if err != nil {
				f.Close()
				writeJSONError(w, http.StatusInternalServerError, "Failed to save image locally")
				return
			}

			if _, err := io.Copy(out, f); err != nil {
				out.Close()
				f.Close()
				writeJSONError(w, http.StatusInternalServerError, "Failed to write image data")
				return
			}
			out.Close()
			f.Close()

			imageURL = "/assets/uploads/" + filename
		}

		currentOrder := maxOrder + 1 + i
		var imgID int
		err = h.DB.QueryRow(`
			INSERT INTO product_images (product_id, url, display_order)
			VALUES ($1, $2, $3)
			RETURNING id
		`, productID, imageURL, currentOrder).Scan(&imgID)

		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to record image in database: %s", err.Error()))
			return
		}

		createdImages = append(createdImages, models.ProductImage{
			ID:           imgID,
			ProductID:    productID,
			URL:          imageURL,
			DisplayOrder: currentOrder,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdImages)
}

func (h *UploadHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	productID := r.PathValue("id")
	imageIdStr := r.PathValue("imageId")
	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil || imageId <= 0 {
		writeJSONError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	var imageURL string
	err = h.DB.QueryRow("SELECT url FROM product_images WHERE id = $1 AND product_id = $2", imageId, productID).Scan(&imageURL)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Image not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	_, err = h.DB.Exec("DELETE FROM product_images WHERE id = $1 AND product_id = $2", imageId, productID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete image record")
		return
	}

	// Delete local file if applicable
	if strings.HasPrefix(imageURL, "/assets/uploads/") {
		fn := filepath.Base(imageURL)
		dest := filepath.Join("frontend", "assets", "uploads", fn)
		_ = os.Remove(dest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image deleted successfully",
	})
}

func (h *UploadHandler) ReorderImages(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	productID := r.PathValue("id")
	var req struct {
		Order []int `json:"order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Order) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Valid image order array is required")
		return
	}

	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil || !exists {
		writeJSONError(w, http.StatusNotFound, "Product not found")
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to begin transaction")
		return
	}
	defer tx.Rollback()

	for idx, imgID := range req.Order {
		_, err = tx.Exec("UPDATE product_images SET display_order = $1 WHERE id = $2 AND product_id = $3", idx, imgID, productID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update image order for id %d: %s", imgID, err.Error()))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to commit image order")
		return
	}

	var updated []models.ProductImage
	rows, err := h.DB.Query("SELECT id, product_id, url, display_order FROM product_images WHERE product_id = $1 ORDER BY display_order ASC, id ASC", productID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var img models.ProductImage
			if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder); err == nil {
				updated = append(updated, img)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func uploadToCloudinary(file multipart.File, filename, cloudName, apiKey, apiSecret string) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	sigStr := fmt.Sprintf("timestamp=%s%s", timestamp, apiSecret)
	h := sha1.New()
	h.Write([]byte(sigStr))
	signature := hex.EncodeToString(h.Sum(nil))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.WriteField("api_key", apiKey)
	writer.WriteField("timestamp", timestamp)
	writer.WriteField("signature", signature)
	writer.Close()

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", cloudName)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var cloudRes CloudinaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudRes); err != nil {
		return "", err
	}

	if cloudRes.Error != nil {
		return "", fmt.Errorf("%s", cloudRes.Error.Message)
	}

	if cloudRes.SecureURL == "" {
		return "", fmt.Errorf("no secure URL returned from Cloudinary")
	}

	return cloudRes.SecureURL, nil
}
