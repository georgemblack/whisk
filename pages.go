package whisk

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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

// createPage from markdown file
func createPage(sourcePath string) {
	// read sample file
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return
	}

	// safely convert to html
	htmlUnsafe := blackfriday.Run(source)
	htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(htmlUnsafe)

	// register
	offset := time.Hour * 1
	page := Page{
		id:         generatePageID(),
		expiration: time.Now().Add(offset).Unix(),
	}
	addToRegister(page)

	// write new file
	output, err := os.Create(pagesDir + page.id + ".html")
	if err != nil {
		log.Printf("Error creating page: %s\n", err)
		return
	}
	output.Write(htmlSafe)
	output.Close()
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
