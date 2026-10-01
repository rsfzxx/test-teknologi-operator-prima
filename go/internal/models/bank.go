package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type BankStatus string

const (
	BankStatusActive   BankStatus = "Active"
	BankStatusInactive BankStatus = "Inactive"
)

func ParseBankStatus(s string) (BankStatus, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "active":
		return BankStatusActive, true
	case "inactive":
		return BankStatusInactive, true
	default:
		return "", false
	}
}

type Bank struct {
	BankUUID  uuid.UUID  `json:"bank_uuid"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Status    BankStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}