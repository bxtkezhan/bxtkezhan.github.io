package main

import (
    "fmt"
    "strconv"
    "io/ioutil"
    "./bklexer"
)

type Node interface {
    GetValue() float64
}

type Block struct {
    statements []Node
}

func NewBlock() *Block {
    return &Block{}
}

func (block *Block) AddStatement(statement Node) {
    block.statements = append(block.statements, statement)
}

func (block *Block) Eval() {
    for i, statement := range block.statements {
        fmt.Printf("out[%d] = %f\n", i, statement.GetValue())
    }
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

func parse(lexer *BKLexer.Lexer) *Block {
    block := NewBlock()
    token := lexer.NextToken()
    for token.TType != BKLexer.TOKEN_TYPE_EOF {
        statement := parse_binary_add(lexer)
        if statement == nil {
            return nil;
        }
        token = lexer.GetToken()
        if token.TType != BKLexer.TOKEN_TYPE_NEWLINE &&
           token.TType != BKLexer.TOKEN_TYPE_EOF {
            return nil;
        }
        block.AddStatement(statement)
        for token.TType == BKLexer.TOKEN_TYPE_NEWLINE {
            token = lexer.NextToken()
        }
    }
    return block
}

func parse_binary_add(lexer *BKLexer.Lexer) Node {
    lhs := parse_binary_mul(lexer)
    if lhs == nil {
        return nil
    }
    token := lexer.GetToken()
    for token.Source == "+" || token.Source == "-" {
        lexer.NextToken()
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
        lexer.NextToken()
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
    token := lexer.GetToken()
    if token.Name == "LPAR" {
        lexer.NextToken()
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
    lexer.AddIgnores("[ \\f\\t]+")
    lexer.AddIgnores("#[^\\r\\n]*")

    bytes, err := ioutil.ReadFile("../test.txt")
    if err != nil {
        fmt.Println("read faild")
        return
    }
    code := string(bytes)
    fmt.Println(code)
    lexer.Build(code)
    result := parse(lexer)
    if result == nil {
        fmt.Println("null result")
        return
    }
    result.Eval()
}
