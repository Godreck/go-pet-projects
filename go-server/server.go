package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// simpleData struct {name string, fArtist string} for storing data about user.
type simpleData struct {
	name    string
	fArtist string
}

// helloHandler func(w http.ResponseWriter, r *http.Request) stands for test how it's possible to move on to the different static pages. Just jun)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	} else {
		fmt.Fprintf(w, "Hello, stranger!")
	}
}

// formHandler func(w http.ResponseWriter, r *http.Request) takes data from form.html and greets the user
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./static/form.html")
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "PareseForm() error: %v", err)
			return
		}
		uData := simpleData{name: r.FormValue("name"), fArtist: r.FormValue("fArtist")}

		fmt.Fprintf(w, "Hello, %v! I heard u like %v", uData.name, uData.fArtist)
	}
}

func main() {
	//static web server. It will take resoursers from http.Dir(path)
	fileServer := http.FileServer(http.Dir("./static"))
	//registrate handlers
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)
	//set up server settings
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//lounching server
	fmt.Println("Starting server at port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Fatal error starting the server:", err)
	}
}
