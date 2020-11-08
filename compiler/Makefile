
.DEFAULT_GOAL := all

.PHONY: all lexer

zoe: cmd/*.go compiler/*.go
	cd ./cmd && go build && mv ./cmd ../zoe

lexer:
	cd compiler && make

all: lexer zoe
