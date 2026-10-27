package service

import (
	"fmt"
	"strings"

	"setho-pymnt/internal/models"
)

type Repository interface {
	ListBusinesses() ([]models.Business, error)
	CreateBusiness(models.BusinessInput) (models.Business, error)
	GetBusiness(id string) (models.Business, error)
	ListTransactions(businessID string) ([]models.Transaction, error)
	CreateTransaction(models.TransactionInput) (models.Transaction, error)
	GetTransaction(id string) (models.Transaction, error)
	ListUsers(businessID string) ([]models.User, error)
	CreateUser(models.UserInput) (models.User, error)
	GetUser(id string) (models.User, error)
	ListProducts(businessID string) ([]models.Product, error)
	CreateProduct(models.ProductInput) (models.Product, error)
	GetProduct(id string) (models.Product, error)
	ListBillings(businessID string) ([]models.Billing, error)
	CreateBilling(models.BillingInput) (models.Billing, error)
	GetBilling(id string) (models.Billing, error)
	Summary() (models.Summary, error)
}

type Service struct {
	Repo Repository
}

func New(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) ListBusinesses() ([]models.Business, error) {
	return s.Repo.ListBusinesses()
}

func (s *Service) CreateBusiness(input models.BusinessInput) (models.Business, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		return models.Business{}, fmt.Errorf("name is required")
	}
	if input.Email == "" {
		return models.Business{}, fmt.Errorf("email is required")
	}

	return s.Repo.CreateBusiness(input)
}

func (s *Service) GetBusiness(id string) (models.Business, error) {
	if strings.TrimSpace(id) == "" {
		return models.Business{}, fmt.Errorf("business id is required")
	}
	return s.Repo.GetBusiness(id)
}

func (s *Service) ListTransactions(businessID string) ([]models.Transaction, error) {
	return s.Repo.ListTransactions(businessID)
}

func (s *Service) CreateTransaction(input models.TransactionInput) (models.Transaction, error) {
	input.BusinessID = strings.TrimSpace(input.BusinessID)
	input.Customer = strings.TrimSpace(input.Customer)
	input.Reference = strings.TrimSpace(input.Reference)
	input.Currency = strings.TrimSpace(input.Currency)
	input.Notes = strings.TrimSpace(input.Notes)

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
	if input.Currency == "" {
		input.Currency = "USD"
	}
	if input.Status == "" {
		input.Status = "pending"
	}

	return s.Repo.CreateTransaction(input)
}

func (s *Service) GetTransaction(id string) (models.Transaction, error) {
	if strings.TrimSpace(id) == "" {
		return models.Transaction{}, fmt.Errorf("transaction id is required")
	}
	return s.Repo.GetTransaction(id)
}

func (s *Service) ListUsers(businessID string) ([]models.User, error) {
	return s.Repo.ListUsers(businessID)
}

func (s *Service) CreateUser(input models.UserInput) (models.User, error) {
	input.BusinessID = strings.TrimSpace(input.BusinessID)
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Role = strings.TrimSpace(input.Role)

	if input.BusinessID == "" {
		return models.User{}, fmt.Errorf("business_id is required")
	}
	if input.Name == "" {
		return models.User{}, fmt.Errorf("name is required")
	}
	if input.Email == "" {
		return models.User{}, fmt.Errorf("email is required")
	}
	if input.Role == "" {
		input.Role = "staff"
	}

	return s.Repo.CreateUser(input)
}

func (s *Service) GetUser(id string) (models.User, error) {
	if strings.TrimSpace(id) == "" {
		return models.User{}, fmt.Errorf("user id is required")
	}
	return s.Repo.GetUser(id)
}

func (s *Service) ListProducts(businessID string) ([]models.Product, error) {
	return s.Repo.ListProducts(businessID)
}

func (s *Service) CreateProduct(input models.ProductInput) (models.Product, error) {
	input.BusinessID = strings.TrimSpace(input.BusinessID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Currency = strings.TrimSpace(input.Currency)

	if input.BusinessID == "" {
		return models.Product{}, fmt.Errorf("business_id is required")
	}
	if input.Name == "" {
		return models.Product{}, fmt.Errorf("name is required")
	}
	if input.Price <= 0 {
		return models.Product{}, fmt.Errorf("price must be greater than zero")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}

	return s.Repo.CreateProduct(input)
}

func (s *Service) GetProduct(id string) (models.Product, error) {
	if strings.TrimSpace(id) == "" {
		return models.Product{}, fmt.Errorf("product id is required")
	}
	return s.Repo.GetProduct(id)
}

func (s *Service) ListBillings(businessID string) ([]models.Billing, error) {
	return s.Repo.ListBillings(businessID)
}

func (s *Service) CreateBilling(input models.BillingInput) (models.Billing, error) {
	input.BusinessID = strings.TrimSpace(input.BusinessID)
	input.UserID = strings.TrimSpace(input.UserID)
	input.ProductID = strings.TrimSpace(input.ProductID)
	input.Reference = strings.TrimSpace(input.Reference)
	input.Currency = strings.TrimSpace(input.Currency)
	input.Notes = strings.TrimSpace(input.Notes)

	if input.BusinessID == "" {
		return models.Billing{}, fmt.Errorf("business_id is required")
	}
	if input.Reference == "" {
		return models.Billing{}, fmt.Errorf("reference is required")
	}
	if input.Amount <= 0 {
		return models.Billing{}, fmt.Errorf("amount must be greater than zero")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}
	if input.Status == "" {
		input.Status = "pending"
	}

	return s.Repo.CreateBilling(input)
}

func (s *Service) GetBilling(id string) (models.Billing, error) {
	if strings.TrimSpace(id) == "" {
		return models.Billing{}, fmt.Errorf("billing id is required")
	}
	return s.Repo.GetBilling(id)
}

func (s *Service) Summary() (models.Summary, error) {
	return s.Repo.Summary()
}
