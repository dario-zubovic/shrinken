// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"shrinken/SDDL/token"
)

const (
	NoState    = -1
	NumStates  = 155
	NumSymbols = 189
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = new(token.Token)
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '-'
1: '0'
2: '-'
3: '"'
4: '"'
5: 'p'
6: 'a'
7: 'c'
8: 'k'
9: 'a'
10: 'g'
11: 'e'
12: 'u'
13: 's'
14: 'e'
15: 'c'
16: 'l'
17: 'a'
18: 's'
19: 's'
20: '{'
21: '}'
22: ':'
23: 's'
24: 't'
25: 'r'
26: 'u'
27: 'c'
28: 't'
29: 'e'
30: 'n'
31: 'u'
32: 'm'
33: 'i'
34: 'n'
35: 't'
36: 'i'
37: 'n'
38: 't'
39: '3'
40: '2'
41: 'i'
42: 'n'
43: 't'
44: '6'
45: '4'
46: 'l'
47: 'o'
48: 'n'
49: 'g'
50: 's'
51: 'h'
52: 'o'
53: 'r'
54: 't'
55: 'u'
56: 'i'
57: 'n'
58: 't'
59: 'u'
60: 'i'
61: 'n'
62: 't'
63: '3'
64: '2'
65: 'u'
66: 'i'
67: 'n'
68: 't'
69: '6'
70: '4'
71: 'u'
72: 'l'
73: 'o'
74: 'n'
75: 'g'
76: 'u'
77: 's'
78: 'h'
79: 'o'
80: 'r'
81: 't'
82: 'b'
83: 'y'
84: 't'
85: 'e'
86: 'b'
87: 'o'
88: 'o'
89: 'l'
90: 's'
91: 't'
92: 'r'
93: 'i'
94: 'n'
95: 'g'
96: 'c'
97: 'h'
98: 'a'
99: 'r'
100: 'f'
101: 'l'
102: 'o'
103: 'a'
104: 't'
105: 'd'
106: 'o'
107: 'u'
108: 'b'
109: 'l'
110: 'e'
111: '['
112: ']'
113: '['
114: ']'
115: ','
116: '@'
117: 'r'
118: 'a'
119: 'n'
120: 'g'
121: 'e'
122: 'e'
123: 'x'
124: 'p'
125: 'o'
126: 'r'
127: 't'
128: 'A'
129: 's'
130: 'p'
131: 'r'
132: 'e'
133: 'c'
134: 'i'
135: 's'
136: 'i'
137: 'o'
138: 'n'
139: 'v'
140: 'e'
141: 'r'
142: 's'
143: 'i'
144: 'o'
145: 'n'
146: 'm'
147: 'e'
148: 's'
149: 's'
150: 'a'
151: 'g'
152: 'e'
153: '>'
154: '<'
155: 'p'
156: 'i'
157: 'e'
158: '+'
159: '-'
160: '*'
161: '/'
162: '^'
163: 's'
164: 'q'
165: 'r'
166: 't'
167: '('
168: ')'
169: '('
170: '/'
171: '/'
172: '\n'
173: '/'
174: '*'
175: '*'
176: '*'
177: '/'
178: '.'
179: '_'
180: ' '
181: '\t'
182: '\n'
183: '\r'
184: '0'-'9'
185: '1'-'9'
186: 'a'-'z'
187: 'A'-'Z'
188: .
*/
