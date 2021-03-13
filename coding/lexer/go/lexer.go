package main

import (
	"fmt"
	"regexp"
	"unicode/utf8"
	"os"
)

var exprs = []string{"\\d+", "[\\p{L}\\d_]+", "[\\+-=]"}
var names = []string{"INT",  "NAME",         "SYMBOL"}

func main() {
	rules := []*regexp.Regexp{}
	for i, expr := range exprs {
		rule, _ := regexp.Compile("^" + expr)
		rules = append(rules, rule)
		fmt.Println(names[i], rule)
	}

	fmt.Println("--------------------------------")
	for row, code := range os.Args[1:] {
		position := 0
		col := 0
		for true {
			for position < len(code) && (code[position] == ' ' || code[position] == '\t') {
				position += 1
				col += 1
			}
			if position >= len(code) {
				break
			}
			source := ""
			tokenType := -1
			for i, rule := range rules {
				source = rule.FindString(code[position:])
				if source != "" {
					tokenType = i
					break
				}
			}
			if tokenType >= 0 {
				fmt.Printf("%s\t`%s`\t%d\t%d\n", names[tokenType], source, row, col)
				position += len(source)
				col += utf8.RuneCountInString(source)
			} else {
				fmt.Printf("error in: %d, %d\n", row, col)
				break
			}
		}
	}

}
