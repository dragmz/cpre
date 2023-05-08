package cpre

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dragmz/cpre/eval"
)

type blockKind int

const (
	blockTypeNone blockKind = iota
	blockTypeConditional
	blockTypeUnconditional
)

type block struct {
	parent *block
	kind   blockKind
	skip   bool
	value  bool
}

type Includer func(filePath string, global bool) (id string, source []byte, err error)

type Preprocessor struct {
	defines map[string]string
	stack   *block

	include  Includer
	includes map[string]bool
}

type PreprocessorConfig struct {
	Include Includer
}

func NewPreprocessor(config PreprocessorConfig) *Preprocessor {
	return &Preprocessor{
		defines: make(map[string]string),
		stack:   &block{},

		include:  config.Include,
		includes: map[string]bool{},
	}
}

func (p *Preprocessor) Define(id, value string) {
	p.defines[id] = value
}

func (p *Preprocessor) Undefine(id string) {
	delete(p.defines, id)
}

func (p *Preprocessor) Process(source string) string {
	source = strings.ReplaceAll(source, "\r\n", "\n")

	bs := []byte(source)

	s := &state{
		p:     p,
		s:     bs,
		start: 0,
		end:   0,
	}

	return s.process()
}

func (p *Preprocessor) push() *block {
	previous := p.stack

	p.stack = &block{
		parent: previous,
		skip:   previous.skip,
	}

	return previous
}

func (p *Preprocessor) pop() *block {
	previous := p.stack
	p.stack = p.stack.parent
	return previous
}

type state struct {
	p *Preprocessor
	s []byte

	start int
	end   int

	once bool
}

func (p *state) skipWhitespace() {
	for p.end < len(p.s) {
		r, w := utf8.DecodeRune(p.s[p.end:])
		if r == ' ' || r == '\t' {
			p.end += w
		} else {
			p.start = p.end
			break
		}
	}
}

func (p *state) readInt() (int64, bool) {
	p.start = p.end

	r, w := utf8.DecodeRune(p.s[p.end:])
	if r == '-' {
		p.end += w
	}

loop:

	for p.end < len(p.s) {
		r, w := utf8.DecodeRune(p.s[p.end:])
		if !(r >= '0' && r <= '9') {
			switch r {
			case ' ', '\t', '\n':
				break loop
			default:
				return 0, false
			}
		}
		p.end += w
	}

	v := string(p.s[p.start:p.end])

	iv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return iv, true
}

func (p *state) readID() string {
	for p.end < len(p.s) {
		r, w := utf8.DecodeRune(p.s[p.end:])

		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (r >= '0' && r <= '9') {
			p.end += w
		} else {
			break
		}
	}

	idbs := p.s[p.start:p.end]
	id := string(idbs)

	return id
}

func (p *state) readToEOL() string {
	for p.end < len(p.s) {
		r, w := utf8.DecodeRune(p.s[p.end:])

		if r == '\n' {
			break
		}
		p.end += w
	}

	return string(p.s[p.start:p.end])
}

