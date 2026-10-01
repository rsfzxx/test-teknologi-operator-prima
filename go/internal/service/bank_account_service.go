package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/repository"
)

const (
	maxBankCodeLength      = 20
	maxAccountNumberLength = 30
	maxAccountNameLength   = 150
	maxReasonLength        = 500

	msgDuplicateAccount = "Bank account with the same account number and bank already exists"
	msgDeletedAccount   = "There is an issue with your data; please contact our customer service to resolve the problem."
)

var accountNamePattern = regexp.MustCompile(`^[A-Za-z .,&'-]+$`)

type BankAccountListResult struct {
	Accounts []models.BankAccount
	Page     int
	PerPage  int
	Total    int64
}

type CreateBankAccountResult struct {
	Account  *models.BankAccount
	Restored bool
}

type BankAccountService interface {
	List(ctx context.Context, q dto.ListBankAccountsQuery) (*BankAccountListResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)
	Create(ctx context.Context, req dto.CreateBankAccountRequest) (*CreateBankAccountResult, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateBankAccountRequest) (*models.BankAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type bankAccountService struct {
	repo     repository.BankAccountRepository
	bankRepo repository.BankRepository
}

func NewBankAccountService(repo repository.BankAccountRepository, bankRepo repository.BankRepository) BankAccountService {
	return &bankAccountService{repo: repo, bankRepo: bankRepo}
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
		return nil, mapBankAccountError(err)
	}
	return account, nil
}

func (s *bankAccountService) Create(ctx context.Context, req dto.CreateBankAccountRequest) (*CreateBankAccountResult, error) {
	var errs fieldErrors
	in := validateAccountFields(&errs, req.Bank, string(req.AccountNumber), req.AccountName)
	if len(errs) > 0 {
		return nil, apperror.Validation(errs...)
	}

	bank, err := s.resolveBank(ctx, in.BankID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByNumberAndBank(ctx, in.AccountNumber, bank.Name)
	switch {
	case err == nil:
		if existing.DeletedAt == nil {
			return nil, apperror.Conflict(msgDuplicateAccount)
		}
		restored, err := s.repo.Restore(ctx, existing.BankAccountUUID)
		if err != nil {
			return nil, mapBankAccountError(err)
		}
		return &CreateBankAccountResult{Account: restored, Restored: true}, nil
	case !errors.Is(err, repository.ErrNotFound):
		return nil, apperror.Internal(err)
	}

	created, err := s.repo.Create(ctx, repository.CreateBankAccountParams{
		AccountNumber: in.AccountNumber,
		BankName:      bank.Name,
		AccountName:   in.AccountName,
		BankCode:      bank.Code,
	})
	if err != nil {
		return nil, mapBankAccountError(err)
	}
	return &CreateBankAccountResult{Account: created}, nil
}

func (s *bankAccountService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateBankAccountRequest) (*models.BankAccount, error) {
	var errs fieldErrors
	in := validateAccountFields(&errs, req.Bank, string(req.AccountNumber), req.AccountName)
	status, reason := validateStatusAndReason(&errs, req.Status, req.Reason)
	if len(errs) > 0 {
		return nil, apperror.Validation(errs...)
	}

	existing, err := s.repo.FindByIDIncludingDeleted(ctx, id)
	if err != nil {
		return nil, mapBankAccountError(err)
	}
	if existing.DeletedAt != nil {
		return nil, apperror.New(http.StatusConflict, msgDeletedAccount)
	}

	bank, err := s.resolveBank(ctx, in.BankID)
	if err != nil {
		return nil, err
	}

	other, err := s.repo.FindByNumberAndBank(ctx, in.AccountNumber, bank.Name)
	switch {
	case err == nil:
		if other.BankAccountUUID != id {
			return nil, apperror.Conflict(msgDuplicateAccount)
		}
	case !errors.Is(err, repository.ErrNotFound):
		return nil, apperror.Internal(err)
	}

	updated, err := s.repo.Update(ctx, id, repository.UpdateBankAccountParams{
		AccountNumber: in.AccountNumber,
		BankName:      bank.Name,
		AccountName:   in.AccountName,
		BankCode:      bank.Code,
		Status:        status,
		Reason:        reason,
	})
	if err != nil {
		return nil, mapBankAccountError(err)
	}
	return updated, nil
}

func (s *bankAccountService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return mapBankAccountError(err)
	}
	return nil
}

func (s *bankAccountService) resolveBank(ctx context.Context, id uuid.UUID) (*models.Bank, error) {
	bank, err := s.bankRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Validation(apperror.FieldError{Field: "bank", Message: "bank not found"})
		}
		return nil, apperror.Internal(err)
	}
	return bank, nil
}

type fieldErrors []apperror.FieldError

func (f *fieldErrors) add(field, message string) {
	*f = append(*f, apperror.FieldError{Field: field, Message: message})
}

type accountInput struct {
	BankID        uuid.UUID
	AccountNumber string
	AccountName   string
}

