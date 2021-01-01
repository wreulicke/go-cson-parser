package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	buf := bytes.NewBufferString("")
	l := NewLexer(buf)
	assertToken(t, l, EOF)
}

func TestLexerSimple(t *testing.T) {
	text := `xxx: "xxx"
yyy: 2
`
	buf := bytes.NewBufferString(text)
	l := NewLexer(buf)
	expects := []TokenType{
		IDENTIFIER,
		COLON,
		STRING,
		IDENTIFIER,
		COLON,
		NUMBER,
		EOF,
	}
	for _, e := range expects {
		assertToken(t, l, e)
	}
}
func TestLexerNested(t *testing.T) {
	text := `a: 
  b: "b"
c: "c"
`
	buf := bytes.NewBufferString(text)
	l := NewLexer(buf)
	expects := []TokenType{
		IDENTIFIER,
		COLON,
		INDENT,
		IDENTIFIER,
		COLON,
		STRING,
		DEDENT,
		IDENTIFIER, COLON, STRING,
		EOF,
	}
	for _, e := range expects {
		assertToken(t, l, e)
	}
}

func TestLexerIndent(t *testing.T) {
	text := `
foo:"foo"
bar: "bar"

baz:         "baz"
xxx:
  test: "xxx"
  v: 
    x: "xxx"
`
	buf := bytes.NewBufferString(text)
	l := NewLexer(buf)
	expects := []TokenType{
		IDENTIFIER,
		COLON,
		STRING,
		IDENTIFIER,
		COLON,
		STRING,
		IDENTIFIER,
		COLON,
		STRING,
		IDENTIFIER,
		COLON,
		INDENT,
		IDENTIFIER,
		COLON,
		STRING,
		DEDENT,
		INDENT,
		IDENTIFIER,
		COLON,
		DEDENT,
		INDENT,
		INDENT,
		IDENTIFIER,
		COLON,
		STRING,
		DEDENT,
		DEDENT,
		EOF,
	}
	for i, e := range expects {
		t.Run(fmt.Sprintf("token[%d]", i), func(t *testing.T) {
			assertToken(t, l, e)
		})
	}
}

func assertToken(t *testing.T, l *Lexer, expect TokenType) {
	if tok := l.NextToken(); tok.Type != expect {
		t.Errorf("current token is unexpected. expected: %s, actual type: %s, position: %d:%d, literal:\n  %s ",
			TokenName(expect), TokenName(tok.Type), l.position.line, l.position.column, tok.Literal)
	}
}
