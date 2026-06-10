// This defines your core entity and the Repository Interface
// (the contract that makes your databases swappable).

package domain

import "context"

// Todo represents the core business model.
type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoRepository defines the database operations required by the application.
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id string) (*Todo, error)
	GetAll(ctx context.Context) ([]*Todo, error)
}