func validateAccountFields(errs *fieldErrors, bank, accountNumber, accountName string) accountInput {
	var in accountInput

	if strings.TrimSpace(bank) == "" {
		errs.add("bank", "is required")
	} else if id, err := uuid.Parse(strings.TrimSpace(bank)); err != nil {
		errs.add("bank", "must be a valid UUID")
	} else {
		in.BankID = id
	}

	switch {
	case accountNumber == "":
		errs.add("account_number", "is required")
	case !isDigits(accountNumber):
		errs.add("account_number", "must contain digits only")
	case len(accountNumber) > maxAccountNumberLength:
		errs.add("account_number", fmt.Sprintf("must be at most %d digits", maxAccountNumberLength))
	default:
		in.AccountNumber = accountNumber
	}

	name := strings.TrimSpace(accountName)
	switch {
	case name == "":
		errs.add("account_name", "is required")
	case utf8.RuneCountInString(name) > maxAccountNameLength:
		errs.add("account_name", fmt.Sprintf("must be at most %d characters", maxAccountNameLength))
	case !accountNamePattern.MatchString(name):
		errs.add("account_name", "may only contain letters, spaces, and . , & ' -")
	default:
		in.AccountName = name
	}

	return in
}

func validateStatusAndReason(errs *fieldErrors, rawStatus string, rawReason *string) (models.BankAccountStatus, *string) {
	if strings.TrimSpace(rawStatus) == "" {
		errs.add("status", "is required")
		return "", nil
	}
	status, ok := models.ParseBankAccountStatus(rawStatus)
	if !ok {
		errs.add("status", "must be one of: Accepted, Review, Rejected")
		return "", nil
	}

	reason := ""
	if rawReason != nil {
		reason = strings.TrimSpace(*rawReason)
	}

	if status == models.BankAccountRejected {
		switch {
		case reason == "":
			errs.add("reason", "is required when status is Rejected")
		case utf8.RuneCountInString(reason) > maxReasonLength:
			errs.add("reason", fmt.Sprintf("must be at most %d characters", maxReasonLength))
		default:
			return status, &reason
		}
		return status, nil
	}

	if reason != "" {
		errs.add("reason", "must be empty unless status is Rejected")
	}
	return status, nil
}

func mapBankAccountError(err error) error {
	var invalid *repository.InvalidDataError

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return apperror.NotFound("Bank account not found")
	case errors.Is(err, repository.ErrConflict):
		return apperror.Conflict(msgDuplicateAccount)
	case errors.As(err, &invalid):
		return apperror.Validation(invalidDataField(invalid.Constraint))
	default:
		return apperror.Internal(err)
	}
}

func invalidDataField(constraint string) apperror.FieldError {
	switch constraint {
	case "chk_ba_account_number":
		return apperror.FieldError{Field: "account_number", Message: "must contain digits only"}
	case "chk_ba_account_name":
		return apperror.FieldError{Field: "account_name", Message: "contains characters that are not allowed"}
	case "chk_ba_bank_name", "chk_ba_bank_code":
		return apperror.FieldError{Field: "bank", Message: "selected bank has data that is not accepted for bank accounts"}
	default:
		return apperror.FieldError{Field: "data", Message: "violates a database rule"}
	}
}

func buildBankAccountFilter(q dto.ListBankAccountsQuery) (repository.BankAccountFilter, []apperror.FieldError) {
	var errs fieldErrors

	if q.Page < 1 || q.Page > maxPage {
		errs.add("page", fmt.Sprintf("must be between 1 and %d", maxPage))
	}
	if q.PerPage < 1 || q.PerPage > dto.MaxPerPage {
		errs.add("per_page", fmt.Sprintf("must be between 1 and %d", dto.MaxPerPage))
	}

	search := strings.TrimSpace(q.Search)
	if utf8.RuneCountInString(search) > maxSearchLength {
		errs.add("search", fmt.Sprintf("must be at most %d characters", maxSearchLength))
	}

	var status *models.BankAccountStatus
	if raw := strings.TrimSpace(q.Status); raw != "" {
		parsed, ok := models.ParseBankAccountStatus(raw)
		if !ok {
			errs.add("status", "must be one of: Accepted, Review, Rejected")
		} else {
			status = &parsed
		}
	}

	bankCode := strings.TrimSpace(q.BankCode)
	if bankCode != "" && (len(bankCode) > maxBankCodeLength || !isDigits(bankCode)) {
		errs.add("bank_code", "must contain digits only")
	}

	sortBy := strings.ToLower(strings.TrimSpace(q.SortBy))
	usingDefaultSort := sortBy == ""
	if usingDefaultSort {
		sortBy = repository.BankAccountSortCreatedAt
	}
	if !repository.IsValidBankAccountSortField(sortBy) {
		errs.add("sort_by", "must be one of: "+strings.Join(repository.BankAccountSortFields, ", "))
	}

	desc := usingDefaultSort
	switch strings.ToLower(strings.TrimSpace(q.Order)) {
	case "":
	case "asc":
		desc = false
	case "desc":
		desc = true
	default:
		errs.add("order", "must be one of: asc, desc")
	}

	if len(errs) > 0 {
		return repository.BankAccountFilter{}, errs
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