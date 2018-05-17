package gobot

import (
	"io/ioutil"
	"fmt"
)

func ReadFileBytes(fileName string) []byte {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func ReadFileString(fileName string) string {
	return string(ReadFileBytes(fileName))
}