package flyline

import (
	"math"
	"unsafe"
)

type entry struct {
	value interface{}
}

func newArray(capacity int64) (a *array) {
	if capacity > 0 && (capacity&(capacity-1)) != 0 {
		panic("The array capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
		return
	}
	var items []*entry = nil
	align := int64(unsafe.Alignof(items))
	mask := int64(capacity - 1)
	shift := int64(math.Log2(float64(capacity)))
	size := int64(capacity) * align
	items = make([]*entry, size)
	itemBasePtr := uintptr(unsafe.Pointer(&items[0]))
	itemMSize := unsafe.Sizeof(items[0])
	for i := int64(0); i < size; i++ {
		items[i&mask*align] = &entry{}
	}
	return &array{
		capacity:    capacity,
		size:        size,
		shift:       shift,
		align:       align,
		mask:        mask,
		items:       items,
		itemBasePtr: itemBasePtr,
		itemMSize:   itemMSize,
	}
}

// shared array
type array struct {
	capacity    int64
	size        int64
	shift       int64
	align       int64
	mask        int64
	items       []*entry
	itemBasePtr uintptr
	itemMSize   uintptr
}

func (a *array) elementAt(seq int64) (e *entry) {
	mask := a.mask
	align := a.align
	basePtr := a.itemBasePtr
	mSize := a.itemMSize
	entryPtr := basePtr + uintptr(seq&mask*align)*mSize
	e = *((*(*entry))(unsafe.Pointer(entryPtr)))
	return
}

func (a *array) set(seq int64, v interface{}) {
	a.elementAt(seq).value = v
}

func (a *array) get(seq int64) (v interface{}) {
	v = a.elementAt(seq).value
	return
}
