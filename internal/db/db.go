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

func (s *Store) Close() error {
	return s.DB.Close()
}

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
	err := s.DB.QueryRow(
		`SELECT id, name, email, status, created_at FROM businesses WHERE id = ?`,
		id,
	).Scan(&b.ID, &b.Name, &b.Email, &b.Status, &b.CreatedAt)
	if err != nil {
		return models.Business{}, err
	}
	return b, nil
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
	if err := s.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE status = 'paid'`).Scan(&summary.TotalRevenue); err != nil {
		return models.Summary{}, err
	}
	if err := s.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE status = 'pending'`).Scan(&summary.PendingAmount); err != nil {
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

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
