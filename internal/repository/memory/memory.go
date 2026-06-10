// An In-Memory implementation of the TodoRepository interface using a thread-safe sync.Map.
package memory

import (
	"context"
	"errors"
	"sync"
	"todo-api/internal/handler/domain"
)

type MemoryTodoRepository struct {
	store sync.Map
}

func NewMemoryTodoRepository() domain.TodoRepository {
	return &MemoryTodoRepository{}
}

func (r *MemoryTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	r.store.Store(todo.ID, todo)
	return nil
}

func (r *MemoryTodoRepository) GetByID(ctx context.Context, id string) (*domain.Todo, error) {
	val, exists := r.store.Load(id)
	if !exists {
		return nil, errors.New("todo not found in-memory")
	}
	return val.(*domain.Todo), nil
}

func (r *MemoryTodoRepository) GetAll(ctx context.Context) ([]*domain.Todo, error) {
	var todos []*domain.Todo
	r.store.Range(func(key, value interface{}) bool { // TODO: Need to Understand this part
		todos = append(todos, value.(*domain.Todo))
		return true
	})
	return todos, nil
}
