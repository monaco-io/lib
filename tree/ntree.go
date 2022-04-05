package tree

type NTree[T any] struct {
	value    *NTree[T]
	children []*NTree[T]
}
