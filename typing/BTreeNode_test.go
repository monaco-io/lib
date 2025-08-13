package typing

import (
	"testing"
)

func TestBTreeString(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	peopleList := &BTreeNode[Person]{value: Person{"root", 1}}
	peopleList.left = &BTreeNode[Person]{value: Person{"Fred", 23}}
	peopleList.right = &BTreeNode[Person]{value: Person{"Joan", 30}}
	t.Log(peopleList)
}

func TestBTreeDeserialize(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	lsit := []Person{
		{"root", 1}, {"Fred", 23}, {"Joan", 30}, {"Marry", 25}, {"Mark", 21}, {"Jack", 25},
	}
	tn := new(BTreeNode[Person])
	tn.Deserialize(lsit)
	t.Log(tn)
}
