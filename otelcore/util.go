package otelcore

import (
	"reflect"

	"github.com/hyperledger-labs/yui-relayer/core"
)

func As(v any, target any) bool {
	var fieldName string
	switch v.(type) {
	case core.Chain:
		fieldName = "Chain"
	case core.Prover:
		fieldName = "Prover"
	default:
		return false
	}

	targetType := reflect.TypeOf(target)

	rv := reflect.ValueOf(v)
	for {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		field := rv.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}

		fieldValue := field.Interface()
		rv = reflect.ValueOf(fieldValue)
		if reflect.TypeOf(fieldValue).AssignableTo(targetType) {
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			reflect.ValueOf(target).Elem().Set(rv)
			return true
		}
	}
}
