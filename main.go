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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// root context & graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx = context.WithValue(ctx, todo.TraceIDKey, uuid.NewString())

	// load initial data then start actor
	initial, _ := todo.LoadTodos(ctx, todo.FileName)
	readCh, writeCh := todo.StartActor(ctx, todo.FileName, initial)

	// HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/create", todo.TraceMiddleware(todo.CreateHandler(readCh, writeCh)))
	mux.Handle("/get", todo.TraceMiddleware(todo.GetHandler(readCh)))
	mux.Handle("/update", todo.TraceMiddleware(todo.UpdateHandler(readCh, writeCh)))
	mux.Handle("/delete", todo.TraceMiddleware(todo.DeleteHandler(writeCh)))

	// static & dynamic pages
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/about/", http.StripPrefix("/about", fs))
	mux.Handle("/list", todo.TraceMiddleware(todo.ListHandler(readCh)))

	go func() {
		slog.Info("HTTP server :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil && ctx.Err() == nil {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done() // wait Ctrl-C
	slog.Info("server shutdown")
}
