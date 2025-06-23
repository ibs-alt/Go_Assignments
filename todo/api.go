package todo

import (
	"encoding/json"
	"net/http"
)

// TraceMiddleware unchanged
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handler constructors that capture actor channels
func CreateHandler(reads chan<- readReq, writes chan<- writeReq) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Description string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", 400)
			return
		}

		t := NewTodo(body.Description)
		resp := make(chan error)
		writes <- writeReq{op: "add", todo: t, resp: resp}
		if err := <-resp; err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetHandler(reads chan<- readReq) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := readReq{resp: make(chan []Todo)}
		reads <- req
		todos := <-req.resp
		json.NewEncoder(w).Encode(todos)
	}
}

func UpdateHandler(reads chan<- readReq, writes chan<- writeReq) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ ID, Status string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		resp := make(chan error)
		writes <- writeReq{op: "update", todo: Todo{ID: body.ID, Status: body.Status, Description: ""}, resp: resp}
		if err := <-resp; err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteHandler(writes chan<- writeReq) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct{ ID string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		resp := make(chan error)
		writes <- writeReq{op: "delete", todo: Todo{ID: body.ID}, resp: resp}
		if err := <-resp; err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
