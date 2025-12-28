package models

// ==================== AUTH MODELS ====================

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email" example:"user@example.com"`
	Password    string `json:"password" binding:"required,min=6" example:"password123"`
	Username    string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	FullName    string `json:"full_name" binding:"required" example:"John Doe"`
	PhoneNumber string `json:"phone_number" example:"+1234567890"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Message string   `json:"message" example:"login successful"`
	Data    AuthData `json:"data"`
}

// AuthData contains token and user info
type AuthData struct {
	Token string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *UserInfo `json:"user"`
}

// ==================== USER MODELS ====================

// UserInfo represents user information
type UserInfo struct {
	ID          uint   `json:"id" example:"1"`
	Email       string `json:"email" example:"user@example.com"`
	Username    string `json:"username" example:"johndoe"`
	FullName    string `json:"full_name" example:"John Doe"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+1234567890"`
	Role        string `json:"role" example:"user"`
	Status      string `json:"status" example:"ACTIVE"`
}

// UpdateProfileRequest represents the request body for updating profile
type UpdateProfileRequest struct {
	FullName    string `json:"full_name,omitempty" example:"John Doe Updated"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+1234567890"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpassword123"`
}

// ==================== ADDRESS MODELS ====================

// CreateAddressRequest represents the request body for creating an address
type CreateAddressRequest struct {
	Street     string `json:"street" binding:"required" example:"123 Main St"`
	City       string `json:"city" binding:"required" example:"Ho Chi Minh City"`
	State      string `json:"state" binding:"required" example:"Ho Chi Minh"`
	PostalCode string `json:"postal_code" binding:"required" example:"70000"`
	Country    string `json:"country" binding:"required" example:"Vietnam"`
	IsDefault  bool   `json:"is_default" example:"false"`
}

// UpdateAddressRequest represents the request body for updating an address
type UpdateAddressRequest struct {
	Street     string `json:"street,omitempty" example:"123 Main St"`
	City       string `json:"city,omitempty" example:"Ho Chi Minh City"`
	State      string `json:"state,omitempty" example:"Ho Chi Minh"`
	PostalCode string `json:"postal_code,omitempty" example:"70000"`
	Country    string `json:"country,omitempty" example:"Vietnam"`
	IsDefault  *bool  `json:"is_default,omitempty" example:"false"`
}

// AddressInfo represents address information
type AddressInfo struct {
	ID         uint   `json:"id" example:"1"`
	UserID     uint   `json:"user_id" example:"1"`
	Street     string `json:"street" example:"123 Main St"`
	City       string `json:"city" example:"Ho Chi Minh City"`
	State      string `json:"state" example:"Ho Chi Minh"`
	PostalCode string `json:"postal_code" example:"70000"`
	Country    string `json:"country" example:"Vietnam"`
	IsDefault  bool   `json:"is_default" example:"true"`
}

// ==================== PRODUCT MODELS ====================

// Product represents a product
type Product struct {
	ID          uint     `json:"id" example:"1"`
	Name        string   `json:"name" example:"iPhone 15 Pro"`
	Description string   `json:"description" example:"Latest iPhone with A17 Pro chip"`
	Price       float64  `json:"price" example:"999.99"`
	SKU         string   `json:"sku" example:"IPH15P-001"`
	CategoryID  *uint    `json:"category_id,omitempty" example:"1"`
	Category    *Category `json:"category,omitempty"`
	Status      string   `json:"status" example:"ACTIVE"`
	Images      []string `json:"images,omitempty"`
	Stock       int      `json:"stock" example:"50"`
	IsActive    bool     `json:"is_active" example:"true"`
	CreatedAt   string   `json:"created_at" example:"2025-12-23T10:00:00Z"`
	UpdatedAt   string   `json:"updated_at" example:"2025-12-23T10:00:00Z"`
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required" example:"iPhone 15 Pro"`
	Description string   `json:"description" binding:"required" example:"Latest iPhone with A17 Pro chip"`
	Price       float64  `json:"price" binding:"required,gt=0" example:"999.99"`
	SKU         string   `json:"sku" binding:"required" example:"IPH15P-001"`
	CategoryID  *uint    `json:"category_id,omitempty" example:"1"`
	Status      string   `json:"status" example:"ACTIVE"`
	Images      []string `json:"images,omitempty"`
	Stock       int      `json:"stock" binding:"gte=0" example:"50"`
	IsActive    bool     `json:"is_active" example:"true"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name,omitempty" example:"iPhone 15 Pro Updated"`
	Description string   `json:"description,omitempty" example:"Updated description"`
	Price       *float64 `json:"price,omitempty" example:"1099.99"`
	CategoryID  *uint    `json:"category_id,omitempty" example:"1"`
	Status      string   `json:"status,omitempty" example:"ACTIVE"`
	Images      []string `json:"images,omitempty"`
	Stock       *int     `json:"stock,omitempty" example:"60"`
	IsActive    *bool    `json:"is_active,omitempty" example:"true"`
}

// ProductsResponse represents paginated products response
type ProductsResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total" example:"100"`
	Page     int       `json:"page" example:"1"`
	Limit    int       `json:"limit" example:"20"`
}

// ==================== CATEGORY MODELS ====================

// Category represents a category
type Category struct {
	ID          uint       `json:"id" example:"1"`
	Name        string     `json:"name" example:"Electronics"`
	Slug        string     `json:"slug" example:"electronics"`
	Description string     `json:"description,omitempty" example:"Electronic devices and gadgets"`
	ParentID    *uint      `json:"parent_id,omitempty" example:"0"`
	Parent      *Category  `json:"parent,omitempty"`
	Children    []Category `json:"children,omitempty"`
	CreatedAt   string     `json:"created_at" example:"2025-12-23T10:00:00Z"`
	UpdatedAt   string     `json:"updated_at" example:"2025-12-23T10:00:00Z"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required" example:"Electronics"`
	Slug        string `json:"slug,omitempty" example:"electronics"`
	Description string `json:"description,omitempty" example:"Electronic devices and gadgets"`
	ParentID    *uint  `json:"parent_id,omitempty" example:"0"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty" example:"Electronics Updated"`
	Slug        string `json:"slug,omitempty" example:"electronics-updated"`
	Description string `json:"description,omitempty" example:"Updated description"`
	ParentID    *uint  `json:"parent_id,omitempty" example:"0"`
}

// ==================== COMMON MODELS ====================

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad request"`
	Message string `json:"message,omitempty" example:"Invalid input"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ==================== SEARCH MODELS ====================

// SearchResponse represents search results with pagination
type SearchResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total" example:"100"`
	Page     int       `json:"page" example:"1"`
	Limit    int       `json:"limit" example:"20"`
}


