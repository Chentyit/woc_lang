package ast

import (
	"woc_lang/token"
)

// IdentExpression 标识符节点，例如 `var x = 5;` 中的 x
// 标识符之所以定义为表达式，是因为当一个值绑定到 `x` 上后，
// x 就指向了这个值，举个栗子：var other_ident = x; 这个语句中，x 就作为表达式
// 将 x 指向的 5 赋值给了 other_ident
type IdentExpression struct {
	Token token.Token
	Value string
}

func (ie *IdentExpression) eNode() {}

func (ie *IdentExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IdentExpression) String() string {
	return ie.Value
}

// IntegerLiteral 整型字面量
type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) eNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
