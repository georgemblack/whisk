package whisk

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const dataFilePath = "register.whisk"

var register map[string]Page

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
	register[page.id] = page
}

// removeFromRegister a single page matching id
func removeFromRegister(id string) {
	delete(register, id)
}

// pageInRegister returns true if id is in register
func pageInRegister(id string) bool {
	_, ok := register[id]
	return ok
}

// writeRegister to file
func writeRegister() {
	output, err := os.Create(dataFilePath)
	if err != nil {
		log.Printf("Error creating %s: %s\n", dataFilePath, err)
		return
	}
	for _, v := range register {
		data := []byte(v.id + "|" + strconv.FormatInt(v.expiration, 10) + "\n")
		output.Write(data)
	}
	output.Close()
}

// sweepRegister removes pages that have expired
func sweepRegister() {
	currTime := time.Now().Unix()
	for _, v := range register {
		if currTime > v.expiration {
			removePage(v.id)
			removeFromRegister(v.id)
		}
	}
}
