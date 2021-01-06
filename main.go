package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// NONE means no more characters were found in input.
const NONE = 255
const (
	// INTEGER //
	INTEGER = "INTEGER"
	// PLUS //
	PLUS = "PLUS"
	// MINUS //
	MINUS = "MINUS"
	// MUL //
	MUL = "MUL"
	// DIV //
	DIV = "DIV"
	// LPAREN //
	LPAREN = "("
	// RPAREN //
	RPAREN = ")"
	// EOF //
	EOF = "EOF"
)

// TokenType //
type TokenType string

// Token //
type Token struct {
	Type    TokenType
	Value   interface{}
	Literal string
}

// Lexer //
type Lexer struct {
	text        string
	pos         int
	currentChar byte
}

// NewLexer //
func NewLexer(text string) *Lexer {
	l := &Lexer{text: text}
	l.pos = 0
	l.currentChar = l.text[l.pos]
	return l
}

// lexerError
func (l *Lexer) lexerError() {
	fmt.Printf(`Unknown character %c`, l.currentChar)
	os.Exit(1)
}

// advance
func (l *Lexer) advance() {
	l.pos++
	if l.pos > len(l.text)-1 {
		l.currentChar = NONE
	} else {
		l.currentChar = l.text[l.pos]
	}
}

// skipWhiteSpace
func (l *Lexer) skipWhiteSpace() {
	for l.currentChar != NONE && l.isSpace(l.currentChar) {
		l.advance()
	}
}

// integer
func (l *Lexer) integer() int {
	var result string
	for l.currentChar != NONE && l.isDigit(l.currentChar) {
		result += string(l.currentChar)
		l.advance()
	}
	num, _ := strconv.Atoi(result)
	return num
}

// getNextToken
func (l *Lexer) getNextToken() Token {
	for l.currentChar != NONE {
		if l.isSpace(l.currentChar) {
			l.skipWhiteSpace()
		}
		if l.isDigit(l.currentChar) {
			tokenValue := l.integer()
			return Token{Type: INTEGER, Value: tokenValue, Literal: fmt.Sprint(tokenValue)}
		}
		if l.currentChar == '(' {
			l.advance()
			return Token{Type: LPAREN, Value: "(", Literal: "("}
		}
		if l.currentChar == ')' {
			l.advance()
			return Token{Type: RPAREN, Value: ")", Literal: ")"}
		}
		if l.currentChar == '+' {
			l.advance()
			return Token{Type: PLUS, Value: "+", Literal: "+"}
		}
		if l.currentChar == '-' {
			l.advance()
			return Token{Type: MINUS, Value: "-", Literal: "-"}
		}
		if l.currentChar == '*' {
			l.advance()
			return Token{Type: MUL, Value: "*", Literal: "*"}
		}
		if l.currentChar == '/' {
			l.advance()
			return Token{Type: DIV, Value: "/", Literal: "/"}
		}
		l.lexerError()
	}
	return Token{Type: EOF, Value: NONE, Literal: " "}
}

// isSpace lexer helper method
func (l *Lexer) isSpace(ch byte) bool {
	return ch == ' '
}

// isDigit lexer helper method
func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Node interface
type Node interface {
	String() string
}

// BinOp Node
type BinOp struct {
	Left  Node
	Op    Token
	Right Node
}

// String from Node interface.
func (b *BinOp) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(b.Op.Literal)
	out.WriteString(" " + b.Left.String() + " " + b.Right.String() + " ")
	out.WriteString(")")

	return out.String()
}

// Num Node
type Num struct {
	Token Token
	Value int
}

func (n *Num) String() string { return n.Token.Literal }

// Parser //
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

// parseError
func (p *Parser) parseError() {
	fmt.Printf("Syntax error")
	os.Exit(1)
}

// eat
func (p *Parser) eat(tokenType TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.getNextToken()
	} else {
		p.parseError()
	}
}

// factor
func (p *Parser) factor() Node {
	var nodeFactor Node
	if p.currentToken.Type == LPAREN {
		p.eat(LPAREN)
		nodeFactor = p.expr()
		p.eat(RPAREN)
	} else {
		nodeFactor = &Num{Token: p.currentToken, Value: p.currentToken.Value.(int)}
		p.eat(INTEGER)
	}
	return nodeFactor
}

// term
func (p *Parser) term() Node {
	nodeTerm := p.factor()
	for p.currentToken.Type == MUL || p.currentToken.Type == DIV {
		curToken := p.currentToken
		p.eat(curToken.Type)
		nodeTerm = &BinOp{Left: nodeTerm, Op: curToken, Right: p.factor()}
	}
	return nodeTerm
}

// expr
func (p *Parser) expr() Node {
	nodeExpr := p.term()
	for p.currentToken.Type == PLUS || p.currentToken.Type == MINUS {
		curToken := p.currentToken
		p.eat(curToken.Type)
		nodeExpr = &BinOp{Left: nodeExpr, Op: curToken, Right: p.term()}
	}
	return nodeExpr
}

// parse
func (p *Parser) parse() Node {
	return p.expr()
}

// NewParser create an instance of Parser
func NewParser(lexer *Lexer) *Parser {
	p := Parser{lexer: lexer}
	p.currentToken = p.lexer.getNextToken()
	return &p
}

// VisitorNode //
type VisitorNode struct {
	Parser *Parser
}

// NewVisitor //
func NewVisitor(parser *Parser) *VisitorNode {
	v := &VisitorNode{Parser: parser}
	return v
}

/**
* Visitor Node
 */
func (v *VisitorNode) visit(node Node) string {
	switch node := node.(type) {
	case *BinOp:
		return v.visitBinOp(node)
	case *Num:
		return v.visitNum(node)
	default:
		fmt.Printf("Unknown Node: %-v\n", node)
		os.Exit(1)
	}
	return ""
}

func (v *VisitorNode) visitBinOp(node *BinOp) string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(node.Op.Literal)
	out.WriteString(" ")
	out.WriteString(v.visit(node.Left))
	out.WriteString(" ")
	out.WriteString(v.visit(node.Right))
	out.WriteString(")")

	return out.String()
}

func (v *VisitorNode) visitNum(node *Num) string {
	return fmt.Sprintf("%d", node.Value)
}

func (v *VisitorNode) interpret() string {
	tree := v.Parser.parse()
	return v.visit(tree)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("goCalc> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		text := scanner.Text()
		lexer := NewLexer(text)
		parser := NewParser(lexer)
		visitor := NewVisitor(parser)
		result := visitor.interpret()
		fmt.Printf("%s\n", result)
	}
}
