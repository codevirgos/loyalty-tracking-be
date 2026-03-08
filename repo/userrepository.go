package repo

import (
	"Payback_BE/models"
	"context"
	"database/sql"
	"errors"
)

type UserRepo interface {
	SaveUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUser(ctx context.Context, id int) (*models.User, error)
	FindUserFromNum(ctx context.Context, number string) (*models.User, error)
}

type pgUserRepo struct {
	db *sql.DB
}

func NewPgUserRepo(db *sql.DB) *pgUserRepo {
	return &pgUserRepo{
		db: db,
	}
}

// non implementation specific
var (
	ErrInternal = errors.New("internal db error")
	ErrNotFound = errors.New("user not found")
)

// -------------------- Queries ------------------------//
func (pr *pgUserRepo) FindUser(ctx context.Context, id int) (*models.User, error) {
	query := `Select id, nkey, user_name FROM users WHERE id = $1`

	var user models.User
	err := pr.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Number,
		&user.Name,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (pr *pgUserRepo) FindUserFromNum(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `Select id, nkey, user_name FROM users WHERE nkey = $1`

	var user models.User
	err := pr.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Number,
		&user.Name,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (pr *pgUserRepo) SaveUser(ctx context.Context, user *models.User) (*models.User, error) {
	var id int
	err := pr.db.QueryRow("Insert into users (nkey, user_name) VALUES ($1, $2) RETURNING id", user.Number, user.Name).Scan(&id)

	if err != nil {
		return nil, ErrInternal
	}

	u := &models.User{
		ID:     id,
		Number: user.Number,
		Name:   user.Name,
	}
	return u, nil
}
