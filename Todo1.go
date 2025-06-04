package main

import (
	"bufio"   // For reading user input
	"fmt"     // For formatted I/O operations
	"os"      // For standard input/output
	"strings" // For string manipulation
)

// Todo represents a single task item with a description and status
type Todo struct {
	Description string // The task description
	Status      string // Current status of the task
}

func main() {
	todos := []Todo{}                   // Initialize empty slice to store tasks
	reader := bufio.NewReader(os.Stdin) // Create new reader for user input

	// Main program loop
	for {
		// Display menu options
		fmt.Println("\n1. Add Task 2. List Tasks 3. Exit")
		fmt.Print("Choose an option: ")

		// Read and process user input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Handle different menu options
		switch input {
		case "1":
			// Add new task
			fmt.Print("Enter Task: ")
			text, _ := reader.ReadString('\n')
			todos = append(todos, Todo{
				Description: strings.TrimSpace(text),
				Status:      "Not started",
			})
			fmt.Println("Task added.")

		case "2":
			// List all tasks
			if len(todos) == 0 {
				fmt.Println("No tasks available.")
				continue
			}
			fmt.Println("Tasks:")
			// Display each task with its number, description and status
			for i, todo := range todos {
				fmt.Printf("%d. %s [%s]\n", i+1, todo.Description, todo.Status)
			}

		case "3":
			// Exit the program
			fmt.Println("Exiting...")
			return

		default:
			// Handle invalid input
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
