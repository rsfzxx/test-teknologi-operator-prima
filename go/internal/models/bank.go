package models

import (
	"time"

	"github.com/google/uuid"
)

type BankStatus string

const (
	BankStatusActive   BankStatus = "Active"
	BankStatusInactive BankStatus = "Inactive"
)

type Bank struct {
	BankUUID  uuid.UUID  `json:"bank_uuid"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Status    BankStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}