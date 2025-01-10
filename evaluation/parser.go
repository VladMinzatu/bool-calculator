package evaluation

import (
	"fmt"
	"strings"
)

type Expression interface {
	numOutputs() int
	numInputs() int
}

func ParseExpression(input string) (Expression, error) {
	tokens, err := ParseTokens(input)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tokens from input: %w", err)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression cannot be evaluated")
	}

	if _, ok := gateTokens[tokens[0].tokenType]; !ok && len(tokens) > 1 {
		return nil, fmt.Errorf("Expression must either start with a gate name or contain exactly one literal or variable name")
	}

	parser := parser{tokens: tokens, pos: -1}
	return parser.parse()
}

var gateTokens map[TokenType]bool = map[TokenType]bool{
	TokenNand: true,
	TokenNot:  true,
	TokenAnd:  true,
	TokenOr:   true,
	TokenXor:  true,
	TokenMux:  true,
	TokenDmux: true,
}

type parser struct {
	tokens []Token
	pos    int
}

func (p *parser) parse() (Expression, error) {
	p.pos++
	tok := p.tokens[p.pos]

	switch tok.tokenType {
	case TokenValue:
		value := tok.literal == "1"
		return &LiteralExpression{value: value}, nil
	case TokenVariable:
		return &VariableExpression{variableName: tok.literal}, nil
	case TokenNot:
		exprs, err := p.parseArgs(1)
		if err != nil {
			return nil, err
		}
		return &NotExpression{expression: exprs[0]}, nil
	case TokenNand, TokenAnd, TokenOr, TokenXor:
		exprs, err := p.parseArgs(2)
		if err != nil {
			return nil, err
		}
		return &BinaryExpression{tok.tokenType, exprs}, nil
	case TokenMux:
		exprs, err := p.parseArgs(3)
		if err != nil {
			return nil, err
		}
		return &MuxExpression{exprs}, nil
	case TokenDmux:
		exprs, err := p.parseArgs(2)
		if err != nil {
			return nil, err
		}
		return &DmuxExpression{exprs}, nil
	default:
		return nil, fmt.Errorf("invalid single token expression")
	}
}

func (p *parser) parseArgs(expectedInputs int) ([]Expression, error) {
	err := p.expect(TokenLparan)
	if err != nil {
		return nil, err
	}

	result := []Expression{}
	for expectedInputs > 0 {
		expr, err := p.parse()
		if err != nil {
			return nil, err
		}
		result = append(result, expr)
		expectedInputs -= expr.numOutputs()

		if expectedInputs < 0 {
			return nil, fmt.Errorf("too many inputs for gate")
		}

		if expectedInputs == 0 {
			break // enough inputs, we should expect ')' now
		}
		err = p.expect(TokenComma)
		if err != nil {
			return nil, err
		}
	}
	err = p.expect(TokenRparan)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *parser) expect(expected TokenType) error {
	p.pos++
	if p.pos >= len(p.tokens) {
		return fmt.Errorf("expected token %v, but reached end of string", expected.String())
	}
	tok := p.tokens[p.pos]
	if tok.tokenType != expected {
		return fmt.Errorf("expected token %v, but found %v at %s^", expected.String(), p.tokens[p.pos].literal, tokenString(p.tokens, p.pos))
	}
	return nil
}

func tokenString(tokens []Token, until int) string {
	var sb strings.Builder
	for i := 0; i <= until; i++ {
		sb.WriteString(tokens[i].literal)
	}
	return sb.String()
}

type LiteralExpression struct {
	value bool
}

func (e *LiteralExpression) numOutputs() int {
	return 1
}

func (e *LiteralExpression) numInputs() int {
	return 0
}

type VariableExpression struct {
	variableName string
}

func (e *VariableExpression) numOutputs() int {
	return 1
}

func (e *VariableExpression) numInputs() int {
	return 0
}

type NotExpression struct {
	expression Expression
}

func (e *NotExpression) numOutputs() int {
	return 1
}

func (e *NotExpression) numInputs() int {
	return 1
}

type BinaryExpression struct {
	op          TokenType
	expressions []Expression
}

func (e *BinaryExpression) numOutputs() int {
	return 1
}

func (e *BinaryExpression) numInputs() int {
	return 2
}

type MuxExpression struct {
	expressions []Expression
}

func (e *MuxExpression) numOutputs() int {
	return 1
}

func (e *MuxExpression) numInputs() int {
	return 3
}

type DmuxExpression struct {
	expressions []Expression
}

func (e *DmuxExpression) numOutputs() int {
	return 2
}

func (e *DmuxExpression) numInputs() int {
	return 2
}
