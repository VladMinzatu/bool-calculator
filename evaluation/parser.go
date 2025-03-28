package evaluation

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: revise error handling in general
// but in particular, evaluation returns only one kind of error to be handled by client outside the package (use case for sentinel errors?)

// TODO: maybe refactoring to split evaluation and parsing??

type Expression interface {
	NumOutputs() int
	Evaluate(args map[string]bool) ([]bool, error)
}

type VariableSet map[string]struct{}

func ParseExpression(input string) (Expression, VariableSet, error) {
	tokens, err := ParseTokens(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract tokens from input: %w", err)
	}

	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("empty expression cannot be evaluated")
	}

	gateTokens := map[TokenType]struct{}{
		TokenNand: {},
		TokenNot:  {},
		TokenAnd:  {},
		TokenOr:   {},
		TokenXor:  {},
		TokenMux:  {},
		TokenDmux: {},
	}

	if _, ok := gateTokens[tokens[0].tokenType]; !ok && len(tokens) > 1 {
		return nil, nil, fmt.Errorf("Expression must either start with a gate name or contain exactly one literal or variable name")
	}

	parser := parser{tokens: tokens, pos: -1}
	variableSet := map[string]struct{}{}
	expression, err := parser.parse(variableSet, true)
	if err != nil {
		return nil, nil, err
	}
	return expression, variableSet, nil
}

type parser struct {
	tokens []Token
	pos    int
}

func (p *parser) parse(variableCollector VariableSet, isRoot bool /*Sorry, Uncle Bob*/) (Expression, error) {
	p.pos++
	if p.pos >= len(p.tokens) {
		return nil, errors.New("unexpected end of string encountered")
	}
	tok := p.tokens[p.pos]

	switch tok.tokenType {
	case TokenValue:
		value := tok.literal == "1"
		return &LiteralExpression{value: value}, nil
	case TokenVariable:
		variableCollector[tok.literal] = struct{}{}
		return &VariableExpression{variableName: tok.literal}, nil
	case TokenNot:
		exprs, err := p.parseArgs(1, variableCollector, isRoot)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &NotExpression{expression: exprs[0]}, nil
	case TokenNand, TokenAnd, TokenOr, TokenXor:
		exprs, err := p.parseArgs(2, variableCollector, isRoot)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &BinaryExpression{tok.tokenType, exprs}, nil
	case TokenMux:
		exprs, err := p.parseArgs(3, variableCollector, isRoot)
		if err != nil {
			return nil, argsError(tok, err)
		}
		return &MuxExpression{exprs}, nil
	case TokenDmux:
		exprs, err := p.parseArgs(2, variableCollector, isRoot)
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

func (p *parser) parseArgs(expectedInputs int, variableCollector VariableSet, isRoot bool) ([]Expression, error) {
	err := p.expect(TokenLparan)
	if err != nil {
		return nil, err
	}

	result := []Expression{}
	for expectedInputs > 0 {
		expr, err := p.parse(variableCollector, false)
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

	if isRoot && p.pos != len(p.tokens)-1 {
		return nil, errors.New("tokens found after root gate expression ended. Remove text following the root expression's closing paran")
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

func (e *LiteralExpression) Evaluate(args map[string]bool) ([]bool, error) {
	return []bool{e.value}, nil
}

type VariableExpression struct {
	variableName string
}

func (e *VariableExpression) NumOutputs() int {
	return 1
}

func (e *VariableExpression) Evaluate(args map[string]bool) ([]bool, error) {
	val, ok := args[e.variableName]
	if !ok {
		return nil, fmt.Errorf("cannot evaluate expression: no value provided for variable %v", e.variableName)
	}
	return []bool{val}, nil
}

type NotExpression struct {
	expression Expression
}

func (e *NotExpression) NumOutputs() int {
	return 1
}

func (e *NotExpression) Evaluate(args map[string]bool) ([]bool, error) {
	in, err := e.expression.Evaluate(args)
	if err != nil {
		return nil, err
	}
	if len(in) != 1 {
		panic("parser messed up. Not got more than 1 input during evaluation")
	}
	return []bool{Not(in[0])}, nil
}

type BinaryExpression struct {
	op          TokenType
	expressions []Expression
}

func (e *BinaryExpression) NumOutputs() int {
	return 1
}

func (e *BinaryExpression) Evaluate(args map[string]bool) ([]bool, error) {
	in, err := collectInputs(e.expressions, args)
	if err != nil {
		return nil, err
	}
	if len(in) != 2 {
		panic("parser messed up. Binary operator didn't get 2 inputs")
	}
	switch e.op {
	case TokenNand:
		return []bool{Nand(in[0], in[1])}, nil
	case TokenAnd:
		return []bool{And(in[0], in[1])}, nil
	case TokenOr:
		return []bool{Or(in[0], in[1])}, nil
	case TokenXor:
		return []bool{Xor(in[0], in[1])}, nil
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

func (e *MuxExpression) Evaluate(args map[string]bool) ([]bool, error) {
	in, err := collectInputs(e.expressions, args)
	if err != nil {
		return nil, err
	}
	if len(in) != 3 {
		panic("parser messed up. Mux operator didn't get 3 inputs")
	}
	return []bool{Mux(in[0], in[1], in[2])}, nil
}

type DmuxExpression struct {
	expressions []Expression
}

func (e *DmuxExpression) NumOutputs() int {
	return 2
}

func (e *DmuxExpression) Evaluate(args map[string]bool) ([]bool, error) {
	in, err := collectInputs(e.expressions, args)
	if err != nil {
		return nil, err
	}
	if len(in) != 2 {
		panic("parser messed up. Dmux operator didn't get 2 inputs")
	}
	out1, out2 := Dmux(in[0], in[1])
	return []bool{out1, out2}, nil
}

func collectInputs(expressions []Expression, args map[string]bool) ([]bool, error) {
	result := []bool{}
	for _, expr := range expressions {
		exprOuts, err := expr.Evaluate(args)
		if err != nil {
			return nil, err
		}
		result = append(result, exprOuts...)
	}
	return result, nil
}
