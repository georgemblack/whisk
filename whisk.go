package whisk

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	logFilePath   = "log.whisk"
	cleanInterval = time.Minute * 15
	port          = "8081"
)

// Launch starts server, initializes log, register, etc.
func Launch() {

	initializeLog()
	initializeRegister()
	initializeTimer()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		createPage("./sample.md")
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

// initializeLog creates log file if it doesn't exist, sets output
func initializeLog() {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error initializing log file: %s\n", err)
	}
	log.SetOutput(f)
}

// initializeTimer will sweep/write the register on given interval
func initializeTimer() {
	ticker := time.NewTicker(cleanInterval)
	go func() {
		for range ticker.C {
			sweepRegister()
			writeRegister()
		}
	}()
}
