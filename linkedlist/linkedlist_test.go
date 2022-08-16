package linkedlist

import (
	"testing"
)

func TestListNodeString(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var peopleList *ListNode[Person]
	peopleList = peopleList.Append(Person{"Fred", 23})
	peopleList = peopleList.Append(Person{"Joan", 30})
	t.Log(peopleList)
}

func TestListNodeValue(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var peopleList *ListNode[Person]
	peopleList = peopleList.Append(Person{"Fred", 23})
	peopleList = peopleList.Append(Person{"Joan", 30})
	for v := peopleList; v != nil; v = v.next {
		t.Log(v.value)
	}
}

func TestListNodeInsertAt(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var peopleList *ListNode[Person]
	peopleList = peopleList.Append(Person{"Fred", 23})
	peopleList = peopleList.Append(Person{"Joan", 30})
	peopleList.InsertAt(1, Person{"John", 25})
	t.Log(peopleList)
}

func TestListNodeContains(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var peopleList *ListNode[Person]
	peopleList = peopleList.Append(Person{"Fred", 23})
	peopleList = peopleList.Append(Person{"Joan", 30})
	peopleList = peopleList.Append(Person{"John", 25})
	t.Log(peopleList.Contains(Person{"John", 25}))
	t.Log(peopleList.Contains(Person{"John", 21}))
}

// type PP map[string]string

// func (p PP) Equal(p2 PP) bool {
// 	return false
// }
// func (i PP) LessThan(j PP) bool {
// 	return false
// }
// func (a PP) Compare(b PP) int {
// 	return 0
// }
// func (i PP) EqualTo(j PP) bool {
// 	return false
// }
// func TestListNodeContains2(t *testing.T) {

// 	var peopleList *ListNode[PP]
// 	// peopleList = peopleList.Append(PP{"Fred", 23})
// 	// peopleList = peopleList.Append(PP{"Joan", 30})
// 	// peopleList = peopleList.Append(PP{"John", 25})
// 	// t.Log(peopleList.Contains(PP{"John", 25}))
// 	// t.Log(peopleList.Contains(PP{"John", 21}))
// 	t.Log(peopleList)
// }
