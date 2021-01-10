.DEFAULT_GOAL := all

generated.go: genaccessors.ts nodes.go
	tsc --inlineSourceMap --outFile /dev/stdout genaccessors.ts --lib es2017 --target es2017 | node --enable-source-maps > generated.go

lexer.go: lexer.go_
	re2go -c -r --input-encoding utf8 -c -W lexer.go_ -o lexer.go
	echo "const (" >> lexer.go
	echo "  TK_EOF TokenKind = iota" >> lexer.go
	grep -oP 'addToken\(\K[^\)]+' lexer.go_ | grep -v TK_EOF | sort | uniq >> lexer.go
	echo "TK__MAX\n)" >> lexer.go
	echo "var tokstr = [...]string{" >> lexer.go
	grep -oP 'addToken\(\K[^\)]+' lexer.go_ | sort | uniq | sed -E 's/(.*)/\1: "\1",/' >> lexer.go
	echo "}" >> lexer.go

test:
	go test -cpuprofile ~/cpu.prof -memprofile ~/mem.prof -bench .

all: lexer.go generated.go
