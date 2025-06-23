package todo

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

// helpers

func freshActor() (reads chan<- readReq, writes chan<- writeReq, cleanup func()) {
	ctx := context.WithValue(context.Background(), TraceIDKey, "unit-test")
	FileName = "unit_data.json"
	_ = os.Remove(FileName)

	reads, writes = StartActor(ctx, FileName, nil)
	cleanup = func() { _ = os.Remove(FileName) }
	return
}

func ask(reads chan<- readReq) []Todo {
	req := readReq{resp: make(chan []Todo)}
	reads <- req
	return <-req.resp
}

func add(w chan<- writeReq, desc string) {
	wr := writeReq{op: "add", todo: NewTodo(desc), resp: make(chan error)}
	w <- wr
	<-wr.resp
}

// tests

func TestActorAddUpdateDeletePersist(t *testing.T) {
	reads, writes, cleanup := freshActor()
	defer cleanup()

	// add
	add(writes, "alpha")
	todos := ask(reads)
	if len(todos) != 1 || todos[0].Description != "alpha" {
		t.Fatalf("add failed %+v", todos)
	}
	id := todos[0].ID

	// update
	wr := writeReq{op: "update", todo: Todo{ID: id, Status: "Completed"}, resp: make(chan error)}
	writes <- wr
	if err := <-wr.resp; err != nil {
		t.Fatalf("update err: %v", err)
	}
	if ask(reads)[0].Status != "Completed" {
		t.Fatalf("status not updated")
	}

	// delete
	wr = writeReq{op: "delete", todo: Todo{ID: id}, resp: make(chan error)}
	writes <- wr
	<-wr.resp
	if len(ask(reads)) != 0 {
		t.Fatalf("delete failed")
	}

	// persisted JSON should decode to empty slice
	data, _ := os.ReadFile(FileName)
	var slice []Todo
	if err := json.Unmarshal(data, &slice); err != nil || len(slice) != 0 {
		t.Fatalf("file not empty json: %s (err %v)", data, err)
	}
}
