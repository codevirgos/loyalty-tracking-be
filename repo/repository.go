package repo

import (
	"Payback_BE/models"
	"context"
	"database/sql"
	"errors"
	"os"
	"time"
)

type UserRepo interface {
	// todo lookup why to use pointer here
	LookUp(ctx context.Context, phoneNumber int) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) (*models.User, error)
}

type pgUserRepo struct {
	db *sql.DB
}

// non implementation specific
var (
	ErrInternal = errors.New("internal db error")
	ErrNotFound = errors.New("user not found")
)

func InitPostgresDB() (*sql.DB, error) {
	db, err := sql.Open("postgress", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err

	}
	// todo reseaarch this
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	return db, nil
}

func NewPgUserRepo(db *sql.DB) *pgUserRepo {
	return &pgUserRepo{
		db: db,
	}
}

// -------------------- Queries ------------------------//
func (pr *pgUserRepo) LookUp(ctx context.Context, phoneNumber int) (*models.User, error) {
	query := `Select nkey, name, points, visit_count FROM users WHERE nkey = $1`

	var user models.User
	err := pr.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.Number,
		&user.Name,
		&user.TotalPoints,
		&user.Visits,
	)

	// todo be careful with this non-generic error capture
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (pr *pgUserRepo) SaveUser(ctx context.Context, u *models.User) (*models.User, error) {
	_, err := pr.db.Exec("Insert into users (nkey, name, points, visit_count) VALUES ($1, $2, $3, $4)", u.Number, u.Name, u.TotalPoints, u.Visits)

	if err != nil {
		return nil, ErrInternal
	}
	return u, err
}
