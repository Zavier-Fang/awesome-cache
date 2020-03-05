package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := NewCache(0, nil)
	lru.Add("key1", String("123"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "123" {
		t.Fatalf("cache hit key1=123 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemove(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1, value2, value3 := "value1", "value2", "value3"
	lru := NewCache(int64(len(key1+key2+value1+value2)), nil)
	lru.Add(key1, String(value1))
	lru.Add(key2, String(value2))
	lru.Add(key3, String(value3))

	if _, ok := lru.Get(key1); ok || lru.Len() != 2 {
		t.Fatalf("remove oldest key1 failed")
	}
}

func TestEvicted(t *testing.T) {
	evicted := make([]string, 0)
	onEvicted := func(key string, value Value) {
		evicted = append(evicted, key)
	}
	lru := NewCache(10, onEvicted)
	lru.Add("key1", String("value1"))
	lru.Add("key2", String("value2"))
	lru.Add("key3", String("value3"))

	expected := []string{"key1", "key2"}

	if !reflect.DeepEqual(expected, evicted) || lru.Len() != 1 {
		t.Fatalf("evicted failed")
	}
}

func TestAdd(t *testing.T) {
	lru := NewCache(0, nil)
	lru.Add("key1", String("value1"))
	lru.Add("key2", String("value2"))
	lru.Add("key1", String("v3"))

	if lru.nBytes != int64(len("key1"+"v3"+"key2"+"value2")) {
		t.Fatalf("add fail")
	}
}
