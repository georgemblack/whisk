package whisk

import (
	"log"
	"net/http"
)

const (
	idCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idLength     = 6
	pagesDir     = "pages/"
)

// Launch starts server, initializes register, etc.
func Launch() {
	initializeRegister()
	defer writeRegister()

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
	log.Fatal(http.ListenAndServe(":8081", nil))
}
