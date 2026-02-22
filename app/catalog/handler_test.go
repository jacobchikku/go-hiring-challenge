package catalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// MockRepository allows us to simulate the DB logic without a real database
type MockRepository struct {
	models.ProductRepository
	GetProductsFunc    func(limit, offset int, category string, priceLessThan float64) ([]models.Product, int64, error)
	CreateCategoryFunc func(category *models.Category) error
}

func (m *MockRepository) GetProducts(limit, offset int, category string, priceLessThan float64) ([]models.Product, int64, error) {
	return m.GetProductsFunc(limit, offset, category, priceLessThan)
}

func (m *MockRepository) CreateCategory(category *models.Category) error {
	return m.CreateCategoryFunc(category)
}

// TestHandleGet_Features tests pagination and filtering
func TestHandleGet_Features(t *testing.T) {
	t.Run("Default Pagination", func(t *testing.T) {
		mockRepo := &MockRepository{
			GetProductsFunc: func(limit, offset int, category string, priceLessThan float64) ([]models.Product, int64, error) {
				if limit != 10 {
					t.Errorf("expected default limit 10, got %d", limit)
				}
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog", nil)
		rr := httptest.NewRecorder()

		handler.HandleGet(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Filtering Parameters", func(t *testing.T) {
		mockRepo := &MockRepository{
			GetProductsFunc: func(limit, offset int, category string, priceLessThan float64) ([]models.Product, int64, error) {
				if category != "clothing" || priceLessThan != 50.0 {
					t.Errorf("filters not passed correctly: got %s, %f", category, priceLessThan)
				}
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?category=clothing&priceLessThan=50.0", nil)
		rr := httptest.NewRecorder()

		handler.HandleGet(rr, req)
	})
}

// TestHandleCreateCategory tests the POST endpoint
func TestHandleCreateCategory(t *testing.T) {
	mockRepo := &MockRepository{
		CreateCategoryFunc: func(category *models.Category) error {
			if category.Code != "new-cat" {
				t.Errorf("expected code new-cat, got %s", category.Code)
			}
			return nil
		},
	}

	handler := NewCatalogHandler(mockRepo)

	// Create a dummy JSON body
	body := strings.NewReader(`{"code":"new-cat", "name":"New Category"}`)
	req := httptest.NewRequest("POST", "/categories", body)
	rr := httptest.NewRecorder()

	handler.HandleCreateCategory(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
