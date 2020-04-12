package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var state []Todo

func main() {
	r := mux.NewRouter()

	//routes

	r.HandleFunc("/", handleRoot).Methods("GET")

	r.HandleFunc("/todos", handleAddTodo).Methods("POST")

	r.HandleFunc("/todos/{todoId}", handleUpdateTodo).Methods("PUT")

	r.HandleFunc("/todos/{todoId}", handleDeleteTodo).Methods("DELETE")

	fmt.Println("server started, listening on port :80")
	if err := http.ListenAndServe(":80", r); err != nil {
		panic(err)
	}
}

//handleroot loads the root page.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	listOfJSON, _ := json.Marshal(state)

	w.Write(listOfJSON)
}

//handleAddTodo adds a new todo.
func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	json.NewDecoder(r.Body).Decode(&todo)
	state = append(state, todo)

	todoJSON, _ := json.Marshal(todo)
	w.Write(todoJSON)
}

//handleUpdateTodo updates a todo.
func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {

}

//handleDeleteTodo deletes a todo.
func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {

}

//Todo type
type Todo struct {
	Body string `json:"body"`
}
