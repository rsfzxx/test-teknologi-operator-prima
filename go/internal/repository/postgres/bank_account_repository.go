package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/models"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/repository"
)

const bankAccountColumns = `bank_account_uuid, account_number, bank_name, account_name, bank_code, status, reason, created_at, updated_at, deleted_at`

const (
	pgUniqueViolation = "23505"
	pgCheckViolation  = "23514"
)

type BankAccountRepository struct {
	db *pgxpool.Pool
}

var _ repository.BankAccountRepository = (*BankAccountRepository)(nil)

func NewBankAccountRepository(db *pgxpool.Pool) *BankAccountRepository {
	return &BankAccountRepository{db: db}
}

func (r *BankAccountRepository) FindAll(ctx context.Context, f repository.BankAccountFilter) ([]models.BankAccount, int64, error) {
	where, args := buildBankAccountWhere(f)

	var total int64
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM bank_accounts"+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bank accounts: %w", err)
	}
	if total == 0 {
		return []models.BankAccount{}, 0, nil
	}

	direction := "ASC"
	if f.SortDesc {
		direction = "DESC"
	}

	query := fmt.Sprintf(
		"SELECT %s FROM bank_accounts%s ORDER BY %s %s NULLS LAST, bank_account_uuid ASC LIMIT $%d OFFSET $%d",
		bankAccountColumns, where, bankAccountSortColumn(f.SortBy), direction, len(args)+1, len(args)+2,
	)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query bank accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]models.BankAccount, 0, f.Limit)
	for rows.Next() {
		a, err := scanBankAccount(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan bank account: %w", err)
		}
		accounts = append(accounts, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate bank accounts: %w", err)
	}
	return accounts, total, nil
}

func (r *BankAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.BankAccount, error) {
	return r.queryOne(ctx, "find bank account by id",
		"SELECT "+bankAccountColumns+" FROM bank_accounts WHERE bank_account_uuid = $1 AND deleted_at IS NULL",
		id.String())
}

func (r *BankAccountRepository) FindByIDIncludingDeleted(ctx context.Context, id uuid.UUID) (*models.BankAccount, error) {
	return r.queryOne(ctx, "find bank account by id (including deleted)",
		"SELECT "+bankAccountColumns+" FROM bank_accounts WHERE bank_account_uuid = $1",
		id.String())
}

func (r *BankAccountRepository) FindByNumberAndBank(ctx context.Context, accountNumber, bankName string) (*models.BankAccount, error) {
	return r.queryOne(ctx, "find bank account by number and bank",
		"SELECT "+bankAccountColumns+" FROM bank_accounts WHERE account_number = $1 AND bank_name = $2 LIMIT 1",
		accountNumber, bankName)
}

func (r *BankAccountRepository) Create(ctx context.Context, p repository.CreateBankAccountParams) (*models.BankAccount, error) {
	return r.queryOne(ctx, "create bank account",
		`INSERT INTO bank_accounts
		        (account_number, bank_name, account_name, bank_code, status, reason, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NULL, NOW(), NOW())
		 RETURNING `+bankAccountColumns,
		p.AccountNumber, p.BankName, p.AccountName, p.BankCode, string(models.BankAccountReview))
}

func (r *BankAccountRepository) Update(ctx context.Context, id uuid.UUID, p repository.UpdateBankAccountParams) (*models.BankAccount, error) {
	var reason any
	if p.Reason != nil {
		reason = *p.Reason
	}

	return r.queryOne(ctx, "update bank account",
		`UPDATE bank_accounts
		    SET account_number = $1,
		        bank_name      = $2,
		        account_name   = $3,
		        bank_code      = $4,
		        status         = $5,
		        reason         = $6,
		        updated_at     = NOW()
		  WHERE bank_account_uuid = $7
		    AND deleted_at IS NULL
		RETURNING `+bankAccountColumns,
		p.AccountNumber, p.BankName, p.AccountName, p.BankCode,
		string(p.Status), reason, id.String())
}

func (r *BankAccountRepository) Restore(ctx context.Context, id uuid.UUID) (*models.BankAccount, error) {
	return r.queryOne(ctx, "restore bank account",
		`UPDATE bank_accounts
		    SET deleted_at = NULL
		  WHERE bank_account_uuid = $1
		RETURNING `+bankAccountColumns,
		id.String())
}

func (r *BankAccountRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE bank_accounts
		    SET deleted_at = NOW()
		  WHERE bank_account_uuid = $1
		    AND deleted_at IS NULL`,
		id.String())
	if err != nil {
		return fmt.Errorf("soft delete bank account: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *BankAccountRepository) queryOne(ctx context.Context, op, sql string, args ...any) (*models.BankAccount, error) {
	a, err := scanBankAccount(r.db.QueryRow(ctx, sql, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, mapPgError(err))
	}
	return &a, nil
}

func mapPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgUniqueViolation:
			return repository.ErrConflict
		case pgCheckViolation:
			return &repository.InvalidDataError{Constraint: pgErr.ConstraintName}
		}
	}
	return err
}

func scanBankAccount(row pgx.Row) (models.BankAccount, error) {
	var (
		a      models.BankAccount
		id     string
		status string
	)
	if err := row.Scan(
		&id, &a.AccountNumber, &a.BankName, &a.AccountName, &a.BankCode,
		&status, &a.Reason, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	); err != nil {
		return models.BankAccount{}, err
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return models.BankAccount{}, fmt.Errorf("parse bank account uuid: %w", err)
	}
	a.BankAccountUUID = parsed
	a.Status = models.BankAccountStatus(status)
	return a, nil
}

func buildBankAccountWhere(f repository.BankAccountFilter) (string, []any) {
	conds := []string{"deleted_at IS NULL"}
	var args []any

	if f.Search != "" {
		args = append(args, "%"+escapeLike(f.Search)+"%")
		n := len(args)
		conds = append(conds, fmt.Sprintf(
			"(account_number ILIKE $%d OR account_name ILIKE $%d OR bank_name ILIKE $%d OR bank_code ILIKE $%d)",
			n, n, n, n))
	}
	if f.Status != nil {
		args = append(args, string(*f.Status))
		conds = append(conds, fmt.Sprintf("status = $%d", len(args)))
	}
	if f.BankCode != "" {
		args = append(args, f.BankCode)
		conds = append(conds, fmt.Sprintf("bank_code = $%d", len(args)))
	}

	return " WHERE " + strings.Join(conds, " AND "), args
}

func bankAccountSortColumn(field string) string {
	switch field {
	case repository.BankAccountSortAccountNumber:
		return "account_number"
	case repository.BankAccountSortAccountName:
		return "account_name"
	case repository.BankAccountSortBankName:
		return "bank_name"
	case repository.BankAccountSortBankCode:
		return "bank_code"
	case repository.BankAccountSortStatus:
		return "status"
	case repository.BankAccountSortUpdatedAt:
		return "updated_at"
	default:
		return "created_at"
	}
}