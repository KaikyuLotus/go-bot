package gobot

import (
	"io/ioutil"
	"fmt"
)

func readFileBytes(fileName string) []byte {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func readFileString(fileName string) string {
	return string(readFileBytes(fileName))
}