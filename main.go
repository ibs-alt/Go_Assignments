package main

import (
	"bufio"         // For reading user input
	"encoding/json" // For JSON serialization/deserialization
	"flag"          // For command-line flag parsing
	"fmt"           // For formatted I/O
	"log"           // For error logging
	"os"            // For file operations
	"strconv"       // For string to integer conversion
	"strings"       // For string manipulation
)

// Todo represents a single task item with unique ID, description and status
type Todo struct {
	ID          int    `json:"id"`          // Unique identifier for each task
	Description string `json:"description"` // Task description
	Status      string `json:"status"`      // Current status of the task
}

// fileName stores the path to the JSON file, configurable via command-line flag
var fileName string

// init runs before main(), sets up command-line flags
func init() {
	flag.StringVar(&fileName, "file", "data.json", "Path to the JSON file storing tasks")
	flag.Parse()
}

// loadTodos reads and parses the JSON file containing tasks
// Returns the slice of todos and any error encountered
func loadTodos() ([]Todo, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	if err := json.Unmarshal(file, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

// saveTodos writes the todos slice to the JSON file
// Returns any error encountered during saving
func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

// printTodos displays all todos with their IDs and status
func printTodos(todos []Todo) {
	fmt.Println("\nYour To-Do List:")
	for _, t := range todos {
		fmt.Printf("ID %d: %s [%s]\n", t.ID, t.Description, t.Status)
	}
}

// nextID determines the next available ID for a new task
// Returns the maximum ID found plus one
func nextID(todos []Todo) int {
	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func main() {
	// Load existing todos, handle errors except for non-existent file
	todos, err := loadTodos()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to load todos: %v", err)
	}

	// Create buffered reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Define menu text
	menu := `
Choose an option:
1. Add Task
2. List Tasks
3. Update Task Status
4. Delete Task
5. Exit
-> `

	// Main program loop
	for {
		fmt.Print(menu)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1": // Add new task
			fmt.Print("Enter Task: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			// Create new task with next available ID
			newTask := Todo{ID: nextID(todos), Description: desc, Status: "Not started"}
			todos = append(todos, newTask)
			if err := saveTodos(todos); err != nil {
				log.Printf("Failed to save task: %v", err)
			} else {
				fmt.Println("Task added.")
			}

		case "2": // List all tasks
			printTodos(todos)

		case "3": // Update task status
			printTodos(todos)
			// Get task ID from user
			fmt.Print("Enter ID of task to update: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			id, err := strconv.Atoi(idInput)
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}

			// Find task with given ID
			index := -1
			for i, t := range todos {
				if t.ID == id {
					index = i
					break
				}
			}

			if index == -1 {
				fmt.Println("Task not found.")
				continue
			}

			// Show status options and get user choice
			fmt.Println("Choose new status:")
			fmt.Println("1. Not started")
			fmt.Println("2. Started")
			fmt.Println("3. Completed")
			fmt.Print("-> ")
			statusInput, _ := reader.ReadString('\n')
			statusInput = strings.TrimSpace(statusInput)

			// Update task status based on user choice
			switch statusInput {
			case "1":
				todos[index].Status = "Not started"
			case "2":
				todos[index].Status = "Started"
			case "3":
				todos[index].Status = "Completed"
			default:
				fmt.Println("Invalid choice.")
				continue
			}

			// Save updated todos
			if err := saveTodos(todos); err != nil {
				log.Printf("Failed to update task: %v", err)
			} else {
				fmt.Println("Task updated.")
			}

		case "4": // Delete task
			printTodos(todos)
			// Get task ID from user
			fmt.Print("Enter ID of task to delete: ")
			idInput, _ := reader.ReadString('\n')
			idInput = strings.TrimSpace(idInput)
			id, err := strconv.Atoi(idInput)
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}

			// Find task with given ID
			index := -1
			for i, t := range todos {
				if t.ID == id {
					index = i
					break
				}
			}

			if index == -1 {
				fmt.Println("Task not found.")
				continue
			}

			// Remove task and save updated todos
			todos = append(todos[:index], todos[index+1:]...)
			if err := saveTodos(todos); err != nil {
				log.Printf("Failed to delete task: %v", err)
			} else {
				fmt.Println("Task deleted.")
			}

		case "5": // Exit program
			fmt.Println("Goodbye!")
			return

		default: // Handle invalid input
			fmt.Println("Invalid option.")
		}
	}
}
