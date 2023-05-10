package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type evalTest struct {
	source  string
	defines map[string]string
	want    bool
}

func runEvalTest(t *testing.T, tt evalTest) {
	actual := Evaluate(tt.source, tt.defines)
	assert.Equal(t, tt.want, actual)
}

func TestEvaluateEmpty(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "",
		defines: map[string]string{},
		want:    false,
	})
}

func TestEvaluateWhitespace(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  " ",
		defines: map[string]string{},
		want:    false,
	})
}

func TestEvaluateIDUndefined(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "abc",
		defines: map[string]string{},
		want:    false,
	})
}

func TestEvaluateIDDefined(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "abc",
		defines: map[string]string{"abc": "1"},
		want:    true,
	})
}

func TestEvaluateIDDoubleResolve(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "abc",
		defines: map[string]string{"abc": "def", "def": "1"},
		want:    true,
	})
}

func TestEvaluateIDDoubleResolveUndefined(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "abc",
		defines: map[string]string{"abc": "def", "def": "ghi"},
		want:    false,
	})
}

func TestEvaluateCircular(t *testing.T) {
	runEvalTest(t, evalTest{
		source:  "abc",
		defines: map[string]string{"abc": "def", "def": "abc"},
		want:    false,
	})
}

func TestEvaluateAndCondition(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "1",
			"b": "1",
		},
		want: true,
	})
}

func TestEvaluateAndConditionFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "0",
			"b": "1",
		},
		want: false,
	})
}

func TestEvaluateOrCondition(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "1",
			"b": "0",
		},
		want: true,
	})
}

func TestEvaluateOrConditionBothTrue(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "1",
			"b": "1",
		},
		want: true,
	})
}

func TestEvaluateOrConditionOneFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "0",
			"b": "1",
		},
		want: true,
	})
}

func TestEvaluateOrConditionBothFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "0",
			"b": "0",
		},
		want: false,
	})
}

func TestEvaluateAndOrConditionAllTrue(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
		},
		want: true,
	})
}

func TestEvaluateAndOrConditionOneFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "1",
			"b": "0",
			"c": "1",
		},
		want: true,
	})
}

func TestEvaluateAndOrConditionTwoFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "0",
			"b": "1",
			"c": "1",
		},
		want: true,
	})
}

func TestEvaluateAndOrConditionOneAndOrFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "0",
			"b": "0",
			"c": "1",
		},
		want: true,
	})
}

func TestEvaluateAndOrConditionOneTrueTwoFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "0",
			"b": "0",
			"c": "0",
		},
		want: false,
	})
}

func TestEvaluateAndOrConditionTwoTrueOneFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "0",
		},
		want: true,
	})
}

func TestEvaluateAndOrConditionOneTrueOneFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "1",
			"b": "0",
			"c": "0",
		},
		want: false,
	})
}

func TestEvaluateAndOrConditionOneAndOneOrFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "0",
			"b": "1",
			"c": "0",
		},
		want: false,
	})
}

func TestEvaluateAndOrConditionAllFalse(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b || c",
		defines: map[string]string{
			"a": "0",
			"b": "0",
			"c": "0",
		},
		want: false,
	})
}

func TestAndLBOrRB(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && (b || c)",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
		},
		want: true,
	})

	runEvalTest(t, evalTest{
		source: "a && (b || c)",
		defines: map[string]string{
			"a": "1",
			"b": "0",
			"c": "0",
		},
		want: false,
	})
}

func TestOrLBAndRB(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || (b && c)",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
		},
		want: true,
	})

	runEvalTest(t, evalTest{
		source: "a || (b && c)",
		defines: map[string]string{
			"a": "1",
			"b": "0",
			"c": "0",
		},
		want: true,
	})
}

func TestLBOrRBAndLBOrRB(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "(a || b) && (c || d)",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
			"d": "1",
		},
		want: true,
	})

	runEvalTest(t, evalTest{
		source: "(a || b) && (c || d)",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "0",
			"d": "0",
		},
		want: false,
	})

	runEvalTest(t, evalTest{
		source: "(a || b) && (c || d)",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
			"d": "0",
		},
		want: true,
	})
}

func TestLBLBOrRBAndLBOrRBRb(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "((a || b) && (c || d))",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
			"d": "1",
		},
		want: true,
	})

	runEvalTest(t, evalTest{
		source: "((a || b) && (c || d))",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "0",
			"d": "0",
		},
		want: false,
	})
}

func TestA1AndB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "1",
			"b": "0",
		},
		want: false,
	})
}

func TestA1AndB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "1",
			"b": "1",
		},
		want: true,
	})
}

func TestA0AndB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "0",
			"b": "0",
		},
		want: false,
	})
}

func TestA0AndB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b",
		defines: map[string]string{
			"a": "0",
			"b": "1",
		},
		want: false,
	})
}

func TestA1OrB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "1",
			"b": "0",
		},
		want: true,
	})
}

func TestA1OrB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "1",
			"b": "1",
		},
		want: true,
	})
}

func TestA0OrB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "0",
			"b": "0",
		},
		want: false,
	})
}

func TestA0OrB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b",
		defines: map[string]string{
			"a": "0",
			"b": "1",
		},
		want: true,
	})
}

func TestA1EqualsB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b",
		defines: map[string]string{
			"a": "1",
			"b": "1",
		},
		want: true,
	})
}

func TestA1EqualsB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b",
		defines: map[string]string{
			"a": "1",
			"b": "0",
		},
		want: false,
	})
}

func TestA0EqualsB1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b",
		defines: map[string]string{
			"a": "0",
			"b": "1",
		},
		want: false,
	})
}

func TestA0EqualsB0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b",
		defines: map[string]string{
			"a": "0",
			"b": "0",
		},
		want: true,
	})
}

func TestA0OrB0OrC0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b || c",
		defines: map[string]string{
			"a": "0",
			"b": "0",
			"c": "0",
		},
		want: false,
	})
}

func TestA0OrB0OrC1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a || b || c",
		defines: map[string]string{
			"a": "0",
			"b": "0",
			"c": "1",
		},
		want: true,
	})
}

func TestA1AndB1AndC1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b && c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
		},
		want: true,
	})
}

func TestA1AndB1AndC0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a && b && c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "0",
		},
		want: false,
	})
}

func TestA1EqB1AndB1EqC1(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b && b == c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "1",
		},
		want: true,
	})
}

func TestA1EqB1AndB1EqC0(t *testing.T) {
	runEvalTest(t, evalTest{
		source: "a == b && b == c",
		defines: map[string]string{
			"a": "1",
			"b": "1",
			"c": "0",
		},
		want: false,
	})
}
