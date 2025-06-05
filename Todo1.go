// Package main implements a simple command-line todo list application
package main

import (
	"bufio"   // For reading user input
	"fmt"     // For formatted I/O operations
	"os"      // For standard input/output
	"strings" // For string manipulation
)

// Todo represents a single task item with a description and status
// Each task has a description and a status
type Todo struct {
	Description string // Stores the task description provided by user
	Status      string // Current status of the task
}

func main() {
	// Initialize an empty slice to store todo items
	todos := []Todo{}

	// Create a new buffered reader for reading user input
	reader := bufio.NewReader(os.Stdin)

	// Main application loop - continues until user chooses to exit
	for {
		// Display menu options
		fmt.Println("\n1. Add Task 2. List Tasks 3. Delete Task 4. Exit")
		fmt.Print("Choose an option: ")

		// Read user input and remove any whitespace/newlines
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Process user input using switch statement
		switch input {
		case "1":
			// Add new task
			fmt.Print("Enter Task: ")
			text, _ := reader.ReadString('\n')
			// Create and append new todo item to the slice
			todos = append(todos, Todo{
				Description: strings.TrimSpace(text),
				Status:      "Not started", // Default status for new tasks
			})
			fmt.Println("Task added.")

		case "2":
			// Handle listing all tasks
			if len(todos) == 0 {
				fmt.Println("No tasks available.")
				continue
			}
			fmt.Println("Tasks:")
			// Iterate through and display all tasks with their indices
			for i, todo := range todos {
				fmt.Printf("%d. %s [%s]\n", i+1, todo.Description, todo.Status)
			}

		case "3":
			// Handle deleting a task
			if len(todos) == 0 {
				fmt.Println("No tasks available to delete.")
				continue
			}
			// Display current tasks for user reference
			fmt.Println("Tasks:")
			for i, todo := range todos {
				fmt.Printf("%d. %s [%s]\n", i+1, todo.Description, todo.Status)
			}
			// Get task number to delete
			fmt.Print("Enter task number to delete: ")
			var taskNum int
			fmt.Scanln(&taskNum)
			// Validate task number
			if taskNum < 1 || taskNum > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			// Remove the selected task using slice operations
			todos = append(todos[:taskNum-1], todos[taskNum:]...)
			fmt.Println("Task deleted.")

		case "4":
			// Handle program exit
			fmt.Println("Exiting...")
			return

		default:
			// Handle invalid input
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
