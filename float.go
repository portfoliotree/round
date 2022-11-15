package round

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Decimal rounds a floating point number to an expected decimal rounded floating point number.
func Decimal(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	shifted := num * shift
	rounded := int(shifted + math.Copysign(0.5, shifted))
	return float64(rounded) / shift
}

// Recursive navigates ptr using reflection to round floats appropriately
// It defaults to using the given precision; however, it consults a struct tag
// as it navigates the fields. It uses the struct tag value for all values in a field
// unless it encounters another struct tag.
//
//		type Data struct {
//	      // Field1 will be rounded to the value "precision" passed to Recursive
//		  Field1 float64
//
//	      // Field2 will be rounded to 2 decimal places
//		  Field2 float64 `precision:"2"`
//
//	      // Field3 will have the slice elements rounded to 3 decimal places
//		  Field3 []float64 `precision:"3"`
//
//	      // Percent will be multiplied by 100 and then rounded to 2 decimal places
//	      // for example, "0.123401" would become "12.34"
//		  Field3 []float64 `precision:"2"`
//		}
//
// The value of ptr must be a pointer.
func Recursive(ptr interface{}, precision int) error {
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
