package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

var state map[uuid.UUID]Todo

func main() {
	r := mux.NewRouter()

	//routes

	r.HandleFunc("/todos", handleGetAllTodos).Methods("GET")

	r.HandleFunc("/todos", handleAddTodo).Methods("POST")

	r.HandleFunc("/todos/{todoId}", handleGetTodo).Methods("GET")

	r.HandleFunc("/todos/{todoId}", handleUpdateTodo).Methods("PUT")

	r.HandleFunc("/todos/{todoId}", handleDeleteTodo).Methods("DELETE")

	fmt.Println("server started, listening on port :80")
	if err := http.ListenAndServe(":80", r); err != nil {
		panic(err)
	}
}

//handleGetAllTodos loads the root page.
func handleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	if len(state) <= 0 {
		state = make(map[uuid.UUID]Todo)
	}

	listOfJSON, _ := json.Marshal(state)

	w.Write(listOfJSON)
}

func handleGetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["todoId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, isPresent := state[id]
	if !isPresent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeResponse(todo, w)
}

//handleAddTodo adds a new todo.
func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	if len(state) <= 0 {
		state = make(map[uuid.UUID]Todo)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var todo Todo
	if err = json.Unmarshal(body, &todo); err != nil {
		w.WriteHeader(400)
		return
	}
	todo.ID, _ = uuid.NewRandom()
	state[todo.ID] = todo

	writeResponse(todo, w)
}

//handleUpdateTodo updates a todo.
func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["todoId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, isPresent := state[id]
	if !isPresent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var update Todo
	if err = json.Unmarshal(body, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	update.ID = todo.ID
	state[todo.ID] = update
	writeResponse(update, w)
}

//handleDeleteTodo deletes a todo.
func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["todoId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, isPresent := state[id]; !isPresent {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(state, id)
	w.WriteHeader(http.StatusNoContent)
}

//Todo type
type Todo struct {
	ID   uuid.UUID `json:"id"`
	Body string    `json:"body"`
}

// helper function to write response
func writeResponse(todo Todo, w http.ResponseWriter) {
	todoJSON, _ := json.Marshal(todo)
	w.Write(todoJSON)
}
