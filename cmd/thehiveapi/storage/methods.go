package storage

import (
	"fmt"
	"time"
)

// GetObjects все существующие объекты
func (s *StorageFoundObjects) GetObjects() map[string][]byte {
	s.storages.mutex.RLock()
	defer s.storages.mutex.RUnlock()

	objects := make(map[string][]byte, len(s.storages.foundObjects))
	for k, v := range s.storages.foundObjects {
		objects[k] = v.object
	}

	return objects
}

// GetObject получить объект по ключу
func (s *StorageFoundObjects) GetObject(key string) ([]byte, bool) {
	s.storages.mutex.RLock()
	defer s.storages.mutex.RUnlock()

	data, ok := s.storages.foundObjects[key]

	return data.object, ok
}

// SetObject добавить объект по ключу
func (s *StorageFoundObjects) SetObject(key string, object []byte) error {
	s.storages.mutex.Lock()
	defer s.storages.mutex.Unlock()

	if len(s.storages.foundObjects) >= s.maxSize {
		//удаляем самый старый объект
		if err := s.deleteOldestObjectFromCache(); err != nil {
			return err
		}
	}

	s.storages.foundObjects[key] = foundObject{
		object:     object,
		timeExpiry: time.Now().Add(s.maxTtl),
	}

	return nil
}

// DeleteForTimeExpiryObjectFromCache удаляет все объекты у которых истекло время жизни, без учета других параметров
func (s *StorageFoundObjects) DeleteForTimeExpiryObjectFromCache() {
	s.storages.mutex.Lock()
	defer s.storages.mutex.Unlock()

	for key, object := range s.storages.foundObjects {
		if object.timeExpiry.Before(time.Now()) {
			delete(s.storages.foundObjects, key)
		}
	}
}

// deleteOldestObjectFromCache поиск и удаление самого старого объекта
func (s *StorageFoundObjects) deleteOldestObjectFromCache() error {
	//получаем самый старый объект в кэше
	index := s.getOldestObjectFromCache()
	if _, ok := s.storages.foundObjects[index]; ok {
		delete(s.storages.foundObjects, index)
	} else {
		return fmt.Errorf("the object with id '%s' cannot be deleted, it may be in progress", index)
	}

	return nil
}

// getOldestObjectFromCache возвращает индекс самого старого объекта
func (s *StorageFoundObjects) getOldestObjectFromCache() string {
	var (
		index      string
		timeExpiry time.Time
	)

	for k, v := range s.storages.foundObjects {
		if index == "" {
			index = k
			timeExpiry = v.timeExpiry

			continue
		}

		if v.timeExpiry.Before(timeExpiry) {
			index = k
			timeExpiry = v.timeExpiry
		}
	}

	return index
}