func (s *state) process() string {
	bol := true

	for s.end < len(s.s) {
		r, w := utf8.DecodeRune(s.s[s.end:])

		switch r {
		case '\n':
			bol = true
			s.end += w
		case '#':
			start := s.end
			s.end += w
			if s.end == len(s.s) {
				break
			}

			if !bol {
				break
			}

			clear := func() {
				s.s = append(s.s[:start], s.s[s.end:]...)
				s.end = start
			}

			s.skipWhitespace()
			directive := s.readID()
			s.skipWhitespace()

			switch directive {
			case "pragma":
				if s.p.stack.skip {
					s.readToEOL()
					clear()
					break
				}

				id := s.readID()

				switch id {
				case "once":
					s.once = true
					clear()
				}

			case "define":
				if s.p.stack.skip {
					s.readToEOL()
					clear()
					break
				}

				id := s.readID()

				s.skipWhitespace()

				s.start = s.end

				for s.end < len(s.s) {
					r, w := utf8.DecodeRune(s.s[s.end:])
					if r == '\n' {
						break
					} else {
						if r == '/' {
							r, w2 := utf8.DecodeRune(s.s[s.end+w:])
							if r == '/' {
								break
							}

							s.end += w

							if r == '*' {
								s.end += w2
								for s.end < len(s.s) {
									r, w := utf8.DecodeRune(s.s[s.end:])
									s.end += w

									if r == '*' {
										r, w := utf8.DecodeRune(s.s[s.end:])
										s.end += w
										if r == '/' {
											break
										}
									}
								}
							}
						} else {
							s.end += w
						}
					}
				}

				value := string(s.s[s.start:s.end])
				s.p.defines[id] = value
				clear()
			case "undef":
				if s.p.stack.skip {
					s.readToEOL()
					clear()
					break
				}

				id := s.readID()
				delete(s.p.defines, id)
				clear()
			case "if":
				s.skipWhitespace()

				s.p.push()

				if !s.p.stack.skip {
					value := s.readToEOL()
					s.p.stack.value = eval.Evaluate(value, s.p.defines)
				}

				s.p.stack.skip = !s.p.stack.value

				clear()
			case "ifdef":
				s.p.push()

				s.skipWhitespace()
				value := s.readToEOL()

				if !s.p.stack.skip {
					l := eval.NewLexer([]byte(value))
					t := l.Read()
					switch t.Kind {
					case eval.TokenKindID:
						_, ok := s.p.defines[value[t.Start:t.End]]
						s.p.stack.value = ok
					}
				}

				s.p.stack.skip = !s.p.stack.value

				clear()
			case "ifndef":
				s.p.push()

				s.skipWhitespace()
				value := s.readToEOL()

				if !s.p.stack.skip {
					l := eval.NewLexer([]byte(value))
					t := l.Read()
					switch t.Kind {
					case eval.TokenKindID:
						_, ok := s.p.defines[value[t.Start:t.End]]
						s.p.stack.value = !ok
					}
				}

				s.p.stack.skip = !s.p.stack.value

				clear()
			case "else":
				prev := s.p.pop()

				s.p.push()
				s.p.stack.value = !prev.value
				s.p.stack.skip = !s.p.stack.value
				clear()
			case "elif":
				prev := s.p.pop()
				s.p.push()

				s.skipWhitespace()
				value := s.readToEOL()

				if !s.p.stack.skip {
					if !prev.value {
						s.p.stack.value = eval.Evaluate(value, s.p.defines)
					}
				}

				s.p.stack.skip = !s.p.stack.value

				clear()
			case "endif":
				s.p.pop()
				clear()
			case "include":
				if s.p.stack.skip {
					s.readToEOL()
					clear()
					break
				}

				s.skipWhitespace()

				r, w := utf8.DecodeRune(s.s[s.end:])
				s.end += w

				var path string
				global := false

				if r == '"' {
					s.start = s.end
					for s.end < len(s.s) {
						r, w := utf8.DecodeRune(s.s[s.end:])

						if r == '"' {
							path = string(s.s[s.start:s.end])
							s.end += w
							break
						}

						s.end += w
					}
				} else if r == '<' {
					global = true

					s.start = s.end
					for s.end < len(s.s) {
						r, w := utf8.DecodeRune(s.s[s.end:])

						if r == '>' {
							path = string(s.s[s.start:s.end])
							s.end += w
							break
						}

						s.end += w
					}
				}

				id, bs, err := s.p.include(path, global)
				if err != nil {
					clear()
					break
				}

				source := strings.ReplaceAll(string(bs), "\r\n", "\n")

				is := &state{
					p: s.p,
					s: []byte(source),
				}

				processed := is.process()

				if is.once && s.p.includes[id] {
					clear()
					break
				}

				s.p.includes[id] = true

				s.s = append(s.s[:start], append([]byte(processed), s.s[s.end:]...)...)
				s.end = start + len(processed)
			default:
				// not a preprocessor directive
				continue
			}
		case '/':
			start := s.end
			s.end += w
			if s.end == len(s.s) {
				break
			}

			r, w := utf8.DecodeRune(s.s[s.end:])

			switch r {
			case '/':
				bol = false
				s.end += w
				for s.end < len(s.s) {
					r, w := utf8.DecodeRune(s.s[s.end:])
					s.end += w

					if r == '\n' {
						break
					}

					if s.end == len(s.s) {
						break
					}
				}
				if s.p.stack.skip {
					s.s = append(s.s[:start], s.s[s.end:]...)
					s.end = start
				}
			case '*':
				s.end += w

				for s.end < len(s.s) {
					r, w := utf8.DecodeRune(s.s[s.end:])
					s.end += w

					if r == '*' {
						r, w := utf8.DecodeRune(s.s[s.end:])
						s.end += w

						if r == '/' {
							break
						}
					}
				}
				if s.p.stack.skip {
					s.s = append(s.s[:start], s.s[s.end:]...)
					s.end = start
				}
			}
		default:
			if s.p.stack.skip {
				s.start = s.end
				s.readToEOL()
				s.s = append(s.s[:s.start], s.s[s.end:]...)
				s.end = s.start
			} else {
				s.start = s.end
				s.end += w

				visited := map[string]bool{}
				for {
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' {
						id := s.readID()

						v, ok := s.p.defines[id]

						if ok {
							if visited[id] {
								break
							}

							s.s = append(s.s[:s.start], append([]byte(v), s.s[s.end:]...)...)
							s.end = s.start

							visited[id] = true

							r, _ = utf8.DecodeRune(s.s[s.end:])
							continue
						}
					}
					break
				}
			}
		}
	}

	result := string(s.s)
	return result
}
