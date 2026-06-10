// The Core Business Logic orchestration layer. Notice how it depends purely on the abstractions
// (domain.TodoRepository and domain.CacheRepository).

package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"todo-api/internal/handler/domain"
)

type TodoService struct {
	repo  domain.TodoRepository
	cache domain.CacheRepository
}

func NewTodoService(repo domain.TodoRepository, cache domain.CacheRepository) *TodoService {
	return &TodoService{
		repo:  repo,
		cache: cache,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	// Generate a simple unique ID
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	todo.ID = fmt.Sprintf("%x", b)

	// Save to whatever repository was injected
	return s.repo.Create(ctx, todo)
}

func (s *TodoService) GetTodo(ctx context.Context, id string) (*domain.Todo, error) {
	// 1. You could check your cache repository layer first here if required
	// 2. Fetch from underlying structural repository
	return s.repo.GetByID(ctx, id)
}

func (s *TodoService) ListTodos(ctx context.Context) ([]*domain.Todo, error) {
	return s.repo.GetAll(ctx)
}
