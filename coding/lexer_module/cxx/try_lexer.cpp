#include <locale>
#include "bklexer/lexer.hpp"


int main() {
    std::locale old;
    std::locale::global(std::locale("en_US.UTF-8"));

	printf("Test Code:\n");
    std::string code = "声明 变量 = PI * 100 - fda\n1024 * 4 * 3.14 ### \n123";
	printf("%s\n", code.c_str());
	printf("--------------------------------\n");

    auto lexer = bklexer::Lexer();
    lexer.add_rule("\\d+\\.\\d*", "FLOAT");
    lexer.add_rule("\\d+", "INT");
    lexer.add_rule("\\w+", "NAME");
    lexer.add_rule("\\+", "PLUS");
    lexer.add_rule("\\-", "MINUS");
    lexer.add_rule("\\*", "MUL");
    lexer.add_rule("/", "DIV");
    lexer.add_rule("=", "ASSIGN");
    lexer.add_rule("#[^\\r\\n]*", "COMMENT");
    lexer.add_ignores("[ \\f\\t]+");

    lexer.build(code);
    while (true) {
        auto token = lexer.next_token();
        if (token.ttype != TOKEN_TYPE_EOF) {
            printf("%s\t%s\tt%d\t%d\t%d,%d\n",
            token.name.c_str(), bklexer::quote(token.source).c_str(),
            token.ttype, token.position, token.row, token.col);
        }
        if (token.ttype == TOKEN_TYPE_EOF || token.ttype == TOKEN_TYPE_ERROR) {
            break;
        }
    }

    std::locale::global(old);
    return 0;
}
