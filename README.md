To-Do List (HTTP + Web UI, Actor Safe Concurrency)

A Go application that exposes a JSON REST API and simple web pages for managing a to-do list.


Features:
  1. Add tasks with a description
  2. View all current tasks with their status
  3. Update a task status(Not started, Started and Completed)
  4. Delete Tasks by ID
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

Curl:

To Create a Task:
$ curl -X POST -H "Content-Type: application/json" -d '{"description": "INSERT TEST NAME HERE"}' http://localhost:8080/create

To Update a Task:
$ curl -X PUT -H "Content-Type: application/json" -d '{"description": "INSERT TEST NAME HERE", "id": "INSERT ID HERE", "status": "INSERT UPDATED STATUS HERE"}' http://localhost:8080/update 

To Delete a Task:
$ curl -X DELETE -H "Content-Type: application/json" -d '{"id": "INSERT ID HERE "}' http://localhost:8080/delete

To View all Tasks:
curl http://localhost:8080/get


Tests:

Functional and Concurrency Tests:
Go test -v ./todo

Race Detector (Needs gcc / CGO):
CGO_ENABLED=1 go test -race ./todo


Actor_test.go fires many concurrent requests to demonstrate the actor layer is race free
Todo_test.go checks add/update/delete JSON persistence


---
Author

Created by Ibrahim Jhetam


  
  
