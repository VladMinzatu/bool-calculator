package evaluation

import (
	"reflect"
	"testing"
)

func TestParsing(t *testing.T) {
	exp, err := ParseExpression("mux(Y,not(X), sel)")
	t.Logf("%v - %v", reflect.TypeOf(exp), exp)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

}
