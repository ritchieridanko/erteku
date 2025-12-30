package databases

import (
	"context"
	"errors"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/database"
	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type AuthDatabase interface {
	Insert(ctx context.Context, data *models.CreateAuth) (auth *models.Auth, err *ce.Error)
	GetByEmail(ctx context.Context, email string) (auth *models.Auth, err *ce.Error)
	IsEmailRegistered(ctx context.Context, email string) (registered bool, err *ce.Error)
}

type authDatabase struct {
	database *database.Database
}

func NewAuthDatabase(db *database.Database) AuthDatabase {
	return &authDatabase{database: db}
}

func (d *authDatabase) Insert(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	query := `
		INSERT INTO auth (email, password)
		VALUES ($1, $2)
		RETURNING id, email, email_verified_at
	`

	row := d.database.Query(ctx, query, data.Email, data.Password)

	var a models.Auth
	if err := row.Scan(&a.ID, &a.Email, &a.EmailVerifiedAt); err != nil {
		return nil, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err)
	}

	return &a, nil
}

func (d *authDatabase) GetByEmail(ctx context.Context, email string) (*models.Auth, *ce.Error) {
	query := `
		SELECT id, email, password, email_verified_at
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := d.database.Query(ctx, query, email)

	var a models.Auth
	if err := row.Scan(&a.ID, &a.Email, &a.Password, &a.EmailVerifiedAt); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err)
	}

	return &a, nil
}

func (d *authDatabase) IsEmailRegistered(ctx context.Context, email string) (bool, *ce.Error) {
	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := d.database.Query(ctx, query, email)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err)
	}

	return true, nil
}
