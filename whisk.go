package whisk

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	dataDir       = "data"
	logFilePath   = dataDir + "/log.whisk"
	cleanInterval = time.Minute * 15
	port          = "8081"
)

var confirmPage *template.Template // confirmation page

// Launch starts server, initializes log, register, etc.
func Launch() {
	initDataDir()
	initLog()
	initHome()
	initPages()
	initRegister()
	startSweep()
	startServer()
}

// Cleanup sweeps/writes the page register
func Cleanup() {
	sweepRegister()
	writeRegister()
}

// initDataDir creates data directory if it doesn't exist
func initDataDir() {
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %s\n", err)
	}
}

// initLog creates log file if it doesn't exist, sets output
func initLog() {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error initializing log file: %s\n", err)
	}
	log.SetOutput(f)
}

// initHome initializes templates used on homepage
func initHome() {
	var err error
	confirmPage, err = template.ParseFiles("resources/confirm-template.html")
	if err != nil {
		log.Fatalf("Error initializing template: %s\n", err)
	}
}

// startSweep will sweep/write the register on given interval
func startSweep() {
	ticker := time.NewTicker(cleanInterval)
	go func() {
		for range ticker.C {
			Cleanup()
		}
	}()
}

// startServer sets up http routes, then listens/serves
func startServer() {
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// HTML forms use '\r\n' instead of '\n'
			r.ParseForm()
			body := strings.Replace(r.Form["body"][0], "\r\n", "\n", -1)

			// make new page, redirect to confirmation
			newPage, err := createPage([]byte(body))
			if err != nil {
				http.ServeFile(w, r, "resources/error.html")
				return
			}
			confirmPage.Execute(w, newPage)
		} else {
			http.ServeFile(w, r, "resources/404.html")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Path[1:]
		if len(address) == 0 {
			http.ServeFile(w, r, "resources/home.html")
			return
		}
		if len(address) != idLength || !pageInRegister(address) {
			http.ServeFile(w, r, "resources/404.html")
			return
		}
		http.ServeFile(w, r, pagesDir+address+".html")
	})

	log.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
