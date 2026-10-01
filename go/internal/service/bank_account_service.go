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

const maxBankCodeLength = 20

type BankAccountListResult struct {
	Accounts []models.BankAccount
	Page     int
	PerPage  int
	Total    int64
}

type BankAccountService interface {
	List(ctx context.Context, q dto.ListBankAccountsQuery) (*BankAccountListResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
}

type bankAccountService struct {
	repo repository.BankAccountRepository
}

func NewBankAccountService(repo repository.BankAccountRepository) BankAccountService {
	return &bankAccountService{repo: repo}
}

func (s *bankAccountService) List(ctx context.Context, q dto.ListBankAccountsQuery) (*BankAccountListResult, error) {
	filter, details := buildBankAccountFilter(q)
	if len(details) > 0 {
		return nil, apperror.BadRequest("Invalid query parameters", details...)
	}

	accounts, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &BankAccountListResult{
		Accounts: accounts,
		Page:     q.Page,
		PerPage:  q.PerPage,
		Total:    total,
	}, nil
}

func (s *bankAccountService) GetByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error) {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NotFound("Bank account not found")
		}
		return nil, apperror.Internal(err)
	}
	return account, nil
}

func buildBankAccountFilter(q dto.ListBankAccountsQuery) (repository.BankAccountFilter, []apperror.FieldError) {
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

	var status *models.BankAccountStatus
	if raw := strings.TrimSpace(q.Status); raw != "" {
		parsed, ok := models.ParseBankAccountStatus(raw)
		if !ok {
			add("status", "must be one of: Accepted, Review, Rejected")
		} else {
			status = &parsed
		}
	}

	bankCode := strings.TrimSpace(q.BankCode)
	if bankCode != "" {
		if len(bankCode) > maxBankCodeLength || !isDigits(bankCode) {
			add("bank_code", "must contain digits only")
		}
	}

	sortBy := strings.ToLower(strings.TrimSpace(q.SortBy))
	usingDefaultSort := sortBy == ""
	if usingDefaultSort {
		sortBy = repository.BankAccountSortCreatedAt
	}
	if !repository.IsValidBankAccountSortField(sortBy) {
		add("sort_by", "must be one of: "+strings.Join(repository.BankAccountSortFields, ", "))
	}

	desc := usingDefaultSort
	switch strings.ToLower(strings.TrimSpace(q.Order)) {
	case "":
	case "asc":
		desc = false
	case "desc":
		desc = true
	default:
		add("order", "must be one of: asc, desc")
	}

	if len(details) > 0 {
		return repository.BankAccountFilter{}, details
	}

	return repository.BankAccountFilter{
		Search:   search,
		Status:   status,
		BankCode: bankCode,
		SortBy:   sortBy,
		SortDesc: desc,
		Limit:    q.PerPage,
		Offset:   (q.Page - 1) * q.PerPage,
	}, nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return s != ""
}