package BKLexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Token struct {
	Position int
	Source string
	TType int
	Row int
	Col int
	Name string
}

type Lexer struct {
	code string
	position int
	source string
	ttype int
	row int
	col int
	exprs []*regexp.Regexp
	names []string
    reserves []string
	ignores *regexp.Regexp
}

const (
	TOKEN_TYPE_ERROR = -2
	TOKEN_TYPE_EOF = -1
	TOKEN_TYPE_NEWLINE = 0
)

func NewToken(positon int, source string, ttype int, row int, col int, name string) *Token {
	token := new(Token)
	token.TType = ttype
	token.Name = name
	token.Source = source
	token.Position = positon
	token.Row = row
	token.Col = col
	return token
}

func NewLexer() *Lexer {
	lexer := new(Lexer)
	lexer.AddRule("\n", "NEWLINE")
	return lexer
}

func (lexer *Lexer) AddRule(expr string, name string) {
	rule, err := regexp.Compile("^" + expr)
	if (err != nil) {
		fmt.Println(err)
	} else {
		lexer.exprs = append(lexer.exprs, rule)
		lexer.names = append(lexer.names, name)
	}
}

func (lexer *Lexer) AddIgnores(expr string) {
    if lexer.ignores != nil {
        expr = lexer.ignores.String() + "|(^" + expr + ")"
    } else {
        expr = "(^" + expr + ")"
    }
	rule, err := regexp.Compile(expr)
	if (err != nil) {
		fmt.Println(err)
	} else {
		lexer.ignores = rule
	}
}

func (lexer *Lexer) AddReserve(reserve string) {
    lexer.reserves = append(lexer.reserves, reserve);
}

func (lexer *Lexer) Build(code string) {
	lexer.code = code
	lexer.position = 0
	lexer.source = ""
	lexer.ttype = 0
	lexer.row  = 0
	lexer.col  = 0
}

func (lexer *Lexer) GetTokenName() string {
	switch lexer.ttype {
		case -1: return "EOF"
		case -2: return "ERROR"
		default: return lexer.names[lexer.ttype]
	}
}

func (lexer *Lexer) GetToken() *Token {
    name := lexer.GetTokenName()
    ttype := lexer.ttype
    if name == "NAME" {
        for i, reserved := range lexer.reserves {
            if lexer.source == reserved {
                name = strings.ToUpper(reserved)
                ttype = len(lexer.names) + i
                break
            }
        }
    }
	return NewToken(
		lexer.position, lexer.source, ttype,
		lexer.row, lexer.col, name)
}

func (lexer *Lexer) NextToken() *Token {
	lexer.position += len(lexer.source);
	if lexer.source != "\n" {
		lexer.col += utf8.RuneCountInString(lexer.source)
	} else {
		lexer.row += 1
		lexer.col = 0
	}

    for true {
        source := lexer.ignores.FindString(lexer.code[lexer.position:])
        if source == "" {
            break
        }
        lexer.position += len(source)
        lexer.col += utf8.RuneCountInString(source)
    }
	if lexer.position >= len(lexer.code) {
		lexer.ttype = -1
		return lexer.GetToken()
	}
	lexer.ttype = -2
	for i, rule := range lexer.exprs {
		source := rule.FindString(lexer.code[lexer.position:])
		if source != "" {
			lexer.source = source;
			lexer.ttype = i
			break
		}
	}
	return lexer.GetToken()
}
