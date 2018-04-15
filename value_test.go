package flyline

import "testing"

func TestValue_Scan_Basic(t *testing.T) {
	var err error = nil
	vInt := &Value{src: int64(1)}
	i := int64(0)
	err = vInt.Scan(&i)
	t.Logf("int64: %v, %v", i == vInt.src, err)
	vBool := &Value{src: true}
	b := false
	err = vBool.Scan(&b)
	t.Logf("bool: %v, %v", b == vBool.src, err)
}

func TestValue_Scan_Slice(t *testing.T) {
	var err error = nil
	v := &Value{src: []int64{1, 2, 3, 4}}
	s := make([]int64, 0, 4)
	err = v.Scan(&s)
	t.Logf("slice: src = %v, dest = %v, err = %v", v.src, s, err)
}

type P struct {
	V int64
}

func TestValue_Scan_Interface(t *testing.T) {
	var err error = nil
	v1 := &Value{src: &P{1}}
	s1 := P{}
	err = v1.Scan(&s1)
	t.Logf("slice: src = %v, dest = %v, err = %v", v1.src, s1, err)

	v2 := &Value{src: P{1}}
	s2 := new(P)
	err = v2.Scan(s2)
	t.Logf("slice: src = %v, dest = %v, err = %v", v2.src, s2, err)

}
