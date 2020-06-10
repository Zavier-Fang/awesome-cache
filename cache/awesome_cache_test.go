package cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	getter := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := getter.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("getter func failed")
	}
}

var db = map[string]string{
	"A": "a",
	"B": "b",
	"C": "c",
}

func TestGroup_Get(t *testing.T) {
	loadCount := make(map[string]int)
	group := NewGroup("test", 1024, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[db] search key", key)
		if v, ok := db[key]; ok {
			loadCount[key] ++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Fatalf("fail to get key [%s] from db", k)
		}
		if _, err := group.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("key [%s] miss cache", k)
		}
	}

	if view, err := group.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
