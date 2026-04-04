// Package postgres implements domain.UserRepository using pgx/v5.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourorg/go-user-crud/internal/domain"
)

// NewPool creates a connection pool with sane defaults.
func NewPool(dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	cfg.MaxConns = 25
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

// UserRepository is the Postgres-backed implementation.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository constructs a UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID fetches one user by primary key.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	user, err := scanUser(r.db.QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("GetByID query: %w", err)
	}
	return user, nil
}

// GetByEmail fetches one user by email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	user, err := scanUser(r.db.QueryRow(ctx, q, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("GetByEmail query: %w", err)
	}
	return user, nil
}

// List returns a page of users with optional search, plus the total count.
func (r *UserRepository) List(ctx context.Context, f domain.ListFilter) ([]*domain.User, int64, error) {
	// Build query dynamically — only add WHERE when searching.
	where := "deleted_at IS NULL"
	args := pgx.NamedArgs{}

	if f.Search != "" {
		where += " AND (name ILIKE @search OR email ILIKE @search)"
		args["search"] = "%" + f.Search + "%"
	}

	// Total count (same filter, no pagination).
	var total int64
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", where)
	if err := r.db.QueryRow(ctx, countQ, args).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("List count: %w", err)
	}

	args["limit"] = f.PerPage
	args["offset"] = f.Offset()

	dataQ := fmt.Sprintf(`
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE %s
		ORDER BY id DESC
		LIMIT @limit OFFSET @offset`, where)

	rows, err := r.db.Query(ctx, dataQ, args)
	if err != nil {
		return nil, 0, fmt.Errorf("List query: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("List scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("List rows: %w", err)
	}

	return users, total, nil
}

// Create inserts a new user and returns it with generated fields.
func (r *UserRepository) Create(ctx context.Context, in domain.CreateUserInput) (*domain.User, error) {
	const q = `
		INSERT INTO users (name, email, role)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, role, created_at, updated_at`

	user, err := scanUser(r.db.QueryRow(ctx, q, in.Name, in.Email, in.Role))
	if err != nil {
		return nil, fmt.Errorf("Create query: %w", err)
	}
	return user, nil
}

// Update applies a partial update and returns the updated user.
func (r *UserRepository) Update(ctx context.Context, id int64, in domain.UpdateUserInput) (*domain.User, error) {
	// Build SET clause dynamically so we only touch provided fields.
	setClauses := []string{"updated_at = NOW()"}
	args := []any{}
	i := 1

	if in.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", i))
		args = append(args, *in.Name)
		i++
	}
	if in.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", i))
		args = append(args, *in.Email)
		i++
	}
	if in.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", i))
		args = append(args, *in.Role)
		i++
	}

	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, name, email, role, created_at, updated_at`,
		strings.Join(setClauses, ", "), i)

	user, err := scanUser(r.db.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("Update query: %w", err)
	}
	return user, nil
}

// Delete soft-deletes a user by setting deleted_at.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("Delete exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// ── Scan helpers ──────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func scanUserRow(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
