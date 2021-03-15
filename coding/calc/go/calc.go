package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "./bklexer"
)

type Node interface {
    GetValue() float64
}

type Number struct {
    value float64
}

func NewNumber(token *BKLexer.Token) *Number {
    value, _ := strconv.ParseFloat(token.Source, 64)
    return &Number{value: value}
}

func (number *Number) GetValue() float64 {
    return number.value
}

type BinaryOpt struct {
    opt string
    lhs Node
    rhs Node
}

func NewBinaryOpt(token *BKLexer.Token, lhs Node, rhs Node) *BinaryOpt {
    return &BinaryOpt{opt: token.Source, lhs: lhs, rhs: rhs}
}

func (binaryOpt *BinaryOpt) GetValue() float64 {
    lhs, rhs := binaryOpt.lhs, binaryOpt.rhs
    switch binaryOpt.opt {
        case "+": return lhs.GetValue() + rhs.GetValue()
        case "-": return lhs.GetValue() - rhs.GetValue()
        case "*": return lhs.GetValue() * rhs.GetValue()
        case "/": return lhs.GetValue() / rhs.GetValue()
    }
    return 0
}

func parse(lexer *BKLexer.Lexer) Node {
    result := parse_binary_add(lexer)
    token := lexer.GetToken()
    if token.TType != BKLexer.TOKEN_TYPE_EOF {
        return nil
    }
    return result
}

func parse_binary_add(lexer *BKLexer.Lexer) Node {
    lhs := parse_binary_mul(lexer)
    if lhs == nil {
        return nil
    }
    token := lexer.GetToken()
    for token.Source == "+" || token.Source == "-" {
        rhs := parse_binary_mul(lexer)
        if rhs == nil {
            return nil
        }
        lhs = NewBinaryOpt(token, lhs, rhs)
        token = lexer.GetToken()
    }
    return lhs
}

func parse_binary_mul(lexer *BKLexer.Lexer) Node {
    lhs := parse_number(lexer)
    if lhs == nil {
        return nil
    }
    token := lexer.GetToken()
    for token.Source == "*" || token.Source == "/" {
        rhs := parse_number(lexer)
        if rhs == nil {
            return nil
        }
        lhs = NewBinaryOpt(token, lhs, rhs)
        token = lexer.GetToken()
    }
    return lhs
}

func parse_number(lexer *BKLexer.Lexer) Node {
    token := lexer.NextToken()
    if token.Name == "LPAR" {
        expr := parse_binary_add(lexer)
        if expr == nil {
            return nil
        }
        token := lexer.GetToken()
        if token.Name != "RPAR" {
            return nil
        }
        lexer.NextToken()
        return expr
    }
    if token.Name == "NUMBER" {
        number := NewNumber(token)
        lexer.NextToken()
        return number
    }
    return nil
}

func main() {
    fmt.Println("Hello My Calc.")

    lexer := BKLexer.NewLexer()
    lexer.AddRule("\\d+\\.?\\d*", "NUMBER")
    lexer.AddRule("\\+", "PLUS")
    lexer.AddRule("-", "MINUS")
    lexer.AddRule("\\*", "MUL")
    lexer.AddRule("/", "DIV")
    lexer.AddRule("\\(", "LPAR")
    lexer.AddRule("\\)", "RPAR")
    lexer.AddIgnores("[ \\f\\f]+")

    reader := bufio.NewReader(os.Stdin)
    for true {
        fmt.Print("> ")
        inputs, _ := reader.ReadString('\n')
        inputs = strings.Trim(inputs, " \f\t\n")
        if inputs == "quit" {
            break
        }
        if inputs != "" {
            lexer.Build(inputs)
            result := parse(lexer)
            if result == nil {
                positon := lexer.GetToken().Col
                fmt.Println("error in :", positon)
                continue
            }
            fmt.Println("out =", result.GetValue())
        }

    }

    fmt.Println("bye!")
}
