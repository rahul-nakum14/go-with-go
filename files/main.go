package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	var (
		filename = flag.String("file", "", "file to write to")
		append   = flag.Bool("append", false, "append instead of replace")
	)
	flag.Parse()

	if *filename == "" {
		log.Fatal("usage: -file <filename> [-append]")
	}

	flags := os.O_WRONLY | os.O_CREATE
	if *append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC // replace
	}

	file, err := os.OpenFile(*filename, flags, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatalf("write error: %v", err)
	}

	if err := writer.Flush(); err != nil {
		log.Fatalf("flush error: %v", err)
	}
}
