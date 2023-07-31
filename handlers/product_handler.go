package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plsc/golang/types"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ProductHandler struct {
	tx gorm.DB
}

func NewProductHandler(bdConn gorm.DB) *ProductHandler {
	return &ProductHandler{
		tx: bdConn,
	}
}

// getProducts is the HTTP handler for GET /products.
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []types.Product
	result := h.tx.Preload("Category").Find(&products)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(products)
}

// getProduct is the HTTP handler for GET /products/{id}.
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product types.Product
	result := h.tx.First(&product, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// getCategories is the HTTP handler for GET /categories.
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []types.Category
	result := h.tx.Find(&categories)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(categories)
}

// getCategory is the HTTP handler for GET /category/{id}.
func (h *ProductHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	var category types.Category
	result := h.tx.First(&category, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(category)
}

// getRawQuery is the HTTP handler for GET /query.
func (h *ProductHandler) GetRawQuery(w http.ResponseWriter, r *http.Request) {
	var products []types.Product
	result := h.tx.Raw("SELECT * FROM products").Scan(&products)
	if result.Error != nil {
		http.NotFound(w, r)
		fmt.Println(result.Error)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(result)
}
