package dto

// ResponseEnvelope wraps all successful responses.
type ResponseEnvelope struct {
    Status  string      `json:"status"`  // e.g., "success"
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data"`  // Specific type per endpoint
}

// ErrorResponse standardizes all error responses.
type ErrorResponse struct {
    Error struct {
        Message string      `json:"message"`
        Code    string      `json:"code"`  // e.g., "INVALID_INPUT", "UNAUTHORIZED"
        Details interface{} `json:"details,omitempty"`  // e.g., map or []ErrorDetail
    } `json:"error"`
}

// ErrorDetail for validation or field-specific errors.
type ErrorDetail struct {
    Field   string `json:"field,omitempty"`
    Message string `json:"message"`
}

// PaginationMeta for list metadata.
type PaginationMeta struct {
    Total      int `json:"total"`
    Page       int `json:"page"`
    PageSize   int `json:"pageSize"`
    TotalPages int `json:"totalPages"`
}

// PaginatedResponse for list endpoints.
type PaginatedResponse struct {
    Data interface{} `json:"data"`  // Usually []SomeItem
    Meta PaginationMeta `json:"meta"`
}