// The REST Endpoints implementation mapping out incoming payloads
//  to the TodoService calls using gorilla/mux routing logic.

package http

import (
	"encoding/json"
	"net/http"
	"todo-api/internal/handler/domain"
	"todo-api/internal/service"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	svc *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// Create handles the HTTP POST request.
// @Summary      Create a new Todo
// @Description  Create a todo entity inside our abstract datastore layer
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body      domain.Todo  true  "Todo request item payload"
// @Success      210   {object}  domain.Todo
// @Failure      400   {string}  string "Invalid JSON Input"
// @Router       /todos [post]
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Todo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid JSON input payload"}`))
		return
	}

	if err := h.svc.CreateTodo(r.Context(), &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Failed to create data persistence layer object"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(req)
}

// GetByID handles retrieving a specific task item.
// @Summary      Get a single Todo item by identifier string
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo Element ID"
// @Success      200  {object}  domain.Todo
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	todo, err := h.svc.GetTodo(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Todo target entry index element not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}
