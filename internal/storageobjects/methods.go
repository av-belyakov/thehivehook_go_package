package storageobjects

import (
	"context"
	"fmt"
	"time"
)

func (s *StorageObjects) Start(ctx context.Context) {
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

				//проверяем, подошло ли время отправки объекта
				if time.Now().After(s.storage.objects[index].timeSending) {
					//забираем объект из хранилища
					obj := s.getObject(index)

					//отправка объекта в канал после истечения таймаута
					s.chanOut <- StorageObjectData{
						Id:          obj.id,
						Data:        obj.data,
						ObjectType:  obj.objectType,
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
func (s *StorageObjects) AddObject(id, objectType string, data []byte) {
	s.addObject(id, objectType, data)
}

// addObject добавляет объект или обновляет его (если есть объект с таким индексом)
func (s *StorageObjects) addObject(id, objectType string, data []byte) {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	for k, v := range s.storage.objects {
		if v.id == id {
			s.storage.objects[k] = object{
				id:          v.id,
				objectType:  objectType,
				timeExpiry:  v.timeExpiry,
				timeCreated: v.timeCreated,
				timeSending: v.timeSending,
				data:        data,
			}

			return
		}
	}

	s.storage.objects = append(s.storage.objects, object{
		id:          id,
		objectType:  objectType,
		data:        data,
		timeExpiry:  time.Now().Add(s.timeToLive),
		timeSending: time.Now().Add(s.timeDelayToSend),
		timeCreated: time.Now(),
	})
}

// Len объем хранилища
func (s *StorageObjects) Len() int {
	return s.len()
}

// len объем хранилища
func (s *StorageObjects) len() int {
	s.storage.mtx.RLock()
	defer s.storage.mtx.RUnlock()

	return len(s.storage.objects)
}

// GetObjects канал через который можно получать самые старые объекты
func (s *StorageObjects) GetObjects() <-chan StorageObjectData {
	return s.chanOut
}

// getObject возвращает объект по индексу и удаляет его из хранилища
func (s *StorageObjects) getObject(index int) object {
	s.storage.mtx.Lock()
	defer s.storage.mtx.Unlock()

	obj := s.storage.objects[index]

	s.storage.objects[index] = s.storage.objects[len(s.storage.objects)-1]
	s.storage.objects = s.storage.objects[:len(s.storage.objects)-1]

	return obj
}

// getOldestObject возвращает индекс самого старого объекта
func (s *StorageObjects) getOldestIndexObject() int {
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
func (s *StorageObjects) deleteObjectsAfterLifetimeExpired() {
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
