package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/repository"
)

const bankColumns = `bank_uuid, name, code, type, status, created_at, updated_at`

type BankRepository struct {
	db *pgxpool.Pool
}

var _ repository.BankRepository = (*BankRepository)(nil)

func NewBankRepository(db *pgxpool.Pool) *BankRepository {
	return &BankRepository{db: db}
}

func (r *BankRepository) FindAll(ctx context.Context, f repository.BankFilter) ([]models.Bank, int64, error) {
	where, args := buildBankWhere(f)

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM banks"+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count banks: %w", err)
	}
	if total == 0 {
		return []models.Bank{}, 0, nil
	}

	direction := "ASC"
	if f.SortDesc {
		direction = "DESC"
	}

	query := fmt.Sprintf(
		"SELECT %s FROM banks%s ORDER BY %s %s, bank_uuid ASC LIMIT $%d OFFSET $%d",
		bankColumns, where, sortColumn(f.SortBy), direction, len(args)+1, len(args)+2,
	)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query banks: %w", err)
	}
	defer rows.Close()

	banks := make([]models.Bank, 0, f.Limit)
	for rows.Next() {
		b, err := scanBank(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan bank: %w", err)
		}
		banks = append(banks, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate banks: %w", err)
	}
	return banks, total, nil
}

func (r *BankRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Bank, error) {
	row := r.db.QueryRow(ctx,
		"SELECT "+bankColumns+" FROM banks WHERE bank_uuid = $1", id.String())

	b, err := scanBank(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find bank by id: %w", err)
	}
	return &b, nil
}

func (r *BankRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BankStatus) (*models.Bank, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE banks
		    SET status = $1, updated_at = NOW()
		  WHERE bank_uuid = $2
		RETURNING `+bankColumns,
		string(status), id.String())

	b, err := scanBank(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("update bank status: %w", err)
	}
	return &b, nil
}

func scanBank(row pgx.Row) (models.Bank, error) {
	var (
		b      models.Bank
		id     string
		status string
	)
	if err := row.Scan(&id, &b.Name, &b.Code, &b.Type, &status, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return models.Bank{}, err
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return models.Bank{}, fmt.Errorf("parse bank uuid: %w", err)
	}
	b.BankUUID = parsed
	b.Status = models.BankStatus(status)
	return b, nil
}

func buildBankWhere(f repository.BankFilter) (string, []any) {
	var (
		conds []string
		args  []any
	)

	if f.Search != "" {
		args = append(args, "%"+escapeLike(f.Search)+"%")
		n := len(args)
		conds = append(conds, fmt.Sprintf("(name ILIKE $%d OR code ILIKE $%d)", n, n))
	}
	if f.Status != nil {
		args = append(args, string(*f.Status))
		conds = append(conds, fmt.Sprintf("status = $%d", len(args)))
	}
	if f.Type != "" {
		args = append(args, f.Type)
		conds = append(conds, fmt.Sprintf("LOWER(type) = LOWER($%d)", len(args)))
	}

	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func sortColumn(field string) string {
	switch field {
	case repository.BankSortCode:
		return "code"
	case repository.BankSortType:
		return "type"
	case repository.BankSortStatus:
		return "status"
	case repository.BankSortCreatedAt:
		return "created_at"
	case repository.BankSortUpdatedAt:
		return "updated_at"
	default:
		return "name"
	}
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func escapeLike(s string) string { return likeEscaper.Replace(s) }