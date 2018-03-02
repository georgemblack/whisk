package whisk

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const dataFilePath = "register.whisk"

var register map[string]Page
var lock = sync.RWMutex{}

// initializeRegister allocates memory for page register,
// imports existing pages
func initializeRegister() {
	register = make(map[string]Page)

	// open register file, if it exists
	input, err := os.Open(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("File %s not detected, initialized empty page register\n", dataFilePath)
			return // file doesn't exist, we're done
		}
		log.Printf("Error opening %s: %s\n", dataFilePath, err)
		return // other error
	}
	defer input.Close()

	// register each listed page
	numPages := 0
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "|")
		exp, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			log.Printf("Error parsing %s: %s", dataFilePath, err)
			return
		}
		addToRegister(Page{
			id:         line[0],
			expiration: exp,
		})
		numPages++
	}
	log.Printf("Initialized page register with %d items\n", numPages)
}

// addToRegister a single page object
func addToRegister(page Page) {
	lock.Lock()
	register[page.id] = page
	lock.Unlock()
}

// removeFromRegister a single page matching id
func removeFromRegister(id string) {
	lock.Lock()
	delete(register, id)
	lock.Unlock()
}

// pageInRegister returns true if id is in register
func pageInRegister(id string) bool {
	lock.RLock()
	_, ok := register[id]
	lock.RUnlock()
	return ok
}

// writeRegister to file
func writeRegister() {
	output, err := os.Create(dataFilePath)
	if err != nil {
		log.Printf("Error creating %s: %s\n", dataFilePath, err)
		return
	}
	lock.RLock()
	for _, v := range register {
		data := []byte(v.id + "|" + strconv.FormatInt(v.expiration, 10) + "\n")
		output.Write(data)
	}
	lock.RUnlock()
	output.Close()
}

// sweepRegister removes pages that have expired
func sweepRegister() {
	currTime := time.Now().Unix()
	ids := []string{} // pages to remove

	// find pages to remove
	lock.RLock()
	for _, v := range register {
		if currTime > v.expiration {
			ids = append(ids, v.id)
		}
	}
	lock.RUnlock()

	// remove them
	for _, id := range ids {
		removePage(id)
		removeFromRegister(id)
	}
}
