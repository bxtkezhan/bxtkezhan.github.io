import re


class Token:
    TOKEN_TYPE_ERROR = -2
    TOKEN_TYPE_EOF = -1
    TOKEN_TYPE_NEWLINE = 0

    def __init__(self, position:int, source:str, ttype:int, row:int, col:int, name:str):
        self.position = position
        self.source   = source
        self.ttype    = ttype
        self.row      = row
        self.col      = col
        self.name     = name


class Lexer:
    def __init__(self):
        self.code = ""
        self.position = 0
        self.source = ""
        self.ttype = 0
        self.row = 0
        self.col = 0
        self.exprs = []
        self.names = []
        self.ignores:re.Pattern = None

        self.add_rule("\n", "NEWLINE")

    def add_rule(self, expr:str, name:str):
        self.exprs.append(re.compile("^" + expr))
        self.names.append(name)

    def add_ignores(self, expr:str):
        self.ignores = re.compile("^" + expr)

    def build(self, code:str):
        self.code = code
        self.position = 0
        self.source = ""
        self.ttype = 0
        self.row  = 0
        self.col  = 0

    def get_token_name(self):
        if self.ttype == -1:
            return "EOF"
        elif self.ttype == -2:
            return "ERROR"
        return self.names[self.ttype]

    def get_token(self):
        return Token(
            self.position, self.source, self.ttype,
            self.row, self.col, self.get_token_name())

    def next_token(self) -> Token:
        self.position += len(self.source)
        if self.source != "\n":
            self.col += len(self.source)
        else:
            self.row += 1
            self.col = 0

        result = self.ignores.findall(self.code[self.position:])
        source = result[0] if len(result) else ''
        self.position += len(source)
        self.col += len(source)
        if self.position >= len(self.code):
            self.ttype = -1
            return self.get_token()
        self.ttype = -2
        for i, rule in enumerate(self.exprs):
            result = rule.findall(self.code[self.position:])
            source = result[0] if len(result) else ''
            if source != "":
                self.source = source
                self.ttype = i
                break
        return self.get_token()
