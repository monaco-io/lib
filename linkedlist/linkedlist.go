package linkedlist

import "fmt"

type ListNode[T comparable] struct {
	value T
	next  *ListNode[T]
}

func (ll *ListNode[T]) String() string {
	if ll == nil {
		return "nil"
	}
	return fmt.Sprintf("%v->%v", ll.value, ll.next.String())
}

func (ll *ListNode[T]) Len() int {
	count := 0
	for node := ll; node != nil; node = node.next {
		count++
	}
	return count
}

func (ll *ListNode[T]) InsertAt(pos int, value T) *ListNode[T] {
	if ll == nil || pos <= 0 {
		return &ListNode[T]{
			value: value,
			next:  ll,
		}
	}
	ll.next = ll.next.InsertAt(pos-1, value)
	return ll
}

func (ll *ListNode[T]) Append(value T) *ListNode[T] {
	return ll.InsertAt(ll.Len(), value)
}

func (ll *ListNode[T]) Contains(value T) bool {
	for node := ll; node != nil; node = node.next {
		if node.value == value {
			return true
		}
	}
	return false
}
