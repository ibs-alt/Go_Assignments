package todo_test

import (
	"context"
	"os"
	"testing"

	"github.com/ibs-alt/Go_Assignments/todo"
)

func TestAddTodo(t *testing.T) {
	todos := []todo.Todo{}
	ctx := context.WithValue(context.Background(), todo.TraceIDKey, "test-trace-id")

	todos = todo.AddTodo(ctx, todos, "Test Task 1")

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].Description != "Test Task 1" {
		t.Errorf("Expected 'Test Task 1', got '%s'", todos[0].Description)
	}
}

func TestSaveAndLoadTodos(t *testing.T) {
	ctx := context.WithValue(context.Background(), todo.TraceIDKey, "test-trace-id")
	fileName := "testdata.json"
	defer os.Remove(fileName)

	todos := []todo.Todo{
		{ID: "1", Description: "Task A", Status: "Not started"},
	}

	err := todo.SaveTodos(ctx, fileName, todos)
	if err != nil {
		t.Fatalf("Failed to save todos: %v", err)
	}

	loadedTodos, err := todo.LoadTodos(ctx, fileName)
	if err != nil {
		t.Fatalf("Failed to load todos: %v", err)
	}

	if len(loadedTodos) != 1 {
		t.Errorf("Expected 1 loaded todo, got %d", len(loadedTodos))
	}
	if loadedTodos[0].Description != "Task A" {
		t.Errorf("Expected 'Task A', got '%s'", loadedTodos[0].Description)
	}
}
