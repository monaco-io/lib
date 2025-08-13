package typing

import (
	"fmt"
)

type BTreeNode[T comparable] struct {
	value T
	left  *BTreeNode[T]
	right *BTreeNode[T]
}

func (bt *BTreeNode[T]) String() string {
	var (
		answer string
		layers = []*BTreeNode[T]{bt}
	)

	for layers != nil {
		var temp []*BTreeNode[T]
		for i := range layers {
			if layers[i] == nil {
				answer = fmt.Sprintf("%s,null", answer)
				continue
			}
			answer = fmt.Sprintf("%s,%v", answer, layers[i].value)
			temp = append(temp, layers[i].left, layers[i].right)
		}
		layers = temp
	}

	return answer[1:]
}

func (bt *BTreeNode[T]) Deserialize(list []T) {
	if len(list) == 0 {
		return
	}
	*bt = BTreeNode[T]{value: list[0]}
	tnList := []*BTreeNode[T]{bt}
	list = list[1:]

	for len(list) > 0 {
		if len(list) == 1 {
			tnList[0].left = &BTreeNode[T]{value: list[0]}
			return
		}
		tnList[0].left = &BTreeNode[T]{value: list[0]}
		tnList[0].right = &BTreeNode[T]{value: list[1]}
		if tnList[0].left != nil {
			tnList = append(tnList, tnList[0].left)
		}
		if tnList[0].right != nil {
			tnList = append(tnList, tnList[0].right)
		}
		list = list[2:]
		tnList = tnList[1:]
	}
}
