package repositories

import (
	"context"
	"errors"

	"github.com/danilobml/go-cached/internal/db"
	"github.com/danilobml/go-cached/internal/errs"
	"github.com/danilobml/go-cached/internal/models"
	"github.com/jackc/pgx/v4"
)

type PgUserRepository struct {
	db db.DBInterface
}

func NewPgUserRepository(db db.DBInterface) *PgUserRepository {
	return &PgUserRepository{
		db: db,
	}
}

func (ur *PgUserRepository) Create(ctx context.Context, username, email string) error {
	query := `
		INSERT INTO users (username, email) 
		VALUES ($1, $2)
		RETURNING id, username, email
		`

	var newUser models.User

	err := ur.db.QueryRow(ctx, query, username, email).Scan(&newUser.Id, &newUser.Username, &newUser.Email)
	if err != nil {
		return err
	}

	return nil
}

func (ur *PgUserRepository) List(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, email 
				FROM users`
	rows, err := ur.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *PgUserRepository) FindById(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, email
				FROM users
				WHERE id = $1`

	var user models.User

	err := ur.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
