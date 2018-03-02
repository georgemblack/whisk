package whisk

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	logFilePath = "log.whisk"
	port        = "8081"
)

// Launch starts server, initializes log, register, etc.
func Launch() {

	initializeLog()
	initializeRegister()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		createPage("./sample.md")
	}) // temporary
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		writeRegister()
	}) // temporary
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Path[1:]
		if len(address) == 0 {
			http.ServeFile(w, r, "./index.html")
			return
		}
		if len(address) != idLength || !pageInRegister(address) {
			http.ServeFile(w, r, "./404.html")
			return
		}
		http.ServeFile(w, r, pagesDir+address+".html")
	})
	log.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initializeLog() {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error initializing log file: %s\n", err)
	}
	log.SetOutput(f)
}
