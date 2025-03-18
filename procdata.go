package main

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
	//"net/http"
	"regexp"
	//"strings"
	//"golang.org/x/net/html"
)

var (
	startRegex = regexp.MustCompile("<DOCUMENT>")
)

type Document struct {
	Type     string
	Content  string
	ItemsMap map[string]string
}

func process10k(r io.Reader) error {
	//var documents []Document
	//var currentDoc *Document
	//var docBuffer bytes.Buffer
	//var inDocument bool
	count := 0

	reader := bufio.NewReader(r)

	fmt.Println("got a bufio reader")
	for {
		if count >= 5 {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}

		fmt.Println("*", line)
		fmt.Println("count:", count)
		count++
	}
	return nil
}
