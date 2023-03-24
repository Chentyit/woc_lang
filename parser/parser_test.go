package parser

import (
	"fmt"
	"testing"
	"woc_lang/ast"
	"woc_lang/lexer"
)

func TestParsingVarStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"var x = 5;", "x", 5},
		{"var foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		parser := New(l)
		checkParserErrors(t, parser)

		if len(parser.Program.Statements) != 1 {
			t.Fatalf("语句解析错误，解析语句数量为: %d", len(parser.Program.Statements))
		}

		stmt := parser.Program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestAssignStmt(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x = 7;", "x", 7},
		{"foo = bar;", "foo", "bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		p := New(l)
		checkParserErrors(t, p)

		if len(p.Program.Statements) != 1 {
			t.Fatalf("语句解析错误，解析语句数量为: %d", len(p.Program.Statements))
		}

		stmt := p.Program.Statements[0].(*ast.AssignStatement)
		ident := stmt.Ident
		exp := stmt.Exp

		if !testLiteralExpression(t, ident, tt.expectedIdentifier) {
			return
		}

		if !testLiteralExpression(t, exp, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)
		p := New(l)
		checkParserErrors(t, p)

		if len(p.Program.Statements) != 1 {
			t.Fatalf("语句数量错误，实际: %d",
				len(p.Program.Statements))
		}

		stmt := p.Program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt 不是 *ast.ReturnStatement 类型. 实际: %T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral 不是 'return' 关键字, 实际: %q",
				returnStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		Operator string
		value    any
	}{
		{"!foobar;", "!", "foobar"},
		{"-5;", "-", 5},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		parser := New(l)
		checkParserErrors(t, parser)

		if parser.Program == nil {
			t.Fatalf("测试用例未解析到代码")
		}

		if len(parser.Program.Statements) != 1 {
			t.Fatalf("测试用例语法结构与预期不符:\n预期: %d\n实际: %d",
				1, len(parser.Program.Statements))
		}

		stmt, ok := parser.Program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("语句解析错误: %T", parser.Program.Statements[0])
		}

		preExp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("该语句并非前缀表达式: %T", stmt.Expression)
		}

		if preExp.Operator != tt.Operator {
			t.Fatalf("该语句操作符错误\n期望: '%s'\n, 实际得到: '%s'",
				tt.Operator,
				preExp.Operator)
		}

		if !testLiteralExpression(t, preExp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		parser := New(l)
		checkParserErrors(t, parser)

		if len(parser.Program.Statements) != 1 {
			t.Fatalf("program.Statements 期望语句数量: %d，实际获得=%d\n",
				1, len(parser.Program.Statements))
		}

		stmt, ok := parser.Program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] 不是表达式声明语句，实际获得=%T",
				parser.Program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		parser := New(l)
		checkParserErrors(t, parser)

		if len(parser.Program.Statements) != 1 {
			t.Fatalf("语句解析失败，实际解析语句数量: %d", len(parser.Program.Statements))
		}

		stmt, ok := parser.Program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("parser.program.Statements[0] 并非表达式声明语句，解析结果: %T",
				parser.Program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("表达式解析结果不是 ast.BooleanLiteral，解析结果为: %T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("布尔值期望结果: %t, 实际获取: %t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4;",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5));",
		},
		{
			"(5 + 5) * 2 * (5 + 5);",
			"(((5 + 5) * 2) * (5 + 5));",
		},
		{
			"-(5 + 5);",
			"(-(5 + 5));",
		},
		{
			"!(true == true);",
			"(!(true == true));",
		},
		{
			"a + add(b * c) + d;",
			"((a + add((b * c))) + d);",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"add(a + b + c * d / f + g);",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		checkLexerErrors(t, l)

		parser := New(l)
		checkParserErrors(t, parser)

		actual := parser.Program.String()
		if actual != tt.expected {
			t.Errorf("期望结果: %q\n实际获得: %q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `
	if (x < y) {
		x;
	}
	`

	l := lexer.New(input)
	checkLexerErrors(t, l)
	p := New(l)
	checkParserErrors(t, p)

	if len(p.Program.Statements) != 1 {
		t.Fatalf("语句解析失败，实际解析语句数量: %d", len(p.Program.Statements))
	}

	stmt, ok := p.Program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStatement，实际取值: %T",
			p.Program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression 不是 ast.IfExpression，实际取值: %T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("IfExpression 的 consequence 并不是 1 个 statements，实际取值: %d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.ElseExpression != nil {
		t.Errorf("exp.Alternative 不为空，实际取值: %+v", exp.ElseExpression)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `
	if (x < y) {
		x;
	} else {
		y;
	}
	`

	l := lexer.New(input)
	checkLexerErrors(t, l)
	p := New(l)
	checkParserErrors(t, p)

	if len(p.Program.Statements) != 1 {
		t.Fatalf("语句解析失败，实际解析语句数量: %d", len(p.Program.Statements))
	}

	stmt, ok := p.Program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStatement，实际取值: %T",
			p.Program.Statements[0])
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression 不是 ast.IfExpression，实际取值: %T",
			stmt.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		return
	}

	ifConsequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, ifConsequence.Expression, "x") {
		return
	}

	elseExp := ifExp.ElseExpression

	elseConsequence, ok := elseExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, elseConsequence.Expression, "y") {
		return
	}
}

func TestElseIfExpression(t *testing.T) {
	input := `
	if (x < y) {
		x;
	} else if (y < x) {
		y;
	} else {
		z;
	}
	`

	l := lexer.New(input)
	checkLexerErrors(t, l)
	p := New(l)
	checkParserErrors(t, p)

	if len(p.Program.Statements) != 1 {
		t.Fatalf("语句解析失败，实际解析语句数量: %d", len(p.Program.Statements))
	}

	stmt, ok := p.Program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStatement，实际取值: %T",
			p.Program.Statements[0])
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression 不是 ast.IfExpression，实际取值: %T",
			stmt.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		return
	}

	ifConsequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, ifConsequence.Expression, "x") {
		return
	}

	ifExp = ifExp.ElseExpression.NextIfExp

	if !testInfixExpression(t, ifExp.Condition, "y", "<", "x") {
		return
	}

	ifConsequence, ok = ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, ifConsequence.Expression, "y") {
		return
	}

	elseExp := ifExp.ElseExpression

	elseConsequence, ok := elseExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] 不是 ast.ExpressionStatement 类型，实际取值: %T",
			ifExp.Consequence.Statements[0])
	}

	if !testIdentifier(t, elseConsequence.Expression, "z") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "func test(x, y) { x * y; }"

	l := lexer.New(input)
	checkLexerErrors(t, l)
	p := New(l)
	checkParserErrors(t, p)

	if len(p.Program.Statements) != 1 {
		t.Fatalf("p.program.Statements 包含语句数量错误:\n期望: %d\n实际:%d\n",
			1, len(p.Program.Statements))
	}

	stmt, ok := p.Program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("p.program.Statements[0] 并不是 ast.ExpressionStatement 类型，实际类型为: %T",
			p.Program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression 并不是 ast.FunctionLiteral 类型，实际类型为: %T",
			stmt.Expression)
	}

	if function.Name.Value != "test" {
		t.Fatalf("函数名称错误:\n期望: %s\n实际: %s\n",
			"test", function.Name.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("函数字面量的形参数量错误:\n期望: %d\n实际: %d\n",
			2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements 包含语句数量错误. 实际: %d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("函数体类型不是 ast.ExpressionStatement. 实际: %T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "*", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	checkLexerErrors(t, l)
	p := New(l)
	checkParserErrors(t, p)

	if len(p.Program.Statements) != 1 {
		t.Fatalf("语句数量错误，实际: %d", len(p.Program.Statements))
	}

	stmt, ok := p.Program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("p.program.Statements[0] 并不是 ast.ExpressionStatement 类型，实际类型为: %T",
			p.Program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression 并不是 ast.CallExpression 类型，实际类型为: %T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.FunctionName, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("实参数量错误，实际: %d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral 不是 'var', 实际得到: %T", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("语句并非 *ast.LetStatement 类型，实际得到: %T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value 期望得到: '%s'，实际得到: %s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() 期望得到: '%s'，got=%s",
			name, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("无法处理该类型: %T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("传入字面量类型不是 *ast.IntegerLiteral: %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value 期望是 %d. 实际=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral 期望是 %d. 实际=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.IdentLiteral)
	if !ok {
		t.Errorf("exp 实际得到结果类型并不是 *ast.Identifier，实际得到=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value 期望得到 %s. 实际得到: %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral 期望得到 %s. 实际得到: %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("表达式不是 ast.InfixExpression 类型，实际获得: %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.LeftExp, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator 期望操作符为: '%s'，实际获得: '%q'", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.RightExp, right) {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errMessages := p.Errors()
	if len(errMessages) == 0 {
		return
	}

	t.Errorf("语法分析存在错误")
	for _, msg := range errMessages {
		t.Errorf("语法法分析错误: %q", msg)
	}
	t.FailNow()
}

func checkLexerErrors(t *testing.T, l *lexer.Lexer) {
	errTokens := l.Errors()
	if len(errTokens) == 0 {
		return
	}

	t.Errorf("词法分析器存在 %d 个错误", len(errTokens))
	for _, errTok := range errTokens {
		t.Errorf("(%s) 词法分析错误: %q", errTok.Literal, errTok.Msg)
	}
	t.FailNow()
}
