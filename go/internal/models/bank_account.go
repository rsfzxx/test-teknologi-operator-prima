package models

import (
	"time"

	"github.com/google/uuid"
)

type BankAccountStatus string

const (
	BankAccountAccepted BankAccountStatus = "Accepted"
	BankAccountReview   BankAccountStatus = "Review"
	BankAccountRejected BankAccountStatus = "Rejected"
)

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