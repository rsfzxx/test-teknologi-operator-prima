package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type BankAccountStatus string

const (
	BankAccountAccepted BankAccountStatus = "Accepted"
	BankAccountReview   BankAccountStatus = "Review"
	BankAccountRejected BankAccountStatus = "Rejected"
)

func ParseBankAccountStatus(s string) (BankAccountStatus, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "accepted":
		return BankAccountAccepted, true
	case "review":
		return BankAccountReview, true
	case "rejected":
		return BankAccountRejected, true
	default:
		return "", false
	}
}

type BankAccount struct {
	BankAccountUUID uuid.UUID         `json:"bank_account_uuid"`
	AccountNumber   string            `json:"account_number"`
	BankName        string            `json:"bank_name"`
	AccountName     string            `json:"account_name"`
	BankCode        string            `json:"bank_code"`
	Status          BankAccountStatus `json:"status"`
	Reason          *string           `json:"reason"`
	CreatedAt       *time.Time        `json:"created_at"`
	UpdatedAt       *time.Time        `json:"updated_at"`
	DeletedAt       *time.Time        `json:"deleted_at,omitempty"`
}