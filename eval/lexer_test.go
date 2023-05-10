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

func runLexerTest(t *testing.T, source string, expected []Token) {
	l := &lexer{
		s: []byte(source),
	}

	actual := readAllTokens(l)
	assert.Equal(t, expected, actual)
}

func TestLexerNegativeNumber(t *testing.T) {
	runLexerTest(t, "-123", []Token{
		{Kind: TokenKindNumber, Start: 0, End: 4},
	})
}

func TestLexerEmpty(t *testing.T) {
	runLexerTest(t, "", nil)
}

func TestLexerWhitespace(t *testing.T) {
	runLexerTest(t, " ", nil)
}

func TestLexerID(t *testing.T) {
	runLexerTest(t, "abc", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
	})
}

func TestLexerIDWithWhitespace(t *testing.T) {
	runLexerTest(t, " abc ", []Token{
		{Kind: TokenKindID, Start: 1, End: 4},
	})
}

func TestLexerIDWithUnderscore(t *testing.T) {
	runLexerTest(t, "abc_def", []Token{
		{Kind: TokenKindID, Start: 0, End: 7},
	})
}

func TestLexerIDWithNumber(t *testing.T) {
	runLexerTest(t, "abc123", []Token{
		{Kind: TokenKindID, Start: 0, End: 6},
	})
}

func TestLexerIDWithNumberAndUnderscore(t *testing.T) {
	runLexerTest(t, "abc_123", []Token{
		{Kind: TokenKindID, Start: 0, End: 7},
	})
}

func TestLexerIDAndID(t *testing.T) {
	runLexerTest(t, "abc && def", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindAnd, Start: 4, End: 6},
		{Kind: TokenKindID, Start: 7, End: 10},
	})
}

func TestLexerIDOrID(t *testing.T) {
	runLexerTest(t, "abc || def", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindOr, Start: 4, End: 6},
		{Kind: TokenKindID, Start: 7, End: 10},
	})
}

func TestLexerIDAndIDOrID(t *testing.T) {
	runLexerTest(t, "abc && def || ghi", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindAnd, Start: 4, End: 6},
		{Kind: TokenKindID, Start: 7, End: 10},
		{Kind: TokenKindOr, Start: 11, End: 13},
		{Kind: TokenKindID, Start: 14, End: 17},
	})
}

func TestLexerIDAndIDOrIDInParentheses(t *testing.T) {
	runLexerTest(t, "abc && (def || ghi)", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindAnd, Start: 4, End: 6},
		{Kind: TokenKindLeftParen, Start: 7, End: 8},
		{Kind: TokenKindID, Start: 8, End: 11},
		{Kind: TokenKindOr, Start: 12, End: 14},
		{Kind: TokenKindID, Start: 15, End: 18},
		{Kind: TokenKindRightParen, Start: 18, End: 19},
	})
}

func TestLexerIDAndIDOrIDInParenthesesAndID(t *testing.T) {
	runLexerTest(t, "abc && (def || ghi) && jkl", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindAnd, Start: 4, End: 6},
		{Kind: TokenKindLeftParen, Start: 7, End: 8},
		{Kind: TokenKindID, Start: 8, End: 11},
		{Kind: TokenKindOr, Start: 12, End: 14},
		{Kind: TokenKindID, Start: 15, End: 18},
		{Kind: TokenKindRightParen, Start: 18, End: 19},
		{Kind: TokenKindAnd, Start: 20, End: 22},
		{Kind: TokenKindID, Start: 23, End: 26},
	})
}

func TestLexerNumber(t *testing.T) {
	runLexerTest(t, "123", []Token{
		{Kind: TokenKindNumber, Start: 0, End: 3},
	})
}

func TestLexerNumberWithWhitespace(t *testing.T) {
	runLexerTest(t, " 123 ", []Token{
		{Kind: TokenKindNumber, Start: 1, End: 4},
	})
}

func TestLexerIDEqualsNumber(t *testing.T) {
	runLexerTest(t, "abc == 123", []Token{
		{Kind: TokenKindID, Start: 0, End: 3},
		{Kind: TokenKindEquals, Start: 4, End: 6},
		{Kind: TokenKindNumber, Start: 7, End: 10},
	})
}
