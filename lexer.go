package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type Position struct {
	line   int
	column int
}

type Lexer struct {
	input    *bufio.Reader
	buffer   bytes.Buffer
	position *Position
	offset   int
	indent   int
	dedent   bool
	error    error
}

func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(input)}
	l.position = &Position{line: 1}
	return l
}

func (l *Lexer) Error(e string) {
	err := fmt.Errorf("%s in %d:%d", e, (*l).position.line, (*l).position.column)
	l.error = err
}

func (l *Lexer) TokenText() string {
	return l.buffer.String()
}

func (l *Lexer) Next() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	if r == '\n' {
		l.position = &Position{line: l.position.line + 1}
	}
	l.position.column += w
	l.offset += w
	l.buffer.WriteRune(r)
	return r
}

func (l *Lexer) Skip() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	if r == '\n' {
		l.position = &Position{line: l.position.line + 1}
	}
	l.position.column += w
	l.offset += w
	return r
}

func (l *Lexer) Peek() rune {
	lead, err := l.input.Peek(1)
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error(err.Error())
		return 0
	}

	p, err := l.input.Peek(runeLen(lead[0]))

	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error("unexpected input error")
		return 0
	}

	ruNe, _ := utf8.DecodeRune(p)
	return ruNe
}

func (l *Lexer) NextToken() Token {
start:
	l.buffer.Reset()
	if l.dedent {
		if l.indent > 0 {
			l.indent--
			return l.newToken(DEDENT)
		}
		l.dedent = false
	}
	next := l.Peek()
	switch next {
	case '"':
		l.Skip()
		l.readString(next)
		return l.newToken(STRING)
	case '\'':
		l.Skip()
		l.readString(next)
		return l.newToken(STRING)
	case eof:
		return l.newToken(EOF)
	}

	next = l.Next()
	switch next {
	case ':':
		l.skipWhitespace()
		return l.newToken(COLON)
	case '\n':
		l.dedent = true
		goto start
	case ' ':
		if l.Peek() == ' ' {
			l.Next()
			l.indent++
			return l.newToken(INDENT)
		}
		goto start
	case eof:
		l.dedent = true
		goto start
	default:
		if isLetter(next) {
			l.readIdentifier()
			text := l.TokenText()
			return l.newToken(LookupIdent(text))
		} else if isDigit(next) {
			l.readNumber(next)
			return l.newToken(NUMBER)
		}
		return l.newToken(ILLEGAL)
	}
}

func (l *Lexer) readIdentifier() {
	next := l.Peek()
	for unicode.IsLetter(next) {
		l.Next()
		next = l.Peek()
	}
}

func (l *Lexer) readNumber(next rune) {
	if next == '0' && isDigit(l.Peek()) {
		l.Error("unexpected digit '0'")
		return
	} else if isDigit(next) {
		next := l.Peek()
		for {
			if !isDigit(next) {
				break
			}
			l.Next()
			next = l.Peek()
		}
		next = l.Peek()
		if next == '.' {
			l.Next()
			next = l.Peek()
			if !isDigit(next) {
				l.Error("unexpected token: expected digits")
				return
			}
			for {
				if !isDigit(next) {
					break
				}
				l.Next()
				next = l.Peek()
			}
		}
		next = l.Peek()
		if next == 'e' || next == 'E' {
			l.Next()
			next := l.Peek()
			if next == '+' || next == '-' {
				l.Next()
			}
			next = l.Peek()
			if !isDigit(next) {
				l.Error("digit expected for number exponent")
				return
			}
			l.Next()
			next = l.Peek()
			for {
				if !isDigit(next) {
					break
				}
				l.Next()
				next = l.Peek()
			}
		}
	} else {
		l.Error("error")
		return
	}
}

func (l *Lexer) readString(start rune) {
	for {
		next := l.Peek()
		if next == start {
			l.Skip()
			return
		}
		switch {
		case next == '\\':
			l.Skip()
			next := l.Peek()
			if next == start {
				l.Next()
			} else if next == 'b' {
				l.Skip()
				l.buffer.WriteRune('\b')
			} else if next == 'f' {
				l.Skip()
				l.buffer.WriteRune('\f')
			} else if next == 'n' {
				l.Skip()
				l.buffer.WriteRune('\n')
			} else if next == 'r' {
				l.Skip()
				l.buffer.WriteRune('\r')
			} else if next == 't' {
				l.Skip()
				l.buffer.WriteRune('\t')
			} else {
				l.Error("unsupported escape character")
				return
			}
		case unicode.IsControl(next):
			l.Error("cannot contain control characters in strings")
			return
		case next == eof:
			l.Error("unclosed string")
			return
		default:
			l.Next()
		}
	}
}

func (l *Lexer) newToken(typ TokenType) Token {
	return Token{Type: typ, Literal: l.TokenText(), Position: *l.position}
}

func (l *Lexer) skipWhitespace() {
	ruNe := l.Peek()
	for unicode.IsSpace(ruNe) && ruNe != '\n' {
		l.Next()
		ruNe = l.Peek()
	}
	l.buffer.Reset()
}

func runeLen(lead byte) int {
	if lead < 0xC0 {
		return 1
	} else if lead < 0xE0 {
		return 2
	} else if lead < 0xF0 {
		return 3
	}
	return 4
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && 'Z' <= ch || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
