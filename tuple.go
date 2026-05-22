package x

type Tuple[T1, T2 any] struct {
	Val1 T1
	Val2 T2
}

type Tuple3[T1, T2, T3 any] struct {
	Val1 T1
	Val2 T2
	Val3 T3
}

type Tuple4[T1, T2, T3, T4 any] struct {
	Val1 T1
	Val2 T2
	Val3 T3
	Val4 T4
}

type Tuple5[T1, T2, T3, T4, T5 any] struct {
	Val1 T1
	Val2 T2
	Val3 T3
	Val4 T4
	Val5 T5
}

func NewTuple[T1, T2 any](v1 T1, v2 T2) Tuple[T1, T2] {
	return Tuple[T1, T2]{v1, v2}
}

func NewTuple3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{v1, v2, v3}
}

func NewTuple4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4) Tuple4[T1, T2, T3, T4] {
	return Tuple4[T1, T2, T3, T4]{v1, v2, v3, v4}
}

func NewTuple5[T1, T2, T3, T4, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) Tuple5[T1, T2, T3, T4, T5] {
	return Tuple5[T1, T2, T3, T4, T5]{v1, v2, v3, v4, v5}
}
