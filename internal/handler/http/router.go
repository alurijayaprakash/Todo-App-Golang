// Initializes your HTTP routing engine, setups basic middleware chains,
// and injects our custom Swagger UI routes.

package http

import (
	"net/http"
	"todo-api/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(todoService *service.TodoService, enableSwagger bool) *mux.Router {
	r := mux.NewRouter()

	// Global Middlewares using gorilla/mux style
	r.Use(loggingMiddleware)

	if enableSwagger {
		// Attach Swagger UI route entry point
		r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	}

	todoHandler := NewTodoHandler(todoService)

	// API Routing Sub-Router Mapping Definitions
	api := r.PathPrefix("/api/v1").Subrouter()
	healthHandler := NewHealthHandler()
	api.HandleFunc("/health", healthHandler.Ping).Methods(http.MethodGet)
	api.HandleFunc("/todos", todoHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/todos/{id}", todoHandler.GetByID).Methods(http.MethodGet)

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("Incoming Request:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
