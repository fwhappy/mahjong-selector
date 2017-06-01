package tests

import (
	"testing"
	"util"
)

func TestJsonMap(t *testing.T) {
	a := 1
	b := "b"
	sub := 3
	c := map[string]interface{}{"sub": sub}

	m := map[string]interface{}{"a": a, "b": b, "c": c}

	j := util.JsonMap(m)
	_a, _ := j.JsonGetInt("a")
	if _a != a {
		t.Error("JsonGetInt error")
	}

	_b, _ := j.JsonGetString("b")
	if _b != b {
		t.Error("JsonGetInt error")
	}

	_c, _ := j.JsonGetJsonMap("c")
	_sub, _ := _c.JsonGetInt("sub")
	if _sub != sub {
		t.Error("JsonGetJsonMap error")
	}
}
