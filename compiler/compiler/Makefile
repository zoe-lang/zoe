.DEFAULT_GOAL := all

#nodesmethods.go: genaccessors.js nodes.go
	#node genaccessors.js > nodesmethods.go

lexer.go: lexer.go_
	re2go -c -r --input-encoding utf8 -c -W lexer.go_ -o lexer.go
	echo "const (" >> lexer.go
	echo "  TK_EOF TokenKind = iota" >> lexer.go
	grep -oP 'addToken\(\K[^\)]+' lexer.go_ | sort | uniq >> lexer.go
	echo "TK__MAX\n)" >> lexer.go
	echo "var tokstr = [...]string{" >> lexer.go
	grep -oP 'addToken\(\K[^\)]+' lexer.go_ | sort | uniq | sed -E 's/(.*)/\1: "\1",/' >> lexer.go
	echo "}" >> lexer.go

all: lexer.go # nodesmethods.go
