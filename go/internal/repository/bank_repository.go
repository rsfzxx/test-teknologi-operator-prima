package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
)

var ErrNotFound = errors.New("record not found")

const (
	BankSortName      = "name"
	BankSortCode      = "code"
	BankSortType      = "type"
	BankSortStatus    = "status"
	BankSortCreatedAt = "created_at"
	BankSortUpdatedAt = "updated_at"
)

var BankSortFields = []string{
	BankSortName, BankSortCode, BankSortType,
	BankSortStatus, BankSortCreatedAt, BankSortUpdatedAt,
}

func IsValidBankSortField(field string) bool {
	for _, f := range BankSortFields {
		if f == field {
			return true
		}
	}
	return false
}

type BankFilter struct {
	Search   string
	Status   *models.BankStatus
	Type     string
	SortBy   string
	SortDesc bool
	Limit    int
	Offset   int
}

type BankRepository interface {
	FindAll(ctx context.Context, filter BankFilter) ([]models.Bank, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Bank, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BankStatus) (*models.Bank, error)
}