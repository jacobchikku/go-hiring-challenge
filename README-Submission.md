# Submission Notes - Go Hiring Challenge

## Completed Tasks
- **Task 1 & 2:** Implemented Catalog list and Product details with GORM.
- **Task 2.1:** Implemented Price Inheritance using an `AfterFind` GORM hook in the `Product` model.
- **Task 3:** Implemented Category listing and creation endpoints.
- **Task 4 & 5:** Added pagination, filtering, and comprehensive unit tests with 100% coverage on API utilities.

## Design Decisions
- **Repository Pattern:** Used an interface for `ProductRepository` to ensure the code is decoupled and easily testable via mocks.
- **GORM Hooks:** Chose `AfterFind` for price inheritance to centralize business logic, ensuring consistent pricing across all API endpoints (List and Details).
- **Go 1.22 Routing:** Utilized the enhanced `http.ServeMux` for clean path-parameter handling (`{code}`).

## API Documentation
The server runs on port 8484 by default (as configured in .env).
        Method Endpoint Description Query Params / Body
##      GET/catalogPaginated product list with filterslimit, offset, category, priceLessThan
##      GET/catalog/{code}Detailed product view with variants{code} (e.g., PROD001)
##      GET/categoriesList all available categoriesNone
##      POST/categoriesCreate a new categoryBody: {"code": "string", "name": "string"}