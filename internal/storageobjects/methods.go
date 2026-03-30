package storageobjects

import (
	"context"
	"fmt"
	"time"
)

func (s *storageObjects[T]) Start(ctx context.Context) {
	go func() {
		tick := time.NewTicker(s.timeTick)
		tickerDel := time.NewTicker(1 * time.Second)
		defer func() {
			tick.Stop()
			tickerDel.Stop()

			close(s.chanOut)
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				if s.Len() == 0 {
					continue
				}

				//ищем самый старый объект
				index := s.getOldestIndexObject()
				//если ничего не было найдено, вероятно пустое хранилище
				if index == -1 {
					continue
				}

				fmt.Println("___ storageObjects.Start, time:", time.Now())

				//проверяем, подошло ли время отправки объекта
				if time.Now().After(s.storage.objects[index].timeSending) {
					//забираем объект из хранилища
					obj := s.getObject(index)

					fmt.Println("sending object ->")

					//отправка объекта в канал после истечения таймаута
					s.chanOut <- StorageObjectData[T]{
						Id:          obj.id,
						Data:        obj.element,
						TimeCreated: obj.timeCreated,
					}
				}

			case <-tickerDel.C:
				//удаление всех объектов у которых истекло время жизни
				s.deleteObjectsAfterLifetimeExpired()
			}
		}
	}()
}

// AddObject добавляет объект или обновляет его (если есть объект с таким индексом)
func (s *storageObjects[T]) AddObject(id string, element T) {
	s.addObject(id, element)
}

// addObject добавляет объект или обновляет его (если есть объект с таким индексом)
func (s *storageObjects[T]) addObject(id string, element T) {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	for k, v := range s.storage.objects {
		if v.id == id {
			s.storage.objects[k] = object[T]{
				id:          v.id,
				timeExpiry:  v.timeExpiry,
				timeCreated: v.timeCreated,
				timeSending: v.timeSending,
				element:     element,
			}

			return
		}
	}

	s.storage.objects = append(s.storage.objects, object[T]{
		id:          id,
		element:     element,
		timeExpiry:  time.Now().Add(s.timeToLive),
		timeSending: time.Now().Add(s.timeWaitingToSend),
		timeCreated: time.Now(),
	})
}

// Len объем хранилища
func (s *storageObjects[T]) Len() int {
	return s.len()
}

// len объем хранилища
func (s *storageObjects[T]) len() int {
	s.storage.mtx.RLock()
	defer s.storage.mtx.RUnlock()

	return len(s.storage.objects)
}

// GetObjects канал через который можно получать самые старые объекты
func (s *storageObjects[T]) GetObjects() <-chan StorageObjectData[T] {
	return s.chanOut
}

// getObject возвращает объект по индексу и удаляет его из хранилища
func (s *storageObjects[T]) getObject(index int) object[T] {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	obj := s.storage.objects[index]

	s.storage.objects[index] = s.storage.objects[len(s.storage.objects)-1]
	s.storage.objects = s.storage.objects[:len(s.storage.objects)-1]

	return obj
}

// getOldestObject возвращает индекс самого старого объекта
func (s *storageObjects[T]) getOldestIndexObject() int {
	var timeCreated time.Time
	index := -1

	s.storage.mtx.RLock()
	defer s.storage.mtx.RUnlock()

	for k, v := range s.storage.objects {
		if index == -1 {
			index = k
			timeCreated = v.timeCreated

			continue
		}

		if v.timeCreated.Before(timeCreated) {
			index = k
			timeCreated = v.timeCreated
		}
	}

	return index
}

// deleteObjectAfterLifetimeExpired удаляет все объекты у которых истекло время жизни
func (s *storageObjects[T]) deleteObjectsAfterLifetimeExpired() {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	index := -1
	for k, v := range s.storage.objects {
		if v.timeExpiry.Before(time.Now()) {
			index = k

			break
		}
	}

	if index != -1 {
		fmt.Println("func 'deleteObjectsAfterLifetimeExpired', delete object with index:", index)
		s.storage.objects[index] = s.storage.objects[len(s.storage.objects)-1]
		s.storage.objects = s.storage.objects[:len(s.storage.objects)-1]
	}
}
