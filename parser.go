package lexer

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

type OF struct {
	OFType    string
	variables []Variable
}
type Constraint struct {
	name     string
	LH       []Variable
	Operator string
	RH       Variable
}

type Variable struct {
	name        string
	coefficient float64
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

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*LinearProblem, error) {
	lp := &LinearProblem{}

	// First token should be a "MAX" or "MIN" keyword.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != TOKEN_MAX && tok != TOKEN_MIN {
		return nil, fmt.Errorf("found %q, expected MAX or MIN", lit)
	}
	lp.ObjectiveFunction.OFType = lit
	// Next we should loop over the objective function.
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != TOKEN_VARIABLE && tok != TOKEN_NUMBER {
			return nil, fmt.Errorf("found %q, expected coefficient or variable", lit)
		}
		if tok == TOKEN_VARIABLE {
			lp.ObjectiveFunction.variables = append(lp.ObjectiveFunction.variables, Variable{lit, 1.0})
		} else if tok == TOKEN_NUMBER {
			var temp Variable
			if num, err := strconv.ParseFloat(lit, 64); err != nil {
				temp.coefficient = num
			} else {
				return nil, fmt.Errorf("Error converting %q to float", lit)
			}

			tok, lit := p.scanIgnoreWhitespace()

			if tok == TOKEN_VARIABLE {
				temp.name = lit
				//append only variables, because constants do not affect the optimal solution
				lp.ObjectiveFunction.variables = append(lp.ObjectiveFunction.variables, temp)
			} else {
				p.unscan()
			}

		}

		// If the next token is ST, BREAK THE LOOP.
		if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_NUMBER && tok != TOKEN_VARIABLE {
			p.unscan()
			break
		}
	}

	//CHECK ST TOKEN
	if tok, lit := p.scanIgnoreWhitespace(); tok != TOKEN_ST {
		return nil, fmt.Errorf("found %q, expected ST/SUBJECT TO/... ", lit)
	}
	// Next we should loop over the constraints.
	for {
		//loop the variables on the left hand sideo of the equation
		for {
			tok, lit := p.scanIgnoreWhitespace()
			var cons Constraint
			if tok != TOKEN_VARIABLE && tok != TOKEN_NUMBER {
				return nil, fmt.Errorf("found %q, expected coefficient or variable", lit)
			}
			if tok == TOKEN_VARIABLE {
				cons.LH = append(cons.LH, Variable{lit, 1.0})
			} else if tok == TOKEN_NUMBER {
				var temp Variable
				if num, err := strconv.ParseFloat(lit, 64); err != nil {
					temp.coefficient = num
				} else {
					return nil, fmt.Errorf("Error converting %q to float", lit)
				}

				tok, lit := p.scanIgnoreWhitespace()
				if tok == TOKEN_VARIABLE {
					temp.name = lit
				} else {
					return nil, fmt.Errorf("found constant %q on the lefthand side of a constraint", lit)
					p.unscan()
				}
				cons.LH = append(cons.LH, temp)
				// If the next token is not a number of variable.
				if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_NUMBER && tok != TOKEN_VARIABLE {
					p.unscan()
					break
				}
			}
			//expects operator

			if tok, lit := p.scanIgnoreWhitespace(); tok == TOKEN_OPERATOR {
				cons.Operator = lit
			} else {
				return nil, fmt.Errorf("found constant %q instead of an operator", lit)
			}

			//right hand side of the equation must be a constant

			if tok, lit := p.scanIgnoreWhitespace(); tok == TOKEN_NUMBER {
				lp.Constraints = append(lp.Constraints, cons)
			} else {
				return nil, fmt.Errorf("found %q instead of a constant", lit)
			}

		}

		// If the next token is END, BREAK THE LOOP.
		if tok, _ := p.scanIgnoreWhitespace(); tok != TOKEN_NUMBER && tok != TOKEN_VARIABLE {
			p.unscan()
			break
		}
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
