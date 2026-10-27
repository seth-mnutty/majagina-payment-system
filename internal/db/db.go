package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"setho-pymnt/internal/models"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func New(path string) (*Store, error) {
	if path == "" {
		path = "data/payment.db"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	dsn := path
	if !strings.HasPrefix(dsn, "file:") {
		dsn = "file:" + dsn
	}
	dsn += "?_foreign_keys=1"

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Store{DB: conn}, nil
}

func (s *Store) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS businesses (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			business_id TEXT NOT NULL,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'staff',
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,
			UNIQUE (business_id, email)
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			business_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'USD',
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS billings (
			id TEXT PRIMARY KEY,
			business_id TEXT NOT NULL,
			user_id TEXT,
			product_id TEXT,
			reference TEXT NOT NULL UNIQUE,
			amount REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'USD',
			status TEXT NOT NULL DEFAULT 'pending',
			due_date TEXT,
			paid_at TEXT,
			notes TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			business_id TEXT NOT NULL,
			customer_name TEXT NOT NULL,
			amount REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'USD',
			status TEXT NOT NULL DEFAULT 'pending',
			reference TEXT NOT NULL UNIQUE,
			notes TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			entity_type TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			action TEXT NOT NULL,
			details TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := s.DB.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Close() error { return s.DB.Close() }

func (s *Store) ListBusinesses() ([]models.Business, error) {
	rows, err := s.DB.Query(`SELECT id, name, email, status, created_at FROM businesses ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	businesses := []models.Business{}
	for rows.Next() {
		var b models.Business
		if err := rows.Scan(&b.ID, &b.Name, &b.Email, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, b)
	}
	return businesses, rows.Err()
}

func (s *Store) CreateBusiness(input models.BusinessInput) (models.Business, error) {
	business := models.Business{
		ID:        newID("biz"),
		Name:      strings.TrimSpace(input.Name),
		Email:     strings.TrimSpace(input.Email),
		Status:    "active",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	_, err := s.DB.Exec(
		`INSERT INTO businesses (id, name, email, status, created_at) VALUES (?, ?, ?, ?, ?)`,
		business.ID, business.Name, business.Email, business.Status, business.CreatedAt,
	)
	if err != nil {
		return models.Business{}, err
	}
	if err := s.LogAction("business", business.ID, "created", fmt.Sprintf("Business %s created", business.Name)); err != nil {
		return models.Business{}, err
	}
	return business, nil
}

func (s *Store) GetBusiness(id string) (models.Business, error) {
	var b models.Business
	err := s.DB.QueryRow(`SELECT id, name, email, status, created_at FROM businesses WHERE id = ?`, id).Scan(&b.ID, &b.Name, &b.Email, &b.Status, &b.CreatedAt)
	if err != nil {
		return models.Business{}, err
	}
	return b, nil
}

func (s *Store) ListUsers(businessID string) ([]models.User, error) {
	query := `SELECT id, business_id, name, email, role, status, created_at FROM users`
	args := []any{}
	if businessID != "" {
		query += ` WHERE business_id = ?`
		args = append(args, businessID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.BusinessID, &u.Name, &u.Email, &u.Role, &u.Status, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *Store) CreateUser(input models.UserInput) (models.User, error) {
	if _, err := s.GetBusiness(input.BusinessID); err != nil {
		return models.User{}, fmt.Errorf("business not found: %w", err)
	}
	user := models.User{
		ID:         newID("usr"),
		BusinessID: input.BusinessID,
		Name:       strings.TrimSpace(input.Name),
		Email:      strings.TrimSpace(input.Email),
		Role:       input.Role,
		Status:     "active",
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	if user.Role == "" {
		user.Role = "staff"
	}

	_, err := s.DB.Exec(
		`INSERT INTO users (id, business_id, name, email, role, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.BusinessID, user.Name, user.Email, user.Role, user.Status, user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	if err := s.LogAction("user", user.ID, "created", fmt.Sprintf("User %s created", user.Name)); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Store) GetUser(id string) (models.User, error) {
	var u models.User
	err := s.DB.QueryRow(`SELECT id, business_id, name, email, role, status, created_at FROM users WHERE id = ?`, id).Scan(&u.ID, &u.BusinessID, &u.Name, &u.Email, &u.Role, &u.Status, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (s *Store) ListProducts(businessID string) ([]models.Product, error) {
	query := `SELECT id, business_id, name, description, price, currency, is_active, created_at FROM products`
	args := []any{}
	if businessID != "" {
		query += ` WHERE business_id = ?`
		args = append(args, businessID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		var isActive int
		if err := rows.Scan(&p.ID, &p.BusinessID, &p.Name, &p.Description, &p.Price, &p.Currency, &isActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.IsActive = isActive == 1
		products = append(products, p)
	}
	return products, rows.Err()
}

func (s *Store) CreateProduct(input models.ProductInput) (models.Product, error) {
	if _, err := s.GetBusiness(input.BusinessID); err != nil {
		return models.Product{}, fmt.Errorf("business not found: %w", err)
	}
	product := models.Product{
		ID:          newID("prod"),
		BusinessID:  input.BusinessID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Price:       input.Price,
		Currency:    input.Currency,
		IsActive:    true,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	if product.Currency == "" {
		product.Currency = "USD"
	}

	_, err := s.DB.Exec(
		`INSERT INTO products (id, business_id, name, description, price, currency, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.BusinessID, product.Name, product.Description, product.Price, product.Currency, boolToInt(product.IsActive), product.CreatedAt,
	)
	if err != nil {
		return models.Product{}, err
	}
	if err := s.LogAction("product", product.ID, "created", fmt.Sprintf("Product %s created", product.Name)); err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (s *Store) GetProduct(id string) (models.Product, error) {
	var p models.Product
	var isActive int
	err := s.DB.QueryRow(`SELECT id, business_id, name, description, price, currency, is_active, created_at FROM products WHERE id = ?`, id).Scan(&p.ID, &p.BusinessID, &p.Name, &p.Description, &p.Price, &p.Currency, &isActive, &p.CreatedAt)
	if err != nil {
		return models.Product{}, err
	}
	p.IsActive = isActive == 1
	return p, nil
}

func (s *Store) ListBillings(businessID string) ([]models.Billing, error) {
	query := `SELECT b.id, b.business_id, b.user_id, b.product_id, b.reference, b.amount, b.currency, b.status, b.due_date, b.paid_at, b.notes, b.created_at, bus.name, p.name, u.name
		FROM billings b
		LEFT JOIN businesses bus ON bus.id = b.business_id
		LEFT JOIN products p ON p.id = b.product_id
		LEFT JOIN users u ON u.id = b.user_id`
	args := []any{}
	if businessID != "" {
		query += ` WHERE b.business_id = ?`
		args = append(args, businessID)
	}
	query += ` ORDER BY b.created_at DESC`

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	billings := []models.Billing{}
	for rows.Next() {
		var bill models.Billing
		if err := rows.Scan(&bill.ID, &bill.BusinessID, &bill.UserID, &bill.ProductID, &bill.Reference, &bill.Amount, &bill.Currency, &bill.Status, &bill.DueDate, &bill.PaidAt, &bill.Notes, &bill.CreatedAt, &bill.BusinessName, &bill.ProductName, &bill.UserName); err != nil {
			return nil, err
		}
		billings = append(billings, bill)
	}
	return billings, rows.Err()
}

func (s *Store) CreateBilling(input models.BillingInput) (models.Billing, error) {
	if _, err := s.GetBusiness(input.BusinessID); err != nil {
		return models.Billing{}, fmt.Errorf("business not found: %w", err)
	}
	if input.UserID != "" {
		if _, err := s.GetUser(input.UserID); err != nil {
			return models.Billing{}, fmt.Errorf("user not found: %w", err)
		}
	}
	if input.ProductID != "" {
		if _, err := s.GetProduct(input.ProductID); err != nil {
			return models.Billing{}, fmt.Errorf("product not found: %w", err)
		}
	}

	billing := models.Billing{
		ID:         newID("bill"),
		BusinessID: input.BusinessID,
		UserID:     input.UserID,
		ProductID:  input.ProductID,
		Reference:  strings.TrimSpace(input.Reference),
		Amount:     input.Amount,
		Currency:   input.Currency,
		Status:     input.Status,
		DueDate:    input.DueDate,
		Notes:      input.Notes,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	if billing.Currency == "" {
		billing.Currency = "USD"
	}
	if billing.Status == "" {
		billing.Status = "pending"
	}
	if billing.PaidAt == "" && billing.Status == "paid" {
		billing.PaidAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.DB.Exec(
		`INSERT INTO billings (id, business_id, user_id, product_id, reference, amount, currency, status, due_date, paid_at, notes, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		billing.ID, billing.BusinessID, billing.UserID, billing.ProductID, billing.Reference, billing.Amount, billing.Currency,
		billing.Status, billing.DueDate, billing.PaidAt, billing.Notes, billing.CreatedAt,
	)
	if err != nil {
		return models.Billing{}, err
	}
	if err := s.LogAction("billing", billing.ID, "created", fmt.Sprintf("Billing %s created", billing.Reference)); err != nil {
		return models.Billing{}, err
	}
	return billing, nil
}

func (s *Store) GetBilling(id string) (models.Billing, error) {
	var bill models.Billing
	err := s.DB.QueryRow(`SELECT b.id, b.business_id, b.user_id, b.product_id, b.reference, b.amount, b.currency, b.status, b.due_date, b.paid_at, b.notes, b.created_at, bus.name, p.name, u.name
		FROM billings b
		LEFT JOIN businesses bus ON bus.id = b.business_id
		LEFT JOIN products p ON p.id = b.product_id
		LEFT JOIN users u ON u.id = b.user_id
		WHERE b.id = ?`, id).Scan(&bill.ID, &bill.BusinessID, &bill.UserID, &bill.ProductID, &bill.Reference, &bill.Amount, &bill.Currency, &bill.Status, &bill.DueDate, &bill.PaidAt, &bill.Notes, &bill.CreatedAt, &bill.BusinessName, &bill.ProductName, &bill.UserName)
	if err != nil {
		return models.Billing{}, err
	}
	return bill, nil
}

func (s *Store) CreateTransaction(input models.TransactionInput) (models.Transaction, error) {
	if input.BusinessID == "" {
		return models.Transaction{}, fmt.Errorf("business_id is required")
	}
	if input.Customer == "" {
		return models.Transaction{}, fmt.Errorf("customer_name is required")
	}
	if input.Reference == "" {
		return models.Transaction{}, fmt.Errorf("reference is required")
	}
	if input.Amount <= 0 {
		return models.Transaction{}, fmt.Errorf("amount must be greater than zero")
	}

	if _, err := s.GetBusiness(input.BusinessID); err != nil {
		return models.Transaction{}, fmt.Errorf("business not found: %w", err)
	}

	status := input.Status
	if status == "" {
		status = "pending"
	}
	currency := input.Currency
	if currency == "" {
		currency = "USD"
	}

	transaction := models.Transaction{
		ID:         newID("txn"),
		BusinessID: input.BusinessID,
		Customer:   strings.TrimSpace(input.Customer),
		Amount:     input.Amount,
		Currency:   currency,
		Status:     status,
		Reference:  strings.TrimSpace(input.Reference),
		Notes:      input.Notes,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	_, err := s.DB.Exec(
		`INSERT INTO transactions (id, business_id, customer_name, amount, currency, status, reference, notes, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.ID, transaction.BusinessID, transaction.Customer, transaction.Amount, transaction.Currency,
		transaction.Status, transaction.Reference, transaction.Notes, transaction.CreatedAt,
	)
	if err != nil {
		return models.Transaction{}, err
	}
	if err := s.LogAction("transaction", transaction.ID, "created", fmt.Sprintf("Payment %s recorded for %s", transaction.Reference, transaction.Customer)); err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func (s *Store) ListTransactions(businessID string) ([]models.Transaction, error) {
	query := `SELECT t.id, t.business_id, t.customer_name, t.amount, t.currency, t.status, t.reference, t.notes, t.created_at, b.name
		FROM transactions t
		LEFT JOIN businesses b ON b.id = t.business_id`
	args := []any{}
	if businessID != "" {
		query += ` WHERE t.business_id = ?`
		args = append(args, businessID)
	}
	query += ` ORDER BY t.created_at DESC`

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.BusinessID, &tx.Customer, &tx.Amount, &tx.Currency, &tx.Status, &tx.Reference, &tx.Notes, &tx.CreatedAt, &tx.BusinessName); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, rows.Err()
}

func (s *Store) GetTransaction(id string) (models.Transaction, error) {
	var tx models.Transaction
	err := s.DB.QueryRow(
		`SELECT t.id, t.business_id, t.customer_name, t.amount, t.currency, t.status, t.reference, t.notes, t.created_at, b.name
		FROM transactions t
		LEFT JOIN businesses b ON b.id = t.business_id
		WHERE t.id = ?`,
		id,
	).Scan(&tx.ID, &tx.BusinessID, &tx.Customer, &tx.Amount, &tx.Currency, &tx.Status, &tx.Reference, &tx.Notes, &tx.CreatedAt, &tx.BusinessName)
	if err != nil {
		return models.Transaction{}, err
	}
	return tx, nil
}

func (s *Store) Summary() (models.Summary, error) {
	var summary models.Summary
	if err := s.DB.QueryRow(`SELECT COUNT(*) FROM businesses`).Scan(&summary.BusinessCount); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&summary.TransactionCount); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&summary.UserCount); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&summary.ProductCount); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE status = 'paid'`).Scan(&summary.TotalRevenue); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE status = 'pending'`).Scan(&summary.PendingAmount); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM billings WHERE status IN ('pending', 'overdue')`).Scan(&summary.OutstandingBills); err != nil {
		return models.Summary{}, err
	}
	return summary, nil
}

func (s *Store) LogAction(entityType, entityID, action, details string) error {
	_, err := s.DB.Exec(
		`INSERT INTO audit_logs (id, entity_type, entity_id, action, details, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		newID("log"), entityType, entityID, action, details, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
