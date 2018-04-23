package flyline

import (
	"testing"
	"fmt"
)

func TestArray_Simple(t *testing.T)  {
	a := newArray(4)
	for i := int64(0) ; i < int64(10) ; i ++ {
		a.set(i, i)
		fmt.Printf("[%v] %v, want %v \n", i, a.get(i), i)
	}
}
