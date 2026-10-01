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

const bankAccountColumns = `bank_account_uuid, account_number, bank_name, account_name, bank_code, status, reason, created_at, updated_at`

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
	row := r.db.QueryRow(ctx,
		"SELECT "+bankAccountColumns+" FROM bank_accounts WHERE bank_account_uuid = $1 AND deleted_at IS NULL",
		id.String())

	a, err := scanBankAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find bank account by id: %w", err)
	}
	return &a, nil
}

func scanBankAccount(row pgx.Row) (models.BankAccount, error) {
	var (
		a      models.BankAccount
		id     string
		status string
	)
	if err := row.Scan(
		&id, &a.AccountNumber, &a.BankName, &a.AccountName, &a.BankCode,
		&status, &a.Reason, &a.CreatedAt, &a.UpdatedAt,
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