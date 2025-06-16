package todo

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateID()
		}
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateHandler(todos *[]Todo) http.HandlerFunc { //
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		t.ID = generateID()
		*todos = append(*todos, t)
		SaveTodos(ctx, FileName, *todos)
		w.WriteHeader(http.StatusCreated)
	}
}

func GetHandler(todos *[]Todo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		json.NewEncoder(w).Encode(*todos)
		slog.Info("Returned todos", "traceID", ctx.Value(TraceIDKey))
	}
}

func UpdateHandler(todos *[]Todo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		for i := range *todos {
			if (*todos)[i].ID == input.ID {
				(*todos)[i].Status = input.Status
				SaveTodos(ctx, FileName, *todos)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}

func DeleteHandler(todos *[]Todo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		for i := range *todos {
			if (*todos)[i].ID == input.ID {
				*todos = append((*todos)[:i], (*todos)[i+1:]...)
				SaveTodos(ctx, FileName, *todos)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
	}
}
