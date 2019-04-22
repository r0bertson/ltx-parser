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

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigitOrDot(ch) {
		s.unread()
		return s.scanNumber()
	} else if isOperator(ch) {
		s.unread()
		return s.scanOperator()
	}

	// Otherwise read the individual character.
	//TODO: HANDLE <= AND >=
	switch ch {
	case eof:
		return TOKEN_EOF, ""
	case '+', '-':
		return TOKEN_SIGN_OPERATOR, string(ch)
	}
	return TOKEN_ILLEGAL, string(ch)
}

func (s *Scanner) scanNumber() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent valid numeric character into the buffer.
	// Anything except for digit characters, dots or commas will cause the loop to exit.
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

	// Otherwise return as a regular identifier.
	return TOKEN_NUMBER, buf.String()
}

func (s *Scanner) scanOperator() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent relational operator caractere.
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
	// Otherwise return as a regular identifier.
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

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
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
	// Otherwise return as a regular identifier with its constraints (no dots inside variables)
	return TOKEN_VARIABLE, strings.Replace(buf.String(), ".", "", -1)
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

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
