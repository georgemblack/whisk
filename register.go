package whisk

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dataFilePath = "register.whisk"

var pageRegister map[string]Page

// initializeRegister allocates memory for page register,
// imports existing pages
func initializeRegister() {
	pageRegister = make(map[string]Page)

	// open register file, if it exists
	input, err := os.Open(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file doesn't exist")
			return // file doesn't exist, we're done
		}
		fmt.Println("other error")
		// other error
		return
	}
	defer input.Close()

	// register each listed page
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "|")
		exp, _ := strconv.ParseInt(line[1], 10, 64)
		addToRegister(Page{
			id:         line[0],
			expiration: exp,
		})
	}
}

// addToRegister a single page object
func addToRegister(page Page) {
	pageRegister[page.id] = page
}

// removeFromRegister a single page matching id
func removeFromRegister(id string) {
	delete(pageRegister, id)
}

// pageInRegister returns true if id is in register
func pageInRegister(id string) bool {
	_, ok := pageRegister[id]
	return ok
}

// writeRegister to file
func writeRegister() {
	output, err := os.Create(dataFilePath)
	if err != nil {
		return
	}
	for _, v := range pageRegister {
		data := []byte(v.id + "|" + strconv.FormatInt(v.expiration, 10) + "\n")
		output.Write(data)
	}
	output.Close()
}

// sweepRegister removes pages that have expired
func sweepRegister() {
	currTime := time.Now().Unix()
	for _, v := range pageRegister {
		if currTime > v.expiration {
			removePage(v.id)
			removeFromRegister(v.id)
		}
	}
}
