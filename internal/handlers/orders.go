package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"beads-haven/internal/models"
)

type OrderHandler struct {
	DB *sql.DB
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var req models.Order
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order payload")
		return
	}

	if req.CustomerName == "" || req.CustomerPhone == "" || req.DeliveryLocation == "" {
		writeJSONError(w, http.StatusBadRequest, "Name, phone, and delivery location are required")
		return
	}

	if len(req.Items) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Order must contain at least one item")
		return
	}

	orderID := fmt.Sprintf("ORD-%d-%04d", time.Now().Year(), rand.Intn(10000))
	req.ID = orderID
	req.Status = "pending"

	tx, err := h.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO orders (id, customer_name, customer_phone, customer_email, delivery_location, payment_method, mpesa_phone, subtotal, delivery_fee, total_amount, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, req.ID, req.CustomerName, req.CustomerPhone, req.CustomerEmail, req.DeliveryLocation, req.PaymentMethod, req.MpesaPhone, req.Subtotal, req.DeliveryFee, req.TotalAmount, req.Status)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %s", err.Error()))
		return
	}

	for _, item := range req.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, size, quantity, price)
			VALUES ($1, $2, $3, $4, $5)
		`, req.ID, item.ProductID, item.Size, item.Quantity, item.Price)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert order item: %s", err.Error()))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"orderId": orderID,
		"message": "Order placed successfully",
	})
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "Order ID required")
		return
	}

	var ord models.Order
	err := h.DB.QueryRow(`
		SELECT id, customer_name, customer_phone, customer_email, delivery_location, 
		       payment_method, mpesa_phone, subtotal, delivery_fee, total_amount, status, created_at
		FROM orders
		WHERE id = $1
	`, id).Scan(
		&ord.ID, &ord.CustomerName, &ord.CustomerPhone, &ord.CustomerEmail, &ord.DeliveryLocation,
		&ord.PaymentMethod, &ord.MpesaPhone, &ord.Subtotal, &ord.DeliveryFee, &ord.TotalAmount, &ord.Status, &ord.CreatedAt,
	)

	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Order not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Fetch items
	rows, err := h.DB.Query(`SELECT product_id, size, quantity, price FROM order_items WHERE order_id = $1`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it models.OrderItem
			if err := rows.Scan(&it.ProductID, &it.Size, &it.Quantity, &it.Price); err == nil {
				ord.Items = append(ord.Items, it)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ord)
}
