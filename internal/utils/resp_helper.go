package utils // Or wherever, e.g., internal/utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

// Map maps fields from src to dest using reflection. It supports structs, slices, pointers, and basic type conversions.
// It matches fields by name or json tag. Dest must be a pointer to a struct.
func RespMap(src any, dest any) error {
	srcV := reflect.ValueOf(src)
	if srcV.Kind() == reflect.Ptr {
		srcV = srcV.Elem()
	}
	if srcV.Kind() != reflect.Struct {
		return errors.New("src must be a struct or pointer to struct")
	}

	destV := reflect.ValueOf(dest)
	if destV.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	destV = destV.Elem()
	if destV.Kind() != reflect.Struct {
		return errors.New("dest must point to a struct")
	}

	for i := 0; i < srcV.NumField(); i++ {
		srcF := srcV.Field(i)
		srcT := srcV.Type().Field(i)

		// Skip unexported or ignored fields
		if !srcF.CanInterface() || srcT.Tag.Get("json") == "-" {
			continue
		}

		name := srcT.Name
		// Prefer json tag for matching
		if tag := srcT.Tag.Get("json"); tag != "" {
			name = strings.Split(tag, ",")[0]
		}

		// Find matching field in dest by name or json tag
		destF := destV.FieldByName(srcT.Name)
		if !destF.IsValid() {
			// Search by json tag
			for j := 0; j < destV.NumField(); j++ {
				destTF := destV.Type().Field(j)
				if destTF.Tag.Get("json") == name || destTF.Name == name {
					destF = destV.Field(j)
					break
				}
			}
		}
		if !destF.IsValid() || !destF.CanSet() {
			continue // No matching field or can't set
		}

		// Type matching and setting
		if srcF.Type() == destF.Type() {
			destF.Set(srcF)
		} else {
			switch srcF.Kind() {
			case reflect.Struct:
				if destF.Kind() == reflect.Struct {
					if err := RespMap(srcF.Addr().Interface(), destF.Addr().Interface()); err != nil {
						return err
					}
				}
			case reflect.Slice:
				if destF.Kind() == reflect.Slice {
					destF.Set(reflect.MakeSlice(destF.Type(), srcF.Len(), srcF.Len()))
					for j := 0; j < srcF.Len(); j++ {
						srcElem := srcF.Index(j)
						destElemP := reflect.New(destF.Type().Elem())
						if err := RespMap(srcElem.Interface(), destElemP.Interface()); err != nil {
							return err
						}
						destF.Index(j).Set(destElemP.Elem())
					}
				}
			case reflect.Ptr:
				if !srcF.IsNil() && destF.Kind() == reflect.Ptr {
					destElemP := reflect.New(destF.Type().Elem())
					if err := RespMap(srcF.Elem().Interface(), destElemP.Interface()); err != nil {
						return err
					}
					destF.Set(destElemP)
				}
			default:
				// Basic type conversions
				if srcF.Type().String() == "github.com/shopspring/decimal.Decimal" && destF.Kind() == reflect.Float64 {
					destF.SetFloat(srcF.Interface().(decimal.Decimal).InexactFloat64())
				} else if srcF.CanConvert(destF.Type()) {
					destF.Set(srcF.Convert(destF.Type()))
				}
				// Add more conversions as needed, e.g., time.Time, enums (if int/string)
			}
		}
	}
	return nil
}