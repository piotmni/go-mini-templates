package note

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
)

// noteRow is the database representation of a note.
type noteRow struct {
	ID         string    `db:"id"`
	CategoryID string    `db:"category_id"`
	Title      string    `db:"title"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (r noteRow) toDomain() (Note, error) {
	id, err := ParseID(r.ID)
	if err != nil {
		return Note{}, err
	}
	catID, err := category.ParseID(r.CategoryID)
	if err != nil {
		return Note{}, err
	}
	return Note{
		ID:         id,
		CategoryID: catID,
		Title:      r.Title,
		Content:    r.Content,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}, nil
}

func toRow(n Note) noteRow {
	return noteRow{
		ID:         n.ID.String(),
		CategoryID: n.CategoryID.String(),
		Title:      n.Title,
		Content:    n.Content,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
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

func (r *PostgresRepository) Create(ctx context.Context, n Note) error {
	row := toRow(n)
	_, err := r.db.Exec(ctx,
		`INSERT INTO notes (id, category_id, title, content, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		row.ID, row.CategoryID, row.Title, row.Content, row.CreatedAt, row.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id ID) (Note, error) {
	var row noteRow
	err := r.db.QueryRow(ctx,
		`SELECT id, category_id, title, content, created_at, updated_at FROM notes WHERE id = $1`,
		id.String(),
	).Scan(&row.ID, &row.CategoryID, &row.Title, &row.Content, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return row.toDomain()
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Note, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, category_id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var row noteRow
		if err := rows.Scan(&row.ID, &row.CategoryID, &row.Title, &row.Content, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		n, err := row.toDomain()
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *PostgresRepository) GetByCategory(ctx context.Context, categoryID category.ID) ([]Note, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, category_id, title, content, created_at, updated_at FROM notes WHERE category_id = $1 ORDER BY created_at DESC`,
		categoryID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var row noteRow
		if err := rows.Scan(&row.ID, &row.CategoryID, &row.Title, &row.Content, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		n, err := row.toDomain()
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, n Note) error {
	row := toRow(n)
	result, err := r.db.Exec(ctx,
		`UPDATE notes SET category_id = $2, title = $3, content = $4, updated_at = $5 WHERE id = $1`,
		row.ID, row.CategoryID, row.Title, row.Content, row.UpdatedAt,
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
	result, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id.String())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
