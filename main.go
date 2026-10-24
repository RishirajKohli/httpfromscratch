package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(lines)
		data := make([]byte, 8)
		curLine := ""
		for {
			n, err := f.Read(data)

			if err != nil {
				if err == io.EOF {
					lines <- curLine
					break
				}

				log.Fatal(err)
			}
			parts := strings.Split(string(data[:n]), "\n")
			// fmt.Println(len(parts))
			if len(parts) == 1 {
				curLine += parts[0]
				continue
			}
			lines <- curLine + parts[0]
			for i := 1; i < len(parts)-1; i++ {
				lines <- parts[i]
			}
			curLine = parts[len(parts)-1]
		}
	}()

	return lines
}

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
