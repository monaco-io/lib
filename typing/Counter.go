package typing

type Counter[T comparable] struct {
	Value T      `json:"value"`
	Label string `json:"label"`
}
