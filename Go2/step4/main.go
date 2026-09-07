package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeMap struct {
	m   map[string]interface{}
	mux sync.Mutex
}

func (s *SafeMap) Get(key string) interface{} {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.m[key]
}

func (s *SafeMap) Set(key string, value interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] = value
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		make(map[string]interface{}),
		sync.Mutex{},
	}
}

type Counter struct {
	value int
	mux   sync.RWMutex
}

func (c *Counter) Increment() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.value++
}

func (c *Counter) GetValue() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.value
}

type Count interface {
	Increment()
	GetValue() int
}

type ConcurrentQueue struct {
	queue []interface{}
	mux   sync.Mutex
}

type Queue interface {
	Enqueue(element interface{})
	Dequeue() interface{}
}

func (c *ConcurrentQueue) Enqueue(element interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.queue = append(c.queue, element)
}

func (c *ConcurrentQueue) Dequeue() interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()
	if len(c.queue) == 0 {
		return nil
	}
	element := c.queue[0]
	c.queue = c.queue[1:]
	return element
}

var (
	Buf   []int
	mutex sync.Mutex
)

func Write(num int) {
	mutex.Lock()
	defer mutex.Unlock()
	Buf = append(Buf, num)
}

func Consume() int {
	mutex.Lock()
	defer mutex.Unlock()

	res := Buf[0]
	Buf = Buf[1:]
	return res
}

func main() {

}

func Test() {
	for i := range 10 {
		go func() {
			fmt.Println(i)
		}()
	}
	time.Sleep(time.Second)
}
