package x

import (
	"slices"
	"strings"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// 遍历
func ForEach[S ~[]E, E any](s S, fn func(E)) {
	for i := range s {
		fn(s[i])
	}
}

// 计数
func Count[S ~[]E, E comparable](s S, value E) int {
	ret := 0
	for i := range s {
		if s[i] == value {
			ret++
		}
	}
	return ret
}

// 计数
func CountFunc[S ~[]E, E any](s S, fn func(E) bool) int {
	ret := 0
	for i := range s {
		if fn(s[i]) {
			ret++
		}
	}
	return ret
}

// 累加
func Sum[S ~[]E, E Number](s S, value E) E {
	var ret E
	for i := range s {
		ret += s[i]
	}
	return ret
}

// 累加
func SumFunc[S ~[]E, E any, T Number](s S, fn func(E) T) T {
	var ret T
	for i := range s {
		ret += fn(s[i])
	}
	return ret
}

// 去重
func Distinct[S ~[]E, E comparable](s S) S {
	if s == nil {
		return nil
	}
	ret := make(S, 0, len(s))
	m := make(map[E]struct{})
	for i := range s {
		if _, ok := m[s[i]]; !ok {
			ret = append(ret, s[i])
			m[s[i]] = struct{}{}
		}
	}
	return ret
}

// 去重
func DistinctFunc[S ~[]E, E any, C comparable](s S, fn func(E) C) S {
	if s == nil {
		return nil
	}
	ret := make(S, 0, len(s))
	m := make(map[C]struct{})
	for i := range s {
		k := fn(s[i])
		if _, ok := m[k]; !ok {
			ret = append(ret, s[i])
			m[k] = struct{}{}
		}
	}
	return ret
}

// 映射
func Map[S1 ~[]E1, E1 any, S2 ~[]E2, E2 any](s S1, fn func(E1) E2) S2 {
	if s == nil {
		return nil
	}
	ret := make(S2, len(s))
	for i := range s {
		ret[i] = fn(s[i])
	}
	return ret
}

// 过滤
func Filter[S ~[]E, E any](s S, fn func(E) bool) S {
	if s == nil {
		return nil
	}
	ret := make(S, 0, len(s))
	for i := range s {
		if fn(s[i]) {
			ret = append(ret, s[i])
		}
	}
	return ret
}

// 过滤+映射
func FilterMap[S1 ~[]E1, E1 any, S2 ~[]E2, E2 any](s S1, fn func(E1) (E2, bool)) S2 {
	if s == nil {
		return nil
	}
	ret := make(S2, 0, len(s))
	for i := range s {
		if v, ok := fn(s[i]); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// 展开
func Flat[S ~[]E, E any](s S, fn func(E) S) S {
	return FlatMap(s, fn)
}

// 展开+映射
func FlatMap[S1 ~[]E1, E1 any, S2 ~[]E2, E2 any](s S1, fn func(E1) S2) S2 {
	if s == nil {
		return nil
	}
	ret := make(S2, 0)
	for i := range s {
		ret = append(ret, fn(s[i])...)
	}
	return ret
}

// 归约/聚合+映射
func Reduce[S ~[]E, E any, T any](s S, fn func(T, E) T, initial T) T {
	ret := initial
	for i := range s {
		ret = fn(ret, s[i])
	}
	return ret
}

// 分组
func Group[S ~[]E, E any, C comparable](s S, fn func(E) C) map[C]S {
	if s == nil {
		return nil
	}
	ret := make(map[C]S)
	for i := range s {
		k := fn(s[i])
		if _, ok := ret[k]; !ok {
			ret[k] = make(S, 0, 2)
		}
		ret[k] = append(ret[k], s[i])
	}
	return ret
}

// 转map
func Key[S ~[]E, E any, C comparable](s S, fn func(E) C) map[C]E {
	if s == nil {
		return nil
	}
	ret := make(map[C]E, len(s))
	for i := range s {
		ret[fn(s[i])] = s[i]
	}
	return ret
}

func Zip[S1 ~[]E1, E1 any, S2 ~[]E2, E2 any](s1 S1, s2 S2) []Tuple[E1, E2] {
	l := min(len(s1), len(s2))
	ret := make([]Tuple[E1, E2], l)
	for i := range l {
		ret[i].Val1, ret[i].Val2 = s1[i], s2[i]
	}
	return ret
}

// 并集
func Union[S ~[]E, E comparable](s1, s2 S) S {
	return Distinct(slices.Concat(s1, s2))
}

// 并集
func UnionFunc[S ~[]E, E any, C comparable](s1, s2 S, fn func(E) C) S {
	return DistinctFunc(slices.Concat(s1, s2), fn)
}

// 交集
func Intersection[S ~[]E, E comparable](s1, s2 S) S {
	if len(s1) == 0 || len(s2) == 0 {
		return make(S, 0)
	}
	set := make(map[E]struct{}, len(s2))
	for i := range s2 {
		set[s2[i]] = struct{}{}
	}
	ret := make(S, 0, min(len(s1), len(s2)))
	seen := make(map[E]struct{})
	for i := range s1 {
		if _, ok := set[s1[i]]; ok {
			if _, dup := seen[s1[i]]; !dup {
				ret = append(ret, s1[i])
				seen[s1[i]] = struct{}{}
			}
		}
	}
	return ret
}

// 交集
func IntersectionFunc[S ~[]E, E any, C comparable](s1, s2 S, fn func(E) C) S {
	if len(s1) == 0 || len(s2) == 0 {
		return make(S, 0)
	}
	set := make(map[C]struct{}, len(s2))
	for i := range s2 {
		set[fn(s2[i])] = struct{}{}
	}
	ret := make(S, 0, min(len(s1), len(s2)))
	seen := make(map[C]struct{})
	for i := range s1 {
		v := fn(s1[i])
		if _, ok := set[v]; ok {
			if _, dup := seen[v]; !dup {
				ret = append(ret, s1[i])
				seen[v] = struct{}{}
			}
		}
	}
	return ret
}

// 差集（s1-s2，属于s1但不属于s2）
func Difference[S ~[]E, E comparable](s1, s2 S) S {
	if len(s1) == 0 {
		return make(S, 0)
	}
	if len(s2) == 0 {
		return Distinct(s1)
	}
	set := make(map[E]struct{}, len(s2))
	for i := range s2 {
		set[s2[i]] = struct{}{}
	}
	ret := make(S, 0, len(s1))
	seen := make(map[E]struct{})
	for i := range s1 {
		if _, ok := set[s1[i]]; !ok {
			if _, dup := seen[s1[i]]; !dup {
				ret = append(ret, s1[i])
				seen[s1[i]] = struct{}{}
			}
		}
	}
	return ret
}

// 差集（s1-s2，属于s1但不属于s2）
func DifferenceFunc[S ~[]E, E any, C comparable](s1, s2 S, fn func(E) C) S {
	if len(s1) == 0 {
		return make(S, 0)
	}
	if len(s2) == 0 {
		return DistinctFunc(s1, fn)
	}
	set := make(map[C]struct{}, len(s2))
	for i := range s2 {
		set[fn(s2[i])] = struct{}{}
	}
	ret := make(S, 0, len(s1))
	seen := make(map[C]struct{})
	for i := range s1 {
		v := fn(s1[i])
		if _, ok := set[v]; !ok {
			if _, dup := seen[v]; !dup {
				ret = append(ret, s1[i])
				seen[v] = struct{}{}
			}
		}
	}
	return ret
}

// 对称差（并集-交集，仅属于s1或仅属于s2）
func SymDifference[S ~[]E, E comparable](s1, s2 S) S {
	diff1 := Difference(s1, s2)
	diff2 := Difference(s2, s1)
	return append(diff1, diff2...)
}

// 对称差（并集-交集，仅属于s1或仅属于s2）
func SymDifferenceFunc[S ~[]E, E any, C comparable](s1, s2 S, fn func(E) C) S {
	diff1 := DifferenceFunc(s1, s2, fn)
	diff2 := DifferenceFunc(s2, s1, fn)
	return append(diff1, diff2...)
}

// 转字符串
func JoinFunc[S ~[]E, E any](s S, sep string, fn func(E) string) string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	if sep == "" {
		for i := range s {
			b.WriteString(fn(s[i]))
		}
	} else {
		for i := range s {
			if i > 0 {
				b.WriteString(sep)
			}
			b.WriteString(fn(s[i]))
		}
	}
	return b.String()
}
