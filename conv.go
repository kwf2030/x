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
	I64 | I32 | I16 | I8 | ~int | ~uint
}

type I64 interface {
	~int64 | ~uint64
}

type I32 interface {
	~int32 | ~uint32
}

type I16 interface {
	~int16 | ~uint16
}

type I8 interface {
	~int8 | ~uint8
}

func I64tob[T I64](i T) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func I32tob[T I32](i T) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

func I16tob[T I16](i T) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	return b
}

func F64tob(f float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	return b
}

func F32tob(f float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	return b
}

func Btoi64[T I64](b []byte) T {
	return T(binary.BigEndian.Uint64(b))
}

func Btoi32[T I32](b []byte) T {
	return T(binary.BigEndian.Uint32(b))
}

func Btoi16[T I16](b []byte) T {
	return T(binary.BigEndian.Uint16(b))
}

func Btof64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

func Btof32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

func Bool(v any) bool {
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

func Int(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case uint:
		return int(val)
	case string:
		if i, e := strconv.Atoi(val); e == nil {
			return i
		}
	case []byte:
		switch len(val) {
		case 8:
			return int(Btoi64[int64](val))
		case 4:
			return int(Btoi32[int32](val))
		case 2:
			return int(Btoi16[int16](val))
		case 1:
			return int(val[0])
		}
	case int64:
		return int(val)
	case uint64:
		return int(val)
	case float64:
		return int(val)
	case int32:
		return int(val)
	case uint32:
		return int(val)
	case float32:
		return int(val)
	case int16:
		return int(val)
	case uint16:
		return int(val)
	case int8:
		return int(val)
	case uint8:
		return int(val)
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func Int64(v any) int64 {
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
			return Btoi64[int64](val)
		case 4:
			return int64(Btoi32[int32](val))
		case 2:
			return int64(Btoi16[int16](val))
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

func Uint(v any) uint {
	switch val := v.(type) {
	case uint:
		return val
	case int:
		return uint(val)
	case string:
		if i, e := strconv.ParseUint(val, 10, 0); e == nil {
			return uint(i)
		}
	case []byte:
		switch len(val) {
		case 8:
			return uint(Btoi64[uint64](val))
		case 4:
			return uint(Btoi32[uint32](val))
		case 2:
			return uint(Btoi16[uint16](val))
		case 1:
			return uint(val[0])
		}
	case uint64:
		return uint(val)
	case int64:
		return uint(val)
	case float64:
		return uint(val)
	case uint32:
		return uint(val)
	case int32:
		return uint(val)
	case float32:
		return uint(val)
	case uint16:
		return uint(val)
	case int16:
		return uint(val)
	case uint8:
		return uint(val)
	case int8:
		return uint(val)
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func Uint64(v any) uint64 {
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
			return Btoi64[uint64](val)
		case 4:
			return uint64(Btoi32[uint32](val))
		case 2:
			return uint64(Btoi16[uint16](val))
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

func Float64(v any) float64 {
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
			return Btof64(val)
		case 4:
			return float64(Btof32(val))
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

func String(v any) string {
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

func Bytes(v any) []byte {
	switch val := v.(type) {
	case []byte:
		return val
	case string:
		return []byte(val)
	case int:
		switch strconv.IntSize {
		case 64:
			return I64tob(int64(val))
		case 32:
			return I32tob(int32(val))
		}
	case uint:
		switch strconv.IntSize {
		case 64:
			return I64tob(uint64(val))
		case 32:
			return I32tob(uint32(val))
		}
	case int64:
		return I64tob(val)
	case uint64:
		return I64tob(val)
	case float64:
		return F64tob(val)
	case int32:
		return I32tob(val)
	case uint32:
		return I32tob(val)
	case float32:
		return F32tob(val)
	case int16:
		return I16tob(val)
	case uint16:
		return I16tob(val)
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
