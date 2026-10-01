package handler

import (
	"net/http"
	"strconv"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/response"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/service"
)

type BankHandler struct {
	service service.BankService
}

func NewBankHandler(s service.BankService) *BankHandler {
	return &BankHandler{service: s}
}

func (h *BankHandler) List(w http.ResponseWriter, r *http.Request) {
	q, err := parseListBanksQuery(r)
	if err != nil {
		response.Error(w, err)
		return
	}

	result, err := h.service.List(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}

	perPage := int64(result.PerPage)
	response.Paginated(w, "Banks retrieved successfully", dto.NewBankResponses(result.Banks), response.Meta{
		Page:       result.Page,
		PerPage:    result.PerPage,
		Total:      result.Total,
		TotalPages: int((result.Total + perPage - 1) / perPage),
	})
}

func (h *BankHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	bank, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Bank retrieved successfully", dto.NewBankResponse(*bank))
}

func (h *BankHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var req dto.UpdateBankStatusRequest
	if err := decodeJSON(w, r, &req); err != nil {
		response.Error(w, err)
		return
	}

	bank, err := h.service.UpdateStatus(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Bank status updated successfully", dto.NewBankResponse(*bank))
}

func parseListBanksQuery(r *http.Request) (dto.ListBanksQuery, error) {
	v := r.URL.Query()

	q := dto.ListBanksQuery{
		Page:    dto.DefaultPage,
		PerPage: dto.DefaultPerPage,
		Search:  v.Get("search"),
		Status:  v.Get("status"),
		Type:    v.Get("type"),
		SortBy:  v.Get("sort_by"),
		Order:   v.Get("order"),
	}

	var details []apperror.FieldError
	parseInt := func(key string, dst *int) {
		raw := v.Get(key)
		if raw == "" {
			return
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			details = append(details, apperror.FieldError{Field: key, Message: "must be a number"})
			return
		}
		*dst = n
	}
	parseInt("page", &q.Page)
	parseInt("per_page", &q.PerPage)

	if len(details) > 0 {
		return q, apperror.BadRequest("Invalid query parameters", details...)
	}
	return q, nil
}