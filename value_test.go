package flyline

import "testing"

func TestValue_Scan_Basic(t *testing.T) {
	var err error = nil
	vInt := int64(1)
	i := int64(0)
	err = ValueScan(vInt, &i)
	t.Logf("int64: %v, %v", i == vInt, err)
	vBool := true
	b := false
	err = ValueScan(vBool, &b)
	t.Logf("bool: %v, %v", b == vBool, err)
}

func TestValue_Scan_Slice(t *testing.T) {
	var err error = nil
	v := []int64{1, 2, 3, 4}
	s := make([]int64, 0, 4)
	err = ValueScan(v, &s)
	t.Logf("slice: src = %v, dest = %v, err = %v", v, s, err)
}

type P struct {
	V int64
}

func TestValue_Scan_Interface(t *testing.T) {
	var err error = nil
	v1 := &P{1}
	s1 := P{}
	err = ValueScan(v1, &s1)
	t.Logf("slice: src = %v, dest = %v, err = %v", v1, s1, err)

	v2 := P{1}
	s2 := new(P)
	err = ValueScan(v2, s2)
	t.Logf("slice: src = %v, dest = %v, err = %v", v2, s2, err)

}
