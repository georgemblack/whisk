package whisk

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"
)

const (
	idCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idLength     = 6
	pagesDir     = "pages/"
)

// Page represents single page
type Page struct {
	id         string
	expiration int64 // Unix timestamp
}

// page content
type pageData struct {
	Title string
	Body  string
}

// createPageFromFile
func createPageFromFile(sourcePath string) {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return
	}
	createPage(source)
}

// createPage from source
func createPage(source []byte) {
	// safely convert to html
	htmlUnsafe := blackfriday.Run(source)
	htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(htmlUnsafe)

	// register
	offset := time.Minute * 2
	page := Page{
		id:         generatePageID(),
		expiration: time.Now().Add(offset).Unix(),
	}
	addToRegister(page)

	// create new file
	output, err := os.Create(pagesDir + page.id + ".html")
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return
	}
	defer output.Close()

	// build template and write
	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		log.Printf("Error parsing template: %s\n", err)
		return
	}
	tmpl.Execute(output, pageData{Title: "Whisk Page", Body: string(htmlSafe)})
}

// removePage from file system and register
func removePage(id string) {
	err := os.Remove(pagesDir + id + ".html")
	if err != nil {
		log.Printf("Error removing page: %s\n", err)
	}
}

// generatePageID that is unique
func generatePageID() string {
	bytes := make([]byte, idLength)
	id := ""
	for id == "" || pageInRegister(id) {
		for i := range bytes {
			bytes[i] = idCharacters[rand.Intn(len(idCharacters))]
		}
		id = string(bytes)
	}
	return id
}
