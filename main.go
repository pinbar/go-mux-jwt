package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := ConfigureRouter()
	log.Fatal(http.ListenAndServe(":3001", handlers.LoggingHandler(os.Stdout, router)))
}

//ConfigureRouter setup the router
func ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/metacortex", metacortexHandler)
	router.HandleFunc("/agents/{name}", agentsHandler)

	router.HandleFunc("/authenticate", authenticate)

	router.Handle("/api/megacity", authMiddleware(megacityHandler))

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Matrix!"))
}
func metacortexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mr Anderson's not so secure workplace!"))
}
func agentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("My name is agent " + vars["name"]))
}

var megacityHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Megacity!"))
})
