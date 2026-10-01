package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/repository"
)

const (
	maxPage         = 1_000_000
	maxSearchLength = 100
)

type BankListResult struct {
	Banks   []models.Bank
	Page    int
	PerPage int
	Total   int64
}

type BankService interface {
	List(ctx context.Context, q dto.ListBanksQuery) (*BankListResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Bank, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req dto.UpdateBankStatusRequest) (*models.Bank, error)
}

type bankService struct {
	repo repository.BankRepository
}

func NewBankService(repo repository.BankRepository) BankService {
	return &bankService{repo: repo}
}

func (s *bankService) List(ctx context.Context, q dto.ListBanksQuery) (*BankListResult, error) {
	filter, details := buildBankFilter(q)
	if len(details) > 0 {
		return nil, apperror.BadRequest("Invalid query parameters", details...)
	}

	banks, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &BankListResult{
		Banks:   banks,
		Page:    q.Page,
		PerPage: q.PerPage,
		Total:   total,
	}, nil
}

func (s *bankService) GetByID(ctx context.Context, id uuid.UUID) (*models.Bank, error) {
	bank, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return bank, nil
}

func (s *bankService) UpdateStatus(ctx context.Context, id uuid.UUID, req dto.UpdateBankStatusRequest) (*models.Bank, error) {
	status, ok := models.ParseBankStatus(req.Status)
	if !ok {
		return nil, apperror.Validation(apperror.FieldError{
			Field:   "status",
			Message: "must be one of: Active, Inactive",
		})
	}

	bank, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return bank, nil
}

func mapRepoError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return apperror.NotFound("Bank not found")
	}
	return apperror.Internal(err)
}

func buildBankFilter(q dto.ListBanksQuery) (repository.BankFilter, []apperror.FieldError) {
	var details []apperror.FieldError
	add := func(field, message string) {
		details = append(details, apperror.FieldError{Field: field, Message: message})
	}

	if q.Page < 1 || q.Page > maxPage {
		add("page", fmt.Sprintf("must be between 1 and %d", maxPage))
	}
	if q.PerPage < 1 || q.PerPage > dto.MaxPerPage {
		add("per_page", fmt.Sprintf("must be between 1 and %d", dto.MaxPerPage))
	}

	search := strings.TrimSpace(q.Search)
	if utf8.RuneCountInString(search) > maxSearchLength {
		add("search", fmt.Sprintf("must be at most %d characters", maxSearchLength))
	}

	var status *models.BankStatus
	if raw := strings.TrimSpace(q.Status); raw != "" {
		parsed, ok := models.ParseBankStatus(raw)
		if !ok {
			add("status", "must be one of: Active, Inactive")
		} else {
			status = &parsed
		}
	}

	sortBy := strings.ToLower(strings.TrimSpace(q.SortBy))
	if sortBy == "" {
		sortBy = repository.BankSortName
	}
	if !repository.IsValidBankSortField(sortBy) {
		add("sort_by", "must be one of: "+strings.Join(repository.BankSortFields, ", "))
	}

	desc := false
	switch strings.ToLower(strings.TrimSpace(q.Order)) {
	case "", "asc":
	case "desc":
		desc = true
	default:
		add("order", "must be one of: asc, desc")
	}

	if len(details) > 0 {
		return repository.BankFilter{}, details
	}

	return repository.BankFilter{
		Search:   search,
		Status:   status,
		Type:     strings.TrimSpace(q.Type),
		SortBy:   sortBy,
		SortDesc: desc,
		Limit:    q.PerPage,
		Offset:   (q.Page - 1) * q.PerPage,
	}, nil
}