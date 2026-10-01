package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
)

const (
	BankAccountSortAccountNumber = "account_number"
	BankAccountSortAccountName   = "account_name"
	BankAccountSortBankName      = "bank_name"
	BankAccountSortBankCode      = "bank_code"
	BankAccountSortStatus        = "status"
	BankAccountSortCreatedAt     = "created_at"
	BankAccountSortUpdatedAt     = "updated_at"
)

var BankAccountSortFields = []string{
	BankAccountSortAccountNumber, BankAccountSortAccountName, BankAccountSortBankName,
	BankAccountSortBankCode, BankAccountSortStatus, BankAccountSortCreatedAt,
	BankAccountSortUpdatedAt,
}

func IsValidBankAccountSortField(field string) bool {
	for _, f := range BankAccountSortFields {
		if f == field {
			return true
		}
	}
	return false
}

type BankAccountFilter struct {
	Search   string
	Status   *models.BankAccountStatus
	BankCode string
	SortBy   string
	SortDesc bool
	Limit    int
	Offset   int
}

type BankAccountRepository interface {
	FindAll(ctx context.Context, filter BankAccountFilter) ([]models.BankAccount, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
}