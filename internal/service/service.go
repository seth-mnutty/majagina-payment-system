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

func (s *Service) Summary() (models.Summary, error) {
	return s.Repo.Summary()
}
