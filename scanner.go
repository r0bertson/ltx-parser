package ltxparser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	if isWhitespace(ch) { // If we see whitespace then consume all contiguous whitespace.
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) { // If we see a letter then consume as a variable or reserved word.
		s.unread()
		return s.scanVariable()
	} else if isDigitOrDot(ch) { // If we see a digit or a dot then consume as a number.
		s.unread()
		return s.scanNumber()
	} else if isOperator(ch) { // If we sen a comparison operator, consume acordingly
		s.unread()
		return s.scanOperator()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return TOKEN_EOF, ""
	case '+', '-':
		return TOKEN_SIGN_OPERATOR, string(ch)
	case ')':
		return TOKEN_RIGHT_PARENTHESIS, string(ch)
	}

	return TOKEN_ILLEGAL, string(ch)
}

func (s *Scanner) scanNumber() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent valid numeric character into the buffer.
	// Anything except for digit characters or dots will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return TOKEN_NUMBER, buf.String()
}

func (s *Scanner) scanOperator() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent relational operator character.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isOperator(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "=":
		return TOKEN_OPERATOR, buf.String()
	case "<=", "<":
		return TOKEN_OPERATOR, "<="
	case ">=", ">":
		return TOKEN_OPERATOR, ">="

	}
	// Otherwise return as illegal token.
	return TOKEN_ILLEGAL, buf.String()
}

// scanWhitespace consumes the current rune and all contiguous whitespace
// including space, tab, or newline.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return TOKEN_WS, buf.String()
}

// scanVariable consumes the current rune and all contiguous variable runes until it reaches max size = 8
func (s *Scanner) scanVariable() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	varSize := 1
	// Read every subsequent valid variable character into the buffer.
	// Non-variable characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isVariableChar(ch) || isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
			varSize++ //increment size
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "MAX", "MAXIMIZE", "MAXIMISE":
		return TOKEN_MAX, buf.String()
	case "MIN", "MINIMIZE", "MINIMISE":
		return TOKEN_MIN, buf.String()
	case "ST", "SUBJECT", "TO", "SUCH", "THAT", "S.T.":
		return TOKEN_ST, buf.String()
	case "END":
		return TOKEN_END, buf.String()
	}
	// If variable max size is respected, returns variable token. Otherwise, returns illegal.
	if varSize <= 8 {
		return TOKEN_VARIABLE, buf.String()
	}
	return TOKEN_ILLEGAL, buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isOperator returns true if the rune is a relational operator.
func isOperator(ch rune) bool { return (ch == '<' || ch == '=' || ch == '>') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// isDigit returns true if the rune is a digit.
func isDigitOrDot(ch rune) bool { return (ch >= '0' && ch <= '9') || ch == '.' }

// isVariableChar returns true if the rune is a valid caracter variable.
func isVariableChar(ch rune) bool {
	switch ch {
	case '!', ')', '+', '-':
		return false
	}
	return !isOperator(ch)
}

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
