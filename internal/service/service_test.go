package service

import (
	"path/filepath"
	"testing"

	"setho-pymnt/internal/db"
	"setho-pymnt/internal/models"
)

func TestServiceBusinessAndTransactionLifecycle(t *testing.T) {
	store, err := db.New(filepath.Join(t.TempDir(), "payment.db"))
	if err != nil {
		t.Fatalf("db.New() error = %v", err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	svc := New(store)

	business, err := svc.CreateBusiness(models.BusinessInput{
		Name:  "Acme Labs",
		Email: "billing@acme.com",
	})
	if err != nil {
		t.Fatalf("CreateBusiness() error = %v", err)
	}
	if business.ID == "" {
		t.Fatal("CreateBusiness() returned empty ID")
	}

	transaction, err := svc.CreateTransaction(models.TransactionInput{
		BusinessID: business.ID,
		Customer:   "Jane Doe",
		Amount:     199.99,
		Currency:   "USD",
		Status:     "paid",
		Reference:  "INV-1001",
		Notes:      "Monthly subscription",
	})
	if err != nil {
		t.Fatalf("CreateTransaction() error = %v", err)
	}
	if transaction.Reference != "INV-1001" {
		t.Fatalf("expected reference INV-1001, got %s", transaction.Reference)
	}

	summary, err := svc.Summary()
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if summary.BusinessCount != 1 {
		t.Fatalf("expected 1 business, got %d", summary.BusinessCount)
	}
	if summary.TransactionCount != 1 {
		t.Fatalf("expected 1 transaction, got %d", summary.TransactionCount)
	}
	if summary.TotalRevenue != 199.99 {
		t.Fatalf("expected total revenue 199.99, got %v", summary.TotalRevenue)
	}
}

func TestServiceUserProductAndBillingLifecycle(t *testing.T) {
	store, err := db.New(filepath.Join(t.TempDir(), "payment.db"))
	if err != nil {
		t.Fatalf("db.New() error = %v", err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	svc := New(store)

	business, err := svc.CreateBusiness(models.BusinessInput{
		Name:  "Setho Pty Ltd",
		Email: "ops@setho.com",
	})
	if err != nil {
		t.Fatalf("CreateBusiness() error = %v", err)
	}

	user, err := svc.CreateUser(models.UserInput{
		BusinessID: business.ID,
		Name:       "Jane Smith",
		Email:      "jane@setho.com",
		Role:       "admin",
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if user.ID == "" {
		t.Fatal("CreateUser() returned empty ID")
	}

	product, err := svc.CreateProduct(models.ProductInput{
		BusinessID:  business.ID,
		Name:        "Pro Plan",
		Description: "Premium account",
		Price:       49.99,
		Currency:    "USD",
	})
	if err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}
	if product.ID == "" {
		t.Fatal("CreateProduct() returned empty ID")
	}

	billing, err := svc.CreateBilling(models.BillingInput{
		BusinessID: business.ID,
		UserID:     user.ID,
		ProductID:  product.ID,
		Reference:  "BILL-001",
		Amount:     49.99,
		Currency:   "USD",
		Status:     "pending",
	})
	if err != nil {
		t.Fatalf("CreateBilling() error = %v", err)
	}
	if billing.Reference != "BILL-001" {
		t.Fatalf("expected reference BILL-001, got %s", billing.Reference)
	}

	summary, err := svc.Summary()
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if summary.UserCount != 1 {
		t.Fatalf("expected 1 user, got %d", summary.UserCount)
	}
	if summary.ProductCount != 1 {
		t.Fatalf("expected 1 product, got %d", summary.ProductCount)
	}
	if summary.OutstandingBills != 49.99 {
		t.Fatalf("expected outstanding bills 49.99, got %v", summary.OutstandingBills)
	}
}
