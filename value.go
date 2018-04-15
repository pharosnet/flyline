package flyline

import (
	"reflect"
	"fmt"
)

// Value Scanner
type Value struct {
	v interface{}
}

/*
 * Scan value
 * @param i ptr object
 *
 */
func (v *Value) Scan(i interface{}) (err error) {
	hosted := reflect.ValueOf(v.v)
	promised := reflect.ValueOf(i)
	if hosted.Kind() == reflect.Ptr {
		if hosted.Kind() == reflect.Ptr {
			promised.Elem().Set(hosted.Elem())
			return
		}
		promised.Elem().Set(hosted)
		return
	} else if promised.Kind() == reflect.Slice || promised.Kind() == reflect.Map {
		promised.Set(hosted)
		return
	}
	return fmt.Errorf("not match, v = %v, i = %v", hosted.Kind(), promised.Kind())
}
