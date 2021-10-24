// Package mody allows to modify fields in an object.
package mody

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const delimiter = "."

// Update sets the value val of the named field inside
// the passed struct instance.
func Update(v interface{}, field string, val interface{}) (err error) {
	return update(reflect.ValueOf(v), field, val, true)
}

// UpdateJson unmarshals valStr into a value and tries
// to apply it to the named field in v.
func UpdateJson(v interface{}, field, valStr string) (err error) {
	var val interface{}
	if err = json.Unmarshal([]byte(valStr), &val); err != nil {
		return
	}
	return Update(v, field, val)
}

func update(
	v reflect.Value,
	field string, val interface{},
	isTopLvl bool,
) (err error) {
	split := strings.SplitN(field, delimiter, 2)

	isPtr := v.Kind() == reflect.Ptr

	if isTopLvl && !isPtr {
		err = ErrTypeMustBePointer
		return
	}

	if isPtr {
		v = v.Elem()
	}

	fv := v.FieldByName(split[0])
	if !fv.IsValid() {
		err = ErrFieldNotExistent
		return
	}

	if len(split) == 2 {
		return update(fv, split[1], val, false)
	}

	valT := reflect.ValueOf(val)

	if valT.Kind() != fv.Kind() {
		if fv.Kind() == reflect.String || !valT.CanConvert(fv.Type()) {
			err = ErrTypeMissmatch
			return
		}
		valT = valT.Convert(fv.Type())
	}

	fv.Set(valT)

	return
}

// Catch executes the passed function and recovers
// any occuring panic inside it which is then
// returned as error.
func Catch(fn func()) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()
	fn()
	return
}
