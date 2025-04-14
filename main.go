package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	books   = []Book{}
	nextID  = 1
	bookMux sync.Mutex
)


func listBooks(w http.ResponseWriter, r *http.Request) {
	bookMux.Lock()
	defer bookMux.Unlock()
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b Book
	json.NewDecoder(r.Body).Decode(&b)
	bookMux.Lock()
	b.ID = nextID
	nextID++
	books = append(books, b)
	bookMux.Unlock()
	json.NewEncoder(w).Encode(b)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	bookMux.Lock()
	defer bookMux.Unlock()
	for i, b := range books {
		if fmt.Sprint(b.ID) == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var updated Book
	json.NewDecoder(r.Body).Decode(&updated)

	bookMux.Lock()
	defer bookMux.Unlock()

	for i, b := range books {
		if b.ID == updated.ID {
			books[i] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}


func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.HandleFunc("/", homeHandler) // handle "/" route
	fmt.Println("Server is running on http://localhost:8080")
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listBooks(w, r)
		case "POST":
			createBook(w, r)
		case "DELETE":
			deleteBook(w, r)
		case "PUT":
			updateBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Book API!"))
}