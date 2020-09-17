package cmd

import (
	"io/ioutil"
	"strings"
)

//GETAPIID Gets the API ID
func GETAPIID() string {

	path := originalDir + "/ID.txt"

	file := readFile(path)

	return "Client-ID " + file[0]
}

func readFile(path string) []string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		//Do something
	}
	return strings.Split(string(content), "\n")
}
