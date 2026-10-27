package models

type Business struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type BusinessInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type UserInput struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type Product struct {
	ID          string  `json:"id"`
	BusinessID  string  `json:"business_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
}

type ProductInput struct {
	BusinessID  string  `json:"business_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
}

type Billing struct {
	ID           string  `json:"id"`
	BusinessID   string  `json:"business_id"`
	UserID       string  `json:"user_id,omitempty"`
	ProductID    string  `json:"product_id,omitempty"`
	Reference    string  `json:"reference"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
	DueDate      string  `json:"due_date,omitempty"`
	PaidAt       string  `json:"paid_at,omitempty"`
	Notes        string  `json:"notes,omitempty"`
	CreatedAt    string  `json:"created_at"`
	BusinessName string  `json:"business_name,omitempty"`
	ProductName  string  `json:"product_name,omitempty"`
	UserName     string  `json:"user_name,omitempty"`
}

type BillingInput struct {
	BusinessID string  `json:"business_id"`
	UserID     string  `json:"user_id,omitempty"`
	ProductID  string  `json:"product_id,omitempty"`
	Reference  string  `json:"reference"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Status     string  `json:"status"`
	DueDate    string  `json:"due_date,omitempty"`
	Notes      string  `json:"notes,omitempty"`
}

type Transaction struct {
	ID           string  `json:"id"`
	BusinessID   string  `json:"business_id"`
	Customer     string  `json:"customer_name"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
	Reference    string  `json:"reference"`
	Notes        string  `json:"notes,omitempty"`
	CreatedAt    string  `json:"created_at"`
	BusinessName string  `json:"business_name,omitempty"`
}

type TransactionInput struct {
	BusinessID string  `json:"business_id"`
	Customer   string  `json:"customer_name"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Status     string  `json:"status"`
	Reference  string  `json:"reference"`
	Notes      string  `json:"notes,omitempty"`
}

type AuditLog struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Details    string `json:"details,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type Summary struct {
	BusinessCount    int     `json:"business_count"`
	TransactionCount int     `json:"transaction_count"`
	TotalRevenue     float64 `json:"total_revenue"`
	PendingAmount    float64 `json:"pending_amount"`
	ProductCount     int     `json:"product_count"`
	UserCount        int     `json:"user_count"`
	OutstandingBills float64 `json:"outstanding_bills"`
}
