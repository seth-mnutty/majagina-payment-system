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
}
