package lexer

import (
	"unicode"

	"github.com/Skarlso/horcsog/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char) PEEK
	ch           byte // current char under examination // change this to RUNE to support UTF8
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// give us the next character and advance our position in teh input string
// this read only supports ASCII characters.
// for unicode this would have to be rune and `++` for position wouldn't work
// anymore and some other methods.
func (l *Lexer) readChar() {
	// reached end of the input
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NUL // we haven't read anything yet or end of file
	} else {
		l.ch = l.input[l.readPosition] // next character
	}

	l.position = l.readPosition // update current location to next location
	l.readPosition++            // advance next location
}

// Look at the current character and return its token representation.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// save the current character because we move ahead
			// normally, this is clever, because we might abstract
			// this method to some kind of double character lookup.
			ch := l.ch
			l.readChar()                                        // move ahead
			literal := string(ch) + string(l.ch)                // add = to = that we read before
			tok = token.Token{Type: token.EQ, Literal: literal} // could be simplified with just Literal: "==".
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar() // move one ahead
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier(isLetter)
			tok.Type = token.LookupIdent(tok.Literal) // look up the type

			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Type = token.INT
			tok.Literal = l.readIdentifier(isNumber)

			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

// or eatWhitespace or consumeWhitespace
// read until we skip all whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// // reads in an identifier and advances the pointer until it's not a letter anymore.
// func (l *Lexer) readIdentifier() string {
// 	position := l.position
// 	for isLetter(l.ch) {
// 		l.readChar()
// 	}

// 	return l.input[position:l.position]
// }

// // this could be made better if we pass in the func and call it readSomething
// func (l *Lexer) readNumber() string {
// 	position := l.position
// 	for unicode.IsDigit(rune(l.ch)) {
// 		l.readChar()
// 	}

// 	return l.input[position:l.position]
// }

// this could be made better if we pass in the func and call it readSomething
func (l *Lexer) readIdentifier(fn func(byte) bool) string {
	position := l.position
	for fn(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// take a look at our read position to see what comes next and decide on what to do.
// Sometimes you also have to look backwards!
func (l *Lexer) peekChar() byte {
	// TODO: this would be a seek of some sort for the reader...
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// will ultimately identify what we tolerate in var names like ~ or ! or ? even.
func isLetter(ch byte) bool {
	// the _ allows identifiers such as foo_bar.
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// will ultimately identify what we tolerate in var names like ~ or ! or ? even.
// TODO: Extend this to float numbers.
// We are going to need modulo as well.
func isNumber(ch byte) bool {
	// the _ allows identifiers such as foo_bar.
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
