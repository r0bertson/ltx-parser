package ltxparser_test

import (
	"reflect"
	"strings"
	"testing"

	ltx "github.com/r0bertson/ltx-parser"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s   string
		lp  *ltx.LinearProblem
		err string
	}{
		// Original
		{
			s: `MAX X1 + 2 X2 SUCH THAT X1 + 2 X2 < 6 3 X1 + 2X2 <= 18 END`,
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MAX", Variables: []ltx.Variable{{"X1", 1}, {"X2", 2}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X1", 1}, {"X2", 2}}, Operator: "<", RH: 6.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X1", 3}, {"X2", 2}}, Operator: "<=", RH: 18.0}},
			},
			err: "",
		},
		{
			s: `MIN X234Y + 2Y 
				SUBJECT TO 
					X234Y + 2Y <= 6 
					3 X234Y + 2Y >= 18 
				END`,
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MIN", Variables: []ltx.Variable{{"X234Y", 1}, {"Y", 2}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X234Y", 1}, {"Y", 2}}, Operator: "<=", RH: 6.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X234Y", 3}, {"Y", 2}}, Operator: ">=", RH: 18.0}},
			},
			err: "",
		},
		{
			s: `MAXIMIZE X + 2Y 
				SUCH THAT 
					3X + Y = 6 
					3 X - 2Y < 18 
				END`,
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MAX", Variables: []ltx.Variable{{"X", 1}, {"Y", 2}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X", 3}, {"Y", 1}}, Operator: "=", RH: 6.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X", 3}, {"Y", -2}}, Operator: "<", RH: 18.0}},
			},
			err: "",
		},
	}

	for i, tt := range tests {
		lp, err := ltx.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.lp, lp) {
			t.Errorf("%d. %q\n\nlp mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.lp, lp)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
