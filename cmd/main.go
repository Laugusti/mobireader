package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/Laugusti/gopl/ch5/exercise_5_7/htmlprint"
	"github.com/Laugusti/mobireader"
)

var filename = flag.String("filename", "data.mobi", "name of mobi file")

func main() {
	flag.Parse()

	// open file
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	// create MOBIFile from file
	mobi, err := mobireader.Create(file)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	// get text from MOBIFile
	data, err := mobi.Text()
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	// parse text as html
	node, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	// pretty print html
	fmt.Println(htmlprint.PrettyPrint(node))
}
