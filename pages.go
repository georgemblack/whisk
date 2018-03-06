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
	pagesDir     = "data/pages/"
)

// Page represents single page
type page struct {
	ID         string
	Expiration int64 // Unix timestamp
}

// page content
type pageData struct {
	Title string
	Body  string
}

// createPageFromFile
func createPageFromFile(sourcePath string) (page, error) {
	var newPage page

	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return newPage, err
	}
	newPage, err = createPage(source)
	return newPage, err
}

// createPage from source
func createPage(source []byte) (page, error) {
	var newPage page

	// safely convert to html
	htmlUnsafe := blackfriday.Run(source)
	htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(htmlUnsafe)

	// register
	offset := time.Minute * 2
	newPage = page{
		ID:         generatePageID(),
		Expiration: time.Now().Add(offset).Unix(),
	}
	addToRegister(newPage)

	// create new file
	output, err := os.Create(pagesDir + newPage.ID + ".html")
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return newPage, err
	}
	defer output.Close()

	// build template and write
	tmpl, err := template.ParseFiles("resources/page-template.html")
	if err != nil {
		log.Printf("Error parsing template: %s\n", err)
		return newPage, err
	}
	tmpl.Execute(output, pageData{Title: "Whisk Page", Body: string(htmlSafe)})
	return newPage, nil
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
