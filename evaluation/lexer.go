package evaluation

import "fmt"

type TokenType int

const (
	TokenValue TokenType = iota
	TokenVariable

	// Gates
	TokenNand
	TokenNot
	TokenAnd
	TokenOr
	TokenXor
	TokenMux
	TokenDmux

	TokenLparan
	TokenRparan
	TokenComma

	tokenEOF // used only internally, won't be returned by our parser
)

func (t TokenType) String() string {
	switch t {
	case TokenLparan:
		return "("
	case TokenRparan:
		return ")"
	case TokenComma:
		return ","
	default:
		return "UNHANDLED"
	}
}

var keywords = map[string]TokenType{
	"nand": TokenNand,
	"not":  TokenNot,
	"and":  TokenAnd,
	"or":   TokenOr,
	"xor":  TokenXor,
	"mux":  TokenMux,
	"dmux": TokenDmux,
}

type Token struct {
	tokenType TokenType
	literal   string
}

func ParseTokens(text string) ([]Token, error) {
	result := []Token{}
	for tok, idx, err := nextToken(text, 0); tok.tokenType != tokenEOF; tok, idx, err = nextToken(text, idx) {
		if err != nil {
			return nil, err
		}
		result = append(result, tok)
	}
	return result, nil
}

func nextToken(text string, index int) (Token, int, error) {
	var token Token
	currentIndex := index

	// Skip whitespace
	for currentIndex < len(text) && isWhitespace(text[currentIndex]) {
		currentIndex++
	}

	// Return EOF if no more input
	if currentIndex >= len(text) {
		return Token{tokenType: tokenEOF, literal: ""}, currentIndex, nil
	}

	// Get first character
	ch := text[currentIndex]
	currentIndex++

	switch ch {
	case '(':
		token = Token{tokenType: TokenLparan, literal: string(ch)}
	case ')':
		token = Token{tokenType: TokenRparan, literal: string(ch)}
	case ',':
		token = Token{tokenType: TokenComma, literal: string(ch)}
	case '0', '1':
		token = Token{tokenType: TokenValue, literal: string(ch)}
	default:
		// Handle identifiers (variables and keywords)
		if isLetter(ch) {
			identifier := string(ch)
			for currentIndex < len(text) && isLetter(text[currentIndex]) {
				identifier += string(text[currentIndex])
				currentIndex++
			}

			// Check if identifier is a keyword
			if tokenType, isKeyword := keywords[identifier]; isKeyword {
				token = Token{tokenType: tokenType, literal: identifier}
			} else {
				token = Token{tokenType: TokenVariable, literal: identifier}
			}
		} else {
			return token, currentIndex, fmt.Errorf("invalid character encountered: %c", ch)
		}
	}

	return token, currentIndex, nil
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
