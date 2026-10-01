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

type CreateBankAccountParams struct {
	AccountNumber string
	BankName      string
	AccountName   string
	BankCode      string
}

type UpdateBankAccountParams struct {
	AccountNumber string
	BankName      string
	AccountName   string
	BankCode      string
	Status        models.BankAccountStatus
	Reason        *string
}

type BankAccountRepository interface {
	FindAll(ctx context.Context, filter BankAccountFilter) ([]models.BankAccount, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
	FindByIDIncludingDeleted(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
	FindByNumberAndBank(ctx context.Context, accountNumber, bankName string) (*models.BankAccount, error)
	Create(ctx context.Context, p CreateBankAccountParams) (*models.BankAccount, error)
	Update(ctx context.Context, id uuid.UUID, p UpdateBankAccountParams) (*models.BankAccount, error)
	Restore(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}