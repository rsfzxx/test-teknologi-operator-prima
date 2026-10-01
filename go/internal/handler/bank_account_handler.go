package handler

import (
	"net/http"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/response"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/service"
)

const bankAccountsPath = "/api/v1/bank-accounts"

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

func (h *BankAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBankAccountRequest
	if err := decodeJSON(w, r, &req); err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	data := dto.NewBankAccountResponse(*result.Account)
	if result.Restored {
		response.Success(w, http.StatusOK, "Bank account restored successfully", data)
		return
	}

	w.Header().Set("Location", bankAccountsPath+"/"+result.Account.BankAccountUUID.String())
	response.Success(w, http.StatusCreated, "Bank account created successfully", data)
}

func (h *BankAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var req dto.UpdateBankAccountRequest
	if err := decodeJSON(w, r, &req); err != nil {
		response.Error(w, err)
		return
	}

	account, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Bank account updated successfully", dto.NewBankAccountResponse(*account))
}

func (h *BankAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Bank account deleted successfully", nil)
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