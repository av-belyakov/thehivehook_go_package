package storage

import (
	"fmt"
	"time"
)

// GetObjects все существующие объекты
func (s *StorageAcceptedCommands) GetObjects() map[string]Object {
	s.storages.mutex.RLock()
	defer s.storages.mutex.RUnlock()

	objects := make(map[string]Object, len(s.storages.objects))
	for k, v := range s.storages.objects {
		objects[k] = v
	}

	return objects
}

// GetObject получить объект по ключу
func (s *StorageAcceptedCommands) GetObject(key string) (Object, bool) {
	s.storages.mutex.RLock()
	defer s.storages.mutex.RUnlock()

	data, ok := s.storages.objects[key]

	return data, ok
}

// SetObject добавить объект по ключу
func (s *StorageAcceptedCommands) SetObject(key string, data []byte) error {
	s.storages.mutex.Lock()
	defer s.storages.mutex.Unlock()

	if len(s.storages.objects) >= s.maxSize {
		//удаляем самый старый объект
		if err := s.deleteOldestObjectFromCache(); err != nil {
			return err
		}
	}

	s.storages.objects[key] = Object{
		Data:       data,
		timeExpiry: time.Now().Add(s.maxTtl),
	}

	return nil
}

// DeleteForTimeExpiryObjectFromCache удаляет все объекты у которых истекло время жизни, без учета других параметров
func (s *StorageAcceptedCommands) DeleteForTimeExpiryObjectFromCache() {
	s.storages.mutex.Lock()
	defer s.storages.mutex.Unlock()

	for key, object := range s.storages.objects {
		if object.timeExpiry.Before(time.Now()) {
			delete(s.storages.objects, key)
		}
	}
}

// deleteOldestObjectFromCache поиск и удаление самого старого объекта
func (s *StorageAcceptedCommands) deleteOldestObjectFromCache() error {
	//получаем самый старый объект в кэше
	index := s.getOldestObjectFromCache()
	if _, ok := s.storages.objects[index]; ok {
		delete(s.storages.objects, index)
	} else {
		return fmt.Errorf("the object with id '%s' cannot be deleted, it may be in progress", index)
	}

	return nil
}

// getOldestObjectFromCache возвращает индекс самого старого объекта
func (s *StorageAcceptedCommands) getOldestObjectFromCache() string {
	var (
		index      string
		timeExpiry time.Time
	)

	for k, v := range s.storages.objects {
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
