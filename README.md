Simple command line To-Do List application. It lets you add, view, update and delete tasks with persistant storage using a local JSON file.

Features:
  1. Add tasks with a description
  2. View all curret tasks with their status
  3. Update a task status(Not started, In Progress and Completed)
  4. Delete Tasks by number
  5. Automatically saves and loads tasks from data.json

Prerequisites:
  1. Go installed (Version 1.20+ recommended)
  2. A terminal (VS Code, Powershell, Windows Command Prompt)

How To Run:
  1. Clone or Download the Repo or if you're just using it locally create a folder and place main.go and data.json into it.
  2. Open the folder in your terminal
  3. Initialise Go module (only need to complete this once)
     go mod init
  4. Run the app
     go run .

How to use:

Once running you will see the following menu:

  Choose an option:
  1. Add Task
  2. List Tasks
  3. Update Task Status
  4. Delete Task
  5. Exit
     
     ->     (Enter the number for the action you want here and follow the prompts)

  Your tasks are saved automatically to the data.json in the same folder

---
Author

Created by Ibrahim Jhetam


  
  
