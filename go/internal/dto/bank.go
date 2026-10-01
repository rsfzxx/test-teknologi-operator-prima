package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type ListBanksQuery struct {
	Page    int
	PerPage int
	Search  string
	Status  string
	Type    string
	SortBy  string
	Order   string
}

type UpdateBankStatusRequest struct {
	Status string `json:"status"`
}

type BankResponse struct {
	BankUUID  uuid.UUID         `json:"bank_uuid"`
	Name      string            `json:"name"`
	Code      string            `json:"code"`
	Type      string            `json:"type"`
	Status    models.BankStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewBankResponse(b models.Bank) BankResponse {
	return BankResponse{
		BankUUID:  b.BankUUID,
		Name:      b.Name,
		Code:      b.Code,
		Type:      b.Type,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func NewBankResponses(banks []models.Bank) []BankResponse {
	out := make([]BankResponse, 0, len(banks))
	for _, b := range banks {
		out = append(out, NewBankResponse(b))
	}
	return out
}