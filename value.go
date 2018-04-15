package flyline

import (
	"fmt"
	"reflect"
)

// Value Scanner
type Value struct {
	src interface{}
}

/*
 * Scan value
 * @param dest ptr object
 *
 */
func (v *Value) Scan(dest interface{}) (err error) {
	if dest == nil {
		err = fmt.Errorf("dest is nil, dest = %v", dest)
		return
	}
	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		err = fmt.Errorf("dest's type is not Ptr, dest = %v", dpv.Kind())
		return
	}
	sv := reflect.ValueOf(v.src)
	dv := reflect.Indirect(dpv)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		dv.Set(sv)
		return
	}
	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return
	}
	err = fmt.Errorf("scan failed, not match, src = %v, dest = %v", sv.Type(), dv.Type())
	return
}
