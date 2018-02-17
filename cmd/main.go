package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Laugusti/mobireader"
)

var filename = flag.String("filename", "data.mobi", "name of mobi file")

func main() {
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	mobi, err := mobireader.Create(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	fmt.Print(mobi.Text())
}
