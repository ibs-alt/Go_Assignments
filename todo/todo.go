package todo

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type contextKey string

const TraceIDKey contextKey = "traceID"

var FileName = "data.json"

type Todo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func LoadTodos(ctx context.Context, file string) ([]Todo, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var todos []Todo
	return todos, json.Unmarshal(data, &todos)
}

func SaveTodos(_ context.Context, file string, todos []Todo) error {
	data, _ := json.MarshalIndent(todos, "", "  ")
	return os.WriteFile(file, data, 0644)
}

// Convenience: create brand-new Todo
func NewTodo(desc string) Todo {
	return Todo{ID: uuid.NewString(), Description: desc, Status: "Not started"}
}
