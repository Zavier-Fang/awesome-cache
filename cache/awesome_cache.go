package cache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (getterFunc GetterFunc) Get(key string) ([]byte, error) {
	return getterFunc(key)
}

var (
	mtx    sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter function nil")
	}
	mtx.Lock()
	defer mtx.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheByte: cacheBytes},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mtx.RLock()
	defer mtx.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == ""{
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok{
		log.Println("[cache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error){
	return g.loadLocal(key)
}

func (g *Group) loadLocal(key string) (ByteView, error){
	byteView, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(byteView)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView){
	g.mainCache.add(key, value)
}
