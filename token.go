package main

type Token struct {
	Type     TokenType
	Literal  string
	Position Position
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	INDENT
	DEDENT

	IDENTIFIER
	NUMBER
	STRING

	COMMA
	COLON

	LBRACE
	RBRACE

	TRUE
	FALSE
)

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

func TokenName(typ TokenType) string {
	switch typ {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case INDENT:
		return "INDENT"
	case DEDENT:
		return "DEDENT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"

	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"

	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"

	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	default:
		panic("unknown token type")
	}
}
