package round

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func Decimal(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	shifted := num * shift
	rounded := int(shifted + math.Copysign(0.5, shifted))
	return float64(rounded) / shift
}

func FloatsRecursively(ptr interface{}, precision int) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		panic("v must be a pointer")
	}
	return floatsRecursively(val, precision, false)
}

func floatsRecursively(v reflect.Value, parentPrecision int, isPercent bool) error {
	switch v.Kind() {
	case reflect.Ptr:
		return floatsRecursively(v.Elem(), parentPrecision, isPercent)
	case reflect.Struct:
		valueType := v.Type()

		for i := 0; i < v.NumField(); i++ {
			precision := parentPrecision
			ip := isPercent
			field := v.Field(i)
			if !field.CanSet() {
				continue
			}

			fieldType := valueType.Field(i)

			precisionTag := fieldType.Tag.Get("precision")
			if strings.HasSuffix(precisionTag, ",percent") {
				ip = true
				precisionTag = strings.TrimSuffix(precisionTag, ",percent")
			}
			if precisionTag != "" {
				pt, err := strconv.Atoi(precisionTag)
				if err != nil {
					return fmt.Errorf(".%s: failed to parse precision tag: %w", fieldType.Name, err)
				}
				precision = pt
			}

			if err := floatsRecursively(field, precision, ip); err != nil {
				return fmt.Errorf(".%s%s", fieldType.Name, err)
			}
		}
	case reflect.Map:
		r := v.MapRange()
		tmp := reflect.New(v.Type().Elem())
		for r.Next() {
			tmp.Elem().Set(r.Value())
			if err := floatsRecursively(tmp, parentPrecision, isPercent); err != nil {
				return fmt.Errorf("[%v]%s", r.Key().Interface(), err)
			}
			v.SetMapIndex(r.Key(), tmp.Elem())
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := floatsRecursively(v.Index(i), parentPrecision, isPercent); err != nil {
				return fmt.Errorf("[%d]%s", i, err)
			}
		}
	case reflect.Float64:
		float := v.Float()

		if math.IsNaN(float) {
			return errors.New(" is not a number (but should be)")
		}

		if math.IsInf(float, 0) {
			return nil
		}

		if isPercent {
			float *= 100
		}
		float = Decimal(float, parentPrecision)

		v.SetFloat(float)
	}

	return nil
}
