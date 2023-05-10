package eval

import (
	"fmt"
	"strconv"
)

func eval(l *lexer, defines map[string]string, visited map[string]bool) bool {
	return orExpr(l, defines, visited)
}

func orExpr(l *lexer, defines map[string]string, visited map[string]bool) bool {
	result := andExpr(l, defines, visited)

	for {
		t := l.Peek()
		switch t.Kind {
		case TokenKindOr:
			l.Read()
			right := andExpr(l, defines, visited)
			result = result || right
		default:
			return result
		}
	}
}

func equalityExpr(l *lexer, defines map[string]string, visited map[string]bool) bool {
	left := primaryExpr(l, defines, visited)

	t := l.Peek()
	if t.Kind == TokenKindEquals {
		l.Read()
		right := primaryExpr(l, defines, visited)
		return left == right
	} else {
		return left
	}
}

func andExpr(l *lexer, defines map[string]string, visited map[string]bool) bool {
	left := equalityExpr(l, defines, visited)

	for {
		t := l.Peek()
		switch t.Kind {
		case TokenKindAnd:
			l.Read()
			right := equalityExpr(l, defines, visited)
			left = left && right
		default:
			return left
		}
	}
}
func primaryExpr(l *lexer, defines map[string]string, visited map[string]bool) bool {
	t := l.Read()
	switch t.Kind {
	case TokenKindID:
		str := string(l.s[t.Start:t.End])
		if visited[str] {
			return false
		}

		v, ok := defines[str]
		if !ok {
			return false
		}

		visited[str] = true

		subLexer := NewLexer([]byte(v))
		result := eval(subLexer, defines, visited)
		visited[str] = false

		return result
	case TokenKindNumber:
		str := string(l.s[t.Start:t.End])

		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false
		}
		return v != 0
	case TokenKindLeftParen:
		result := orExpr(l, defines, visited)

		t = l.Read()
		if t.Kind != TokenKindRightParen {
			return false
		}

		return result
	default:
		fmt.Printf("Unexpected token in primaryExpr: %v\n", t)
		return false
	}
}

func Evaluate(source string, defines map[string]string) bool {
	l := &lexer{
		s: []byte(source),
	}

	visited := map[string]bool{}

	return eval(l, defines, visited)
}
