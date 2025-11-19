package randx

import "testing"

func TestGetRand1(t *testing.T) {
	num1 := GenIntWithRange(-484, -387)
	t.Logf("%d", num1)
}
