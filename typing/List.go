package typing

import (
	"github.com/samber/lo"
)

type List[T any] []T

func (l List[T]) ForEach(f func(item T, index int)) {
	lo.ForEach(l, f)
}
