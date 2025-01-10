package evaluation

import "testing"

var unaryInput = []bool{false, true}

var binaryInput = []struct {
	a bool
	b bool
}{
	{false, false},
	{false, true},
	{true, false},
	{true, true},
}

var ternaryInput = []struct {
	a bool
	b bool
	c bool
}{
	{false, false, false},
	{false, false, true},
	{false, true, false},
	{false, true, true},
	{true, false, false},
	{true, false, true},
	{true, true, false},
	{true, true, true},
}

func TestNand(t *testing.T) {
	expected := []bool{true, true, true, false}

	for idx, tc := range binaryInput {
		verifyEquality(t, Nand(tc.a, tc.b), expected[idx])
	}
}

func TestNot(t *testing.T) {
	expected := []bool{true, false}

	for idx, tc := range unaryInput {
		verifyEquality(t, Not(tc), expected[idx])
	}
}

func TestAnd(t *testing.T) {
	expected := []bool{false, false, false, true}

	for idx, tc := range binaryInput {
		verifyEquality(t, And(tc.a, tc.b), expected[idx])
	}
}

func TestOr(t *testing.T) {
	expected := []bool{false, true, true, true}

	for idx, tc := range binaryInput {
		verifyEquality(t, Or(tc.a, tc.b), expected[idx])
	}
}

func TestXor(t *testing.T) {
	expected := []bool{false, true, true, false}

	for idx, tc := range binaryInput {
		verifyEquality(t, Xor(tc.a, tc.b), expected[idx])
	}
}

func TestMux(t *testing.T) {
	expected := []bool{false, false, true, false, false, true, true, true}

	for idx, tc := range ternaryInput {
		verifyEquality(t, Mux(tc.a, tc.b, tc.c), expected[idx])
	}
}

func TestDmux(t *testing.T) {
	type result struct {
		a bool
		b bool
	}
	expected := []result{
		{false, false},
		{false, false},
		{true, false},
		{false, true},
	}

	for idx, tc := range binaryInput {
		a, b := Dmux(tc.a, tc.b)
		verifyEquality(t, result{a, b}, expected[idx])
	}
}

func verifyEquality[T comparable](t *testing.T, actual, expected T) {
	if actual != expected {
		t.Errorf("got %+v, wanted %+v", actual, expected)
	}
}
