package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
)

type ListBankAccountsQuery struct {
	Page     int
	PerPage  int
	Search   string
	Status   string
	BankCode string
	SortBy   string
	Order    string
}

type BankAccountResponse struct {
	BankAccountUUID uuid.UUID                `json:"bank_account_uuid"`
	AccountNumber   string                   `json:"account_number"`
	BankName        string                   `json:"bank_name"`
	AccountName     string                   `json:"account_name"`
	BankCode        string                   `json:"bank_code"`
	Status          models.BankAccountStatus `json:"status"`
	Reason          *string                  `json:"reason"`
	CreatedAt       *time.Time               `json:"created_at"`
	UpdatedAt       *time.Time               `json:"updated_at"`
}

func NewBankAccountResponse(a models.BankAccount) BankAccountResponse {
	return BankAccountResponse{
		BankAccountUUID: a.BankAccountUUID,
		AccountNumber:   a.AccountNumber,
		BankName:        a.BankName,
		AccountName:     a.AccountName,
		BankCode:        a.BankCode,
		Status:          a.Status,
		Reason:          a.Reason,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func NewBankAccountResponses(accounts []models.BankAccount) []BankAccountResponse {
	out := make([]BankAccountResponse, 0, len(accounts))
	for _, a := range accounts {
		out = append(out, NewBankAccountResponse(a))
	}
	return out
}