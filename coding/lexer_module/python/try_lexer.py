from bklexer import Token, Lexer


def main():
    print("Test Code:")
    code = "声明 变量 = PI * 100 - fda\n1024 * 4 * 3.14 ### \n123"
    print(code)
    print("--------------------------------")

    lexer = Lexer()
    lexer.add_rule("\\d+\\.\\d*", "FLOAT")
    lexer.add_rule("\\d+", "INT")
    lexer.add_rule("\\w+", "NAME")
    lexer.add_rule("\\+", "PLUS")
    lexer.add_rule("\\-", "MINUS")
    lexer.add_rule("\\*", "MUL")
    lexer.add_rule("/", "DIV")
    lexer.add_rule("=", "ASSIGN")
    lexer.add_rule("#[^\\r\\n]*", "COMMENT")
    lexer.add_ignores("[ \\f\\t]+")

    lexer.build(code)
    while True:
        token = lexer.next_token()
        if token.ttype != Token.TOKEN_TYPE_EOF:
            print(
                f"{token.name}\t{repr(token.source)}\tt{token.ttype}\t{token.position}\t{token.row},{token.col}")
        if token.ttype == Token.TOKEN_TYPE_EOF or token.ttype == Token.TOKEN_TYPE_ERROR:
            break


if __name__ == "__main__":
    main()
