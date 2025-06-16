package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/ibs-alt/Go_Assignments/todo"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create base context with trace ID
	traceID := uuid.NewString()
	ctx = context.WithValue(ctx, todo.TraceIDKey, traceID)

	// Load todos from file
	todos, err := todo.LoadTodos(ctx, todo.FileName)
	if err != nil {
		slog.Error("Failed to load todos", "traceID", traceID, "error", err)
		os.Exit(1)
	}

	// Start CLI menu in a goroutine
	go todo.StartCLIMenu(ctx, &todos, todo.FileName)

	// Setup ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/create", todo.CreateHandler(&todos))
	mux.HandleFunc("/get", todo.GetHandler(&todos))
	mux.HandleFunc("/update", todo.UpdateHandler(&todos))
	mux.HandleFunc("/delete", todo.DeleteHandler(&todos))

	// Static /about page
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/about/", http.StripPrefix("/about", fs))

	// Dynamic /list page
	mux.HandleFunc("/list", todo.ListHandler(&todos))

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: todo.TraceMiddleware(mux),
	}

	go func() {
		slog.Info("Starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown
	<-ctx.Done()
	slog.Info("Shutting down Program...")
	server.Shutdown(context.Background())
	todo.SaveTodos(ctx, todo.FileName, todos)
	slog.Info("Shutdown complete.")
}
