package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"concurrent/models"
)

var cache = map[int]models.Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		wg.Add(2)
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex) {
			defer wg.Done()
			if b, ok := queryCache(id, m); ok {
				fmt.Println("from cache\n", b)
			}
		}(id, wg, m)
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex) {
			defer wg.Done()
			if b, ok := queryDatabase(id, m); ok {
				fmt.Println("from database\n", b)
			}
		}(id, wg, m)
		time.Sleep(150 * time.Millisecond)
	}
	wg.Wait()
}

func queryCache(id int, m *sync.RWMutex) (models.Book, bool) {
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	return b, ok
}

func queryDatabase(id int, m *sync.RWMutex) (models.Book, bool) {
	time.Sleep(100 * time.Millisecond)
	for _, b := range models.Books {
		if b.ID == id {
			m.Lock()
			cache[id] = b
			m.Unlock()
			return b, true
		}
	}
	return models.Book{}, false
}
