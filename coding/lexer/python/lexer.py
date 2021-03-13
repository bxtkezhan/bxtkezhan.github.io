import re
import sys


exprs = ['\\d+', '\\w+', '[\\+-=]']
names = ['INT',  'NAME', 'SYMBOL']


def main():
    rules = []
    for i, expr in enumerate(exprs):
        rules.append(re.compile('^' + expr))
        print(names[i], rules[-1].pattern)

    print('-' * 32)
    for row, code in enumerate(sys.argv[1:]):
        position = 0
        while True:
            while position < len(code) and (code[position] == ' ' or code[position] == '\t'):
                position += 1
            if position >= len(code):
                break

            source = ''
            tokenType = -1
            for i, rule in enumerate(rules):
                result = rule.findall(code[position:])
                if len(result) > 0:
                    source = result[0]
                    tokenType = i
                    break
            if tokenType >= 0:
                print(f'{names[tokenType]}\t`{source}`\t{row}\t{position}')
                position += len(source)
            else:
                print(f'error in {row}, {position}')


if __name__ == "__main__":
    main()
