package main

func ValueToPointer[T any](value T) *T {
	return &value
}

func GetPointerValue[T any](ptr *T) T {
	if ptr == nil {
		var x interface{}
		ZeroValue, _ := x.(T)
		return ZeroValue
	}
	return *ptr
}
