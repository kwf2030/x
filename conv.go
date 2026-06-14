package x

import (
	"encoding/binary"
	"math"
	"strconv"
)

func I64ToBytes[T ~int64 | ~uint64](i T) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func I32ToBytes[T ~int32 | ~uint32](i T) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

func I16ToBytes[T ~int16 | ~uint16](i T) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	return b
}

func F64ToBytes(f float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	return b
}

func F32ToBytes(f float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	return b
}

func BytesToI64[T ~int64 | ~uint64](b []byte) T {
	return T(binary.BigEndian.Uint64(b))
}

func BytesToI32[T ~int32 | ~uint32](b []byte) T {
	return T(binary.BigEndian.Uint32(b))
}

func BytesToI16[T ~int16 | ~uint16](b []byte) T {
	return T(binary.BigEndian.Uint16(b))
}

func BytesToF64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

func BytesToF32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

func ToBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	case int:
		return val != 0
	case uint:
		return val != 0
	case int64:
		return val != 0
	case uint64:
		return val != 0
	case int32:
		return val != 0
	case uint32:
		return val != 0
	case int16:
		return val != 0
	case uint16:
		return val != 0
	case int8:
		return val != 0
	case uint8:
		return val != 0
	case float64:
		return val != 0
	case float32:
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
	case string:
		if i, e := strconv.ParseInt(val, 10, 64); e == nil {
			return i
		}
	case int:
		return int64(val)
	case uint:
		return int64(val)
	case int64:
		return val
	case uint64:
		return int64(val)
	case int32:
		return int64(val)
	case uint32:
		return int64(val)
	case int16:
		return int64(val)
	case uint16:
		return int64(val)
	case int8:
		return int64(val)
	case uint8:
		return int64(val)
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	case bool:
		if val {
			return 1
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToI64[int64](val)
		case 4:
			return int64(BytesToI32[int32](val))
		case 2:
			return int64(BytesToI16[int16](val))
		case 1:
			return int64(val[0])
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
	case string:
		if i, e := strconv.ParseUint(val, 10, 64); e == nil {
			return i
		}
	case int:
		return uint64(val)
	case uint:
		return uint64(val)
	case int64:
		return uint64(val)
	case uint64:
		return val
	case int32:
		return uint64(val)
	case uint32:
		return uint64(val)
	case int16:
		return uint64(val)
	case uint16:
		return uint64(val)
	case int8:
		return uint64(val)
	case uint8:
		return uint64(val)
	case float64:
		return uint64(val)
	case float32:
		return uint64(val)
	case bool:
		if val {
			return 1
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToI64[uint64](val)
		case 4:
			return uint64(BytesToI32[uint32](val))
		case 2:
			return uint64(BytesToI16[uint16](val))
		case 1:
			return uint64(val[0])
		}
	}
	return 0
}

func ToFloat32(v any) float32 {
	return float32(ToFloat64(v))
}

func ToFloat64(v any) float64 {
	switch val := v.(type) {
	case string:
		if f, e := strconv.ParseFloat(val, 64); e == nil {
			return f
		}
	case int:
		return float64(val)
	case uint:
		return float64(val)
	case int64:
		return float64(val)
	case uint64:
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
	case float64:
		return val
	case float32:
		return float64(val)
	case bool:
		if val {
			return 1
		}
	case []byte:
		switch len(val) {
		case 8:
			return BytesToF64(val)
		case 4:
			return float64(BytesToF32(val))
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
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case float64:
		return strconv.FormatFloat(val, 'f', 2, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', 2, 32)
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
			return I64ToBytes(int64(val))
		case 32:
			return I32ToBytes(int32(val))
		}
	case uint:
		switch strconv.IntSize {
		case 64:
			return I64ToBytes(uint64(val))
		case 32:
			return I32ToBytes(uint32(val))
		}
	case int64:
		return I64ToBytes(val)
	case uint64:
		return I64ToBytes(val)
	case int32:
		return I32ToBytes(val)
	case uint32:
		return I32ToBytes(val)
	case int16:
		return I16ToBytes(val)
	case uint16:
		return I16ToBytes(val)
	case int8:
		return []byte{byte(val)}
	case uint8:
		return []byte{val}
	case float64:
		return F64ToBytes(val)
	case float32:
		return F32ToBytes(val)
	case bool:
		if val {
			return []byte("true")
		}
	}
	return nil
}
