package main

import (
	"fmt"
	"strconv"
	"./bklexer"
)

func main() {
	fmt.Println("Test Code:")
	code := "声明 变量 = PI * 100 - fda\n1024 * 4 * 3.14 ### \n123"
	fmt.Println(code)
	fmt.Println("--------------------------------")

	lexer := BKLexer.NewLexer()
	lexer.AddRule("\\d+\\.\\d*", "FLOAT")
	lexer.AddRule("\\d+", "INT")
	lexer.AddRule("[\\p{L}\\d_]+", "NAME")
	lexer.AddRule("\\+", "PLUS")
	lexer.AddRule("\\-", "MINUS")
	lexer.AddRule("\\*", "MUL")
	lexer.AddRule("/", "DIV")
	lexer.AddRule("=", "ASSIGN")
	lexer.AddRule("#[^\\r\\n]*", "COMMENT")
	lexer.AddIgnores("[ \\f\\t]+")

	lexer.Build(code)
	for true {
		token := lexer.NextToken()
		if (token.TType != BKLexer.TOKEN_TYPE_EOF) {
			fmt.Printf("%s\t%s\tt%d\t%d\t%d,%d\n",
			token.Name, strconv.Quote(token.Source), token.TType, token.Position, token.Row, token.Col)
		}
		if (token.TType == BKLexer.TOKEN_TYPE_EOF || token.TType == BKLexer.TOKEN_TYPE_ERROR) {
			break
		}
	}
}
