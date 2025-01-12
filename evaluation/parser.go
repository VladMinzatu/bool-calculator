package evaluation

import (
	"errors"
	"fmt"
	"strings"
)

type Expression interface {
	NumOutputs() int
	Evaluate() []bool
}

type VariableSet map[string]bool

func ParseExpression(input string) (Expression, VariableSet, error) {
	tokens, err := ParseTokens(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract tokens from input: %w", err)
	}

	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("empty expression cannot be evaluated")
	}

	gateTokens := map[TokenType]bool{
		TokenNand: true,
		TokenNot:  true,
		TokenAnd:  true,
		TokenOr:   true,
		TokenXor:  true,
		TokenMux:  true,
		TokenDmux: true,
	}

	if _, ok := gateTokens[tokens[0].tokenType]; !ok && len(tokens) > 1 {
		return nil, nil, fmt.Errorf("Expression must either start with a gate name or contain exactly one literal or variable name")
	}

	parser := parser{tokens: tokens, pos: -1}
	variableSet := map[string]bool{}
	expression, err := parser.parse(variableSet)
	if err != nil {
		return nil, nil, err
	}
	return expression, variableSet, nil
}

type parser struct {
	tokens []Token
	pos    int
}

func (p *parser) parse(variableCollector VariableSet) (Expression, error) {
	p.pos++
	tok := p.tokens[p.pos]

	switch tok.tokenType {
	case TokenValue:
		value := tok.literal == "1"
		return &LiteralExpression{value: value}, nil
	case TokenVariable:
		variableCollector[tok.literal] = true
		return &VariableExpression{variableName: tok.literal}, nil
	case TokenNot:
		exprs, err := p.parseArgs(1, variableCollector)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &NotExpression{expression: exprs[0]}, nil
	case TokenNand, TokenAnd, TokenOr, TokenXor:
		exprs, err := p.parseArgs(2, variableCollector)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &BinaryExpression{tok.tokenType, exprs}, nil
	case TokenMux:
		exprs, err := p.parseArgs(3, variableCollector)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &MuxExpression{exprs}, nil
	case TokenDmux:
		exprs, err := p.parseArgs(2, variableCollector)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &DmuxExpression{exprs}, nil
	default:
		errorString := fmt.Sprintf("invalid token type: %v", tok)
		return nil, errors.New(errorString)
	}
}

func argsError(tok Token, err error) error {
	return fmt.Errorf("error parsing arguments for %s gate: %w", tok.literal, err)
}

func (p *parser) parseArgs(expectedInputs int, variableCollector VariableSet) ([]Expression, error) {
	err := p.expect(TokenLparan)
	if err != nil {
		return nil, err
	}

	result := []Expression{}
	for expectedInputs > 0 {
		expr, err := p.parse(variableCollector)
		if err != nil {
			return nil, err
		}
		result = append(result, expr)
		expectedInputs -= expr.NumOutputs()

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

func (e *LiteralExpression) NumOutputs() int {
	return 1
}

func (e *LiteralExpression) Evaluate() []bool {
	return []bool{e.value}
}

type VariableExpression struct {
	variableName string
}

func (e *VariableExpression) NumOutputs() int {
	return 1
}

func (e *VariableExpression) Evaluate() []bool {
	return []bool{true} //TODO: implement proper variable handling in evaluation
}

type NotExpression struct {
	expression Expression
}

func (e *NotExpression) NumOutputs() int {
	return 1
}

func (e *NotExpression) Evaluate() []bool {
	in := e.expression.Evaluate()
	if len(in) != 1 {
		panic("parser messed up. Not got more than 1 input during evaluation")
	}
	return []bool{Not(in[0])}
}

type BinaryExpression struct {
	op          TokenType
	expressions []Expression
}

func (e *BinaryExpression) NumOutputs() int {
	return 1
}

func (e *BinaryExpression) Evaluate() []bool {
	in := collectInputs(e.expressions)

	if len(in) != 2 {
		panic("parser messed up. Binary operator didn't get 2 inputs")
	}
	switch e.op {
	case TokenNand:
		return []bool{Nand(in[0], in[1])}
	case TokenAnd:
		return []bool{And(in[0], in[1])}
	case TokenOr:
		return []bool{Or(in[0], in[1])}
	case TokenXor:
		return []bool{Xor(in[0], in[1])}
	default:
		errorString := fmt.Sprintf("evaluation of binary expression %d not implemented", e.op)
		panic(errorString)
	}
}

type MuxExpression struct {
	expressions []Expression
}

func (e *MuxExpression) NumOutputs() int {
	return 1
}

func (e *MuxExpression) Evaluate() []bool {
	in := collectInputs(e.expressions)
	if len(in) != 3 {
		panic("parser messed up. Mux operator didn't get 3 inputs")
	}
	return []bool{Mux(in[0], in[1], in[2])}
}

type DmuxExpression struct {
	expressions []Expression
}

func (e *DmuxExpression) NumOutputs() int {
	return 2
}

func (e *DmuxExpression) Evaluate() []bool {
	in := collectInputs(e.expressions)
	if len(in) != 2 {
		panic("parser messed up. Dmux operator didn't get 2 inputs")
	}
	out1, out2 := Dmux(in[0], in[1])
	return []bool{out1, out2}
}

func collectInputs(expressions []Expression) []bool {
	result := []bool{}
	for _, expr := range expressions {
		exprOuts := expr.Evaluate()
		result = append(result, exprOuts...)
	}
	return result
}
