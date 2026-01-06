package category

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// categoryRow is the database representation of a category.
type categoryRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r categoryRow) toDomain() (Category, error) {
	id, err := ParseID(r.ID)
	if err != nil {
		return Category{}, err
	}
	return Category{
		ID:        id,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func toRow(c Category) categoryRow {
	return categoryRow{
		ID:        c.ID.String(),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	db *pgx.Conn
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, c Category) error {
	row := toRow(c)
	_, err := r.db.Exec(ctx,
		`INSERT INTO categories (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		row.ID, row.Name, row.CreatedAt, row.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id ID) (Category, error) {
	var row categoryRow
	err := r.db.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`,
		id.String(),
	).Scan(&row.ID, &row.Name, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Category{}, ErrNotFound
		}
		return Category{}, err
	}
	return row.toDomain()
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, created_at, updated_at FROM categories ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var row categoryRow
		if err := rows.Scan(&row.ID, &row.Name, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		c, err := row.toDomain()
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, c Category) error {
	row := toRow(c)
	result, err := r.db.Exec(ctx,
		`UPDATE categories SET name = $2, updated_at = $3 WHERE id = $1`,
		row.ID, row.Name, row.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id ID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id.String())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
