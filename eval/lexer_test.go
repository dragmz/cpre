package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func readAllTokens(l *lexer) []Token {
	var tokens []Token
	for {
		t := l.Read()
		if t.Kind == TokenKindNone {
			break
		}
		tokens = append(tokens, t)
	}
	return tokens
}

func TestLexer(t *testing.T) {
	type test struct {
		name   string
		source string
		tokens []Token
	}

	tests := []test{
		{
			name:   "negative number",
			source: "-123",
			tokens: []Token{
				{Kind: TokenKindNumber, Start: 0, End: 4},
			},
		},
		{
			name:   "empty",
			source: "",
		},
		{
			name:   "whitespace",
			source: " ",
		},
		{
			name:   "id",
			source: "abc",
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
			}},
		{
			name:   "id with whitespace",
			source: " abc ",
			tokens: []Token{
				{Kind: TokenKindID, Start: 1, End: 4},
			},
		},
		{
			name:   "id with underscore",
			source: "abc_def",
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 7},
			},
		},
		{
			name:   "id with number",
			source: "abc123",
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 6},
			},
		},
		{
			name:   "id with number and underscore",
			source: "abc_123",
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 7},
			},
		},
		{
			name:   "id && id",
			source: `abc && def`,
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
				{Kind: TokenKindAnd, Start: 4, End: 6},
				{Kind: TokenKindID, Start: 7, End: 10},
			},
		},
		{
			name:   "id || id",
			source: `abc || def`,
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
				{Kind: TokenKindOr, Start: 4, End: 6},
				{Kind: TokenKindID, Start: 7, End: 10},
			},
		},
		{
			name:   "id && id || id",
			source: `abc && def || ghi`,
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
				{Kind: TokenKindAnd, Start: 4, End: 6},
				{Kind: TokenKindID, Start: 7, End: 10},
				{Kind: TokenKindOr, Start: 11, End: 13},
				{Kind: TokenKindID, Start: 14, End: 17},
			},
		},
		{
			name:   "id && (id || id)",
			source: `abc && (def || ghi)`,
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
				{Kind: TokenKindAnd, Start: 4, End: 6},
				{Kind: TokenKindLeftParen, Start: 7, End: 8},
				{Kind: TokenKindID, Start: 8, End: 11},
				{Kind: TokenKindOr, Start: 12, End: 14},
				{Kind: TokenKindID, Start: 15, End: 18},
				{Kind: TokenKindRightParen, Start: 18, End: 19},
			},
		},
		{
			name:   "id && (id || id) && id",
			source: `abc && (def || ghi) && jkl`,
			tokens: []Token{
				{Kind: TokenKindID, Start: 0, End: 3},
				{Kind: TokenKindAnd, Start: 4, End: 6},
				{Kind: TokenKindLeftParen, Start: 7, End: 8},
				{Kind: TokenKindID, Start: 8, End: 11},
				{Kind: TokenKindOr, Start: 12, End: 14},
				{Kind: TokenKindID, Start: 15, End: 18},
				{Kind: TokenKindRightParen, Start: 18, End: 19},
				{Kind: TokenKindAnd, Start: 20, End: 22},
				{Kind: TokenKindID, Start: 23, End: 26},
			},
		},
		{
			name:   "number",
			source: "123",
			tokens: []Token{
				{Kind: TokenKindNumber, Start: 0, End: 3},
			},
		},
		{
			name:   "number with whitespace",
			source: " 123 ",
			tokens: []Token{
				{Kind: TokenKindNumber, Start: 1, End: 4},
			},
		},
	}

	for _, test := range tests {
		l := &lexer{
			s: []byte(test.source),
		}

		actual := readAllTokens(l)
		assert.Equal(t, test.tokens, actual, test.name)
	}
}
