package postgres

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

// TODO: public?
func NewUsersRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db}
}

// TODO: public?
func (r *UsersRepository) Create(ctx context.Context, login string, password string) (domain.User, error) {
	var newUser domain.User
	err := r.db.QueryRow(ctx,
		"INSERT INTO users(login, password) VALUES ($1, $2) RETURNING *",
		login, password).Scan(
		&newUser.ID, &newUser.Login, &newUser.Password)

	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (r *UsersRepository) FindByLogin(ctx context.Context, login string) (domain.User, error) {
	var user domain.User

	err := r.db.QueryRow(ctx,
		"SELECT * FROM users WHERE login = $1", login).Scan(
		&user.ID, &user.Login, &user.Password)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
