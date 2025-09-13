package impl

import (
	"context"

	"evently/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) model.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt)

	return err
}

func (r *userRepositoryImpl) GetByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at
		FROM users WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at
		FROM users WHERE email = $1`

	user := &model.User{}
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
