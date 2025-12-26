// package main

// import (
// 	"fmt"
// 	"os"
// )

// func main(){

// 	f,err := os.Open("example.txt")

// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}	

// 	defer f.Close()
// // 	fileInfo,err := f.Stat()
// // 	if err != nil {
// // 		fmt.Println("Error getting file info:", err)
// // 		return
// // 	}

// // 	fmt.Println("File Name:", fileInfo.Name())
// // 	fmt.Println("File Size:", fileInfo.Size(), "bytes")
// // 	fmt.Println("Last Modified:", fileInfo.Mode())

// 	data := make([]byte, 100)
// 	bytesRead, err := f.Read(data)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}	
// 	fmt.Printf("Read %d bytes: %s\n", bytesRead, string(data[:bytesRead]))
// }



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
