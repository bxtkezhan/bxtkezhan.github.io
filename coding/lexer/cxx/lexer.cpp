#include <locale>
#include <regex>
#include <string>
#include <vector>
#include <codecvt>


std::vector<std::wstring> exprs{L"\\d+", L"\\w+", L"[\\+\\-=]"};
std::vector<std::string> names{"INT",  "NAME", "SYMBOL"};


int main(int argc, char *argv[]) {
    std::locale old;
    std::locale::global(std::locale("en_US.UTF-8"));
    std::wstring_convert<std::codecvt_utf8<wchar_t>> codecvt_utf8;

    std::vector<std::wregex> rules;
    for (size_t i = 0, count = exprs.size(); i < count; ++i) {
        rules.push_back(std::wregex(L"^" + exprs[i]));
        printf("%s ^%s\n", names[i].c_str(), codecvt_utf8.to_bytes(exprs[i]).c_str());
    }

    printf("--------------------------------\n");
    for (int row = 0; row < argc - 1; ++row) {
        std::wstring code = codecvt_utf8.from_bytes(argv[row + 1]);
        size_t position = 0;
        while (true) {
            while (position < code.size() && (code[position] == L' ' || code[position] == L'\t'))
                position += 1;
            if (position >= code.size())
                break;

            auto subcode = code.substr(position);
            std::wsmatch match;
            int tokenType = -1;
            for (size_t i = 0, count = rules.size(); i < count; ++i) {
                if (std::regex_search(subcode, match, rules[i])) {
                    tokenType = i;
                    break;
                }
            }

            if (tokenType >= 0) {
                auto source = match.str(0);
                printf("%s\t`%s`\t%d\t%ld\n",
                    names[tokenType].c_str(), codecvt_utf8.to_bytes(source).c_str(), row, position);
                position += source.size();
            } else {
                printf("error in: %d, %ld\n", row, position);
                break;
            }
        }
    }

    std::locale::global(old);
    return 0;
}
