package eval

import "strconv"

func eval(l *lexer, defines map[string]string, visited map[string]bool) bool {
	for {
		t := l.Read()

		switch t.Kind {
		// TODO: handle && and ||
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

			l.s = []byte(v)
			l.start = 0
			l.end = 0
		case TokenKindNumber:
			str := string(l.s[t.Start:t.End])

			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return false
			}
			return v != 0
		case TokenKindNone:
			return false
		}
	}
}

func Evaluate(source string, defines map[string]string) bool {
	l := &lexer{
		s: []byte(source),
	}

	visited := map[string]bool{}

	return eval(l, defines, visited)
}
