package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data.mobi")
	if err != nil {
		log.Fatal(err)
	}
	format, err := readPalmDatabaseFormat(file)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(*format, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pdb format: %s\n", data)

	pdHeader, err := readPalmDocHeader(file)
	if err != nil {
		log.Fatal(err)
	}
	data, err = json.MarshalIndent(*pdHeader, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("palmdoc header %s\n", data)

	mHeader, err := readMobiHeader(file)
	if err != nil {
		log.Fatal(err)
	}
	data, err = json.MarshalIndent(*mHeader, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mobi header %s\n", data)

	exthHeader, err := readExthHeader(file)
	if err != nil {
		log.Fatal(err)
	}
	data, err = json.MarshalIndent(*exthHeader, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exth header %s\n", data)
}
