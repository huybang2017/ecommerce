package handler

// ErrorResponse defines a standard error response for Swagger
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"description of the error"`
}

// SuccessResponse defines a standard success message response
// @name SuccessResponse
type SuccessResponse struct {
	Message string `json:"message" example:"operation completed successfully"`
}

// OrdersListResponse is used for documenting the list orders response
// @name OrdersListResponse
type OrdersListResponse struct {
	Orders []*interface{} `json:"orders"` // Use interface to reference domain.Order via handler docs if needed
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
