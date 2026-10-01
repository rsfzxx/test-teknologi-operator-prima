package handler

import (
	"net/http"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/response"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/service"
)

type BankAccountHandler struct {
	service service.BankAccountService
}

func NewBankAccountHandler(s service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{service: s}
}

func (h *BankAccountHandler) List(w http.ResponseWriter, r *http.Request) {
	q, err := parseListBankAccountsQuery(r)
	if err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.service.List(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Paginated(w, "Bank accounts retrieved successfully",
		dto.NewBankAccountResponses(result.Accounts),
		response.NewMeta(result.Page, result.PerPage, result.Total))
}

func (h *BankAccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	account, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Bank account retrieved successfully", dto.NewBankAccountResponse(*account))
}

func parseListBankAccountsQuery(r *http.Request) (dto.ListBankAccountsQuery, error) {
	v := r.URL.Query()

	page, perPage, details := parsePagination(v)
	if len(details) > 0 {
		return dto.ListBankAccountsQuery{}, apperror.BadRequest("Invalid query parameters", details...)
	}

	return dto.ListBankAccountsQuery{
		Page:     page,
		PerPage:  perPage,
		Search:   v.Get("search"),
		Status:   v.Get("status"),
		BankCode: v.Get("bank_code"),
		SortBy:   v.Get("sort_by"),
		Order:    v.Get("order"),
	}, nil
}