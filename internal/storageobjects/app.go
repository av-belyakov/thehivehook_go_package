package storageobjects

import (
	"errors"
	"time"
)

// NewStorageObjects инициализация хранилища объектов
func NewStorageObjects[T any](opts ...cacheOptions[T]) (*storageObjects[T], error) {
	storageObjects := &storageObjects[T]{
		timeTick:          time.Duration(1 * time.Second),    // 1 секунда
		timeToLive:        time.Duration(3600 * time.Second), // 1 час
		timeWaitingToSend: time.Duration(10 * time.Second),   // 10 секунд
		storage: objects[T]{
			objects: []object[T]{},
		},
		chanOutSize: 10,
	}

	for _, opt := range opts {
		if err := opt(storageObjects); err != nil {
			return nil, err
		}
	}

	storageObjects.chanOut = make(chan StorageObjectData[T], storageObjects.chanOutSize)

	return storageObjects, nil
}

// WithTimeTick устанавливает интервал времени, заданное время такта, по истечении которого
// запускается новый виток автоматической обработки содержимого хранилища,
// интервал значений должен быть в диапазоне от 1 до 180 секунд
func WithTimeTick[T any](v int) cacheOptions[T] {
	return func(s *storageObjects[T]) error {
		if v < 1 || v > 180 {
			return errors.New("the set clock cycle time should not be less than 1 seconds or more than 180 seconds")
		}

		s.timeTick = time.Duration(v) * time.Second

		return nil
	}
}

// WithTimeToLive устанавливает максимальное время жизни объекта, по истечении которого он будет удален,
// допустимый интервал от 10 до 43200 секунд (12 часов)
func WithTimeToLive[T any](v int) cacheOptions[T] {
	return func(s *storageObjects[T]) error {
		if v < 10 || v > 43200 {
			return errors.New("the maximum time after which an entry in the storage will be deleted should not be less than 10 seconds or more than 12 hours (43200 seconds)")
		}

		s.timeToLive = time.Duration(v) * time.Second

		return nil
	}
}

// WithTimeWaitingToSend устанавливает время по истечении которого объект из записи будет передан в канал,
// а сама запись будет удалена, допустимый интервал от 1 до 180 секунд (от 1 секунды до 3 минут)
func WithTimeWaitingToSend[T any](v int) cacheOptions[T] {
	return func(s *storageObjects[T]) error {
		if v < 1 || v > 180 {
			return errors.New("the time after which the object will be sent to the channel should not be less than 1 second or more than 180 seconds")
		}

		s.timeWaitingToSend = time.Duration(v) * time.Second

		return nil
	}
}

// WithChannelSize устанавливает размер буферизованного канала,
// допустимый интервал от 1 до 1000 элементов
func WithChannelSize[T any](v int) cacheOptions[T] {
	return func(s *storageObjects[T]) error {
		if v < 1 || v > 1_000 {
			return errors.New("the channel size cannot be less than 1 element or more than 1000 elements")
		}

		s.chanOutSize = v

		return nil
	}
}
