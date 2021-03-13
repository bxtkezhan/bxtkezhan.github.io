#ifndef BXTKEZHAN_LEXER_HEADER
#define BXTKEZHAN_LEXER_HEADER 1

#include <regex>
#include <string>
#include <vector>
#include <codecvt>


namespace bklexer {

#define TOKEN_TYPE_ERROR -2
#define TOKEN_TYPE_EOF -1
#define TOKEN_TYPE_NEWLINE 0

class Token {
    public:
        int position;
        std::string source;
        int ttype;
        int row;
        int col;
        std::string name;

        Token(int position, const std::string &source, int ttype, int row, int col, const std::string &name):
        position(position), source(source), ttype(ttype), row(row), col(col), name(name){};
};


std::wstring_convert<std::codecvt_utf8<wchar_t>> codecvt_utf8;

class Lexer {
    private:
        std::wstring code;
        int position;
        std::string source;
        int ttype;
        int row;
        int col;
        std::vector<std::wregex> exprs;
        std::vector<std::string> names;
        std::wregex ignores;

    public:
        Lexer() {
            add_rule("\n", "NEWLINE");
        }

        void add_rule(const std::string &expr, const std::string name) {
            this->exprs.push_back(std::wregex(codecvt_utf8.from_bytes("^" + expr)));
            this->names.push_back(name);
        }

        void add_ignores(const std::string &expr) {
            this->ignores = codecvt_utf8.from_bytes("^" + expr);
        }

        void build(std::string &code) {
            this->code = codecvt_utf8.from_bytes(code);
            this->position = 0;
            this->source = "";
            this->ttype = 0;
            this->row  = 0;
            this->col  = 0;
        }

        std::string get_token_name() {
            switch (ttype) {
                case -1:
                    return "EOF";
                case -2:
                    return "ERROR";
                default:
                    return names[ttype];
            }
        }

        Token get_token() {
            return Token(position, source, ttype, row, col, get_token_name());
        }

        Token next_token() {
            auto wsource = codecvt_utf8.from_bytes(source);
            position += wsource.size();
            if (source != "\n") {
                col += wsource.size();
            } else {
                row += 1;
                col = 0;
            }

            auto subcode = code.substr(position);
            std::wsmatch match;
            if (std::regex_search(subcode, match, ignores)) {
                int match_size = match.str(0).size();
                position += match_size;
                col += match_size;
            }
            if (position >= (int)code.size()) {
                ttype = -1;
                return get_token();
            }
            ttype = -2;
            subcode = code.substr(position);
            for (size_t i = 0, count = exprs.size(); i < count; ++i) {
                if (std::regex_search(subcode, match, exprs[i])) {
                    source = codecvt_utf8.to_bytes(match.str(0));
                    ttype = i;
                    break;
                }
            }
            return get_token();
        }
};


std::string quote(const std::string &src) {
    const char before_escape[] = {'\"', '\\', 't', 'r', 'n', 'b', 'f'};
    const char after_escape[] = {'\"', '\\', '\t', '\r', '\n', '\b', '\f'};
    const size_t escape_list_size = sizeof(before_escape);
    std::string dst = "\"";
	char inchar;
    for (size_t i = 0, length = src.size(); i < length; ++i) {
        inchar = src[i];
        for (size_t j = 0; j < escape_list_size; ++j) {
            if (inchar == after_escape[j]) {
                dst += '\\';
                inchar = before_escape[j];
                break;
            }
        }
        dst += inchar;
    }
    dst += '\"';
	return dst;
}

}

#endif /* ifndef BXTKEZHAN_LEXER_HEADER */
