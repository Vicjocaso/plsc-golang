package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plsc/golang/lib"
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
func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []types.Product
	result := p.tx.Preload("Category").Find(&products)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lib.EncodeOkReponse(w, products)
}

// getProduct is the HTTP handler for GET /products/{id}.
func (p *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product types.Product

	result := p.tx.First(&product, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	lib.EncodeOkReponse(w, product)
}

// getCategories is the HTTP handler for GET /categories.
func (p *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []types.Category
	result := p.tx.Find(&categories)
	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lib.EncodeOkReponse(w, categories)
}

// getCategory is the HTTP handler for GET /category/{id}.
func (p *ProductHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	var category types.Category
	result := p.tx.First(&category, mux.Vars(r)["id"])
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	lib.EncodeOkReponse(w, category)
}

// getRawQuery is the HTTP handler for GET /query.
func (p *ProductHandler) GetRawQuery(w http.ResponseWriter, r *http.Request) {
	var products []types.Product
	result := p.tx.Raw("SELECT * FROM products").Scan(&products)
	if result.Error != nil {
		http.NotFound(w, r)
		fmt.Println(result.Error)
		return
	}

	lib.EncodeOkReponse(w, products)
}

// AddProduct is the HTTP handler for POST /add/products.
func (p *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct types.Product

	json.NewDecoder(r.Body).Decode(&newProduct)

	product := p.tx.Create(&newProduct)
	if product.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lib.EncodeOkReponse(w, newProduct)

}
