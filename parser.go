package ltxparser

import (
	"fmt"
	"io"
	"strconv"
)

// LinearProblem represents a Linear Problem.
type LinearProblem struct {
	ObjectiveFunction OF
	Constraints       []Constraint
}

// OF represents an objective function
type OF struct {
	OFType    string
	Variables []Variable
}

//Constraint represent an LP constraint
type Constraint struct {
	Name     string
	LH       []Variable
	Operator string
	RH       float64
}

//Variable holds the variable name and its coefficient
type Variable struct {
	Name        string
	Coefficient float64
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a LINDO's text file (*.ltx) Linear Problem.
func (p *Parser) Parse() (*LinearProblem, error) {
	lp := &LinearProblem{}

	// First token should be a "MAX" or "MIN" keyword.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != TOKEN_MAX && tok != TOKEN_MIN {
		return nil, fmt.Errorf("found %q, expected MAX or MIN", lit)
	}
	lp.ObjectiveFunction.OFType = handleObjectiveFunctionType(lit)
	// Next we should loop over the objective function.
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != TOKEN_VARIABLE && tok != TOKEN_NUMBER && tok != TOKEN_SIGN_OPERATOR {
			return nil, fmt.Errorf("found %q, expected sign operator, coefficient or variable", lit)
		}
		//IF THERE IS AN EXPLICIT SIGN
		var sign = "+"
		if tok == TOKEN_SIGN_OPERATOR {
			sign = lit
			tok, lit = p.scanIgnoreWhitespace()
		}

		if tok == TOKEN_VARIABLE {
			lp.ObjectiveFunction.Variables = append(lp.ObjectiveFunction.Variables, Variable{lit, handleSign(sign)}) // place 1 or -1 as coefficient
		} else if tok == TOKEN_NUMBER {
			var temp Variable
			if num, err := strconv.ParseFloat(sign+lit, 64); err == nil {
				temp.Coefficient = num
			} else {
				return nil, fmt.Errorf("Error converting %q to float", lit)
			}

			tok, lit := p.scanIgnoreWhitespace()

			if tok == TOKEN_VARIABLE {
				temp.Name = lit
				//append only variables, because constants do not affect the optimal solution
				lp.ObjectiveFunction.Variables = append(lp.ObjectiveFunction.Variables, temp)
			} else {
				p.unscan()
			}
		} else {
			return nil, fmt.Errorf("found %q, expected coefficient or variable", lit)
		}

		// If the next token is not a sign operator, breaks loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_SIGN_OPERATOR {
			p.unscan()
			break
		}
		p.unscan()
	}

	// Check if exists a ST token
	if tok, lit := p.scanIgnoreWhitespace(); tok != TOKEN_ST {
		return nil, fmt.Errorf("found %q, expected ST/SUBJECT TO/... ", lit)
	} else {
		// Handles compound keywords (SUBJECT TO or SUCH THAT)
		if lit == "SUBJECT" {
			if tok, lit := p.scanIgnoreWhitespace(); tok != TOKEN_ST || lit != "TO" {
				return nil, fmt.Errorf("found %q, expected TO ", lit)
			}
		} else if lit == "SUCH" {
			if tok, lit := p.scanIgnoreWhitespace(); tok != TOKEN_ST || lit != "THAT" {
				return nil, fmt.Errorf("found %q, expected THAT ", lit)
			}
		}
	}
	// Next we should loop over the constraints.
	for {
		var cons Constraint
		// Loop through variables on the left hand side of the equation
		for {
			tok, lit := p.scanIgnoreWhitespace()

			if tok != TOKEN_VARIABLE && tok != TOKEN_NUMBER && tok != TOKEN_SIGN_OPERATOR {
				return nil, fmt.Errorf("found %q, expected sign operator, coefficient or variable", lit)
			}
			// Handle explicit sign operator
			var sign = "+"
			if tok == TOKEN_SIGN_OPERATOR {
				sign = lit
				tok, lit = p.scanIgnoreWhitespace()
			}

			if tok == TOKEN_VARIABLE { // No explicit coefficient
				cons.LH = append(cons.LH, Variable{lit, handleSign(sign)}) // place 1 or -1 as coefficient
			} else if tok == TOKEN_NUMBER { // Handle explicit coefficient
				var temp Variable
				if num, err := strconv.ParseFloat(sign+lit, 64); err == nil {
					temp.Coefficient = num
				} else {
					return nil, fmt.Errorf("Error converting %q to float", lit)
				}

				tok, lit := p.scanIgnoreWhitespace()
				if tok == TOKEN_VARIABLE {
					temp.Name = lit
				} else {
					return nil, fmt.Errorf("found constant %q on the lefthand side of a constraint", lit)
				}
				cons.LH = append(cons.LH, temp)
			} else {
				return nil, fmt.Errorf("found %q, expected coefficient or variable", lit)
			}
			// If the next token is not a sign operator, breaks the loop.
			if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_SIGN_OPERATOR {
				p.unscan()
				break
			}
			p.unscan()
		}

		// Expects operator splitting LH and RH

		if tok, lit := p.scanIgnoreWhitespace(); tok == TOKEN_OPERATOR {
			cons.Operator = lit
		} else {
			return nil, fmt.Errorf("found constant %q instead of an operator", lit)
		}

		//Right hand side of the equation must be a constant (sign is not mandatory)

		var sign = "+"
		tok, lit := p.scanIgnoreWhitespace()
		if tok == TOKEN_SIGN_OPERATOR {
			sign = lit
			tok, lit = p.scanIgnoreWhitespace()
		}

		if tok == TOKEN_NUMBER {
			if num, err := strconv.ParseFloat(sign+lit, 64); err == nil {
				cons.RH = num
			} else {
				return nil, fmt.Errorf("Error converting %q to float", lit)
			}
			lp.Constraints = append(lp.Constraints, cons)
		} else {
			return nil, fmt.Errorf("found %q instead of a constant", lit)
		}

		// If the next token is END, BREAK THE LOOP.
		if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_NUMBER && tok != TOKEN_VARIABLE && tok != TOKEN_SIGN_OPERATOR {
			p.unscan()
			break
		}
		p.unscan()
	}

	// Must have END token at the end of the file
	if tok, lit := p.scanIgnoreWhitespace(); tok != TOKEN_END {
		return nil, fmt.Errorf("found %q, expected END", lit)
	}

	return lp, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == TOKEN_WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// handleSign returns a numeric coefficient based on a variable sign operator
func handleSign(signal string) float64 {
	if signal == "-" {
		return -1.0
	}
	return 1.0
}

// handleOFType returns the most concise alias (MIN/MAX) for an objective function type.
func handleObjectiveFunctionType(token string) string {
	if token == "MIN" || token == "MINIMIZE" || token == "MINIMISE" {
		return "MIN"
	}
	return "MAX"
}
