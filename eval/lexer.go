package eval

import "unicode/utf8"

type lexer struct {
	s     []byte
	start int
	end   int
}

type TokenKind int

const (
	TokenKindNone TokenKind = iota
	TokenKindID
	TokenKindNumber
	TokenKindLeftParen
	TokenKindRightParen
	TokenKindAnd
	TokenKindOr
	TokenKindOther
)

type Token struct {
	Kind  TokenKind
	Start int
	End   int
}

func NewLexer(source []byte) *lexer {
	return &lexer{
		s: source,
	}
}

func (l *lexer) Read() Token {
	l.skipWhitespace()

	if l.end >= len(l.s) {
		l.start = l.end
		return Token{}
	}

	r, w := utf8.DecodeRune(l.s[l.end:])
	if r == '_' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
		return l.readID()
	}

	if r == '-' || ('0' <= r && r <= '9') {
		return l.readNumberOrOther()
	}

	if r == '&' && l.end+w < len(l.s) {
		r2, w2 := utf8.DecodeRune(l.s[l.end+w:])
		if r2 == '&' {
			l.end += w + w2

			start := l.start
			l.start = l.end

			return Token{
				Kind:  TokenKindAnd,
				Start: start,
				End:   l.end,
			}
		}
	}

	if r == '|' && l.end+w < len(l.s) {
		r2, w2 := utf8.DecodeRune(l.s[l.end+w:])
		if r2 == '|' {
			l.end += w + w2

			start := l.start
			l.start = l.end

			return Token{
				Kind:  TokenKindOr,
				Start: start,
				End:   l.end,
			}
		}
	}

	switch r {
	case '(':
		l.end += w

		start := l.start
		l.start = l.end

		return Token{
			Kind:  TokenKindLeftParen,
			Start: start,
			End:   l.end,
		}
	case ')':
		l.end += w

		start := l.start
		l.start = l.end

		return Token{
			Kind:  TokenKindRightParen,
			Start: start,
			End:   l.end,
		}
	default:
		return l.readOther()
	}
}

func (l *lexer) skipWhitespace() {
	for l.end < len(l.s) {
		r, w := utf8.DecodeRune(l.s[l.end:])
		if r == ' ' || r == '\t' {
			l.end += w
		} else {
			l.start = l.end
			break
		}
	}
}

func (l *lexer) readID() Token {
	for l.end < len(l.s) {
		r, w := utf8.DecodeRune(l.s[l.end:])

		if r == '_' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') {
			l.end += w
		} else {
			break
		}
	}

	start := l.start
	l.start = l.end

	return Token{
		Kind:  TokenKindID,
		Start: start,
		End:   l.end,
	}
}

func (l *lexer) readOther() Token {
loop:
	for l.end < len(l.s) {
		r, w := utf8.DecodeRune(l.s[l.end:])

		switch r {
		case ' ', '\t', '\n':
			break loop
		}

		l.end += w
	}

	start := l.start
	l.start = l.end

	return Token{
		Kind:  TokenKindOther,
		Start: start,
		End:   l.end,
	}
}

func (l *lexer) readNumberOrOther() Token {
	// read first rune and make sure it is a digit or -

	r, w := utf8.DecodeRune(l.s[l.end:])
	if r == '-' {
		// read next and make sure it's digit
		r, w = utf8.DecodeRune(l.s[l.end+w:])
		if r < '0' || '9' < r {
			return l.readOther()
		}
		l.end += w
	}

	// Read the remaining digits
	for l.end < len(l.s) {
		r, w := utf8.DecodeRune(l.s[l.end:])
		if '0' <= r && r <= '9' {
			l.end += w
		} else if r == ' ' || r == '\t' || r == '\n' {
			break
		} else {
			return l.readOther()
		}
	}

	start := l.start
	l.start = l.end

	return Token{
		Kind:  TokenKindNumber,
		Start: start,
		End:   l.end,
	}
}
