package main

import (
    "fmt"
    "strconv"
    "io/ioutil"
    "./bklexer"
)

var ValueDict map[string]float64

type Node interface {
    Eval() float64
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
    for _, statement := range block.statements {
        statement.Eval()
    }
}

type Number struct {
    value float64
}

func NewNumber(token *BKLexer.Token) *Number {
    value, _ := strconv.ParseFloat(token.Source, 64)
    return &Number{value: value}
}

func (number *Number) Eval() float64 {
    return number.value
}

type Name struct {
    name string
}

func NewName(token *BKLexer.Token) *Name {
    return &Name{name: token.Source}
}

func (name *Name) Eval() float64 {
    if value, found := ValueDict[name.name]; found {
        return value;
    }
    return 0.
}

type BinaryOpt struct {
    opt string
    lhs Node
    rhs Node
}

func NewBinaryOpt(token *BKLexer.Token, lhs Node, rhs Node) *BinaryOpt {
    return &BinaryOpt{opt: token.Source, lhs: lhs, rhs: rhs}
}

func (binaryOpt *BinaryOpt) Eval() float64 {
    lhs, rhs := binaryOpt.lhs, binaryOpt.rhs
    switch binaryOpt.opt {
        case "+": return lhs.Eval() + rhs.Eval()
        case "-": return lhs.Eval() - rhs.Eval()
        case "*": return lhs.Eval() * rhs.Eval()
        case "/": return lhs.Eval() / rhs.Eval()
    }
    return 0
}

type Assign struct {
    name string
    value Node
}

func NewAssign(token *BKLexer.Token, value Node) *Assign {
    return &Assign{name: token.Source, value: value}
}

func (assign *Assign) Eval() float64 {
    value := assign.value.Eval()
    ValueDict[assign.name] = value
    return value
}

type Echo struct {
    value Node
}

func NewEcho(value Node) *Echo {
    return &Echo{value: value}
}

func (echo *Echo) Eval() float64 {
    value := echo.value.Eval()
    fmt.Println(":=", value)
    return value
}

func parse(lexer *BKLexer.Lexer) *Block {
    block := NewBlock()
    token := lexer.NextToken()
    for token.TType == BKLexer.TOKEN_TYPE_NEWLINE {
        token = lexer.NextToken()
    }
    for token.TType != BKLexer.TOKEN_TYPE_EOF {
        statement := parse_statement(lexer)
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

func parse_statement(lexer *BKLexer.Lexer) Node {
    token := lexer.GetToken()
    if token.Name == "SET" {
        name := lexer.NextToken()
        if name.Name != "NAME" {
            return nil
        }
        token = lexer.NextToken()
        if token.Name != "ASSIGN" {
            return nil
        }
        lexer.NextToken()
        value := parse_binary_add(lexer)
        if value == nil {
            return nil
        }
        return NewAssign(name, value)
    } else if token.Name == "ECHO" {
        lexer.NextToken()
        value := parse_binary_add(lexer)
        if (value == nil) {
            return nil
        }
        return NewEcho(value)
    }
    return parse_binary_add(lexer)
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
    lhs := factor(lexer)
    if lhs == nil {
        return nil
    }
    token := lexer.GetToken()
    for token.Source == "*" || token.Source == "/" {
        lexer.NextToken()
        rhs := factor(lexer)
        if rhs == nil {
            return nil
        }
        lhs = NewBinaryOpt(token, lhs, rhs)
        token = lexer.GetToken()
    }
    return lhs
}

func factor(lexer *BKLexer.Lexer) Node {
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
    if token.Name == "NAME" {
        name := NewName(token)
        lexer.NextToken()
        return name
    }
    return nil
}

func main() {
    lexer := BKLexer.NewLexer()
    lexer.AddRule("\\d+\\.?\\d*", "NUMBER")
    lexer.AddRule("[\\p{L}\\d_]+", "NAME")
    lexer.AddRule("\\+", "PLUS")
    lexer.AddRule("-", "MINUS")
    lexer.AddRule("\\*", "MUL")
    lexer.AddRule("/", "DIV")
    lexer.AddRule("\\(", "LPAR")
    lexer.AddRule("\\)", "RPAR")
    lexer.AddRule("=", "ASSIGN")
    lexer.AddIgnores("[ \\f\\t]+")
    lexer.AddIgnores("#[^\\r\\n]*")
    lexer.AddReserve("set")
    lexer.AddReserve("echo")

    bytes, err := ioutil.ReadFile("../test.txt")
    if err != nil {
        fmt.Println("read faild")
        return
    }
    code := string(bytes)
    lexer.Build(code)
    result := parse(lexer)
    if result == nil {
        fmt.Println("null result")
        return
    }
    ValueDict = make(map[string]float64)
    result.Eval()
}
