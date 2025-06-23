package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// helper: create isolated actor + HTTP server
func localServer() (*httptest.Server, chan<- readReq, chan<- writeReq) {
	ctx := context.WithValue(context.Background(), TraceIDKey, "parallel-test")

	reads, writes := StartActor(ctx, "inmem.json", nil)

	mux := http.NewServeMux()
	mux.Handle("/create", CreateHandler(reads, writes))
	mux.Handle("/get", GetHandler(reads))
	mux.Handle("/update", UpdateHandler(reads, writes))
	mux.Handle("/delete", DeleteHandler(writes))

	return httptest.NewServer(mux), reads, writes
}

// Parent test groups sub-tests; each sub-test runs truly in parallel

func TestAPIParallelScenarios(t *testing.T) {

	// sub-test: 100 concurrent creates
	t.Run("create-100", func(t *testing.T) {
		t.Parallel()
		ts, reads, _ := localServer()
		defer ts.Close()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				body := fmt.Sprintf(`{"description":"task %d"}`, n)
				http.Post(ts.URL+"/create", "application/json", strings.NewReader(body))
			}(i)
		}
		wg.Wait()

		// verify length == 100
		rq := readReq{resp: make(chan []Todo)}
		reads <- rq
		if got := len(<-rq.resp); got != 100 {
			t.Fatalf("want 100 tasks, got %d", got)
		}
	})

	// sub-test: concurrent update storm
	t.Run("update-burst", func(t *testing.T) {
		t.Parallel()
		ts, reads, _ := localServer()
		defer ts.Close()

		// seed one task
		http.Post(ts.URL+"/create", "application/json",
			strings.NewReader(`{"description":"demo"}`))

		// grab its ID
		res, _ := http.Get(ts.URL + "/get")
		var list []Todo
		_ = json.NewDecoder(res.Body).Decode(&list)
		id := list[0].ID

		statuses := []string{"Started", "Completed", "Not started"}
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			for _, st := range statuses {
				wg.Add(1)
				go func(s string) {
					defer wg.Done()
					body := fmt.Sprintf(`{"id":"%s","status":"%s"}`, id, s)
					http.Post(ts.URL+"/update", "application/json", strings.NewReader(body))
				}(st)
			}
		}
		wg.Wait()

		// still exactly one task
		rq := readReq{resp: make(chan []Todo)}
		reads <- rq
		if len(<-rq.resp) != 1 {
			t.Fatalf("task duplicated during updates")
		}
	})

	// sub-test: concurrent deletes
	t.Run("delete-all", func(t *testing.T) {
		t.Parallel()
		ts, reads, writes := localServer()
		defer ts.Close()

		// create 10 tasks
		for i := 0; i < 10; i++ {
			body := fmt.Sprintf(`{"description":"d%d"}`, i)
			http.Post(ts.URL+"/create", "application/json", strings.NewReader(body))
		}

		// fetch ids
		res, _ := http.Get(ts.URL + "/get")
		var todos []Todo
		_ = json.NewDecoder(res.Body).Decode(&todos)

		var wg sync.WaitGroup
		for _, td := range todos {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				body := fmt.Sprintf(`{"id":"%s"}`, id)
				http.Post(ts.URL+"/delete", "application/json", strings.NewReader(body))
			}(td.ID)
		}
		wg.Wait()

		// ensure slice empty via read channel
		rq := readReq{resp: make(chan []Todo)}
		reads <- rq
		if len(<-rq.resp) != 0 {
			t.Fatalf("tasks not fully deleted")
		}

		// also verify write-only channel works compile-time
		wr := writeReq{op: "add", todo: NewTodo("x"), resp: make(chan error)}
		writes <- wr
		<-wr.resp
	})
}
