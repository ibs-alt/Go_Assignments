package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"bufio"

	"github.com/google/uuid"
)

type contextKey string // Define a custom type for context keys to avoid collisions

const TraceIDKey contextKey = "traceID" // Key for storing trace ID in context

var FileName = "data.json" // Default file name for storing todos

type Todo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func LoadTodos(ctx context.Context, file string) ([]Todo, error) { // LoadTodos reads todos from a JSON file
	traceID := ctx.Value(TraceIDKey)
	data, err := os.ReadFile(file)
	if err != nil {
		slog.Error("Failed to read file", "traceID", traceID, "error", err)
		return nil, err
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		slog.Error("Failed to unmarshal JSON", "traceID", traceID, "error", err)
		return nil, err
	}
	return todos, nil
}

func AddTodo(ctx context.Context, todos []Todo, desc string) []Todo { // AddTodo creates a new todo item and appends it to the list
	traceID := ctx.Value(TraceIDKey)
	id := uuid.NewString()
	todo := Todo{ID: id, Description: desc, Status: "Not started"}
	todos = append(todos, todo)
	slog.Info("Added new todo", "traceID", traceID, "todo", todo)
	return todos
}

func SaveTodos(ctx context.Context, file string, todos []Todo) { // SaveTodos writes the current list of todos to a JSON file
	traceID := ctx.Value(TraceIDKey)
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal JSON", "traceID", traceID, "error", err)
		return
	}
	if err := os.WriteFile(file, data, 0644); err != nil {
		slog.Error("Failed to write file", "traceID", traceID, "error", err)
		return
	}
	slog.Info("Data saved to disk", "traceID", traceID)
}

func StartCLIMenu(ctx context.Context, todos *[]Todo, file string) { // StartCLIMenu starts a command-line interface for managing todos
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`
Choose an option:
1. Add Task
2. List Tasks
3. Update Task Status
4. Delete Task
5. Exit
-> `)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("Enter Task: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			id := generateID()
			*todos = append(*todos, Todo{ID: id, Description: desc, Status: "Not started"})
			SaveTodos(ctx, file, *todos)
			fmt.Println("Task added.")

		case "2":
			fmt.Println("\nYour To-Do List:")
			for i, t := range *todos {
				fmt.Printf("%d. %s [%s]\n", i+1, t.Description, t.Status)
			}

		case "3":
			for i, t := range *todos {
				fmt.Printf("%d. %s [%s]\n", i+1, t.Description, t.Status)
			}
			fmt.Print("Enter task number to update: ")
			var i int
			fmt.Scanln(&i)
			if i < 1 || i > len(*todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			fmt.Print(` 
Choose new status:
1. Not started
2. Started
3. Completed
-> `)
			var statusChoice int
			fmt.Scanln(&statusChoice)

			statusMap := map[int]string{1: "Not started", 2: "Started", 3: "Completed"} // Map for status choices
			status, ok := statusMap[statusChoice]
			if !ok {
				fmt.Println("Invalid choice.")
				continue
			}
			(*todos)[i-1].Status = status
			SaveTodos(ctx, file, *todos)
			fmt.Println("Task updated.")

		case "4":
			for i, t := range *todos {
				fmt.Printf("%d. %s [%s]\n", i+1, t.Description, t.Status)
			}
			fmt.Print("Enter task number to delete: ")
			var i int
			fmt.Scanln(&i)
			if i < 1 || i > len(*todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			*todos = append((*todos)[:i-1], (*todos)[i:]...)
			SaveTodos(ctx, file, *todos)
			fmt.Println("Task deleted.")

		case "5":
			fmt.Println("Exiting...")
			SaveTodos(ctx, file, *todos)
			os.Exit(0)

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func generateID() string {
	return uuid.NewString()
}
