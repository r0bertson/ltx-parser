package ltxparser_test

import (
	"strings"
	"testing"

	ltxparser "github.com/r0bertson/ltx-parser"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok ltxparser.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: ltxparser.TOKEN_EOF},
		{s: `#`, tok: ltxparser.TOKEN_ILLEGAL, lit: `#`},
		{s: ` `, tok: ltxparser.TOKEN_WS, lit: " "},
		{s: "\t", tok: ltxparser.TOKEN_WS, lit: "\t"},
		{s: "\n", tok: ltxparser.TOKEN_WS, lit: "\n"},
		// Operators characters
		{s: `<`, tok: ltxparser.TOKEN_OPERATOR, lit: "<"},
		{s: `<=`, tok: ltxparser.TOKEN_OPERATOR, lit: "<="},
		{s: `>`, tok: ltxparser.TOKEN_OPERATOR, lit: ">"},
		{s: `>=`, tok: ltxparser.TOKEN_OPERATOR, lit: ">="},
		{s: `=`, tok: ltxparser.TOKEN_OPERATOR, lit: "="},
		{s: `+`, tok: ltxparser.TOKEN_SIGN_OPERATOR, lit: "+"},
		{s: `-`, tok: ltxparser.TOKEN_SIGN_OPERATOR, lit: "-"},
		{s: `==`, tok: ltxparser.TOKEN_ILLEGAL, lit: "=="},
		{s: `=>`, tok: ltxparser.TOKEN_ILLEGAL, lit: "=>"},
		{s: `=<`, tok: ltxparser.TOKEN_ILLEGAL, lit: "=<"},
		{s: `<<`, tok: ltxparser.TOKEN_ILLEGAL, lit: "<<"},
		// Identifiers
		{s: `foo`, tok: ltxparser.TOKEN_VARIABLE, lit: `foo`},
		{s: `ZxU1234X`, tok: ltxparser.TOKEN_VARIABLE, lit: `ZxU1234X`},
		{s: `ZxU1_234X`, tok: ltxparser.TOKEN_VARIABLE, lit: `ZxU1`},
		{s: `X2`, tok: ltxparser.TOKEN_VARIABLE, lit: `X2`},
		// Numbers
		{s: `3`, tok: ltxparser.TOKEN_NUMBER, lit: `3`},
		{s: `3.4`, tok: ltxparser.TOKEN_NUMBER, lit: `3.4`},
		{s: `.4`, tok: ltxparser.TOKEN_NUMBER, lit: `.4`},
		// Keywords
		{s: `MAX`, tok: ltxparser.TOKEN_MAX, lit: "MAX"},
		{s: `MAXIMIZE`, tok: ltxparser.TOKEN_MAX, lit: "MAXIMIZE"},
		{s: `MAXIMISE`, tok: ltxparser.TOKEN_MAX, lit: "MAXIMISE"},
		{s: `MIN`, tok: ltxparser.TOKEN_MIN, lit: "MIN"},
		{s: `MINIMIZE`, tok: ltxparser.TOKEN_MIN, lit: "MINIMIZE"},
		{s: `MINIMISE`, tok: ltxparser.TOKEN_MIN, lit: "MINIMISE"},
		{s: `ST`, tok: ltxparser.TOKEN_ST, lit: "ST"},
		{s: `S.T.`, tok: ltxparser.TOKEN_ST, lit: "S.T."},
		{s: `END`, tok: ltxparser.TOKEN_END, lit: "END"},
		// Compound keywords are handled on parser
		{s: `SUBJECT`, tok: ltxparser.TOKEN_ST, lit: "SUBJECT"},
		{s: `TO`, tok: ltxparser.TOKEN_ST, lit: "TO"},
		{s: `SUCH`, tok: ltxparser.TOKEN_ST, lit: "SUCH"},
		{s: `THAT`, tok: ltxparser.TOKEN_ST, lit: "THAT"},
	}

	for i, tt := range tests {
		s := ltxparser.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
