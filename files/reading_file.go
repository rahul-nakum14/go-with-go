
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func mainx() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <filename>", os.Args[0])
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					fmt.Print(line)
				}
				break
			}
			log.Fatalf("read error: %v", err)
		}

		fmt.Print(line)
	}
}
