package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogHandler struct {
	repo models.ProductRepository
}

func NewCatalogHandler(r models.ProductRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Pagination Parameters (Requirement 4)
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Defaults and Constraints
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// 2. Parse Filter Parameters (Requirement 5)
	category := r.URL.Query().Get("category")
	priceLessThanStr := r.URL.Query().Get("priceLessThan")
	priceLessThan, _ := strconv.ParseFloat(priceLessThanStr, 64)

	// 3. Fetch Data from Repository
	products, total, err := h.repo.GetProducts(limit, offset, category, priceLessThan)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Could not fetch products")
		return
	}

	// 4. Return formatted response using your api utility
	api.OKResponse(w, map[string]any{
		"products": products,
		"total":    total,
	})
}
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	// Extract 'code' from URL (e.g., /catalog/PROD001)
	code := r.PathValue("code")

	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	product, err := h.repo.GetByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// api.OKResponse will automatically include the inherited prices
	// because the AfterFind hook ran when the repo called the DB.
	api.OKResponse(w, product)
}

// HandleGetCategories implements GET /categories (Task 3.1)
func (h *CatalogHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	api.OKResponse(w, categories)
}

// HandleCreateCategory implements POST /categories (Task 3.2)
func (h *CatalogHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	// Decode the JSON body into our model
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation: ensure Code and Name aren't empty
	if category.Code == "" || category.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Category code and name are required")
		return
	}

	if err := h.repo.CreateCategory(&category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	api.OKResponse(w, category)
}
