// A PostgreSQL implementation of the TodoRepository interface using database/sql with pq driver.

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"todo-api/internal/handler/domain"

	_ "github.com/lib/pq"
)

type PostgresTodoRepository struct {
	db *sql.DB
}

// NewPostgresTodoRepository creates a new PostgreSQL repository with the given connection string.
func NewPostgresTodoRepository(connString string) (domain.TodoRepository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresTodoRepository{db: db}, nil
}

// Create inserts a new todo into the database.
func (r *PostgresTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	if todo == nil {
		return errors.New("todo cannot be nil")
	}

	insertSQL := `
	INSERT INTO todos (id, title, completed)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO UPDATE SET title = $2, completed = $3, updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, insertSQL, todo.ID, todo.Title, todo.Completed)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// GetByID retrieves a todo by its ID from the database.
func (r *PostgresTodoRepository) GetByID(ctx context.Context, id string) (*domain.Todo, error) {
	if id == "" {
		return nil, errors.New("todo id cannot be empty")
	}

	selectSQL := `SELECT id, title, completed FROM todos WHERE id = $1`

	row := r.db.QueryRowContext(ctx, selectSQL, id)

	todo := &domain.Todo{}
	err := row.Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

// GetAll retrieves all todos from the database.
func (r *PostgresTodoRepository) GetAll(ctx context.Context) ([]*domain.Todo, error) {
	selectAllSQL := `SELECT id, title, completed FROM todos ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, selectAllSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []*domain.Todo

	for rows.Next() {
		todo := &domain.Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}

// Close closes the database connection.
func (r *PostgresTodoRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
