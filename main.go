package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Todo struct { // Todo represents a single task with a description and status.
	Description string `json:"description"`
	Status      string `json:"status"`
}

const fileName = "data.json"

func loadTodos() []Todo { // loadTodos reads the JSON-encoded tasks from the file
	var todos []Todo
	data, err := os.ReadFile(fileName)
	if err == nil {
		_ = json.Unmarshal(data, &todos)
	}
	return todos
}

func saveTodos(todos []Todo) { // saveTodos into JSON format and writes them to the file.
	data, _ := json.MarshalIndent(todos, "", "  ")
	_ = os.WriteFile(fileName, data, 0644)
}

func printTodos(todos []Todo) { // printTodos prints all tasks with their index and current status
	fmt.Println("\nYour To-Do List:")
	for i, t := range todos {
		fmt.Printf("%d. %s [%s]\n", i+1, t.Description, t.Status)
	}
}

func main() {
	todos := loadTodos()                // Load tasks from file into memory.
	reader := bufio.NewReader(os.Stdin) // Create buffered reader for user input.

	for {
		fmt.Println("\nChoose an option:") // Menu prompt
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Update Task Status")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")
		fmt.Print("-> ")

		input, _ := reader.ReadString('\n') // Read user choice
		input = strings.TrimSpace(input)

		switch input { // Handle user input
		case "1":
			fmt.Print("Enter Task: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			todos = append(todos, Todo{Description: desc, Status: "Not started"})
			saveTodos(todos)
			fmt.Println("Task added.")

		case "2": // List all tasks
			printTodos(todos)

		case "3": // Update task status
			printTodos(todos)
			fmt.Print("Enter task number to update: ")
			var i int
			fmt.Scanln(&i)
			if i < 1 || i > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}

			fmt.Println("Choose new status:")
			fmt.Println("1. Not started")
			fmt.Println("2. In progress")
			fmt.Println("3. Completed")
			fmt.Print("-> ")

			var choice int
			fmt.Scanln(&choice)

			switch choice { // Update status based on user choice
			case 1:
				todos[i-1].Status = "Not started"
			case 2:
				todos[i-1].Status = "In progress"
			case 3:
				todos[i-1].Status = "Completed"
			default:
				fmt.Println("Invalid choice.")
				continue
			}

			saveTodos(todos)
			fmt.Println("Task updated.")

		case "4": // Delete a task
			printTodos(todos)
			fmt.Print("Enter task number to delete: ")
			var i int
			fmt.Scanln(&i)
			if i < 1 || i > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}

			todos = append(todos[:i-1], todos[i:]...) // Remove item by slicing
			saveTodos(todos)
			fmt.Println("Task deleted.")

		case "5":
			fmt.Println("Goodbye!")
			return

		default: // Handle invalid input
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
