
.DEFAULT_GOAL := all

.PHONY: all lexer

zoelsp: lsp/*.go compiler/*.go
	cd ./lsp && go build && mv ./lsp ../zoelsp

zoe: cmd/*.go compiler/*.go
	cd ./cmd && go build && mv ./cmd ../zoe

test: all
	cd compiler && make test

lexer:
	cd compiler && make

all: lexer zoe zoelsp
