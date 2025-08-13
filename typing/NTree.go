package typing

type NTree[T any] struct {
	value    *NTree[T]
	children []*NTree[T]
}
