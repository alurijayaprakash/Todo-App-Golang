// he system bootstrap root application orchestrator.
// This reads our mock environmental toggles and acts as the runtime dependency constructor injector.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "todo-api/docs"
	"todo-api/internal/config"
	"todo-api/internal/handler/domain"
	appHttp "todo-api/internal/handler/http"
	"todo-api/internal/repository/memory"
	"todo-api/internal/repository/postgres"
	"todo-api/internal/service"
)

// @title           Todo Core API Spec Layout
// @version         1.0
// @description     Golang REST Production architecture clean skeleton layout engine.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// --- 1. Load Configurations ---
	cfg := config.Load()

	var activeRepository domain.TodoRepository

	// --- 2. Dynamic Dependency Swapping ---
	if cfg.DBType == "postgres" {
		log.Println("[Init] Swapping implementation to: PostgreSQL Engine")
		var err error
		activeRepository, err = postgres.NewPostgresTodoRepository(cfg.GetDatabaseURL())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("[Init] Swapping implementation to: In-Memory Datastore Engine")
		activeRepository = memory.NewMemoryTodoRepository()
	}

	activeCacheStore := memory.NewMemoryCache()

	// --- 3. Inject Dependencies into Service ---
	todoService := service.NewTodoService(activeRepository, activeCacheStore)

	// --- 4. Build Transport Routers Routing Layer ---
	routerEngine := appHttp.NewRouter(todoService, cfg.IsSwaggerEnabled())

	if cfg.IsSwaggerEnabled() {
		log.Println("[Bootstrap] Swagger route enabled")
	} else {
		log.Println("[Bootstrap] Swagger route disabled for production APP_ENV")
	}

	server := &http.Server{
		Handler:      routerEngine,
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("[Bootstrap] Target Web Server Engine successfully online listening active port bind " + cfg.ServerPort)
	log.Fatal(server.ListenAndServe())
}
