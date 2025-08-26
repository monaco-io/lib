package typing

import (
	"github.com/samber/lo"
)

type (
	List[T any] []T
	ListX       List[any]
)

func (l List[T]) ForEach(f func(item T, index int)) {
	lo.ForEach(l, f)
}
