package todo

import (
	"html/template" // for HTML templating
	"net/http"      // for HTTP server functionality
)

func ListHandler(todos *[]Todo) http.HandlerFunc { // Create a handler function for listing todos
	tmpl := template.Must(template.ParseFiles("templates/list.html")) // Load the HTML template for listing todos
	return func(w http.ResponseWriter, r *http.Request) {             // Handle the HTTP request for listing todos
		tmpl.Execute(w, *todos) // Render the template with the current list of todos
	}
}
