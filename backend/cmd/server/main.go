package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"beads-haven/internal/db"
	"beads-haven/internal/handlers"
	"beads-haven/internal/middleware"
)

func loadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

func main() {
	loadEnvFile(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/beads_haven?sslmode=disable"
	}

	var database *sql.DB
	var err error
	database, err = db.Connect(dbURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to PostgreSQL (%v). Operating in offline/mock mode.", err)
	} else {
		defer database.Close()
		if err := db.EnsureSchema(database); err != nil {
			log.Printf("Warning: Schema check returned: %v", err)
		}
	}

	mux := http.NewServeMux()

	// Initialize Handlers
	prodHandler := &handlers.ProductHandler{DB: database}
	orderHandler := &handlers.OrderHandler{DB: database}
	authHandler := &handlers.AuthHandler{DB: database}
	uploadHandler := &handlers.UploadHandler{DB: database}

	// API Routes (Go 1.22+ routing)
	mux.HandleFunc("GET /api/products", prodHandler.List)
	mux.HandleFunc("GET /api/products/{id}", prodHandler.GetByID)

	mux.HandleFunc("POST /api/admin/products", middleware.RequireAdmin(prodHandler.Create))
	mux.HandleFunc("PUT /api/admin/products/{id}", middleware.RequireAdmin(prodHandler.Update))
	mux.HandleFunc("DELETE /api/admin/products/{id}", middleware.RequireAdmin(prodHandler.Delete))
	mux.HandleFunc("POST /api/admin/products/{id}/images", middleware.RequireAdmin(uploadHandler.UploadImage))

	mux.HandleFunc("POST /api/orders", orderHandler.Create)
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetByID)

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Static Assets & Clean Page Routing
	frontendDir := "frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "."
	}

	fs := http.FileServer(http.Dir(frontendDir))

	// Root & Clean URL Router
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Clean routing aliases
		switch path {
		case "/", "/index.html":
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		case "/bracelets":
			http.Redirect(w, r, "/shop.html?type=bracelets", http.StatusFound)
			return
		case "/necklaces":
			http.Redirect(w, r, "/shop.html?type=necklaces", http.StatusFound)
			return
		case "/rings":
			http.Redirect(w, r, "/shop.html?type=rings", http.StatusFound)
			return
		case "/bags":
			http.Redirect(w, r, "/shop.html?type=bags", http.StatusFound)
			return
		case "/fashion-accessories":
			http.Redirect(w, r, "/shop.html?type=fashion-accessories", http.StatusFound)
			return
		case "/art-decor":
			http.Redirect(w, r, "/shop.html?type=art-decor", http.StatusFound)
			return
		case "/sneakers", "/jerseys":
			http.Redirect(w, r, "/shop.html", http.StatusMovedPermanently)
			return
		case "/shop":
			http.ServeFile(w, r, filepath.Join(frontendDir, "shop.html"))
			return
		case "/cart":
			http.ServeFile(w, r, filepath.Join(frontendDir, "cart.html"))
			return
		case "/checkout":
			http.ServeFile(w, r, filepath.Join(frontendDir, "checkout.html"))
			return
		case "/about":
			http.ServeFile(w, r, filepath.Join(frontendDir, "about.html"))
			return
		case "/contact":
			http.ServeFile(w, r, filepath.Join(frontendDir, "contact.html"))
			return
		case "/admin":
			http.ServeFile(w, r, filepath.Join(frontendDir, "admin.html"))
			return
		}

		if strings.HasPrefix(path, "/product/") {
			id := strings.TrimPrefix(path, "/product/")
			http.Redirect(w, r, "/product.html?id="+id, http.StatusFound)
			return
		}

		// Serve static files from /frontend
		targetPath := filepath.Join(frontendDir, path)
		if fi, err := os.Stat(targetPath); err == nil && !fi.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback to 404 or index
		fs.ServeHTTP(w, r)
	})

	handlerWithCORS := middleware.EnableCORS(mux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Beads Haven server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
