package whisk

import (
	"bytes"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"
)

const (
	pagesDir     = dataDir + "/pages/"
	idCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idLength     = 6
)

// page represents single page
type page struct {
	ID         string
	Expiration int64 // Unix timestamp
}

// page content
type pageData struct {
	Title string
	Body  string
	Theme string
}

var minifier *minify.M
var pageTemplate *template.Template

// initPages
func initPages() {

	// create pages dir if it doesn't exist
	err := os.MkdirAll(pagesDir, 0755)
	if err != nil {
		log.Fatalf("Error creating pages dir: %s\n", err)
	}

	// parse page template
	pageTemplate, err = template.ParseFiles("resources/page-template.html")
	if err != nil {
		log.Fatalf("Error initializing page template: %s\n", err)
	}

	// setup minifier
	minifier = minify.New()
	minifier.AddFunc("text/html", html.Minify)
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
	var buf bytes.Buffer

	// create new page struct
	offset := time.Hour * 1440 // 60 days
	newPage = page{
		ID:         generatePageID(),
		Expiration: time.Now().Add(offset).Unix(),
	}

	// safely convert markdown to html
	htmlUnsafe := blackfriday.Run(source)
	htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(htmlUnsafe)

	// execute template, store in buf
	pageTemplate.Execute(&buf, pageData{Title: "Whisk Page", Body: string(htmlSafe), Theme: "minimal"})

	// minify data
	htmlMin, err := minifier.Bytes("text/html", buf.Bytes())
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return newPage, err
	}

	// write to file
	path := pagesDir + newPage.ID + ".html"
	err = ioutil.WriteFile(path, htmlMin, 0644)
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return newPage, err
	}

	addToRegister(newPage) // register new page
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
