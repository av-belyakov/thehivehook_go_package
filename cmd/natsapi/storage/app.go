package storage

import (
	"context"
	"errors"
	"time"
)

// NewStorageAcceptedCommands инициализация хранилища команд
func NewStorageAcceptedCommands(opts ...cacheOptions) (*StorageAcceptedCommands, error) {
	storage := &StorageAcceptedCommands{
		//значение по умолчанию для интервала автоматической обработки
		timeTick: time.Duration(5 * time.Second),
		//значение по умолчанию для времени жизни объекта
		maxTtl: time.Duration(3600 * time.Second),
		//хранилище объектов
		storages: storage{
			objects: map[string]Object{},
		},
	}

	for _, opt := range opts {
		if err := opt(storage); err != nil {
			return storage, err
		}
	}

	return storage, nil
}

func (s *StorageAcceptedCommands) Start(ctx context.Context) {
	go func() {
		tick := time.NewTicker(s.timeTick)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				//поиск и удаление из хранилища всех объектов у которых истекло время жизни
				s.DeleteForTimeExpiryObjectFromCache()
			}
		}
	}()
}

//***************** методы для настройки хранилища ******************

// WithMaxTtl устанавливает максимальное время, по истечении которого запись в cacheStorages будет
// удалена, допустимый интервал времени хранения записи от 5 до 3600 секунд
func WithMaxTtl(v int) cacheOptions {
	return func(s *StorageAcceptedCommands) error {
		if v < 5 || v > 3600 {
			return errors.New("the maximum time after which an entry in the cache will be deleted should not be less than 300 seconds or more than 24 hours (86400 seconds)")
		}

		s.maxTtl = time.Duration(v) * time.Second

		return nil
	}
}

// WithTimeTick устанавливает интервал времени, заданное время такта, по истечении которого
// запускается новый виток автоматической обработки содержимого кэша, интервал значений должен
// быть в диапазоне от 1 до 120 секунд
func WithTimeTick(v int) cacheOptions {
	return func(s *StorageAcceptedCommands) error {
		if v < 1 || v > 120 {
			return errors.New("the set clock cycle time should not be less than 3 seconds or more than 120 seconds")
		}

		s.timeTick = time.Duration(v) * time.Second

		return nil
	}
}

// WithMaxSize устанавливает максимальный размер кэша, не может быть меньше 3 и больше 1000
func WithMaxSize(v int) cacheOptions {
	return func(s *StorageAcceptedCommands) error {
		if v < 3 || v > 1000 {
			return errors.New("the maximum cache size cannot be less than 3 or more than 1000 objects")
		}

		s.maxSize = v

		return nil
	}
}
