package main

import (
	. "awesome-cache/cache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Alice":  "1",
	"Bob":    "2",
	"Claire": "3",
	"David":  "4",
}

func main() {
	NewGroup("test", 1<<10, GetterFunc(func(key string) ([]byte, error){
		log.Println("[db] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := "localhost:10001"
	peers := NewHttpPool(addr)
	log.Println("awesome-cache running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
