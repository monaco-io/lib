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
