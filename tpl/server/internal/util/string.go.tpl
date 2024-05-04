package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"unsafe"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	// nolint:govet
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

// ToString casts a interface to a string.
func ToString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(i.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(i.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(i.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(i.(int32)), 10)
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case uint:
		return strconv.FormatUint(uint64(i.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(i.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(i.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(i.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(i.(uint64), 10)
	case []byte:
		return string(s)
	case nil:
		return ""
	case error:
		return s.Error()
	case fmt.Stringer:
		return s.String()
	default:
		return fmt.Sprint(i)
	}
}

// Truncate string to specific length
func Truncate(s string, n int) string {
	if n > len(s) {
		return s
	}
	return s[:n]
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

// Random generate a random string with fixed length. [a-z0-9]
func Random(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
