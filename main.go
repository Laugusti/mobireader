package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Laugusti/mobireader/palmdoc"
)

func main() {
	foo()
}

func foo() {
	file, err := os.Open("data.mobi")
	if err != nil {
		log.Fatal(err)
	}
	mobi, err := Create(file)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(*mobi, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pdb format: %s\n", data)
}

func bar() {
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

	start := format.RecordInfoEntries[1].Offset
	end := format.RecordInfoEntries[2].Offset
	_, err = file.Seek(int64(start), io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, end-start+1)
	_, err = io.ReadFull(file, b)
	if err != nil {
		log.Fatal(err)
	}
	res, err := palmdoc.Decompress(b)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("****RESULT:\n%s\n", res)
	}
}
