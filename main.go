package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"Name"`
	LastName string `json:"LastName"`
}

var users = []*User{}

func main() {
	server := createServer(":8080")
	fail := server.ListenAndServe()
	if fail != nil {
		panic(fail)
	}
}

func createServer(address string) *http.Server {
	initRoutes()
	return &http.Server{
		Addr: address,
	}
}

func initRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUser(w, r)
		case http.MethodPost:
			addUser(w, r)
		case http.MethodPut:
			updateUser(w, r)
		case http.MethodDelete:
			deleteUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "NOT ALLOWED")
			return
		}
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey, go to /users")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%v", users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	fail := json.NewDecoder(r.Body).Decode(user)
	if fail != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", fail)
		return
	}
	users = append(users, user)
	fmt.Fprintf(w, "User added succesfully")

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	nf := true
	fail := json.NewDecoder(r.Body).Decode(user)
	if fail != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", fail)
		return
	}
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
			nf = false
			fmt.Fprintf(w, "User updated succesfully")
		}
	}
	if nf {
		fmt.Fprintf(w, "User not found")
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	nf := true
	fail := json.NewDecoder(r.Body).Decode(user)
	if fail != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", fail)
		return
	}
	for i, u := range users {
		if u.ID == user.ID {
			users = append(users[:i], users[i+1:]...)
			nf = false
			fmt.Fprintf(w, "User deleted succesfully")
		}
	}
	if nf {
		fmt.Fprintf(w, "User not found")
	}
}
