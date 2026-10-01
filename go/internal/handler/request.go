package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/dto"
)

const maxBodyBytes = 1 << 20

func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		return uuid.Nil, apperror.BadRequest("Invalid path parameter", apperror.FieldError{
			Field:   name,
			Message: "must be a valid UUID",
		})
	}
	return id, nil
}

func parsePagination(v url.Values) (page, perPage int, details []apperror.FieldError) {
	page, perPage = dto.DefaultPage, dto.DefaultPerPage

	read := func(key string, dst *int) {
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
	read("page", &page)
	read("per_page", &perPage)
	return page, perPage, details
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var (
			syntaxErr   *json.SyntaxError
			typeErr     *json.UnmarshalTypeError
			maxBytesErr *http.MaxBytesError
		)
		switch {
		case errors.Is(err, io.EOF):
			return apperror.BadRequest("Request body is required")
		case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
			return apperror.BadRequest("Request body contains malformed JSON")
		case errors.As(err, &typeErr):
			return apperror.BadRequest("Invalid request body", apperror.FieldError{
				Field:   typeErr.Field,
				Message: fmt.Sprintf("must be of type %s", typeErr.Type),
			})
		case errors.As(err, &maxBytesErr):
			return apperror.New(http.StatusRequestEntityTooLarge, "Request body too large")
		default:
			return apperror.BadRequest("Invalid request body: " + err.Error())
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return apperror.BadRequest("Request body must contain a single JSON object")
	}
	return nil
}