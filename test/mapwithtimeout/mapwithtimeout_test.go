package mapwithtimeout_test

import (
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MyTestStorages[T any] struct {
	mutex    sync.RWMutex
	storages map[string]T
}

type element struct {
	id, name string
}

func NewMyTestStorages[T any]() *MyTestStorages[T] {
	return &MyTestStorages[T]{
		storages: make(map[string]T),
	}
}

func (s *MyTestStorages[T]) Add(key string, v T) {
	fmt.Println("func 'Add', BEFORE SIZE:", s.Size())

	if _, ok := s.storages[key]; !ok {
		s.storages[key] = v
	}

	fmt.Println("func 'Add', AFTER SIZE:", s.Size())
}

func (s *MyTestStorages[T]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.storages)
}

func (s *MyTestStorages[T]) Get() map[string]T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.storages
}

var _ = Describe("Mapwithtimeout", func() {
	Context("Test", func() {
		It("any ", func() {
			listIdOne := []string{
				"0000-0000",
				"1111-1111",
				"2222-2222",
				"3333-3333",
				"4444-4444",
				"5555-5555",
				"6666-6666",
				"7777-7777",
				"8888-8888",
				"9999-9999",
			}

			mts := NewMyTestStorages[element]()

			for _, v := range listIdOne {
				time.Sleep(1 * time.Second)

				mts.Add(v, element{id: v})
			}

			fmt.Println("LIST:", mts.Get())
			Expect(len(mts.Get()), len(listIdOne))
		})
	})
})
