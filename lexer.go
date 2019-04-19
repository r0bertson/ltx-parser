package lexer

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
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	//TODO: HANDLE <= AND >=
	switch ch {
	case eof:
		return TOKEN_EOF, ""
	case '<', '>', '=':
		return TOKEN_OPERATOR, string(ch)
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
		} else if isDigit(ch) && ch != '.' && ch != ',' { //HANDLE COMMA ON PARSING
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return TOKEN_NUMBER, buf.String()
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
		} else if !isLetter(ch) && /*!isDigit(ch) &&*/ ch != '_' {
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
	//TODO: HANDLE 2 SEPARATE WORDS AS KEYWORDS LIKE "SUBJECT TO"
	case "ST":
		return TOKEN_ST, buf.String()
	}

	// Otherwise return as a regular identifier.
	return TOKEN_VARIABLE, buf.String()
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

// is Operator returns true if rune is an operator of addition or subtraction (+ or -)
func isOperator(ch rune) bool { return (ch == '-' || ch == '+') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
