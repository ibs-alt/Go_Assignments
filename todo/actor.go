package todo

import (
	"context"
)

// message types
type readReq struct {
	resp chan []Todo // actor copies todos and replies
}
type writeReq struct {
	op   string     // "add" | "update" | "delete"
	todo Todo       // payload
	resp chan error // success / failure
}

// actor starter
func StartActor(ctx context.Context, filename string, initial []Todo) (chan<- readReq, chan<- writeReq) {
	reads := make(chan readReq)
	writes := make(chan writeReq)

	go func() {
		todos := initial // actor owns the slice
		for {
			select {
			case r := <-reads:
				cp := make([]Todo, len(todos)) // defensive copy for readers
				copy(cp, todos)
				r.resp <- cp

			case w := <-writes:
				switch w.op {
				case "add":
					todos = append(todos, w.todo)
				case "update":
					for i := range todos {
						if todos[i].ID == w.todo.ID {
							if w.todo.Status != "" {
								todos[i].Status = w.todo.Status // update status
							}
						}
					}
				case "delete":
					for i := range todos {
						if todos[i].ID == w.todo.ID {
							todos = append(todos[:i], todos[i+1:]...)
							break
						}
					}
				}
				w.resp <- SaveTodos(ctx, filename, todos) // persist safely

			case <-ctx.Done(): // graceful exit
				_ = SaveTodos(ctx, filename, todos)
				return
			}
		}
	}()

	return reads, writes
}
