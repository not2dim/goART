package goART

import "testing"

func TestTree(t *testing.T) {
	tree := new(Tree)
	keys := []string{"absent", "abandon", "app", "apple", "applied"}
	for val, key := range keys {
		tree.Insert([]byte(key), val)
		t.Logf("tree size %v\n", tree.Size())
	}
	for expected, key := range keys {
		actual, exist := tree.Search([]byte(key))
		if !exist {
			t.Fatalf("key %v not exists\n", key)
		}
		if actual != expected {
			t.Fatalf("key expected %v, actual %v\n", expected, actual)
		}
		t.Logf("key %v, val %v\n", actual, actual)
	}
	for expected, key := range keys {
		actual, deleted := tree.Delete([]byte(key))
		if !deleted {
			t.Fatalf("key %v not exists\n", key)
		}
		if actual != expected {
			t.Fatalf("key expected %v, actual %v\n", expected, actual)
		}
		t.Logf("key %v, val %v\n", actual, actual)
		t.Logf("tree size %v\n", tree.Size())
	}
}
