package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	//"strings"
	//"golang.org/x/net/html"
)

var (
	startRegex = regexp.MustCompile(`<DOCUMENT>`)
	endRegex   = regexp.MustCompile(`</DOCUMENT>`)
	typeRegex  = regexp.MustCompile(`<TYPE>[^\n]+`)
	//itemRegex  = regexp.MustCompile(`(>Item(\s|&#160;|&nbsp;)(1A|1B|7A|7|8)\\.{0,1})|(ITEM\s(1A|1B|7A|7|8))`)
	itemRegex = regexp.MustCompile(`(>(Item|ITEM)(\s|&#160;|&nbsp;)(8|9)\.)`)

	opsRegex  = regexp.MustCompile(`(?i)(CONSOLIDATED\s)?(STATEMENTS? OF\s)(OPERATIONS|EARNINGS)`)
	incRegex  = regexp.MustCompile(`(?i)(CONSOLIDATED\s)?(STATEMENTS? OF\s)(COMPREHENSIVE INCOME)`)
	balRegex  = regexp.MustCompile(`(?i)(CONSOLIDATED\s)(BALANCE SHEETS?)`)
	cashRegex = regexp.MustCompile(`(?i)(CONSOLIDATED\s)?(STATEMENTS? OF\s)(CASH FLOWS?)`)
)

type Document struct {
	Type     string
	Content  string
	ItemsMap map[string]string
}

func process10k(r *http.Response) error {
	defer r.Body.Close()

	var documents []Document
	var currentDoc *Document
	var docBuffer bytes.Buffer
	var inDocument bool

	reader := bufio.NewReader(r.Body)

	count := 0
	for {
		if count >= 10000 {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading response: %w", err)
		}

		if startRegex.MatchString(line) {
			fmt.Printf("*%d match start = %s", count, line)
			inDocument = true
			currentDoc = &Document{
				ItemsMap: make(map[string]string),
			}
			docBuffer.Reset()
		}

		if inDocument {
			docBuffer.WriteString(line)
			if typeRegex.MatchString(line) {
				typePart := typeRegex.FindString(line)
				typeName := strings.TrimPrefix(typePart, "<TYPE>")
				if typeName != "10-K" {
					inDocument = false
					continue
				}
				fmt.Printf("*%d match type = %s", count, line)
				currentDoc.Type = typeName
			}

			if endRegex.MatchString(line) {
				fmt.Printf("*%d match end = %s", count, line)
				inDocument = false
				currentDoc.Content = docBuffer.String()

				extractItems(currentDoc)

				documents = append(documents, *currentDoc)
				currentDoc = nil
				break
			}
		}
		count++
	}
	return nil
}

func extractItems(doc *Document) {
	content := doc.Content
	fmt.Println("cont len:", len(content))

	itemMatches := itemRegex.FindAllStringIndex(content, -1)
	if len(itemMatches) == 0 {
		return
	}
	fmt.Println("match slice:", itemMatches)

	for _, match := range itemMatches {
		fmt.Printf("%v - %s\n", match, content[match[0]:match[1]+80])
		fmt.Println()
	}

	tbls := make(map[string][][]int)
	for i := 2; i < len(itemMatches); i = i + 2 {
		stmts := content[itemMatches[i][1]:itemMatches[i+1][0]]
		opsMatches := opsRegex.FindAllStringIndex(stmts, -1)
		if len(opsMatches) > 0 {
			tbls["OpsEarnings"] = opsMatches
		}
	}
	fmt.Printf("Table headers: %v\n", tbls)

}
