package x

import (
	"encoding/binary"
	"math"
	"strconv"
)

type Number interface {
	Integer | Float
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Int64 | Int32 | Int16 | Int8 | ~int | ~uint
}

type Int64 interface {
	~int64 | ~uint64
}

type Int32 interface {
	~int32 | ~uint32
}

type Int16 interface {
	~int16 | ~uint16
}

type Int8 interface {
	~int8 | ~uint8
}

func Int64ToBytes[T Int64](i T) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func Int32ToBytes[T Int32](i T) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

func Int16ToBytes[T Int16](i T) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	return b
}

func Float64ToBytes(f float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	return b
}

func Float32ToBytes(f float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	return b
}

func BytesToInt64[T Int64](b []byte) T {
	return T(binary.BigEndian.Uint64(b))
}

func BytesToInt32[T Int32](b []byte) T {
	return T(binary.BigEndian.Uint32(b))
}

func BytesToInt16[T Int16](b []byte) T {
	return T(binary.BigEndian.Uint16(b))
}

func BytesToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

func ToBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return val != 0
	}
	return v != nil
}

func ToInt(v any) int {
	return int(ToInt64(v))
}

func ToInt32(v any) int32 {
	return int32(ToInt64(v))
}

func ToInt64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case uint64:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		if i, e := strconv.ParseInt(val, 10, 64); e == nil {
			return i
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToInt64[int64](val)
		case 4:
			return int64(BytesToInt32[int32](val))
		case 2:
			return int64(BytesToInt16[int16](val))
		case 1:
			return int64(val[0])
		}
	case int:
		return int64(val)
	case uint:
		return int64(val)
	case int32:
		return int64(val)
	case uint32:
		return int64(val)
	case float32:
		return int64(val)
	case int16:
		return int64(val)
	case uint16:
		return int64(val)
	case int8:
		return int64(val)
	case uint8:
		return int64(val)
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func ToUint(v any) uint {
	return uint(ToUint64(v))
}

func ToUint32(v any) uint32 {
	return uint32(ToUint64(v))
}

func ToUint64(v any) uint64 {
	switch val := v.(type) {
	case uint64:
		return val
	case int64:
		return uint64(val)
	case float64:
		return uint64(val)
	case string:
		if i, e := strconv.ParseUint(val, 10, 64); e == nil {
			return i
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToInt64[uint64](val)
		case 4:
			return uint64(BytesToInt32[uint32](val))
		case 2:
			return uint64(BytesToInt16[uint16](val))
		case 1:
			return uint64(val[0])
		}
	case uint:
		return uint64(val)
	case int:
		return uint64(val)
	case uint32:
		return uint64(val)
	case int32:
		return uint64(val)
	case float32:
		return uint64(val)
	case uint16:
		return uint64(val)
	case int16:
		return uint64(val)
	case uint8:
		return uint64(val)
	case int8:
		return uint64(val)
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func ToFloat32(v any) float32 {
	return float32(ToFloat64(v))
}

func ToFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		if f, e := strconv.ParseFloat(val, 64); e == nil {
			return f
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToFloat64(val)
		case 4:
			return float64(BytesToFloat32(val))
		}
	case int:
		return float64(val)
	case uint:
		return float64(val)
	case float32:
		return float64(val)
	case int32:
		return float64(val)
	case uint32:
		return float64(val)
	case int16:
		return float64(val)
	case uint16:
		return float64(val)
	case int8:
		return float64(val)
	case uint8:
		return float64(val)
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func ToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', 2, 64)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', 2, 32)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case bool:
		if val {
			return "true"
		}
	}
	return ""
}

func ToBytes(v any) []byte {
	switch val := v.(type) {
	case []byte:
		return val
	case string:
		return []byte(val)
	case int:
		switch strconv.IntSize {
		case 64:
			return Int64ToBytes(int64(val))
		case 32:
			return Int32ToBytes(int32(val))
		}
	case uint:
		switch strconv.IntSize {
		case 64:
			return Int64ToBytes(uint64(val))
		case 32:
			return Int32ToBytes(uint32(val))
		}
	case int64:
		return Int64ToBytes(val)
	case uint64:
		return Int64ToBytes(val)
	case float64:
		return Float64ToBytes(val)
	case int32:
		return Int32ToBytes(val)
	case uint32:
		return Int32ToBytes(val)
	case float32:
		return Float32ToBytes(val)
	case int16:
		return Int16ToBytes(val)
	case uint16:
		return Int16ToBytes(val)
	case int8:
		return []byte{byte(val)}
	case uint8:
		return []byte{val}
	case bool:
		if val {
			return []byte("true")
		}
	}
	return nil
}
