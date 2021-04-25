package parser

import (
	"fmt"
	"strings"

	"github.com/gscienty/causer/expr/ast"
	"github.com/spf13/cast"
)

type parser struct {
	tokens  []Token
	current Token
	pos     int
	err     error
}

type associativity string

const (
	associateLeft  associativity = "left"
	associateRight associativity = "right"
)

type operator struct {
	priority  int
	associate associativity
}

var unaryOp = map[string]operator{
	"not": {5, associateLeft},
	"!":   {5, associateLeft},
	"-":   {10, associateLeft},
	"+":   {10, associateLeft},
}

var binaryOp = map[string]operator{
	"or":  {1, associateLeft},
	"||":  {1, associateLeft},
	"and": {2, associateLeft},
	"&&":  {2, associateLeft},
	"==":  {3, associateLeft},
	"!=":  {3, associateLeft},
	"<":   {3, associateLeft},
	">":   {3, associateLeft},
	"<=":  {3, associateLeft},
	">=":  {3, associateLeft},
	"in":  {3, associateLeft},
	"+":   {4, associateLeft},
	"-":   {4, associateLeft},
	"*":   {6, associateLeft},
	"/":   {6, associateLeft},
	"^":   {7, associateLeft},
}

func Parse(source string) (*ast.Tree, error) {
	tokens, err := Lexer(source)
	if err != nil {
		return nil, err
	}

	p := &parser{tokens: tokens, current: tokens[0]}

	node := p.parse(0)

	if p.current.Kind != TokenKindEOF {
		p.err = fmt.Errorf("unexpected token %v", p.current)
	}

	if p.err != nil {
		return nil, p.err
	}

	return &ast.Tree{Root: node}, nil
}

func (p *parser) parse(priority int) ast.Node {
	nodeLeft := p.parsePrimary()

	token := p.current
	for token.Kind == TokenKindOperator && p.err == nil {
		if op, ok := binaryOp[token.Value]; ok {
			if op.priority >= priority {
				p.next()

				var nodeRight ast.Node
				if op.associate == associateLeft {
					nodeRight = p.parse(op.priority + 1)
				} else {
					nodeRight = p.parse(op.priority)
				}

				nodeLeft = &ast.BinaryNode{
					Operator: token.Value,
					Left:     nodeLeft,
					Right:    nodeRight,
				}

				token = p.current
				continue
			}
		}
		break
	}

	if priority == 0 {
		// TODO
	}

	return nodeLeft
}

func (p *parser) parsePrimary() ast.Node {
	token := p.current

	if token.Kind == TokenKindOperator {
		if op, ok := unaryOp[token.Value]; ok {
			p.next()
			expr := p.parse(op.priority)
			node := &ast.UnaryNode{
				Operator: token.Value,
				Expr:     expr,
			}

			return p.parsePostfix(node)
		}
	}

	if token.Kind == TokenKindBracket && token.Value == "(" {
		p.next()
		expr := p.parse(0)
		p.next()
		return p.parsePostfix(expr)
	}

	switch token.Kind {
	case TokenKindIdentifier:
		p.next()
		switch token.Value {
		case "true":
			return &ast.BoolNode{Value: true}
		case "false":
			return &ast.BoolNode{Value: false}
		case "nil":
			return &ast.NilNode{}
		default:
			node := p.parseIdentifier(token)
			return p.parsePostfix(node)
		}

	case TokenKindNumber:
		p.next()
		value := strings.ReplaceAll(token.Value, "_", "")
		if strings.ContainsAny(value, ".") {
			return &ast.FloatNode{Value: cast.ToFloat64(value)}
		} else {
			return &ast.IntNode{Value: cast.ToInt(value)}
		}

	case TokenKindString:
		p.next()
		return &ast.StringNode{Value: token.Value}

	default:
		// TODO
	}

	return nil
}

func (p *parser) parseIdentifier(token Token) ast.Node {
	if p.current.Kind == TokenKindBracket && p.current.Value == "(" {
		p.next()
		arguments := p.parseArguments()
		return &ast.FunctionNode{
			Name:      token.Value,
			Arguments: arguments,
		}
	} else {
		return &ast.IdentifierNode{Value: token.Value}
	}
}

func (p *parser) next() {
	p.pos++
	if p.pos >= len(p.tokens) {
		p.err = fmt.Errorf("unexpect end of expression")
		return
	}

	p.current = p.tokens[p.pos]
}

func (p *parser) parsePostfix(node ast.Node) ast.Node {
	token := p.current
	for (token.Kind == TokenKindOperator || token.Kind == TokenKindBracket) && p.err == nil {
		if token.Value == "." {
			p.next()
			token = p.current
			p.next()

			if token.Kind != TokenKindIdentifier {
				p.err = fmt.Errorf("expect name")
			}

			if p.current.Kind == TokenKindBracket && p.current.Value == "(" {
				p.next()
				args := p.parseArguments()
				node = &ast.MethodNode{
					Node:      node,
					Method:    token.Value,
					Arguments: args,
				}
			} else {
				node = &ast.PropertyNode{
					Node:     node,
					Property: token.Value,
				}
			}
		} else if p.current.Kind == TokenKindBracket && p.current.Value == "(" {
			p.next()
			args := p.parseArguments()
			node = &ast.FunctionNode{
				Name:      token.Value,
				Arguments: args,
			}
		} else {
			break
		}
		token = p.current
	}

	return node
}

func (p *parser) parseArguments() []ast.Node {
	nodes := make([]ast.Node, 0)
	for !(p.current.Kind == TokenKindBracket && p.current.Value == ")") && p.err == nil {
		if len(nodes) > 0 {
			if !(p.current.Kind == TokenKindOperator && p.current.Value == ",") {
				p.err = fmt.Errorf("invalid token")
			}
		}
		node := p.parse(0)
		nodes = append(nodes, node)
	}
	p.next()

	return nodes
}
