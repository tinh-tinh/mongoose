package mongoose

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// DangerousOperators contains MongoDB query operators that could be exploited
// for NoSQL injection when user input is passed directly to filters.
var DangerousOperators = []string{
	// Comparison operators
	"$gt", "$gte", "$lt", "$lte", "$ne", "$nin", "$in",
	// Logical operators
	"$not", "$nor", "$or", "$and",
	// Element operators
	"$exists", "$type",
	// Evaluation operators (most dangerous)
	"$where", "$expr", "$jsonSchema", "$text", "$mod", "$regex",
	// Server-side JavaScript execution (CVE-2025-10061)
	"$function", "$accumulator",
	// Array operators
	"$all", "$elemMatch", "$size",
}

// ErrDangerousOperator is returned when a dangerous MongoDB operator is detected in a filter.
type ErrDangerousOperator struct {
	Operator string
}

func (e *ErrDangerousOperator) Error() string {
	return fmt.Sprintf("dangerous MongoDB operator detected in filter: %s", e.Operator)
}

// IsDangerousOperator checks if a key starts with '$' and is in the dangerous operators list.
func IsDangerousOperator(key string) bool {
	if !strings.HasPrefix(key, "$") {
		return false
	}
	return slices.Contains(DangerousOperators, key)
}

// SanitizeFilter recursively checks a filter for dangerous MongoDB operators.
// Returns an ErrDangerousOperator if any dangerous operator is found.
// This function should be called before passing user-controlled input to query functions.
func SanitizeFilter(filter any) error {
	if filter == nil {
		return nil
	}

	return sanitizeValue(reflect.ValueOf(filter))
}

func sanitizeValue(v reflect.Value) error {
	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		return sanitizeMap(v)
	case reflect.Slice, reflect.Array:
		return sanitizeSlice(v)
	case reflect.Struct:
		return sanitizeStruct(v)
	default:
		// Primitive types are safe
		return nil
	}
}

func sanitizeStruct(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue // Skip unexported fields
		}
		if err := sanitizeValue(field); err != nil {
			return err
		}
	}
	return nil
}

func sanitizeMap(v reflect.Value) error {
	for _, key := range v.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())

		// Check if the key is a dangerous operator
		if IsDangerousOperator(keyStr) {
			return &ErrDangerousOperator{Operator: keyStr}
		}

		// Recursively check the value
		if err := sanitizeValue(v.MapIndex(key)); err != nil {
			return err
		}
	}
	return nil
}

func sanitizeSlice(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		if err := sanitizeValue(v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

// IsDangerousOperatorError checks if an error is an ErrDangerousOperator.
func IsDangerousOperatorError(err error) bool {
	var opErr *ErrDangerousOperator
	return errors.As(err, &opErr)
}
