package dto

import (
	"bytes"
	"encoding/json"
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

type NumericString string

func (n *NumericString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch {
	case len(b) == 0 || string(b) == "null":
		*n = ""
	case b[0] == '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*n = NumericString(s)
	default:
		*n = NumericString(b)
	}
	return nil
}

type CreateBankAccountRequest struct {
	Bank          string        `json:"bank"`
	AccountNumber NumericString `json:"account_number"`
	AccountName   string        `json:"account_name"`
}

type UpdateBankAccountRequest struct {
	Bank          string        `json:"bank"`
	AccountNumber NumericString `json:"account_number"`
	AccountName   string        `json:"account_name"`
	Status        string        `json:"status"`
	Reason        *string       `json:"reason"`
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